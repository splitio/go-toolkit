package mocks

import (
	"github.com/splitio/go-toolkit/v5/backoff"
	"github.com/stretchr/testify/mock"
	"time"
)

type BackoffMock struct {
	mock.Mock
}

// Next implements backoff.Interface.
func (b *BackoffMock) Next() time.Duration {
    return b.Called().Get(0).(time.Duration)
}

// Reset implements backoff.Interface.
func (b *BackoffMock) Reset() {
    b.Called()
}

var _ backoff.Interface = (*BackoffMock)(nil)
