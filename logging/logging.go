// Package logging ...
// Handles logging within the SDK
package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	skipStackFrameBase = 3 // How many stack frames to skip when logging filename
)

// LoggerOptions ...
// Struct that must be passed to the NewLogger constructor to setup a logger
// CommonWriter and ErrorWriter can be <nil>. In that case they'll default to os.Stdout
type LoggerOptions struct {
	LogLevel            int
	ErrorWriter         io.Writer
	WarningWriter       io.Writer
	InfoWriter          io.Writer
	DebugWriter         io.Writer
	VerboseWriter       io.Writer
	StandardLoggerFlags int
	Prefix              string
	ExtraFramesToSkip   int
}

// Logger struct. Encapsulates four different loggers, each for a different "level",
// and provides Error, Debug, Warning and Info functions, that will forward a message
// to the appropriate logger.
type Logger struct {
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	verboseLogger *log.Logger
	framesToSkip  int
	level         int

	contextData   *ContextData
	stringContext string
}

func (l *Logger) appendContext(msg ...interface{}) []interface{} {
	if l.contextData != nil {
		msg = append([]interface{}{l.stringContext}, msg...)
	}
	return msg
}

// Verbose logs a message with Debug level
func (l *Logger) Verbose(msg ...interface{}) {
	l.verboseLogger.Output(l.framesToSkip, fmt.Sprintln(l.appendContext(msg...)...))
}

// Debug logs a message with Debug level
func (l *Logger) Debug(msg ...interface{}) {
	l.debugLogger.Output(l.framesToSkip, fmt.Sprintln(l.appendContext(msg...)...))
}

// Info logs a message with Info level
func (l *Logger) Info(msg ...interface{}) {
	l.infoLogger.Output(l.framesToSkip, fmt.Sprintln(l.appendContext(msg...)...))
}

// Warning logs a message with Warning level
func (l *Logger) Warning(msg ...interface{}) {
	l.warningLogger.Output(l.framesToSkip, fmt.Sprintln(l.appendContext(msg...)...))
}

// Error logs a message with Error level
func (l *Logger) Error(msg ...interface{}) {
	l.errorLogger.Output(l.framesToSkip, fmt.Sprintln(l.appendContext(msg...)...))
}

func (l *Logger) concatString(f string) string {
	if l.contextData != nil {
		f = l.stringContext + " " + f
	}
	return f
}

// Verbose logs a message with Debug level
func (l *Logger) Verbosef(f string, args ...interface{}) {
	l.verboseLogger.Output(l.framesToSkip, fmt.Sprintf(l.concatString(f), args...))
}

// Debug logs a message with Debug level
func (l *Logger) Debugf(f string, args ...interface{}) {
	l.debugLogger.Output(l.framesToSkip, fmt.Sprintf(l.concatString(f), args...))
}

// Info logs a message with Info level
func (l *Logger) Infof(f string, args ...interface{}) {
	l.infoLogger.Output(l.framesToSkip, fmt.Sprintf(l.concatString(f), args...))
}

// Warning logs a message with Warning level
func (l *Logger) Warningf(f string, args ...interface{}) {
	l.warningLogger.Output(l.framesToSkip, fmt.Sprintf(l.concatString(f), args...))
}

// Error logs a message with Error level
func (l *Logger) Errorf(f string, args ...interface{}) {
	l.errorLogger.Output(l.framesToSkip, fmt.Sprintf(l.concatString(f), args...))
}

// WithContext sums one to the number of frames to skip
func (l *Logger) WithContext(ctx context.Context) LoggerInterface {
	contextData := Merge(l.contextData, GetContext(ctx))
	stringContext := ""
	if contextData != nil {
		stringContext = contextData.String()
	}
	return &Logger{
		debugLogger:   l.debugLogger,
		infoLogger:    l.infoLogger,
		warningLogger: l.warningLogger,
		errorLogger:   l.errorLogger,
		verboseLogger: l.verboseLogger,
		framesToSkip:  l.framesToSkip,

		contextData:   contextData,
		stringContext: stringContext,
	}
}

// AugmentFromContext
func (l *Logger) AugmentFromContext(ctx context.Context, values ...string) (LoggerInterface, context.Context) {
	if len(values)%2 == 1 {
		return nil, nil
	}
	newContextData := NewContext()
	for i := 0; i < len(values); i += 2 {
		newContextData = newContextData.WithTag(values[i], values[i+1])
	}
	contextData := Merge(l.contextData, newContextData)
	ctx = context.WithValue(ctx, ContextKey{}, contextData)

	return l.WithContext(ctx), ctx
}

// Clone creates a new logger with the same options as the current one
func (l *Logger) Clone(options ...LoggerOptions) LoggerInterface {
	opts := &LoggerOptions{
		DebugWriter:         l.debugLogger.Writer(),
		InfoWriter:          l.infoLogger.Writer(),
		WarningWriter:       l.warningLogger.Writer(),
		ErrorWriter:         l.errorLogger.Writer(),
		VerboseWriter:       l.verboseLogger.Writer(),
		LogLevel:            l.level,
		StandardLoggerFlags: l.debugLogger.Flags(),
		ExtraFramesToSkip:   l.framesToSkip - skipStackFrameBase + 1,
	}

	for _, opt := range options {
		if opt.DebugWriter != nil {
			opts.DebugWriter = opt.DebugWriter
		}
		if opt.InfoWriter != nil {
			opts.InfoWriter = opt.InfoWriter
		}
		if opt.WarningWriter != nil {
			opts.WarningWriter = opt.WarningWriter
		}
		if opt.ErrorWriter != nil {
			opts.ErrorWriter = opt.ErrorWriter
		}
		if opt.VerboseWriter != nil {
			opts.VerboseWriter = opt.VerboseWriter
		}
		if opt.StandardLoggerFlags != 0 {
			opts.StandardLoggerFlags = opt.StandardLoggerFlags
		}
		if opt.LogLevel != 0 {
			opts.LogLevel = opt.LogLevel
		}
	}

	return NewLogger(opts)
}

func normalizeOptions(options *LoggerOptions) *LoggerOptions {
	var toRet *LoggerOptions
	if options == nil {
		toRet = &LoggerOptions{}
	} else {
		toRet = options
	}

	if toRet.DebugWriter == nil {
		toRet.DebugWriter = os.Stdout
	}

	if toRet.ErrorWriter == nil {
		toRet.ErrorWriter = os.Stdout
	}

	if toRet.InfoWriter == nil {
		toRet.InfoWriter = os.Stdout
	}

	if toRet.VerboseWriter == nil {
		toRet.VerboseWriter = os.Stdout
	}

	if toRet.WarningWriter == nil {
		toRet.WarningWriter = os.Stdout
	}

	if toRet.StandardLoggerFlags == 0 {
		toRet.StandardLoggerFlags = log.Ldate | log.Ltime
	}

	switch toRet.LogLevel {
	case LevelAll, LevelDebug, LevelError, LevelInfo, LevelNone, LevelVerbose, LevelWarning:
	default:
		toRet.LogLevel = LevelError
	}
	return toRet
}

// newLogger constructor of Logger instance.
//
//	Returns: a pointer to *Logger
//	Assumes: that options are alredy normalized
func newLogger(options *LoggerOptions) *Logger {
	prefix := ""
	if options.Prefix != "" {
		prefix = fmt.Sprintf("%s - ", options.Prefix)
	}
	return &Logger{
		debugLogger:   log.New(options.DebugWriter, fmt.Sprintf("%sDEBUG - ", prefix), options.StandardLoggerFlags),
		infoLogger:    log.New(options.InfoWriter, fmt.Sprintf("%sINFO - ", prefix), options.StandardLoggerFlags),
		warningLogger: log.New(options.WarningWriter, fmt.Sprintf("%sWARNING - ", prefix), options.StandardLoggerFlags),
		errorLogger:   log.New(options.ErrorWriter, fmt.Sprintf("%sERROR - ", prefix), options.StandardLoggerFlags),
		verboseLogger: log.New(options.VerboseWriter, fmt.Sprintf("%sVERBOSE - ", prefix), options.StandardLoggerFlags),
		framesToSkip:  skipStackFrameBase + options.ExtraFramesToSkip,
		stringContext: "",
		level:         options.LogLevel,
	}
}

// NewLogger instantiates a new Logger instance. Requires a pointer to a LoggerOptions struct to be passed.
func NewLogger(options *LoggerOptions) LoggerInterface {
	options = normalizeOptions(options)

	logger := newLogger(options)

	return &LevelFilteredLoggerWrapper{
		delegate: logger,
		level:    options.LogLevel,
	}
}

// NewExtendedLogger instantiates a new Logger instance. Requires a pointer to a LoggerOptions struct to be passed.
func NewExtendedLogger(options *LoggerOptions) ExtendedLoggerInterface {
	options = normalizeOptions(options)

	logger := newLogger(options)

	return &ExtendedLevelFilteredLoggerWrapper{&LevelFilteredLoggerWrapper{
		delegate: logger,
		level:    options.LogLevel,
	}}
}
