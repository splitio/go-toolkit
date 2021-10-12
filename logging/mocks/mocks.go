package mocks

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
