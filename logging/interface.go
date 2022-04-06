package logging

// LoggerInterface ...
// If a custom logger object is to be used, it should comply with the following
// interface. (Standard go-lang library log.Logger.Println method signature)
type LoggerInterface interface {
	Error(msg ...interface{})
	Warning(msg ...interface{})
	Info(msg ...interface{})
	Debug(msg ...interface{})
	Verbose(msg ...interface{})
}

// paramsFn is a function that returns a slice of interface{}
type paramsFn = func() []interface{}

// il alias for interface list = []interface{}
type il = []interface{}

// ExtendedLoggerInterface ...
// If a custom logger object is to be used, it should comply with the following
// interface. (Standard go-lang library log.Logger.Println method signature)
type ExtendedLoggerInterface interface {
	LoggerInterface
	ErrorFn(format string, params paramsFn)
	WarningFn(format string, params paramsFn)
	InfoFn(format string, params paramsFn)
	DebugFn(format string, params paramsFn)
	VerboseFn(format string, params paramsFn)
}
