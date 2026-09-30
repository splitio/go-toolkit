package backoff

import (
	"math"
	"sync/atomic"
	"time"
)

const (
	maxAllowedWait = 30 * 60 * time.Second // half an hour
)

// Interface is the backoff interface
type Interface interface {
	Next() time.Duration
	Reset()
}

// Impl implements the Backoff interface
type Impl struct {
	base       int64
	current    int64
	maxAllowed time.Duration
}

// Next returns how long to wait and updates the current count
func (b *Impl) Next() time.Duration {
	current := atomic.LoadInt64(&b.current)
	atomic.AddInt64(&b.current, 1)

	// Compare in float seconds before converting: time.Duration is an int64 of nanoseconds,
	// so a large exponent would overflow and wrap around to zero or a negative value.
	seconds := math.Pow(float64(b.base), float64(current))
	if seconds >= b.maxAllowed.Seconds() {
		return b.maxAllowed
	}
	return time.Duration(seconds * float64(time.Second))
}

// Reset sets the current count to 0
func (b *Impl) Reset() {
	atomic.StoreInt64(&b.current, 0)
}

// New creates a new Backoffer
func New(base int64, maxAllowed time.Duration) *Impl {
	backoffBase := int64(2)
	backoffMaxAllowed := maxAllowedWait
	if base > 0 {
		backoffBase = base
	}
	if maxAllowed.Seconds() > 0 {
		backoffMaxAllowed = maxAllowed
	}
	return &Impl{base: backoffBase, maxAllowed: backoffMaxAllowed}
}
