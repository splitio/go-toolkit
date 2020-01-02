package common

import (
	"errors"
)

func WithAttempts(attempts int, main func() error) error {
	err := errors.New("")
	remaining := attempts
	for err != nil && remaining > 0 {
		remaining--
		err = main()
	}
	return err
}
