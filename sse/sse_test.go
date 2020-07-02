package sse

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/splitio/go-toolkit/logging"
)

func TestSSEError(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})
	clientErr, err := NewSSEClient("", make(chan int), make(chan struct{}), 120, logger)
	if clientErr != nil {
		t.Error("It should be nil")
	}
	if err == nil || err.Error() != "Status channel should have length" {
		t.Error("It should return err")
	}

	status := make(chan int, 1)
	client, _ := NewSSEClient("", status, make(chan struct{}, 1), 120, logger)
	client.Do(make(map[string]string), func(e map[string]interface{}) { t.Error("It should not execute anything") })

	stats := <-status
	if stats != ErrorRequestPerformed {
		t.Error("Unexpected type of error")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}))
	defer ts.Close()

	mockedClient := SSEClient{
		url:      ts.URL,
		client:   http.Client{},
		status:   status,
		stopped:  make(chan struct{}, 1),
		shutdown: make(chan struct{}, 1),
		logger:   logger,
	}

	mockedClient.Do(make(map[string]string), func(e map[string]interface{}) {
		t.Error("Should not execute callback")
	})

	stats = <-status
	if stats != ErrorConnectToStreaming {
		t.Error("Unexpected type of error")
	}
}

func TestSSE(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, err := w.(http.Flusher)
		if !err {
			t.Error("Unexpected error")
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		fmt.Fprintf(w, "data: %s\n\n", "{\"id\":\"YCh53QfLxO:0:0\",\"data\":\"some\",\"timestamp\":1591911770828}")
		flusher.Flush()
	}))
	defer ts.Close()

	status := make(chan int, 1)
	mockedClient := SSEClient{
		url:      ts.URL,
		client:   http.Client{},
		status:   status,
		stopped:  make(chan struct{}, 1),
		shutdown: make(chan struct{}, 1),
		timeout:  30,
		logger:   logger,
	}

	var result map[string]interface{}
	mutextTest := sync.RWMutex{}
	go mockedClient.Do(make(map[string]string), func(e map[string]interface{}) {
		mutextTest.Lock()
		result = e
		mutextTest.Unlock()
	})

	ready := <-status
	if ready != OK {
		t.Error("It should send ready flag")
	}

	time.Sleep(900 * time.Millisecond)
	mockedClient.Shutdown()
	time.Sleep(900 * time.Millisecond)

	mutextTest.RLock()
	if result["data"] != "some" {
		t.Error("Unexpected result")
	}
	mutextTest.RUnlock()
}

func TestStopBlock(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, err := w.(http.Flusher)
		if !err {
			t.Error("Unexpected error")
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		fmt.Fprintf(w, ":keepalive")
		flusher.Flush()
	}))
	defer ts.Close()

	status := make(chan int, 1)
	stopped := make(chan struct{}, 1)
	mockedClient := SSEClient{
		client:   http.Client{},
		logger:   logger,
		stopped:  stopped,
		shutdown: make(chan struct{}, 1),
		timeout:  30,
		url:      ts.URL,
		status:   status,
	}

	go mockedClient.Do(make(map[string]string), func(e map[string]interface{}) {
	})

	ready := <-status
	if ready != OK {
		t.Error("It should send ready flag")
	}

	var msg struct{}
	mutextTest := sync.RWMutex{}
	go func() {
		defer mutextTest.Unlock()
		mutextTest.Lock()
		msg = <-stopped
	}()

	mockedClient.Shutdown()

	mutextTest.RLock()
	if msg != struct{}{} {
		t.Error("It should receive stop")
	}
	mutextTest.RUnlock()
}
