package sse

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/splitio/go-toolkit/logging"
)

var sseDelimiter [2]byte = [...]byte{':', ' '}
var sseData [4]byte = [...]byte{'d', 'a', 't', 'a'}
var sseKeepAlive [10]byte = [...]byte{':', 'k', 'e', 'e', 'p', 'a', 'l', 'i', 'v', 'e'}

// SSEClient struct
type SSEClient struct {
	url      string
	client   http.Client
	sseReady chan struct{}
	shutdown chan struct{}
	mainWG   sync.WaitGroup
	logger   logging.LoggerInterface
}

// NewSSEClient creates new SSEClient
func NewSSEClient(url string, ready chan struct{}, logger logging.LoggerInterface) *SSEClient {
	return &SSEClient{
		url:      url,
		client:   http.Client{},
		sseReady: ready,
		shutdown: make(chan struct{}, 1),
		mainWG:   sync.WaitGroup{},
		logger:   logger,
	}
}

// Shutdown stops SSE
func (l *SSEClient) Shutdown() {
	select {
	case l.shutdown <- struct{}{}:
	default:
		l.logger.Error("Awaited unexpected event")
	}
	l.mainWG.Wait()
}

func parseData(raw []byte) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return nil, fmt.Errorf("error parsing json: %w", err)
	}
	return data, nil
}

func (l *SSEClient) readEvent(reader *bufio.Reader) (map[string]interface{}, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	if len(line) < 2 {
		return nil, nil
	}

	splitted := bytes.Split(line, sseDelimiter[:])

	if bytes.Contains(splitted[0], sseKeepAlive[:]) {
		data := make(map[string]interface{})
		data["event"] = string(sseKeepAlive[1:])
		return data, nil
	}

	if bytes.Compare(splitted[0], sseData[:]) != 0 {
		return nil, nil
	}

	raw := bytes.TrimSpace(splitted[1])
	data, err := parseData(raw)
	if err != nil {
		l.logger.Error("Error parsing event: ", err)
		return nil, nil
	}

	return data, nil
}

// Do starts streaming
func (l *SSEClient) Do(params map[string]string, callback func(e map[string]interface{})) error {
	l.mainWG.Add(1)
	defer l.mainWG.Done()

	req, err := http.NewRequest("GET", l.url, nil)
	if err != nil {
		return errors.New("Could not create client")
	}

	query := req.URL.Query()

	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "text/event-stream")

	resp, err := l.client.Do(req)
	if err != nil {
		return errors.New("Could not perform request")
	}
	if resp.StatusCode != 200 {
		return errors.New("Could not connect to streaming")
	}

	l.sseReady <- struct{}{}
	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	shouldKeepRunning := true
	activeGoroutines := sync.WaitGroup{}

	for shouldKeepRunning {
		select {
		case <-l.shutdown:
			l.logger.Info("Shutting down listener")
			shouldKeepRunning = false
			break
		default:
			event, err := l.readEvent(reader)
			if err != nil {
				l.logger.Error(err)
				l.Shutdown()
			}

			if event != nil {
				activeGoroutines.Add(1)
				go func() {
					callback(event)
					activeGoroutines.Done()
				}()
			}
		}
	}
	l.logger.Info("SSE streaming exiting")
	activeGoroutines.Wait()
	return nil
}
