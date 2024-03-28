package mocks

import "github.com/stretchr/testify/mock"

type RawEventMock struct {
	mock.Mock
}

func (r *RawEventMock) ID() string {
	return r.Called().String(0)
}

func (r *RawEventMock) Event() string {
	return r.Called().String(0)
}

func (r *RawEventMock) Data() string {
	return r.Called().String(0)
}

func (r *RawEventMock) Retry() int64 {
	return r.Called().Get(0).(int64)
}

func (r *RawEventMock) IsError() bool {
	return r.Called().Bool(0)
}

func (r *RawEventMock) IsEmpty() bool {
	return r.Called().Bool(0)
}
