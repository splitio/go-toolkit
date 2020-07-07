package backoff

import (
	"errors"
	"testing"
	"time"

	"github.com/splitio/go-toolkit/logging"
)

func TestBackOff(t *testing.T) {
	count := 0
	logger := logging.NewLogger(&logging.LoggerOptions{LogLevel: logging.LevelDebug})
	perform := func(logger logging.LoggerInterface) (bool, error) {
		if count == 2 {
			return false, nil
		}
		count++
		return true, nil
	}
	backoff := NewBackOff("Test", perform, 1, 10, logger)
	go func() {
		backoff.Start()
	}()
	time.Sleep(100 * time.Millisecond)
	if !backoff.IsRunning() {
		t.Error("It should be running")
	}
	time.Sleep(7 * time.Second)
	if backoff.IsRunning() {
		t.Error("It should not be running")
	}
}

func TestBackOffShouldStopOnError(t *testing.T) {
	count := 0
	logger := logging.NewLogger(&logging.LoggerOptions{LogLevel: logging.LevelDebug})
	perform := func(logger logging.LoggerInterface) (bool, error) {
		if count == 2 {
			return false, errors.New("some")
		}
		count++
		return true, nil
	}
	backoff := NewBackOff("Test", perform, 1, 10, logger)
	go func() {
		backoff.Start()
	}()
	time.Sleep(100 * time.Millisecond)
	if !backoff.IsRunning() {
		t.Error("It should be running")
	}
	time.Sleep(7 * time.Second)
	if backoff.IsRunning() {
		t.Error("It should not be running")
	}
}

func TestBackoffMaxRetry(t *testing.T) {
	count := 0
	logger := logging.NewLogger(&logging.LoggerOptions{LogLevel: logging.LevelDebug})
	perform := func(logger logging.LoggerInterface) (bool, error) {
		if count == 2 {
			return false, errors.New("some")
		}
		count++
		return true, nil
	}
	backoff := NewBackOff("Test", perform, 1, 3, logger)
	go func() {
		backoff.Start()
	}()
	time.Sleep(100 * time.Millisecond)
	if !backoff.IsRunning() {
		t.Error("It should be running")
	}
	time.Sleep(5 * time.Second)
	if backoff.IsRunning() {
		t.Error("It should not be running")
	}
}

func TestBackOffShouldStop(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{LogLevel: logging.LevelDebug})
	perform := func(logger logging.LoggerInterface) (bool, error) {
		return true, nil
	}
	backoff := NewBackOff("Test", perform, 1, 10, logger)
	go func() {
		backoff.Start()
	}()
	time.Sleep(100 * time.Millisecond)
	if !backoff.IsRunning() {
		t.Error("It should be running")
	}
	time.Sleep(2 * time.Second)
	backoff.Stop(true)
	if backoff.IsRunning() {
		t.Error("It should not be running")
	}
}
