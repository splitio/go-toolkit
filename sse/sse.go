package sse

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/splitio/go-toolkit/v4/logging"
)

const (
	statusIdle = iota
	statusRunning
	statusShuttingDown

	endOfLineChar = '\n'
	endOfLineStr  = "\n"
)

// Client struct
type Client struct {
	url              string
	client           http.Client
	shutdownRequest  chan struct{}
	timeout          time.Duration
	shutdownComplete *sync.Cond
	status           int32
	logger           logging.LoggerInterface
}

// NewClient creates new SSEClient
func NewClient(url string, timeout int, logger logging.LoggerInterface) (*Client, error) {
	if timeout < 1 {
		return nil, errors.New("Timeout should be higher than 0")
	}

	return &Client{
		url:              url,
		client:           http.Client{},
		shutdownRequest:  make(chan struct{}, 1),
		shutdownComplete: sync.NewCond(&sync.Mutex{}),
		status:           statusIdle,
		timeout:          time.Duration(timeout) * time.Second,
		logger:           logger,
	}, nil
}

func (l *Client) readEvents(in *bufio.Reader, out chan<- RawEvent) {
	eventBuilder := NewEventBuilder()
	for {
		line, err := in.ReadString(endOfLineChar)
		l.logger.Debug("Incoming SSE line: ", line)
		if err != nil {
			if atomic.LoadInt32(&l.status) == statusShuttingDown {
				l.logger.Error(err)
			}
			close(out)
			return
		}
		if line != endOfLineStr {
			eventBuilder.AddLine(line)
			continue

		}
		l.logger.Debug("Building SSE event")
		if event := eventBuilder.Build(); event != nil {
			out <- *event
		}
		eventBuilder.Reset()
	}
}

// Do starts streaming
func (l *Client) Do(params map[string]string, callback func(e RawEvent)) error {

	if !atomic.CompareAndSwapInt32(&l.status, statusIdle, statusRunning) {
		return ErrNotIdle
	}

	activeGoroutines := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		l.logger.Info("SSE streaming exiting")
		cancel()
		activeGoroutines.Wait()

		// In the rare case that .Shutdown() was called before the actual SSE connection was established,
		// we ensure that at the end of the method, the Shutdown channel is always cleared.
		// This is done prior to signiling the shutdown caller to avoid race conditions with new .Shutdown calls
		select {
		case <-l.shutdownRequest:
		default:
		}

		atomic.StoreInt32(&l.status, statusIdle)
		l.shutdownComplete.Broadcast()

	}()

	req, err := l.buildCancellableRequest(ctx, params)
	if err != nil {
		return &ErrConnectionFailed{wrapped: fmt.Errorf("error building request: %w", err)}
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return &ErrConnectionFailed{wrapped: fmt.Errorf("error issuing request: %w", err)}
	}
	if resp.StatusCode != 200 {
		return &ErrConnectionFailed{wrapped: fmt.Errorf("sse request status code: %d", resp.StatusCode)}
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	eventChannel := make(chan RawEvent, 1000)
	go l.readEvents(reader, eventChannel)

	// Create timeout timer in case SSE dont receive notifications or keepalive messages
	keepAliveTimer := time.NewTimer(l.timeout)
	defer keepAliveTimer.Stop()

	for {
		select {
		case <-l.shutdownRequest:
			l.logger.Info("Shutting down listener")
			return nil
			// The former `return` causes the following to be executed in this specific order:
			// - Shutdown the keep-alive timer
			// - Close SSE response body
			// - Cancel ongoing request
			// - Wait for any SSE currently being processed to finish
			// - Acknowledge end of function
		case event, ok := <-eventChannel:
			keepAliveTimer.Reset(l.timeout)
			if !ok {
				if atomic.LoadInt32(&l.status) == statusShuttingDown {
					return nil
				}
				return ErrReadingStream
			}

			if event.IsEmpty() {
				continue // don't forward empty/comment events
			}
			activeGoroutines.Add(1)
			go func() {
				defer activeGoroutines.Done()
				callback(event)
			}()
		case <-keepAliveTimer.C: // Timeout
			l.logger.Warning("SSE idle timeout. Restarting connection flow")
			atomic.StoreInt32(&l.status, statusShuttingDown)
			return ErrTimeout
		}
	}
}

// Shutdown stops SSE
func (l *Client) Shutdown(blocking bool) {
	if !atomic.CompareAndSwapInt32(&l.status, statusRunning, statusShuttingDown) {
		l.logger.Info("SSE client stopped or shutdown in progress. Ignoring.")
		return
	}

	l.shutdownRequest <- struct{}{}

	if blocking {
		for atomic.LoadInt32(&l.status) == statusIdle {
			l.shutdownComplete.L.Lock()
			l.shutdownComplete.Wait()
		}
	}
}

func (l *Client) buildCancellableRequest(ctx context.Context, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("GET", l.url, nil)
	if err != nil {
		return nil, fmt.Errorf("error instantiating request: %w", err)
	}
	req = req.WithContext(ctx)
	query := req.URL.Query()

	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "text/event-stream")
	return req, nil
}
