package sse

import (
	"testing"
)

func TestEventBuilder(t *testing.T) {
	builder := NewEventBuilder()
	builder.AddLine("event: message")
	builder.AddLine("data: something")
	builder.AddLine("id: 1234")
	builder.AddLine("retry: 1")
	builder.AddLine(":some Comment")

	e := builder.Build()
	if e.Event() != "message" {
		t.Error("event should be 'message'")
	}
	if e.Data() != "something" {
		t.Error("data should be 'something'")
	}
	if e.ID() != "1234" {
		t.Error("Id should be 1234")
	}
	if e.Retry() != 1 {
		t.Error("retry should be 1234")
	}
	if e.IsEmpty() {
		t.Error("event should not be empty")
	}
	if e.IsError() {
		t.Error("event is not an error")
	}

	builder.Reset()
	builder.AddLine("event: error")
	builder.AddLine("data: someError")
	e2 := builder.Build()
	if !e2.IsError() {
		t.Error("event is an error")
	}
}
