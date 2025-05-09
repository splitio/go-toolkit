package backoff

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestBackoff(t *testing.T) {
	base := int64(10)
	maxAllowed := 60 * time.Second
	backoff := New(base, maxAllowed)
    assert.Equal(t, base, backoff.base)
    assert.Equal(t, maxAllowed, backoff.maxAllowed)
    assert.Equal(t, 1*time.Second, backoff.Next())
    assert.Equal(t, 10*time.Second, backoff.Next())
    assert.Equal(t, 60*time.Second, backoff.Next())

	backoff.Reset()
    assert.Equal(t, int64(0), backoff.current)
    assert.Equal(t, 1*time.Second, backoff.Next())
}
