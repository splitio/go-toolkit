package mocks

import "time"

type BackoffMock struct {
	NextCall  func() time.Duration
	ResetCall func()
}

func (b *BackoffMock) Next() time.Duration {
	return b.NextCall()
}

func (b *BackoffMock) Reset() {
	b.ResetCall()
}
