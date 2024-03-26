package logging

import (
	"testing"

	"github.com/splitio/go-toolkit/v6/logging/mocks"
)

func TestErrorLevel(t *testing.T) {

	delegate := &mocks.MockLogger{}
	delegate.On("Error", "text").Once()
	delegate.On("Errorf", "formatted %d", int(3)).Once()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelError,
	}

	logger.Error("text")
	logger.Errorf("formatted %d", int(3))
	logger.Warning("text")
	logger.Warningf("formatted %d", int(3))
	logger.Info("text")
	logger.Infof("formatted %d", int(3))
	logger.Debug("text")
	logger.Debugf("formatted %d", int(3))
	logger.Verbose("text")
	logger.Verbosef("formatted %d", int(3))

	delegate.AssertExpectations(t)
}

func TestWarningLevel(t *testing.T) {
	delegate := &mocks.MockLogger{}
	delegate.On("Error", "text").Once()
	delegate.On("Errorf", "formatted %d", int(3)).Once()
	delegate.On("Warning", "text").Once()
	delegate.On("Warningf", "formatted %d", int(3)).Once()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelWarning,
	}

	logger.Error("text")
	logger.Errorf("formatted %d", int(3))
	logger.Warning("text")
	logger.Warningf("formatted %d", int(3))
	logger.Info("text")
	logger.Infof("formatted %d", int(3))
	logger.Debug("text")
	logger.Debugf("formatted %d", int(3))
	logger.Verbose("text")
	logger.Verbosef("formatted %d", int(3))

	delegate.AssertExpectations(t)
}

func TestInfoLevel(t *testing.T) {
	delegate := &mocks.MockLogger{}
	delegate.On("Error", "text").Once()
	delegate.On("Errorf", "formatted %d", int(3)).Once()
	delegate.On("Warning", "text").Once()
	delegate.On("Warningf", "formatted %d", int(3)).Once()
	delegate.On("Info", "text").Once()
	delegate.On("Infof", "formatted %d", int(3)).Once()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelInfo,
	}

	logger.Error("text")
	logger.Errorf("formatted %d", int(3))
	logger.Warning("text")
	logger.Warningf("formatted %d", int(3))
	logger.Info("text")
	logger.Infof("formatted %d", int(3))
	logger.Debug("text")
	logger.Debugf("formatted %d", int(3))
	logger.Verbose("text")
	logger.Verbosef("formatted %d", int(3))

	delegate.AssertExpectations(t)
}

func TestDebugLevel(t *testing.T) {
	delegate := &mocks.MockLogger{}
	delegate.On("Error", "text").Once()
	delegate.On("Errorf", "formatted %d", int(3)).Once()
	delegate.On("Warning", "text").Once()
	delegate.On("Warningf", "formatted %d", int(3)).Once()
	delegate.On("Info", "text").Once()
	delegate.On("Infof", "formatted %d", int(3)).Once()
	delegate.On("Debug", "text").Once()
	delegate.On("Debugf", "formatted %d", int(3)).Once()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelDebug,
	}

	logger.Error("text")
	logger.Errorf("formatted %d", int(3))
	logger.Warning("text")
	logger.Warningf("formatted %d", int(3))
	logger.Info("text")
	logger.Infof("formatted %d", int(3))
	logger.Debug("text")
	logger.Debugf("formatted %d", int(3))
	logger.Verbose("text")
	logger.Verbosef("formatted %d", int(3))

	delegate.AssertExpectations(t)
}

func TestVerboseLevel(t *testing.T) {
	delegate := &mocks.MockLogger{}
	delegate.On("Error", "text").Once()
	delegate.On("Errorf", "formatted %d", int(3)).Once()
	delegate.On("Warning", "text").Once()
	delegate.On("Warningf", "formatted %d", int(3)).Once()
	delegate.On("Info", "text").Once()
	delegate.On("Infof", "formatted %d", int(3)).Once()
	delegate.On("Debug", "text").Once()
	delegate.On("Debugf", "formatted %d", int(3)).Once()
	delegate.On("Verbose", "text").Once()
	delegate.On("Verbosef", "formatted %d", int(3)).Once()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelVerbose,
	}

	logger.Error("text")
	logger.Errorf("formatted %d", int(3))
	logger.Warning("text")
	logger.Warningf("formatted %d", int(3))
	logger.Info("text")
	logger.Infof("formatted %d", int(3))
	logger.Debug("text")
	logger.Debugf("formatted %d", int(3))
	logger.Verbose("text")
	logger.Verbosef("formatted %d", int(3))

	delegate.AssertExpectations(t)
}

func TestAllLevel(t *testing.T) {
	delegate := &mocks.MockLogger{}
	delegate.On("Error", "text").Once()
	delegate.On("Errorf", "formatted %d", int(3)).Once()
	delegate.On("Warning", "text").Once()
	delegate.On("Warningf", "formatted %d", int(3)).Once()
	delegate.On("Info", "text").Once()
	delegate.On("Infof", "formatted %d", int(3)).Once()
	delegate.On("Debug", "text").Once()
	delegate.On("Debugf", "formatted %d", int(3)).Once()
	delegate.On("Verbose", "text").Once()
	delegate.On("Verbosef", "formatted %d", int(3)).Once()

	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelAll,
	}

	logger.Error("text")
	logger.Errorf("formatted %d", int(3))
	logger.Warning("text")
	logger.Warningf("formatted %d", int(3))
	logger.Info("text")
	logger.Infof("formatted %d", int(3))
	logger.Debug("text")
	logger.Debugf("formatted %d", int(3))
	logger.Verbose("text")
	logger.Verbosef("formatted %d", int(3))

	delegate.AssertExpectations(t)

}

func TestNoneLevel(t *testing.T) {
	delegate := &mocks.MockLogger{}
	logger := LevelFilteredLoggerWrapper{
		delegate: delegate,
		level:    LevelNone,
	}

	logger.Error("text")
	logger.Errorf("formatted %d", int(3))
	logger.Warning("text")
	logger.Warningf("formatted %d", int(3))
	logger.Info("text")
	logger.Infof("formatted %d", int(3))
	logger.Debug("text")
	logger.Debugf("formatted %d", int(3))
	logger.Verbose("text")
	logger.Verbosef("formatted %d", int(3))

	delegate.AssertExpectations(t)
}

// ---------------------------------

type mockedLogger struct{ msgs map[string]bool }

func (*mockedLogger) Debugf(fmt string, msg ...interface{})   { panic("unimplemented") }
func (*mockedLogger) Errorf(fmt string, msg ...interface{})   { panic("unimplemented") }
func (*mockedLogger) Infof(fmt string, msg ...interface{})    { panic("unimplemented") }
func (*mockedLogger) Verbosef(fmt string, msg ...interface{}) { panic("unimplemented") }
func (*mockedLogger) Warningf(fmt string, msg ...interface{}) { panic("unimplemented") }
func (l *mockedLogger) reset()                                { l.msgs = make(map[string]bool) }
func (l *mockedLogger) Error(msg ...interface{})              { l.msgs["Error"] = true }
func (l *mockedLogger) Warning(msg ...interface{})            { l.msgs["Warning"] = true }
func (l *mockedLogger) Info(msg ...interface{})               { l.msgs["Info"] = true }
func (l *mockedLogger) Debug(msg ...interface{})              { l.msgs["Debug"] = true }
func (l *mockedLogger) Verbose(msg ...interface{})            { l.msgs["Verbose"] = true }

var _ LoggerInterface = (*mockedLogger)(nil)

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
