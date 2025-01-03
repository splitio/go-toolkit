package logging

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestWithoutPrefix(t *testing.T) {
	today := time.Now()
	asFormat := fmt.Sprintf("%d/%02d/%02d", today.Year(), today.Month(), today.Day())
	mockedWriter := &MockWriter{}
	logger := NewLogger(&LoggerOptions{
		StandardLoggerFlags: log.Ldate,
		LogLevel:            5,
		ErrorWriter:         mockedWriter,
		WarningWriter:       mockedWriter,
		InfoWriter:          mockedWriter,
		DebugWriter:         mockedWriter,
		VerboseWriter:       mockedWriter,
		Prefix:              "",
	})
	mockedWriter.On("Write", []byte(fmt.Sprintf("DEBUG - %s some\n", asFormat))).Once().Return(0, nil)
	logger.Debug("some")
	mockedWriter.On("Write", []byte(fmt.Sprintf("ERROR - %s some\n", asFormat))).Once().Return(0, nil)
	logger.Error("some")
	mockedWriter.On("Write", []byte(fmt.Sprintf("INFO - %s some\n", asFormat))).Once().Return(0, nil)
	logger.Info("some")
	mockedWriter.On("Write", []byte(fmt.Sprintf("VERBOSE - %s some\n", asFormat))).Once().Return(0, nil)
	logger.Verbose("some")
	mockedWriter.On("Write", []byte(fmt.Sprintf("WARNING - %s some\n", asFormat))).Once().Return(0, nil)
	logger.Warning("some")
	mockedWriter.AssertExpectations(t)
}

func TestWithPrefix(t *testing.T) {
	today := time.Now()
	asFormat := fmt.Sprintf("%d/%02d/%02d", today.Year(), today.Month(), today.Day())
	mockedWriter := &MockWriter{}
	logger := NewLogger(&LoggerOptions{
		StandardLoggerFlags: log.Ldate,
		LogLevel:            5,
		ErrorWriter:         mockedWriter,
		WarningWriter:       mockedWriter,
		InfoWriter:          mockedWriter,
		DebugWriter:         mockedWriter,
		VerboseWriter:       mockedWriter,
		Prefix:              "prefix",
	})
	mockedWriter.On("Write", []byte(fmt.Sprintf("prefix - DEBUG - %s some\n", asFormat))).Once().Return(0, nil)
	logger.Debug("some")
	mockedWriter.On("Write", []byte(fmt.Sprintf("prefix - ERROR - %s some\n", asFormat))).Once().Return(0, nil)
	logger.Error("some")
	mockedWriter.On("Write", []byte(fmt.Sprintf("prefix - INFO - %s some\n", asFormat))).Once().Return(0, nil)
	logger.Info("some")
	mockedWriter.On("Write", []byte(fmt.Sprintf("prefix - VERBOSE - %s some\n", asFormat))).Once().Return(0, nil)
	logger.Verbose("some")
	mockedWriter.On("Write", []byte(fmt.Sprintf("prefix - WARNING - %s some\n", asFormat))).Once().Return(0, nil)
	logger.Warning("some")
	mockedWriter.AssertExpectations(t)
}
