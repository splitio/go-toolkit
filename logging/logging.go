// Package logging ...
// Handles logging within the SDK
package logging

import (
	"io"
	"log"
	"os"
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
}

// Logger struct. Encapsulates four different loggers, each for a different "level",
// and provides Error, Debug, Warning and Info functions, that will forward a message
// to the appropriate logger.
type Logger struct {
	debugLogger   log.Logger
	infoLogger    log.Logger
	warningLogger log.Logger
	errorLogger   log.Logger
	verboseLogger log.Logger
}

// Verbose logs a message with Debug level
func (l *Logger) Verbose(msg ...interface{}) {
	l.verboseLogger.Println(msg...)
}

// Debug logs a message with Debug level
func (l *Logger) Debug(msg ...interface{}) {
	l.debugLogger.Println(msg...)
}

// Info logs a message with Info level
func (l *Logger) Info(msg ...interface{}) {
	l.infoLogger.Println(msg...)
}

// Warning logs a message with Warning level
func (l *Logger) Warning(msg ...interface{}) {
	l.warningLogger.Println(msg...)
}

// Error logs a message with Error level
func (l *Logger) Error(msg ...interface{}) {
	l.errorLogger.Println(msg...)
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

// NewLogger instantiates a new Logger instance. Requires a pointer to a LoggerOptions struct to be passed.
func NewLogger(options *LoggerOptions) LoggerInterface {

	options = normalizeOptions(options)

	logger := &Logger{
		debugLogger:   *log.New(options.DebugWriter, "DEBUG - ", options.StandardLoggerFlags),
		infoLogger:    *log.New(options.InfoWriter, "INFO - ", options.StandardLoggerFlags),
		warningLogger: *log.New(options.WarningWriter, "WARNING - ", options.StandardLoggerFlags),
		errorLogger:   *log.New(options.ErrorWriter, "ERROR - ", options.StandardLoggerFlags),
		verboseLogger: *log.New(options.VerboseWriter, "VERBOSE - ", options.StandardLoggerFlags),
	}

	if options.Prefix != "" {
		logger.debugLogger.SetPrefix(options.Prefix)
		logger.infoLogger.SetPrefix(options.Prefix)
		logger.warningLogger.SetPrefix(options.Prefix)
		logger.verboseLogger.SetPrefix(options.Prefix)
		logger.errorLogger.SetPrefix(options.Prefix)
	}

	return &LevelFilteredLoggerWrapper{
		delegate: logger,
		level:    options.LogLevel,
	}
}
