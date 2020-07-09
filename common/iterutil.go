package common

import (
	"errors"
	"math"
	"time"
)

// WithAttempts executes a function N times or until no error is returned
func WithAttempts(attempts int, main func() error) error {
	err := errors.New("")
	remaining := attempts
	for err != nil && remaining > 0 {
		remaining--
		err = main()
	}
	return err
}

// WithBackoff wraps the function to add Exponential backoff
func WithBackoff(duration time.Duration, main func() error) func() error {
	var count time.Duration = 1
	return func() error {
		err := main()
		if err != nil {
			time.Sleep(count * duration)
			count++
		} else {
			count = 0
		}
		return main()
	}
}

// WithBackoffCancelling wraps the function to add Exponential backoff
func WithBackoffCancelling(unit time.Duration, max time.Duration, main func() bool) func() {
	cancel := make(chan struct{})
	go func() {
		attempts := 0
		isDone := main()
		for !isDone {
			attempts++
			select {
			case <-cancel:
				return
			case <-time.After(time.Duration(math.Min(math.Pow(2, float64(attempts)), max.Seconds())) * unit):
				isDone = main()
			}
		}
	}()
	return func() {
		select {
		case cancel <- struct{}{}:
			return
		default:
		}
	}
}
