package common

import (
	"testing"
	"time"
)

func TestMinDuration(t *testing.T) {
	m := MinDuration(1*time.Second, 2*time.Second)
	if m.Seconds() != 1 {
		t.Error("Unexpected result")
	}

	m2 := MinDuration(2*time.Second, 1*time.Second)
	if m2.Seconds() != 1 {
		t.Error("Unexpected result")
	}

	m3 := MinDuration(2*time.Minute, 1*time.Second, 2*time.Millisecond)
	if m3.Milliseconds() != 2 {
		t.Error("Unexpected result")
	}
}

func TestMaxDuration(t *testing.T) {
	m := MaxDuration(1*time.Second, 2*time.Second)
	if m.Seconds() != 2 {
		t.Error("Unexpected result")
	}

	m2 := MaxDuration(2*time.Second, 1*time.Second)
	if m2.Seconds() != 2 {
		t.Error("Unexpected result")
	}

	m3 := MaxDuration(4*time.Minute, 1*time.Second, 2*time.Millisecond)
	if m3.Seconds() != 240 {
		t.Error("Unexpected result")
	}
}
