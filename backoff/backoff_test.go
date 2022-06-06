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
