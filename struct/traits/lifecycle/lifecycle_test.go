package lifecycle

import (
	"testing"
	"time"
)

func TestLifecycleManager(t *testing.T) {
	m := Manager{}
	m.Setup()

	if !m.BeginInitialization() {
		t.Error("initialization should begin properly.")
	}

	if m.IsRunning() {
		t.Error("isRunning should be false while initialization is going on")
	}

	if m.BeginInitialization() {
		t.Error("initialization should fail if called more than once.")
	}

	if m.BeginShutdown() {
		t.Error("shutdown cannot be started until the manager is fully running")
	}

	m.InitializationComplete()
	if !m.IsRunning() {
		t.Error("it should be running")
	}

	done := make(chan struct{}, 1)
	go func() {
		defer m.ShutdownComplete()
		defer func() { done <- struct{}{} }()
		for {
			select {
			case <-m.ShutdownRequested():
				time.Sleep(1 * time.Second)
				return
			}
		}
	}()

	if !m.BeginShutdown() {
		t.Error("shutdown should be correctly propagated")
	}
	if m.BeginShutdown() {
		t.Error("once shutdown is started, it should no longer propagate further requests")
	}
	m.AwaitShutdownComplete()
	if m.IsRunning() {
		t.Error("should not be running")
	}
	<-done // ensure that await actually waits

	// Start again

	if !m.BeginInitialization() {
		t.Error("initialization should begin properly.")
	}

	if m.IsRunning() {
		t.Error("isRunning should be false while initialization is going on")
	}

	if m.BeginInitialization() {
		t.Error("initialization should fail if called more than once.")
	}

	if m.BeginShutdown() {
		t.Error("shutdown cannot be started until the manager is fully running")
	}

	m.InitializationComplete()
	if !m.IsRunning() {
		t.Error("it should be running")
	}

	done = make(chan struct{}, 1)
	go func() {
		defer m.ShutdownComplete()
		defer func() { done <- struct{}{} }()
		for {
			select {
			case <-m.ShutdownRequested():
				time.Sleep(1 * time.Second)
				return
			}
		}
	}()

	if !m.BeginShutdown() {
		t.Error("shutdown should be correctly propagated")
	}
	if m.BeginShutdown() {
		t.Error("once shutdown is started, it should no longer propagate further requests")
	}
	m.AwaitShutdownComplete()
	if m.IsRunning() {
		t.Error("should not be running")
	}
	<-done // ensure that await actually waits
}

func TestLifecycleManagerAbnormalShutdown(t *testing.T) {
	m := Manager{}
	m.Setup()

	if !m.BeginInitialization() {
		t.Error("initialization should begin properly.")
	}

	if m.IsRunning() {
		t.Error("isRunning should be false while initialization is going on")
	}

	if m.BeginInitialization() {
		t.Error("initialization should fail if called more than once.")
	}

	if m.BeginShutdown() {
		t.Error("shutdown cannot be started until the manager is fully running")
	}

	m.InitializationComplete()
	if !m.IsRunning() {
		t.Error("it should be running")
	}

	done := make(chan struct{}, 1)
	go func() {
		defer m.ShutdownComplete()
		defer func() { done <- struct{}{} }()
		for {
			select {
			case <-time.After(1 * time.Second):
				m.AbnormalShutdown()
				return
			}
		}
	}()

	m.AwaitShutdownComplete()
	if m.IsRunning() {
		t.Error("should not be running")
	}
	<-done // ensure that await actually waits

	// Start again

	if !m.BeginInitialization() {
		t.Error("initialization should begin properly.")
	}

	if m.IsRunning() {
		t.Error("isRunning should be false while initialization is going on")
	}

	if m.BeginInitialization() {
		t.Error("initialization should fail if called more than once.")
	}

	if m.BeginShutdown() {
		t.Error("shutdown cannot be started until the manager is fully running")
	}

	m.InitializationComplete()
	if !m.IsRunning() {
		t.Error("it should be running")
	}

	done = make(chan struct{}, 1)
	go func() {
		defer m.ShutdownComplete()
		defer func() { done <- struct{}{} }()
		for {
			select {
			case <-m.ShutdownRequested():
				time.Sleep(1 * time.Second)
				return
			}
		}
	}()

	if !m.BeginShutdown() {
		t.Error("shutdown should be correctly propagated")
	}

	if m.BeginShutdown() {
		t.Error("once shutdown is started, it should no longer propagate further requests")
	}
	m.AwaitShutdownComplete()
	if m.IsRunning() {
		t.Error("should not be running")
	}
	<-done // ensure that await actually waits
}
