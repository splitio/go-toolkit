package common

import (
	"errors"
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
