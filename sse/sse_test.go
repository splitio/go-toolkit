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
	"github.com/stretchr/testify/require"
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

type fakeRawEvent struct {
	id int
}

func (f fakeRawEvent) ID() string    { return fmt.Sprintf("%d", f.id) }
func (f fakeRawEvent) Event() string { return "test" }
func (f fakeRawEvent) Data() string  { return "data" }
func (f fakeRawEvent) Retry() int64  { return 0 }
func (f fakeRawEvent) IsError() bool { return false }
func (f fakeRawEvent) IsEmpty() bool { return false }

func TestProcessEventsClosureBugWithInterface(t *testing.T) {
	const n = 200

	events := make([]RawEvent, n)
	for i := 0; i < n; i++ {
		events[i] = fakeRawEvent{id: i}
	}

	received := make([]string, 0, n)
	var mu sync.Mutex

	processEventsBug(events, func(e RawEvent) {
		mu.Lock()
		received = append(received, e.ID())
		mu.Unlock()
	})

	if len(received) != n {
		t.Fatalf("expected %d events, got %d", n, len(received))
	}

	unique := map[string]bool{}
	for _, id := range received {
		unique[id] = true
	}

	if len(unique) != n {
		t.Fatalf(
			"expected %d unique events, got %d (closure bug exposed)",
			n,
			len(unique),
		)
	}
}

func processEventsBug(events []RawEvent, callback func(RawEvent)) {
	var wg sync.WaitGroup

	for _, event := range events {
		wg.Add(1)
		go func(ev RawEvent) {
			defer wg.Done()
			callback(ev)
		}(event)
	}

	wg.Wait()
}

func TestShutdownDoesNotHangWhenSSEIsIdle(t *testing.T) {
	// Fake SSE server: accepts connection, sends headers, then blocks forever
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		require.True(t, ok)
		flusher.Flush()

		// Block until client closes the connection
		<-r.Context().Done()
	}))
	defer server.Close()

	logger := logging.NewLogger(nil)

	client, err := NewClient(
		server.URL,
		70, // keepAlive
		0,  // dialTimeout
		logger,
	)
	require.NoError(t, err)

	done := make(chan struct{})

	// Start streaming
	go func() {
		_ = client.Do(
			map[string]string{"channels": "test"},
			nil,
			func(e RawEvent) {},
		)
		close(done)
	}()

	// Give the client time to connect and block on read
	time.Sleep(100 * time.Millisecond)

	shutdownDone := make(chan struct{})

	go func() {
		client.Shutdown(true)
		close(shutdownDone)
	}()

	select {
	case <-shutdownDone:
		// OK
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Shutdown(true) blocked — SSE reader did not exit")
	}

	// Ensure Do() also returns
	select {
	case <-done:
		// OK
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Do() did not return after shutdown")
	}
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
		err := client.Do
(
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
