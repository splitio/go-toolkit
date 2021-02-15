package mocks

type RawEventMock struct {
	IDCall      func() string
	EventCall   func() string
	DataCall    func() string
	RetryCall   func() int64
	IsErrorCall func() bool
	IsEmptyCall func() bool
}

func (r *RawEventMock) ID() string {
	return r.IDCall()
}

func (r *RawEventMock) Event() string {
	return r.EventCall()
}

func (r *RawEventMock) Data() string {
	return r.DataCall()
}

func (r *RawEventMock) Retry() int64 {
	return r.RetryCall()
}

func (r *RawEventMock) IsError() bool {
	return r.IsErrorCall()
}

func (r *RawEventMock) IsEmpty() bool {
	return r.IsEmptyCall()
}
