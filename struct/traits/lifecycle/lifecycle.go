package lifecycle

import (
	"sync"
	"sync/atomic"
)

// Status constants
const (
	StatusIdle = iota
	StatusStarting
	StatusRunning
	StatusStopping
)

// Status type alias
type Status = int32

// Manager is a trait to be embedded in structs that manage the lifecycle of goroutines.
// The trait enables the struct to easily switch between states and await proper shutdown
type Manager struct {
	status   int32
	c        *sync.Cond
	shutdown chan struct{}
}

// Setup must be called in the struct constructor
func (l *Manager) Setup() {
	l.c = sync.NewCond(&sync.Mutex{})
	l.shutdown = make(chan struct{}, 1)
}

// BeginInitialization should be called in the .Start() method (or whichever begins the async work)
func (l *Manager) BeginInitialization() bool {
	return atomic.CompareAndSwapInt32(&l.status, StatusIdle, StatusStarting)
}

// InitializationComplete should be called just prior to the `go ...` directive starting the async work
func (l *Manager) InitializationComplete() {
	atomic.StoreInt32(&l.status, StatusRunning)
}

// BeginShutdown should be called on the .Stop() method or whichever makes a request for the async work to stop
func (l *Manager) BeginShutdown() bool {
	if !atomic.CompareAndSwapInt32(&l.status, StatusRunning, StatusStopping) {
		return false
	}

	l.shutdown <- struct{}{}
	return true
}

// ShutdownComplete should be called just before the goroutine exits. (ie: it should be the FIRST deferred func)
func (l *Manager) ShutdownComplete() {
	// clean up status channel in case a Stop occurred while the task was exiting on its own
	select {
	case <-l.shutdown:
	default:
	}
	atomic.StoreInt32(&l.status, StatusIdle)
	l.c.Broadcast()
}

// AwaitShutdownComplete can be called in case you need to join against the goroutine's end
func (l *Manager) AwaitShutdownComplete() {
	for atomic.LoadInt32(&l.status) != StatusIdle {
		l.c.L.Lock()
		l.c.Wait()
		l.c.L.Unlock()
	}
}

// ShutdownRequested should be queried in a select statement, which should react by terminating the goroutine
func (l *Manager) ShutdownRequested() <-chan struct{} {
	return l.shutdown
}

// IsRunning returns true if the BG work is still going on
func (l *Manager) IsRunning() bool {
	return atomic.LoadInt32(&l.status) == StatusRunning
}

// AbnormalShutdown should be called when the goroutine exits without Stop being called.:w
func (l *Manager) AbnormalShutdown() {
	atomic.CompareAndSwapInt32(&l.status, StatusRunning, StatusStopping)
}
