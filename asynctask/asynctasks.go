package asynctask

import (
	"fmt"
	"github.com/splitio/go-toolkit/logging"
	"time"
)

// AsyncTask is a struct that wraps tasks that should run periodically and can be remotely stopped & started,
// as well as making it's status (running/stopped) available.
type AsyncTask struct {
	task        func(l logging.LoggerInterface) error
	name        string
	running     bool
	stopChannel chan bool
	period      int
	onInit      func(l logging.LoggerInterface) error
	onStop      func(l logging.LoggerInterface)
	logger      logging.LoggerInterface
}

func (t *AsyncTask) waitForInterrupt() bool {
	select {
	case <-t.stopChannel:
		return true
	case <-time.After(time.Second * time.Duration(t.period)):
		return false
	}
}

// Start initiates the task. It wraps the execution in a closure guarded by a call to recover() in order
// to prevent the main application from crashin if something goes wrong while the sdk interacts with the backend.
func (t *AsyncTask) Start() {

	if t.running {
		if t.logger != nil {
			t.logger.Warning("Task %s is already running. Aborting new execution.", t.name)
		}
		return
	}
	t.running = true

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.running = false
				if t.logger != nil {
					t.logger.Error(fmt.Sprintf(
						"AsyncTask %s is panicking! Delaying execution for %d seconds (1 period)",
						t.name,
						t.period,
					))
				}
				time.Sleep(time.Duration(t.period) * time.Second)
			}
		}()

		// If there's an initialization function, execute it
		if t.onInit != nil {
			err := t.onInit(t.logger)
			if err != nil {
				// If something goes wrong during initialization, abort.
				if t.logger != nil {
					t.logger.Error(err.Error())
				}
				return
			}
		}

		for t.running {
			err := t.task(t.logger)
			if err != nil && t.logger != nil {
				t.logger.Error(err.Error())
			}
			// waitForInterrupt will return true if the task was interrupted and should be aborted
			// false if the period has expired and the task is ready to run again
			t.running = !t.waitForInterrupt()
		}
		if t.onStop != nil {
			t.onStop(t.logger)
		}
	}()
}

// Stop prevents future executions of the task
func (t *AsyncTask) Stop() {
	select {
	case t.stopChannel <- true:
		return
	default:
		if t.logger != nil {
			t.logger.Error(fmt.Sprintf(
				"Cannot stop task %s. A stop signal has already been sent and is yet to be processed",
				t.name,
			))
		}
	}
}

// IsRunning returns true if the task is currently running
func (t *AsyncTask) IsRunning() bool {
	return t.running
}

// NewAsyncTask creates a new task and returns a pointer to it
func NewAsyncTask(
	name string,
	task func(l logging.LoggerInterface) error,
	period int,
	onInit func(l logging.LoggerInterface) error,
	onStop func(l logging.LoggerInterface),
	logger logging.LoggerInterface,
) *AsyncTask {
	t := AsyncTask{
		name:        name,
		task:        task,
		running:     false,
		period:      period,
		onInit:      onInit,
		onStop:      onStop,
		logger:      logger,
		stopChannel: make(chan bool, 1),
	}

	return &t
}
