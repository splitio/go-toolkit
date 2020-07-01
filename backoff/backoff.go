package backoff

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/splitio/go-toolkit/logging"
)

// BackOff struct
type BackOff struct {
	perform    func(l logging.LoggerInterface) (bool, error)
	name       string
	running    atomic.Value
	incoming   chan int
	period     int
	retry      atomic.Value
	logger     logging.LoggerInterface
	finished   bool
	finishChan chan struct{}
	mutex      sync.RWMutex
}

const (
	taskMessageStop = iota
)

func (t *BackOff) _running() bool {
	res, ok := t.running.Load().(bool)
	if !ok {
		t.logger.Error("Error parsing backoff task status flag")
		return false
	}
	return res
}

func (t *BackOff) isFinished() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.finished
}

func (t *BackOff) updateStatus(status bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.finished = status
}

// Start initiates the backoff.
func (t *BackOff) Start() {
	if t._running() {
		if t.logger != nil {
			t.logger.Warning("BackOff %s is already running. Aborting new execution.", t.name)
		}
		return
	}
	t.running.Store(true)

	go func() {
		defer func() {
			t.updateStatus(true)
			t.finishChan <- struct{}{}
		}()

		// Execution
		for t._running() {
			// Run the wrapped function and handle the returned error if any.
			shouldRetry, err := t.perform(t.logger)
			if err != nil && t.logger != nil {
				t.logger.Error(err.Error())
				t.Stop(false)
			} else {
				if shouldRetry {
					t.retry.Store(t.retry.Load().(int) + 1)
				} else {
					t.Stop(false)
				}
			}

			// Wait for either a timeout or an interruption (can be a stop signal)
			select {
			case msg := <-t.incoming:
				switch msg {
				case taskMessageStop:
					t.running.Store(false)
				}
			case <-time.After(time.Second * time.Duration(t.period*int(math.Pow(2, float64(t.retry.Load().(int)))))):
			}
		}

	}()
}

func (t *BackOff) sendSignal(signal int) error {
	select {
	case t.incoming <- signal:
		return nil
	default:
		return fmt.Errorf("Couldn't send message to task %s", t.name)
	}
}

// Stop executes onStop hook if any, blocks until its done (if blocking = true) and prevents future executions of the backoff.
func (t *BackOff) Stop(blocking bool) error {
	if t.isFinished() {
		// BackOff already stopped
		return nil
	}
	if err := t.sendSignal(taskMessageStop); err != nil {
		// If the signal couldnt be sent, return error!
		return err
	}

	if blocking {
		// If blocking was set to true, wait until an empty struct is pushed into the channel
		<-t.finishChan
	}
	return nil
}

// IsRunning returns true if the backoff is currently running
func (t *BackOff) IsRunning() bool {
	return t._running()
}

// NewBackOff creates new backoff task with retries
func NewBackOff(
	name string,
	perform func(l logging.LoggerInterface) (bool, error),
	period int,
	logger logging.LoggerInterface,
) *BackOff {
	t := BackOff{
		name:       name,
		perform:    perform,
		period:     period,
		logger:     logger,
		incoming:   make(chan int, 10),
		finishChan: make(chan struct{}, 1),
		finished:   false,
		mutex:      sync.RWMutex{},
	}
	t.retry.Store(0)
	t.running.Store(false)
	return &t
}
