package mocks

import (
	"github.com/splitio/go-toolkit/v5/logging"

	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	ErrorCall   func(msg ...interface{})
	WarningCall func(msg ...interface{})
	InfoCall    func(msg ...interface{})
	DebugCall   func(msg ...interface{})
	VerboseCall func(msg ...interface{})
}

func (l *MockLogger) Error(msg ...interface{}) {
	l.ErrorCall(msg...)
}

func (l *MockLogger) Warning(msg ...interface{}) {
	l.WarningCall(msg...)
}

func (l *MockLogger) Info(msg ...interface{}) {
	l.InfoCall(msg...)
}

func (l *MockLogger) Debug(msg ...interface{}) {
	l.DebugCall(msg...)
}

func (l *MockLogger) Verbose(msg ...interface{}) {
	l.VerboseCall(msg...)
}

type LoggerMock struct {
	mock.Mock
}

func (l *LoggerMock) Debug(msg ...interface{}) {
	l.Called(msg)
}

func (l *LoggerMock) Info(msg ...interface{}) {
	l.Called(msg)
}

func (l *LoggerMock) Warning(msg ...interface{}) {
	l.Called(msg)
}

func (l *LoggerMock) Error(msg ...interface{}) {
	l.Called(msg)
}

func (l *LoggerMock) Verbose(msg ...interface{}) {
	l.Called(msg)
}

var _ logging.LoggerInterface = (*LoggerMock)(nil)
var _ logging.LoggerInterface = (*MockLogger)(nil)
