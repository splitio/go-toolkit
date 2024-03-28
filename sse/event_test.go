package sse

import (
	"testing"

	"github.com/stretchr/testify/assert"

)

func TestEventBuilder(t *testing.T) {
	builder := NewEventBuilder()
	builder.AddLine("event: message")
	builder.AddLine("data: something")
	builder.AddLine("id: 1234")
	builder.AddLine("retry: 1")
	builder.AddLine(":some Comment")

	e := builder.Build()
    assert.Equal(t, "message", e.Event())
    assert.Equal(t, "something", e.Data())
    assert.Equal(t, "1234", e.ID())
    assert.Equal(t, int64(1), e.Retry())
    assert.False(t, e.IsEmpty())
    assert.False(t, e.IsEmpty())

	builder.Reset()
	builder.AddLine("event: error")
	builder.AddLine("data: someError")
	e2 := builder.Build()
    assert.True(t, e2.IsError())
}
