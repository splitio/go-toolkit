package logging

const (
	// LevelError log level
	LevelError = iota

	// LevelWarning log level
	LevelWarning

	// LevelInfo log level
	LevelInfo

	// LevelDebug log level
	LevelDebug

	// LevelVerbose log level
	LevelVerbose
)

// LevelFilteredLoggerWrapper forwards log message to delegate if level is set higher than incoming message
type LevelFilteredLoggerWrapper struct {
	level    int
	delegate LoggerInterface
}

// Error forwards error logging messages
func (l *LevelFilteredLoggerWrapper) Error(is ...interface{}) {
	l.delegate.Error(is...)
}

// Warning forwards warning logging messages
func (l *LevelFilteredLoggerWrapper) Warning(is ...interface{}) {
	if l.level >= LevelWarning {
		l.delegate.Warning(is...)
	}
}

// Info forwards info logging messages
func (l *LevelFilteredLoggerWrapper) Info(is ...interface{}) {
	if l.level >= LevelInfo {
		l.delegate.Info(is...)
	}
}

// Debug forwards debug logging messages
func (l *LevelFilteredLoggerWrapper) Debug(is ...interface{}) {
	if l.level >= LevelDebug {
		l.delegate.Debug(is...)
	}
}

// Verbose forwards verbose logging messages
func (l *LevelFilteredLoggerWrapper) Verbose(is ...interface{}) {
	if l.level >= LevelVerbose {
		l.delegate.Verbose(is...)
	}
}
