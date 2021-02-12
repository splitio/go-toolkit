package sse

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/splitio/go-toolkit/v3/logging"
	gtSync "github.com/splitio/go-toolkit/v3/sync"
)

// Client struct
type Client struct {
	url              string
	client           http.Client
	shutdown         chan struct{}
	timeout          time.Duration
	shutdownWaiter   sync.WaitGroup
	executing        *gtSync.AtomicBool
	shutdownExpected *gtSync.AtomicBool
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
		shutdown:         make(chan struct{}, 1),
		shutdownExpected: gtSync.NewAtomicBool(false),
		executing:        gtSync.NewAtomicBool(false),
		timeout:          time.Duration(timeout) * time.Second,
		logger:           logger,
	}, nil
}

// Do starts streaming
func (l *Client) Do(params map[string]string, callback func(e RawEvent)) error {

	if !l.executing.TestAndSet() {
		return ErrAlreadyRunning
	}
	l.shutdownExpected.Unset()
	l.shutdownWaiter.Add(1) // keep track of when this function finishes
	activeGoroutines := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		l.logger.Info("SSE streaming exiting")
		cancel()
		activeGoroutines.Wait()
		l.executing.Unset()
		l.shutdownWaiter.Done()
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

	l.shutdownWaiter.Add(1)
	go func() {
		defer l.shutdownWaiter.Done()
		eventBuilder := NewEventBuilder()
		for {
			line, err := reader.ReadString('\n')
			l.logger.Debug("Incoming SSE line: ", line)
			if err != nil {
				if !l.shutdownExpected.IsSet() {
					l.logger.Error(err)
				}
				close(eventChannel)
				return
			}
			if line != "\n" {
				eventBuilder.AddLine(line)
				continue

			}
			l.logger.Debug("Building SSE event")
			if event := eventBuilder.Build(); event != nil {
				eventChannel <- *event
			}
			eventBuilder.Reset()
		}
	}()

	// Create timeout timer in case SSE dont receive notifications or keepalive messages
	keepAliveTimer := time.NewTimer(l.timeout)
	defer keepAliveTimer.Stop()

	for {
		select {
		case <-l.shutdown:
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
				if l.shutdownExpected.IsSet() {
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
			l.shutdownExpected.Set()
			return ErrTimeout
		}
	}
}

// Shutdown stops SSE
func (l *Client) Shutdown(blocking bool) {
	if !l.shutdownExpected.TestAndSet() {
		l.logger.Info("SSE client stopped or shutdown in progress. Ignoring.")
		return
	}

	fmt.Println("shutdown")

	select {
	case l.shutdown <- struct{}{}:
	default:
		l.logger.Error("Shutdown already in progress")
	}

	if blocking {
		l.shutdownWaiter.Wait()
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
