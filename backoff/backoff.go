package backoff

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/splitio/go-toolkit/logging"
)

// BackOff struct
type BackOff struct {
	max        float64
	perform    func(l logging.LoggerInterface) bool
	name       string
	incoming   chan int
	finishChan chan struct{}
	period     int
	retry      atomic.Value
	running    atomic.Value
	finished   atomic.Value
	logger     logging.LoggerInterface
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

// Start initiates the backoff.
func (t *BackOff) Start() {
	if len(t.finishChan) > 0 {
		// Discarding finished signal in case it was pending before
		<-t.finishChan
	}

	if t._running() {
		if t.logger != nil {
			t.logger.Warning("BackOff %s is already running. Aborting new execution.", t.name)
		}
		return
	}
	// Reseting configurations
	t.running.Store(true)
	t.finished.Store(false)
	t.retry.Store(0)

	go func() {
		defer func() {
			t.running.Store(false)
			t.finished.Store(true)
			t.finishChan <- struct{}{}
		}()

		// Execution
		for t._running() {
			// Run the wrapped function and handle the returned error if any.
			shouldRetry := t.perform(t.logger)
			if !shouldRetry {
				t.Stop(false)
				return
			}
			t.retry.Store(t.retry.Load().(int) + 1)

			// Wait for either a timeout or an interruption (can be a stop signal)
			select {
			case msg := <-t.incoming:
				switch msg {
				case taskMessageStop:
					return
				}
			case <-time.After(time.Second * time.Duration(t.period*int(math.Min(math.Pow(2, float64(t.retry.Load().(int))), t.max)))):
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
	if t.finished.Load().(bool) {
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
	perform func(l logging.LoggerInterface) bool,
	period int,
	max float64,
	logger logging.LoggerInterface,
) *BackOff {
	t := BackOff{
		max:        max,
		name:       name,
		perform:    perform,
		period:     period,
		logger:     logger,
		incoming:   make(chan int, 10),
		finishChan: make(chan struct{}, 1),
	}
	t.running.Store(false)
	return &t
}
