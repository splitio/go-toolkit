package sse

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/splitio/go-toolkit/logging"
)

func TestSSEError(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})
	clientErr, err := NewSSEClient("", make(chan int), make(chan struct{}), logger)
	if clientErr != nil {
		t.Error("It should be nil")
	}
	if err == nil || err.Error() != "Status channel should have length" {
		t.Error("It should return err")
	}

	status := make(chan int, 1)
	client, _ := NewSSEClient("", status, make(chan struct{}, 1), logger)
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
		client:   http.Client{},
		status:   status,
		stopped:  make(chan struct{}, 1),
		shutdown: make(chan struct{}, 1),
		url:      ts.URL,
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
		client:   http.Client{},
		logger:   logger,
		stopped:  make(chan struct{}, 1),
		shutdown: make(chan struct{}, 1),
		status:   status,
		url:      ts.URL,
	}

	var result map[string]interface{}
	go mockedClient.Do(make(map[string]string), func(e map[string]interface{}) {
		result = e
	})

	ready := <-status
	if ready != OK {
		t.Error("It should send ready flag")
	}

	mockedClient.Shutdown()

	if result["data"] != "some" {
		t.Error("Unexpected result")
	}
}

func TestSSEKeepAlive(t *testing.T) {
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
	mockedClient := SSEClient{
		client:   http.Client{},
		logger:   logger,
		stopped:  make(chan struct{}, 1),
		shutdown: make(chan struct{}, 1),
		url:      ts.URL,
		status:   status,
	}

	var result map[string]interface{}
	go mockedClient.Do(make(map[string]string), func(e map[string]interface{}) {
		result = e
	})

	ready := <-status
	if ready != OK {
		t.Error("It should send ready flag")
	}

	mockedClient.Shutdown()

	if result["event"] != "keepalive" {
		t.Error("Unexpected result")
	}
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
		url:      ts.URL,
		status:   status,
	}

	go mockedClient.Do(make(map[string]string), func(e map[string]interface{}) {
	})

	ready := <-status
	if ready != OK {
		t.Error("It should send ready flag")
	}

	go func() {
		mockedClient.Shutdown()
	}()

	msg := <-stopped

	if msg != struct{}{} {
		t.Error("It should receive stop")
	}
}
