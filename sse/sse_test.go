package sse

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/splitio/go-toolkit/v6/logging"
	"github.com/stretchr/testify/assert"
)

func TestSSEErrorConnecting(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})
	client, _ := NewClient("", 120, 10, logger)
	err := client.Do(make(map[string]string), make(map[string]string), func(e RawEvent) { t.Error("It should not execute anything") })
	_, ok := err.(*ErrConnectionFailed)
    assert.True(t, ok)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}))
	defer ts.Close()

	mockedClient := Client{
		url:    ts.URL,
		client: http.Client{},
		logger: logger,
	}
	mockedClient.lifecycle.Setup()

	err = mockedClient.Do(make(map[string]string), make(map[string]string), func(e RawEvent) {
		assert.Fail(t, "Should not execute callback")
	})
	_, ok = err.(*ErrConnectionFailed)
    assert.True(t, ok)
}

func TestSSE(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "some", r.Header.Get("some"))
		flusher, ok := w.(http.Flusher)
        assert.True(t, ok)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		fmt.Fprintf(w, "data: %s\n\n", `{"id":"YCh53QfLxO:0:0","data":"some","timestamp":1591911770828}`)
		flusher.Flush()
		time.Sleep(2 * time.Second)
	}))
	defer ts.Close()

	mockedClient := Client{
		url:     ts.URL,
		client:  http.Client{},
		timeout: 30 * time.Second,
		logger:  logger,
	}
	mockedClient.lifecycle.Setup()

	var result RawEvent
	mutextTest := sync.RWMutex{}
	go func() {
		err := mockedClient.Do(nil, map[string]string{"some": "some"}, func(e RawEvent) {
			mutextTest.Lock()
			result = e
			mutextTest.Unlock()
		})
        assert.Nil(t, err)
	}()

	time.Sleep(2 * time.Second)
	mockedClient.Shutdown(true)

	mutextTest.RLock()
    assert.Equal(t, `{"id":"YCh53QfLxO:0:0","data":"some","timestamp":1591911770828}`, result.Data())
	mutextTest.RUnlock()
}

func TestSSENoTimeout(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	mutexTest := sync.RWMutex{}

	mutexTest.Lock()
	finished := false
	mutexTest.Unlock()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
        assert.True(t, ok)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		fmt.Fprintf(w, "data: %s\n\n", `{"id":"YCh53QfLxO:0:0","data":"some","timestamp":1591911770828}`)
		flusher.Flush()
		time.Sleep(2 * time.Second)
		mutexTest.Lock()
		finished = true
		mutexTest.Unlock()
	}))
	defer ts.Close()

	clientSSE, _ := NewClient(ts.URL, 70, 1, logger)

	go func() {
		clientSSE.Do(nil, make(map[string]string), func(e RawEvent) {})
	}()

	time.Sleep(1500 * time.Millisecond)
	mutexTest.RLock()
    assert.False(t, finished)
	mutexTest.RUnlock()
	time.Sleep(1500 * time.Millisecond)
	mutexTest.RLock()
    assert.True(t, finished)
	mutexTest.RUnlock()
	clientSSE.Shutdown(true)
}

func TestStopBlock(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
        assert.True(t, ok)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		fmt.Fprintf(w, ":keepalive")
		flusher.Flush()
		time.Sleep(2 * time.Second)
	}))
	defer ts.Close()

	mockedClient := Client{
		client:  http.Client{},
		logger:  logger,
		timeout: 30 * time.Second,
		url:     ts.URL,
	}
	mockedClient.lifecycle.Setup()

	waiter := make(chan struct{}, 1)
	go func() {
		err := mockedClient.Do(make(map[string]string), make(map[string]string), func(e RawEvent) {})
        assert.Nil(t, err)
		waiter <- struct{}{}
	}()

	time.Sleep(1 * time.Second)
	mockedClient.Shutdown(true)
	<-waiter
}

func TestConnectionEOF(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
        assert.True(t, ok)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		fmt.Fprintf(w, ":keepalive")
		flusher.Flush()
		ts.Listener.Close()
	}))
	defer ts.Close()

	mockedClient := Client{
		client:  http.Client{},
		logger:  logger,
		timeout: 30 * time.Second,
		url:     ts.URL,
	}
	mockedClient.lifecycle.Setup()

	err := mockedClient.Do(make(map[string]string), make(map[string]string), func(e RawEvent) {})
    assert.ErrorIs(t, err, ErrReadingStream)
	mockedClient.Shutdown(true)
}
