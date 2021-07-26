package sse

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/splitio/go-toolkit/v5/logging"
)

func TestSSEErrorConnecting(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})
	client, _ := NewClient("", 120, 10, logger)
	err := client.Do(make(map[string]string), make(map[string]string), func(e RawEvent) { t.Error("It should not execute anything") })
	asErrConecting := &ErrConnectionFailed{}
	if !errors.As(err, &asErrConecting) {
		t.Errorf("Unexpected type of error: %+v", err)
	}

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
		t.Error("Should not execute callback")
	})
	if !errors.As(err, &asErrConecting) {
		t.Errorf("Unexpected type of error: %+v", err)
	}
}

func TestSSE(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("some") != "some" {
			t.Error("It should send header")
		}
		flusher, err := w.(http.Flusher)
		if !err {
			t.Error("Unexpected error")
			return
		}

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
		if err != nil {
			t.Error("sse client ended in error:", err)
		}
	}()

	time.Sleep(2 * time.Second)
	mockedClient.Shutdown(true)

	mutextTest.RLock()
	if result.Data() != `{"id":"YCh53QfLxO:0:0","data":"some","timestamp":1591911770828}` {
		t.Error("Unexpected result: ", result.Data())
	}
	mutextTest.RUnlock()
}

func TestSSENoTimeout(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{})

	mutexTest := sync.RWMutex{}

	mutexTest.Lock()
	finished := false
	mutexTest.Unlock()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, err := w.(http.Flusher)
		if !err {
			t.Error("Unexpected error")
			return
		}

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
	if finished {
		t.Error("It should not be finished")
	}
	mutexTest.RUnlock()
	time.Sleep(1500 * time.Millisecond)
	mutexTest.RLock()
	if !finished {
		t.Error("It should be finished")
	}
	mutexTest.RUnlock()
	clientSSE.Shutdown(true)
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
		if err != nil {
			t.Error("sse client ended in error: ", err)
		}
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
		flusher, err := w.(http.Flusher)
		if !err {
			t.Error("Unexpected error")
			return
		}

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
	if err != ErrReadingStream {
		t.Error("Should have triggered an ErrorReadingStreamError. Got: ", err)
	}

	mockedClient.Shutdown(true)
}

/*
func TestCustom(t *testing.T) {
	url := `https://streaming.split.io/event-stream`
	logger := logging.NewLogger(&logging.LoggerOptions{LogLevel: logging.LevelError, StandardLoggerFlags: log.Llongfile})
	client, _ := NewClient(url, 50, logger)

	ready := make(chan struct{})
	accessToken := ``
	channels := "NzM2MDI5Mzc0_MTgyNTg1MTgwNg==_splits,[?occupancy=metrics.publishers]control_pri,[?occupancy=metrics.publishers]control_sec"
	go func() {
		err := client.Do(
			map[string]string{
				"accessToken": accessToken,
				"v":           "1.1",
				"channel":     channels,
			},
			func(e RawEvent) {
				fmt.Printf("Event: %+v\n", e)
			})
		if err != nil {
			t.Error("sse error:", err)
		}
		ready <- struct{}{}
	}()
	time.Sleep(5 * time.Second)
	client.Shutdown(true)
	<-ready
	fmt.Println(1)
	go func() {
		err := client.Do(
			map[string]string{
				"accessToken": accessToken,
				"v":           "1.1",
				"channel":     channels,
			},
			func(e RawEvent) {
				fmt.Printf("Event: %+v\n", e)
			})
		if err != nil {
			t.Error("sse error:", err)
		}
		ready <- struct{}{}
	}()
	time.Sleep(5 * time.Second)
	client.Shutdown(true)
	<-ready
	fmt.Println(2)

}
*/
