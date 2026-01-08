package sse

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/splitio/go-toolkit/v5/logging"
	"github.com/splitio/go-toolkit/v5/struct/traits/lifecycle"
)

const (
	endOfLineChar = '\n'
	endOfLineStr  = "\n"
)

// Client struct
type Client struct {
	lifecycle lifecycle.Manager
	url       string
	client    http.Client
	timeout   time.Duration
	logger    logging.LoggerInterface
	bodyMu    sync.Mutex
	body      io.ReadCloser
	cancel    context.CancelFunc
}

// NewClient creates new SSEClient
func NewClient(url string, keepAlive int, dialTimeout int, logger logging.LoggerInterface) (*Client, error) {
	if keepAlive < 1 {
		return nil, errors.New("keepAlive timeout should be higher than 0")
	}
	if dialTimeout < 0 {
		dialTimeout = 0
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyFromEnvironment

	client := &Client{
		url:     url,
		client:  http.Client{Transport: transport},
		timeout: time.Duration(keepAlive) * time.Second,
		logger:  logger,
	}
	client.lifecycle.Setup()
	return client, nil
}

func (l *Client) readEvents(ctx context.Context, in *bufio.Reader, out chan<- RawEvent) {
	eventBuilder := NewEventBuilder()
	defer close(out)
	defer l.logger.Info("SSE reader goroutine exited")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := in.ReadString(endOfLineChar)
			l.logger.Debug("Incoming SSE line: ", line)
			if err != nil {
				if l.lifecycle.IsRunning() {
					l.logger.Error(err)
				}
				return
			}

			if line != endOfLineStr {
				eventBuilder.AddLine(line)
				continue
			}

			if event := eventBuilder.Build(); event != nil {
				out <- event
			}
			eventBuilder.Reset()
		}
	}
}

// Do starts streaming
func (l *Client) Do(params map[string]string, headers map[string]string, callback func(e RawEvent)) error {

	if !l.lifecycle.BeginInitialization() {
		return ErrNotIdle
	}

	var activeGoroutines sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	l.bodyMu.Lock()
	l.cancel = cancel
	l.bodyMu.Unlock()

	defer func() {
		l.logger.Info("SSE streaming exiting")

		cancel()

		l.bodyMu.Lock()
		l.cancel = nil
		l.bodyMu.Unlock()

		activeGoroutines.Wait()
		l.lifecycle.ShutdownComplete()
	}()

	req, err := l.buildCancellableRequest(ctx, params, headers)
	if err != nil {
		return &ErrConnectionFailed{wrapped: fmt.Errorf("error building request: %w", err)}
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return &ErrConnectionFailed{wrapped: fmt.Errorf("error issuing request: %w", err)}
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return &ErrConnectionFailed{
			wrapped: fmt.Errorf("sse request status code: %d", resp.StatusCode),
		}
	}

	l.bodyMu.Lock()
	l.body = resp.Body
	l.bodyMu.Unlock()

	if !l.lifecycle.InitializationComplete() {
		return nil
	}

	reader := bufio.NewReader(resp.Body)
	eventChannel := make(chan RawEvent, 1000)

	activeGoroutines.Add(1)
	go func() {
		defer activeGoroutines.Done()
		l.readEvents(ctx, reader, eventChannel)
	}()

	keepAliveTimer := time.NewTimer(l.timeout)
	defer keepAliveTimer.Stop()

	for {
		select {
		case <-l.lifecycle.ShutdownRequested():
			return nil

		case event, ok := <-eventChannel:
			keepAliveTimer.Reset(l.timeout)

			if !ok {
				if l.lifecycle.IsRunning() {
					return ErrReadingStream
				}
				return nil
			}

			if event.IsEmpty() {
				continue
			}

			activeGoroutines.Add(1)
			go func(ev RawEvent) {
				defer activeGoroutines.Done()
				callback(ev)
			}(event)

		case <-keepAliveTimer.C:
			l.lifecycle.AbnormalShutdown()
			return ErrTimeout
		}
	}
}

// Shutdown stops SSE
func (l *Client) Shutdown(blocking bool) {
	if !l.lifecycle.BeginShutdown() {
		l.logger.Info("SSE client stopped or shutdown in progress. Ignoring.")
		return
	}

	l.bodyMu.Lock()
	if l.cancel != nil {
		l.cancel()
		l.cancel = nil
	}
	if l.body != nil {
		_ = l.body.Close()
		l.body = nil
	}
	l.bodyMu.Unlock()

	if blocking {
		l.lifecycle.AwaitShutdownComplete()
	}
}

func (l *Client) buildCancellableRequest(ctx context.Context, params map[string]string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("GET", l.url, nil)
	if err != nil {
		return nil, fmt.Errorf("error instantiating request: %w", err)
	}
	req = req.WithContext(ctx)
	query := req.URL.Query()

	for key, value := range params {
		query.Add(key, value)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "text/event-stream")
	return req, nil
}
