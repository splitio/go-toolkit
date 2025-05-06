package logging

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoggerContext(t *testing.T) {
	today := time.Now()
	asFormat := fmt.Sprintf("%d/%02d/%02d", today.Year(), today.Month(), today.Day())
	mW := &MockWriter{}
	mW.On("Write", []byte(fmt.Sprintf("DEBUG - %s test\n", asFormat))).Once().Return(0, nil)
	mW.On("Write", []byte(fmt.Sprintf("DEBUG - %s [txID=tx] test\n", asFormat))).Once().Return(0, nil)
	mW.On("Write", []byte(fmt.Sprintf("DEBUG - %s [orgID=org, txID=tx, userID=user] test\n", asFormat))).Once().Return(0, nil)
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
	info := NewContext().WithTag("txID", "tx")
	bg2 := context.WithValue(bg, ContextKey{}, info)
	loggerWithContext := logger.WithContext(bg2)
	loggerWithContext.Debug("test")

	logger, ctx := loggerWithContext.AugmentFromContext(bg2, "orgID", "org", "userID", "user")
	logger.Debug("test")
	ctxData := ctx.Value(ContextKey{}).(*ContextData)
	assert.Equal(t, "org", ctxData.Get("orgID"))
	assert.Equal(t, "tx", ctxData.Get("txID"))
	assert.Equal(t, "user", ctxData.Get("userID"))
}

func TestLoggerWrapper(t *testing.T) {
	loggerTest := NewLogger(&LoggerOptions{
		StandardLoggerFlags: log.LUTC | log.Ldate | log.Lmicroseconds | log.Lshortfile,
		LogLevel:            LevelDebug,
		Prefix:              "test",
	})

	bg := context.Background()
	info := NewContext().WithTag("txID", "tx")
	bg2 := context.WithValue(bg, ContextKey{}, info)
	loggerTest.Debug("sd", "aa")
	loggerWithContext := loggerTest.WithContext(bg2)
	loggerWithContext.Debug("sd", "aa")

	info2 := NewContext().WithTag("orgid", "org")
	bg3 := context.WithValue(bg2, ContextKey{}, info2)
	loggerWithContext2 := loggerWithContext.WithContext(bg3)
	loggerWithContext2.Debug("3", "aa")

	other := loggerWithContext.WithContext(bg)
	other.Debugf("PreProcess executed successful %s", "test")
}

func TestClone(t *testing.T) {
	today := time.Now()
	asFormat := fmt.Sprintf("%d/%02d/%02d", today.Year(), today.Month(), today.Day())
	mW := &MockWriter{}
	mW.On("Write", []byte(fmt.Sprintf("DEBUG - %s test\n", asFormat))).Once().Return(0, nil)
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

	mW2 := &MockWriter{}
	mW2.On("Write", []byte(fmt.Sprintf("ERROR - %s test\n", asFormat))).Once().Return(0, nil)
	logger2 := logger.Clone(LoggerOptions{LogLevel: LevelError}, LoggerOptions{ErrorWriter: mW2})
	logger2.Error("test")
	logger2.Debug("test")
	mW.AssertExpectations(t)
	mW2.AssertExpectations(t)
}
