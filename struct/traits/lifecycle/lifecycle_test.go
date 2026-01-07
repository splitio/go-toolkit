package lifecycle

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

	if !m.InitializationComplete() {
		t.Error("should complete initialization correctly")
	}

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

func TestShutdownRequestWhileInitNotComplete(t *testing.T) {
	m := Manager{}
	m.Setup()

	m.BeginInitialization()
	if !m.BeginShutdown() {
		t.Error("should accept the shutdown request")
	}

	if m.InitializationComplete() {
		t.Error("initialization cannot complete.")
	}

	m.ShutdownComplete()

	// Now restart the lifecycle to see if it works properly
	m.BeginInitialization()

	done := make(chan struct{}, 1)
	var executed int32
	go func() {
		defer m.ShutdownComplete()
		defer func() { done <- struct{}{} }()
		if !m.InitializationComplete() {
			return
		}
		atomic.StoreInt32(&executed, 1)
		for {
			select {
			case <-time.After(1 * time.Second):
				m.AbnormalShutdown()
				return
			}
		}
	}()
	m.BeginShutdown()
	m.AwaitShutdownComplete()
	if m.IsRunning() {
		t.Error("should not be running")
	}
	<-done // ensure that await actually waits

	if atomic.LoadInt32(&executed) != 0 {
		t.Error("the goroutine should have not executed further than the InitializationComplete check.")
	}
}

func TestInitializationCancelledEventuallyBecomesIdle(t *testing.T) {
	var m Manager
	m.Setup()

	require.True(t, m.BeginInitialization())
	require.True(t, m.BeginShutdown()) // cancela init

	done := make(chan struct{})
	go func() {
		m.AwaitShutdownComplete()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("initialization cancellation never transitions to Idle")
	}
}

func TestShutdownWithoutWorkerDoesNotHang(t *testing.T) {
	var m Manager
	m.Setup()

	require.True(t, m.BeginInitialization())
	require.True(t, m.InitializationComplete())

	require.True(t, m.BeginShutdown())

	done := make(chan struct{})
	go func() {
		m.AwaitShutdownComplete()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("shutdown hangs forever waiting for ShutdownComplete")
	}
}
