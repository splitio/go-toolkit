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
	maxAllowed time.Duration
	current    int64
}

// Next returns how long to wait and updates the current count
func (b *Impl) Next() time.Duration {
	current := atomic.LoadInt64(&b.current)
	nextWait := time.Duration(float64(b.base)*(math.Pow(2, float64(current)))) * time.Second
	atomic.AddInt64(&b.current, 1)
	if nextWait > b.maxAllowed {
		return b.maxAllowed
	}
	return nextWait
}

// Reset sets the current count to 0
func (b *Impl) Reset() {
	atomic.StoreInt64(&b.current, 0)
}

// New creates a new Backoffer
func New(base *int64, maxAllowed *time.Duration) *Impl {
	backoffBase := int64(2)
	backoffMaxAllowed := maxAllowedWait
	if base != nil && *base > 0 {
		backoffBase = *base
	}
	if maxAllowed != nil && *maxAllowed > 0 {
		backoffMaxAllowed = *maxAllowed
	}
	return &Impl{base: backoffBase, maxAllowed: backoffMaxAllowed}
}
