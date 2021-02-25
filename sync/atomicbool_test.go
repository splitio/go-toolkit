package sync

import (
	"testing"
)

func TestAtomicBool(t *testing.T) {
	a := NewAtomicBool(false)
	if a.IsSet() {
		t.Error("initial value should be false")
	}

	if !a.TestAndSet() {
		t.Error("compare and swap should succeed with no other concurrent access.")
	}

	if a.TestAndSet() {
		t.Error("compare and swap should return false if it didn't change anything.")
	}

	if !a.IsSet() {
		t.Error("should now be true")
	}

	b := NewAtomicBool(true)
	if !b.IsSet() {
		t.Error("initial value should be true")
	}

	if b.TestAndClear() != true {
		t.Error("compare and swap should succeed with no other concurrent access.")
	}

	if !a.TestAndClear() {
		t.Error("compare and swap should return false if it didn't change anything.")
	}

	if b.IsSet() {
		t.Error("should now be false")
	}
}
