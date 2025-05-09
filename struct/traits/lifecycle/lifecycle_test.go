package lifecycle

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLifecycleManager(t *testing.T) {
	m := Manager{}
	m.Setup()

    assert.True(t, m.BeginInitialization())
    assert.False(t, m.IsRunning())
    assert.False(t, m.BeginInitialization())
    assert.True(t, m.InitializationComplete())
    assert.True(t, m.IsRunning())

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

    assert.True(t, m.BeginShutdown())
    assert.False(t, m.BeginShutdown())

	m.AwaitShutdownComplete()

    assert.False(t, m.IsRunning())

	<-done // ensure that await actually waits

	// Start again

    assert.True(t, m.BeginInitialization())
    assert.False(t, m.IsRunning())
    assert.False(t, m.BeginInitialization())
    assert.True(t, m.InitializationComplete())
    assert.True(t, m.IsRunning())

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

    assert.True(t, m.BeginShutdown())
    assert.False(t, m.BeginShutdown())

	m.AwaitShutdownComplete()

    assert.False(t, m.IsRunning())

	<-done // ensure that await actually waits
}

func TestLifecycleManagerAbnormalShutdown(t *testing.T) {
	m := Manager{}
	m.Setup()

    assert.True(t, m.BeginInitialization())
    assert.False(t, m.IsRunning())
    assert.False(t, m.BeginInitialization())
    assert.True(t, m.InitializationComplete())
    assert.True(t, m.IsRunning())

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
    assert.False(t, m.IsRunning())
	<-done // ensure that await actually waits

	// Start again

    assert.True(t, m.BeginInitialization())
    assert.False(t, m.IsRunning())
    assert.False(t, m.BeginInitialization())
    assert.True(t, m.InitializationComplete())
    assert.True(t, m.IsRunning())

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

    assert.True(t, m.BeginShutdown())
    assert.False(t, m.BeginShutdown())

	m.AwaitShutdownComplete()
    assert.False(t, m.IsRunning())

	<-done // ensure that await actually waits
}

func TestShutdownRequestWhileInitNotComplete(t *testing.T) {
	m := Manager{}
	m.Setup()

	assert.True(t, m.BeginInitialization())
    assert.True(t, m.BeginShutdown())
    assert.False(t, m.InitializationComplete())
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

    assert.True(t, m.BeginShutdown())
	m.AwaitShutdownComplete()
    assert.False(t, m.IsRunning())
	<-done // ensure that await actually waits

    assert.Equal(t, int32(0), atomic.LoadInt32(&executed))
}
