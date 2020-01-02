package common

import (
	"errors"
	"testing"
)

func TestWithAttempts(t *testing.T) {
	usedAttempts := 0
	err := WithAttempts(3, func() error { usedAttempts++; return nil })
	if err != nil {
		t.Error("Func Should have returned nil.")
	}
	if usedAttempts != 1 {
		t.Errorf("Func should have succeeded after 1 attempts. It took %d", usedAttempts)

	}

	usedAttempts = 0
	err = WithAttempts(3, func() error { usedAttempts++; return errors.New("someError") })
	if err == nil {
		t.Error("Func Should NOT have returned nil.")
	}
	if usedAttempts != 3 {
		t.Errorf("Func should have failed after 3 attempts. It took %d", usedAttempts)

	}

	usedAttempts = 0
	err = WithAttempts(3, func() error {
		usedAttempts++
		if usedAttempts != 3 {
			return errors.New("someError")
		} else {
			return nil
		}
	})
	if err != nil {
		t.Error("Func Should have returned nil.")
	}
	if usedAttempts != 3 {
		t.Errorf("Func should have succeeded after 3 attempts. It took %d", usedAttempts)
	}
}
