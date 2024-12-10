package logging

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestLoggerContext(t *testing.T) {
	today := time.Now()
	asFormat := fmt.Sprintf("%d/%02d/%02d", today.Year(), today.Month(), today.Day())
	mW := &MockWriter{}
	mW.On("Write", []byte(fmt.Sprintf("DEBUG - %s test\n", asFormat))).Once().Return(0, nil)
	mW.On("Write", []byte(fmt.Sprintf("DEBUG - %s [txID: tx] test\n", asFormat))).Once().Return(0, nil)
	logger := NewLogger(&LoggerOptions{
		StandardLoggerFlags: log.Ldate,
		LogLevel:            LevelVerbose,
		ErrorWriter:         mW,
		WarningWriter:       mW,
		InfoWriter:          mW,
		DebugWriter:         mW,
		VerboseWriter:       mW,
	})
	logger.Debug("test")

	bg := context.Background()
	info := NewContextInformation()
	info.Add("txID", "tx")
	bg2 := context.WithValue(bg, ContextKey{}, info)
	loggerWithContext := logger.WithContext(bg2)
	loggerWithContext.Debug("test")
}

func TestLoggerWrapper(t *testing.T) {
	loggerTest := NewLogger(&LoggerOptions{
		StandardLoggerFlags: log.LUTC | log.Ldate | log.Lmicroseconds | log.Lshortfile,
		LogLevel:            LevelDebug,
		Prefix:              "test",
	})

	bg := context.Background()
	info := NewContextInformation()
	info.Add("txID", "tx")
	bg2 := context.WithValue(bg, ContextKey{}, info)
	loggerTest.Debug("sd", "aa")
	loggerWithContext := loggerTest.WithContext(bg2)
	loggerWithContext.Debug("sd", "aa")

	info2 := NewContextInformation()
	info2.Add("orgid", "org")
	bg3 := context.WithValue(bg2, ContextKey{}, info2)
	loggerWithContext2 := loggerWithContext.WithContext(bg3)
	loggerWithContext2.Debug("3", "aa")

	other := loggerWithContext.WithContext(bg)
	other.Debugf("PreProcess executed successful %s", "test")

}
