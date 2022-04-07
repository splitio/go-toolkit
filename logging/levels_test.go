package logging

import (
	"testing"
)

type mockedLogger struct {
	msgs map[string]bool
}

func (l *mockedLogger) reset() {
	l.msgs = make(map[string]bool)
}

func (l *mockedLogger) Error(msg ...interface{}) {
	l.msgs["Error"] = true
}

func (l *mockedLogger) Warning(msg ...interface{}) {
	l.msgs["Warning"] = true
}

func (l *mockedLogger) Info(msg ...interface{}) {
	l.msgs["Info"] = true
}

func (l *mockedLogger) Debug(msg ...interface{}) {
	l.msgs["Debug"] = true
}

func (l *mockedLogger) Verbose(msg ...interface{}) {
	l.msgs["Verbose"] = true
}

func TestErrorLevel(t *testing.T) {

	delegate := &mockedLogger{}
	delegate.reset()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelError,
	}

	logger.Error("text")
	logger.Warning("text")
	logger.Info("text")
	logger.Debug("text")
	logger.Verbose("text")

	shouldBeCalled := []string{"Error"}
	shouldNotBeCalled := []string{"Warning", "Info", "Debug", "Verbose"}

	for _, level := range shouldBeCalled {
		if !delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should have been forwarded", level)
		}
	}

	for _, level := range shouldNotBeCalled {
		if delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should NOT have been forwarded", level)
		}
	}
}

func TestWarningLevel(t *testing.T) {

	delegate := &mockedLogger{}
	delegate.reset()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelWarning,
	}

	logger.Error("text")
	logger.Warning("text")
	logger.Info("text")
	logger.Debug("text")
	logger.Verbose("text")

	shouldBeCalled := []string{"Error", "Warning"}
	shouldNotBeCalled := []string{"Info", "Debug", "Verbose"}

	for _, level := range shouldBeCalled {
		if !delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should have been forwarded", level)
		}
	}

	for _, level := range shouldNotBeCalled {
		if delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should NOT have been forwarded", level)
		}
	}
}

func TestInfoLevel(t *testing.T) {

	delegate := &mockedLogger{}
	delegate.reset()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelInfo,
	}

	logger.Error("text")
	logger.Warning("text")
	logger.Info("text")
	logger.Debug("text")
	logger.Verbose("text")

	shouldBeCalled := []string{"Error", "Warning", "Info"}
	shouldNotBeCalled := []string{"Debug", "Verbose"}

	for _, level := range shouldBeCalled {
		if !delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should have been forwarded", level)
		}
	}

	for _, level := range shouldNotBeCalled {
		if delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should NOT have been forwarded", level)
		}
	}
}

func TestDebugLevel(t *testing.T) {

	delegate := &mockedLogger{}
	delegate.reset()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelDebug,
	}

	logger.Error("text")
	logger.Warning("text")
	logger.Info("text")
	logger.Debug("text")
	logger.Verbose("text")

	shouldBeCalled := []string{"Error", "Warning", "Info", "Debug"}
	shouldNotBeCalled := []string{"Verbose"}

	for _, level := range shouldBeCalled {
		if !delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should have been forwarded", level)
		}
	}

	for _, level := range shouldNotBeCalled {
		if delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should NOT have been forwarded", level)
		}
	}
}

func TestVerboseLevel(t *testing.T) {

	delegate := &mockedLogger{}
	delegate.reset()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelVerbose,
	}

	logger.Error("text")
	logger.Warning("text")
	logger.Info("text")
	logger.Debug("text")
	logger.Verbose("text")

	shouldBeCalled := []string{"Error", "Warning", "Info", "Debug", "Verbose"}
	shouldNotBeCalled := []string{}

	for _, level := range shouldBeCalled {
		if !delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should have been forwarded", level)
		}
	}

	for _, level := range shouldNotBeCalled {
		if delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should NOT have been forwarded", level)
		}
	}
}

func TestAllLevel(t *testing.T) {

	delegate := &mockedLogger{}
	delegate.reset()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelAll,
	}

	logger.Error("text")
	logger.Warning("text")
	logger.Info("text")
	logger.Debug("text")
	logger.Verbose("text")

	shouldBeCalled := []string{"Error", "Warning", "Info", "Debug", "Verbose"}
	shouldNotBeCalled := []string{}

	for _, level := range shouldBeCalled {
		if !delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should have been forwarded", level)
		}
	}

	for _, level := range shouldNotBeCalled {
		if delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should NOT have been forwarded", level)
		}
	}
}

func TestNoneLevel(t *testing.T) {

	delegate := &mockedLogger{}
	delegate.reset()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelNone,
	}

	logger.Error("text")
	logger.Warning("text")
	logger.Info("text")
	logger.Debug("text")
	logger.Verbose("text")

	shouldNotBeCalled := []string{"Error", "Warning", "Info", "Debug", "Verbose"}
	shouldBeCalled := []string{}

	for _, level := range shouldBeCalled {
		if !delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should have been forwarded", level)
		}
	}

	for _, level := range shouldNotBeCalled {
		if delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should NOT have been forwarded", level)
		}
	}
}

func writelog(logger *ExtendedLevelFilteredLoggerWrapper) {
	logger.ErrorFn("hello %s", func() []interface{} { return []interface{}{"world"} })
	logger.WarningFn("hello %s", func() []interface{} { return []interface{}{"world"} })
	logger.InfoFn("hello %s", func() []interface{} { return []interface{}{"world"} })
	logger.DebugFn("hello %s", func() []interface{} { return []interface{}{"world"} })
	logger.VerboseFn("hello %s", func() []interface{} { return []interface{}{"world"} })
}

func assertWrites(t *testing.T, currentLevel string, delegate *mockedLogger, shouldBeCalled []string, shouldNotBeCalled []string) {
	for _, level := range shouldBeCalled {
		if !delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should have been forwarded, current level=%s", level, currentLevel)
		}
	}

	for _, level := range shouldNotBeCalled {
		if delegate.msgs[level] {
			t.Errorf("Call to log level function \"%s\" should NOT have been forwarded, current level=%s", level, currentLevel)
		}
	}
}

func TestExtendedLevelFilteredLogger(t *testing.T) {

	delegate := &mockedLogger{}
	delegate.reset()

	logger := &ExtendedLevelFilteredLoggerWrapper{&LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelError,
	}}

	t.Run("Error", func(t *testing.T) {
		logger.level = LevelError
		delegate.reset()
		writelog(logger)
		assertWrites(t, "ERROR", delegate, []string{"Error"}, []string{"Warning", "Info", "Debug", "Verbose"})
	})

	t.Run("Waring", func(t *testing.T) {
		logger.level = LevelWarning
		delegate.reset()
		writelog(logger)
		assertWrites(t, "WARNING", delegate, []string{"Error", "Warning"}, []string{"Info", "Debug", "Verbose"})
	})

	t.Run("Info", func(t *testing.T) {
		logger.level = LevelInfo
		delegate.reset()
		writelog(logger)
		assertWrites(t, "INFO", delegate, []string{"Error", "Warning", "Info"}, []string{"Debug", "Verbose"})
	})

	t.Run("Debug", func(t *testing.T) {
		logger.level = LevelDebug
		delegate.reset()
		writelog(logger)
		assertWrites(t, "DEBUG", delegate, []string{"Error", "Warning", "Info", "Debug"}, []string{"Verbose"})
	})

	t.Run("Verbose", func(t *testing.T) {
		logger.level = LevelVerbose
		delegate.reset()
		writelog(logger)
		assertWrites(t, "VERBOSE", delegate, []string{"Error", "Warning", "Info", "Debug", "Verbose"}, []string{})
	})

	t.Run("All", func(t *testing.T) {
		logger.level = LevelAll
		delegate.reset()
		writelog(logger)
		assertWrites(t, "ALL", delegate, []string{"Error", "Warning", "Info", "Debug", "Verbose"}, []string{})
	})

	t.Run("None", func(t *testing.T) {
		logger.level = LevelNone
		delegate.reset()
		writelog(logger)
		assertWrites(t, "NONE", delegate, []string{}, []string{"Error", "Warning", "Info", "Debug", "Verbose"})
	})
}
