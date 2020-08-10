package logging

import (
	"strings"
	"sync"
	"testing"
)

type MockWriter struct {
	mutex  sync.RWMutex
	strMsg string
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.strMsg = string(p[:])
	return 0, nil
}

func (m *MockWriter) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.strMsg = ""
}

func (m *MockWriter) Get() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.strMsg
}

var mW MockWriter

func expectedLogMessage(expectedMessage string, t *testing.T) {
	if !strings.Contains(mW.strMsg, expectedMessage) {
		t.Error("Message error is different from the expected: " + mW.Get())
	}
	mW.Reset()
}

func getMockedLogger(prefix string) LoggerInterface {
	return NewLogger(&LoggerOptions{
		LogLevel:      5,
		ErrorWriter:   &mW,
		WarningWriter: &mW,
		InfoWriter:    &mW,
		DebugWriter:   &mW,
		VerboseWriter: &mW,
		Prefix:        prefix,
	})
}

func TestWithoutPrefix(t *testing.T) {
	logger := getMockedLogger("")
	logger.Debug("some")
	debug := strings.Split(mW.strMsg, "-")
	if debug[0] != "DEBUG " {
		t.Error("It should not include prefix")
	}
	logger.Error("some")
	err := strings.Split(mW.strMsg, "-")
	if err[0] != "ERROR " {
		t.Error("It should not include prefix")
	}
	logger.Info("some")
	info := strings.Split(mW.strMsg, "-")
	if info[0] != "INFO " {
		t.Error("It should not include prefix")
	}
	logger.Verbose("some")
	verbose := strings.Split(mW.strMsg, "-")
	if verbose[0] != "VERBOSE " {
		t.Error("It should not include prefix")
	}
	logger.Warning("some")
	warning := strings.Split(mW.strMsg, "-")
	if warning[0] != "WARNING " {
		t.Error("It should not include prefix")
	}
}

func TestWithPrefix(t *testing.T) {
	logger := getMockedLogger("prefix")
	logger.Debug("some")
	debug := strings.Split(mW.strMsg, "-")
	if debug[0] != "prefix " {
		t.Error("It should include prefix")
	}
	if debug[1] != " DEBUG " {
		t.Error("It should not include prefix")
	}
	logger.Error("some")
	err := strings.Split(mW.strMsg, "-")
	if err[0] != "prefix " {
		t.Error("It should include prefix")
	}
	if err[1] != " ERROR " {
		t.Error("It should not include prefix")
	}
	logger.Info("some")
	info := strings.Split(mW.strMsg, "-")
	if info[0] != "prefix " {
		t.Error("It should include prefix")
	}
	if info[1] != " INFO " {
		t.Error("It should not include prefix")
	}
	logger.Verbose("some")
	verbose := strings.Split(mW.strMsg, "-")
	if verbose[0] != "prefix " {
		t.Error("It should include prefix")
	}
	if verbose[1] != " VERBOSE " {
		t.Error("It should not include prefix")
	}
	logger.Warning("some")
	warning := strings.Split(mW.strMsg, "-")
	if warning[0] != "prefix " {
		t.Error("It should include prefix")
	}
	if warning[1] != " WARNING " {
		t.Error("It should not include prefix")
	}
}
