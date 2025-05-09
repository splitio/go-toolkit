package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicBool(t *testing.T) {
	a := NewAtomicBool(false)
    assert.False(t, a.IsSet())
    assert.True(t, a.TestAndSet())
    assert.False(t, a.TestAndSet())
    assert.True(t, a.IsSet())

	b := NewAtomicBool(true)
    assert.True(t, b.IsSet())
    assert.True(t, b.TestAndClear())
	assert.False(t, b.TestAndClear())
    assert.False(t, b.IsSet())
}
