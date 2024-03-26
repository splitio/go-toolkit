package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

// Debug implements logging.LoggerInterface.
func (l *MockLogger) Debug(msg ...interface{}) {
    l.Called(msg...)
}

// Debugf implements logging.LoggerInterface.
func (l *MockLogger) Debugf(fmt string, msg ...interface{}) {
    l.Called(append([]interface{}{fmt}, msg...)...)
}

// Error implements logging.LoggerInterface.
func (l *MockLogger) Error(msg ...interface{}) {
    l.Called(msg...)
}

// Errorf implements logging.LoggerInterface.
func (l *MockLogger) Errorf(fmt string, msg ...interface{}) {
    l.Called(append([]interface{}{fmt}, msg...)...)
}

// Info implements logging.LoggerInterface.
func (l *MockLogger) Info(msg ...interface{}) {
    l.Called(msg...)
}

// Infof implements logging.LoggerInterface.
func (l *MockLogger) Infof(fmt string, msg ...interface{}) {
    l.Called(append([]interface{}{fmt}, msg...)...)
}

// Verbose implements logging.LoggerInterface.
func (l *MockLogger) Verbose(msg ...interface{}) {
    l.Called(msg...)
}

// Verbosef implements logging.LoggerInterface.
func (l *MockLogger) Verbosef(fmt string, msg ...interface{}) {
    l.Called(append([]interface{}{fmt}, msg...)...)
}

// Warning implements logging.LoggerInterface.
func (l *MockLogger) Warning(msg ...interface{}) {
    l.Called(msg...)
}

// Warningf implements logging.LoggerInterface.
func (l *MockLogger) Warningf(fmt string, msg ...interface{}) {
    l.Called(append([]interface{}{fmt}, msg...)...)
}
