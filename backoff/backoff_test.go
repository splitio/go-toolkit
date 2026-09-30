package backoff

import (
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	base := int64(10)
	maxAllowed := 60 * time.Second
	backoff := New(base, maxAllowed)
	if backoff.base != base {
		t.Error("It should be equals to 10")
	}
	if backoff.maxAllowed != maxAllowed {
		t.Error("It should be equals to 60")
	}
	if backoff.Next() != 1*time.Second {
		t.Error("It should be 1 second")
	}
	if backoff.Next() != 10*time.Second {
		t.Error("It should be 10 seconds")
	}
	if backoff.Next() != 60*time.Second {
		t.Error("It should be 60 seconds")
	}
	backoff.Reset()
	if backoff.current != 0 {
		t.Error("It should be zero")
	}
	if backoff.Next() != 1*time.Second {
		t.Error("It should be 1 second")
	}
}

func TestBackoffDoesNotOverflow(t *testing.T) {
	maxAllowed := 30 * time.Minute
	b := New(2, maxAllowed)
	prev := time.Duration(0)
	for i := 0; i < 10000; i++ {
		got := b.Next()
		if got <= 0 {
			t.Fatalf("attempt %d: wait must be positive, got %v", i, got)
		}
		if got > maxAllowed {
			t.Fatalf("attempt %d: wait must not exceed %v, got %v", i, maxAllowed, got)
		}
		if got < prev {
			t.Fatalf("attempt %d: wait must be non-decreasing, got %v after %v", i, got, prev)
		}
		prev = got
	}
	if prev != maxAllowed {
		t.Errorf("wait should settle at the max, got %v", prev)
	}
}
