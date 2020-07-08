package backoff

import (
	"testing"
	"time"

	"github.com/splitio/go-toolkit/logging"
)

func TestBackOff(t *testing.T) {
	count := 0
	logger := logging.NewLogger(&logging.LoggerOptions{LogLevel: logging.LevelDebug})
	perform := func(logger logging.LoggerInterface) bool {
		if count == 2 {
			return false
		}
		count++
		return true
	}
	backoff := NewBackOff("Test", perform, 1, 10, logger)
	backoff.Start()
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
	perform := func(logger logging.LoggerInterface) bool {
		if count == 2 {
			return false
		}
		count++
		return true
	}
	backoff := NewBackOff("Test", perform, 1, 3, logger)
	backoff.Start()
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
	perform := func(logger logging.LoggerInterface) bool {
		return true
	}
	backoff := NewBackOff("Test", perform, 1, 10, logger)
	backoff.Start()
	time.Sleep(100 * time.Millisecond)
	if !backoff.IsRunning() {
		t.Error("It should be running")
	}
	time.Sleep(100 * time.Millisecond)
	backoff.Stop(true)
	if backoff.IsRunning() {
		t.Error("It should not be running")
	}
}

func TestBackOffStartingAgain(t *testing.T) {
	logger := logging.NewLogger(&logging.LoggerOptions{LogLevel: logging.LevelDebug})
	perform := func(logger logging.LoggerInterface) bool {
		return true
	}
	backoff := NewBackOff("Test", perform, 1, 10, logger)
	backoff.Start()
	time.Sleep(100 * time.Millisecond)
	if !backoff.IsRunning() {
		t.Error("It should be running")
	}
	time.Sleep(100 * time.Millisecond)
	backoff.Stop(false)
	time.Sleep(100 * time.Millisecond)
	if backoff.IsRunning() {
		t.Error("It should not be running")
	}
	backoff.Start()
	if !backoff.IsRunning() {
		t.Error("It should be running")
	}
	time.Sleep(100 * time.Millisecond)
	backoff.Stop(false)
	time.Sleep(100 * time.Millisecond)
	if backoff.IsRunning() {
		t.Error("It should not be running")
	}
}
