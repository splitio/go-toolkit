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

	sseError := make(chan error, 1)
	client := NewSSEClient("", make(chan struct{}, 1), sseError, logger)
	client.Do(make(map[string]string), func(e map[string]interface{}) { t.Error("It should not execute anything") })

	err := <-sseError
	if err == nil || err.Error() != "Could not perform request" {
		t.Error("It should not be nil")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}))
	defer ts.Close()

	mockedClient := SSEClient{
		client:   http.Client{},
		logger:   logger,
		mainWG:   sync.WaitGroup{},
		shutdown: make(chan struct{}, 1),
		sseReady: make(chan struct{}, 1),
		sseError: sseError,
		url:      ts.URL,
	}

	mockedClient.Do(make(map[string]string), func(e map[string]interface{}) {
		t.Error("Should not execute callback")
	})

	err = <-sseError
	if err == nil || err.Error() != "Could not connect to streaming" {
		t.Error("Unexpected error")
	}
}

func TestSSE(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	mockedClient := SSEClient{
		client:   http.Client{},
		logger:   logger,
		mainWG:   sync.WaitGroup{},
		shutdown: make(chan struct{}, 1),
		sseReady: make(chan struct{}, 1),
		sseError: make(chan error, 1),
	}

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

		go func() {
			time.Sleep(100 * time.Millisecond)
			mockedClient.Shutdown()
		}()
	}))
	defer ts.Close()

	mockedClient.url = ts.URL

	var result map[string]interface{}
	mockedClient.Do(make(map[string]string), func(e map[string]interface{}) {
		result = e
	})

	if result["data"] != "some" {
		t.Error("Unexpected result")
	}
}

func TestSSEKeepAlive(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	mockedClient := SSEClient{
		client:   http.Client{},
		logger:   logger,
		mainWG:   sync.WaitGroup{},
		shutdown: make(chan struct{}, 1),
		sseReady: make(chan struct{}, 1),
		sseError: make(chan error, 1),
	}

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

		go func() {
			time.Sleep(100 * time.Millisecond)
			mockedClient.Shutdown()
		}()
	}))
	defer ts.Close()

	mockedClient.url = ts.URL

	var result map[string]interface{}
	mockedClient.Do(make(map[string]string), func(e map[string]interface{}) {
		result = e
	})

	if result["event"] != "keepalive" {
		t.Error("Unexpected result")
	}
}
