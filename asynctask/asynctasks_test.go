package asynctask

import (
	"errors"
	"github.com/splitio/go-toolkit/logging"
	"sync/atomic"
	"testing"
	"time"
)

func TestAsyncTaskNormalOperation(t *testing.T) {

	onInit := atomic.Value{}
	onInit.Store(false)
	onExecution := atomic.Value{}
	onExecution.Store(false)
	onStop := atomic.Value{}
	onStop.Store(false)

	task1 := NewAsyncTask(
		"testTask1",
		func(l logging.LoggerInterface) error { onExecution.Store(true); return nil },
		1,
		func(l logging.LoggerInterface) error { onInit.Store(true); return nil },
		func(l logging.LoggerInterface) { onStop.Store(true) },
		logging.NewLogger(nil),
	)

	task1.Start()
	if !task1.IsRunning() {
		t.Error("Task should be running")
	}
	time.Sleep(1 * time.Second)

	task1.Stop()
	time.Sleep(1 * time.Second)
	if task1.IsRunning() {
		t.Error("Task should be stopped")
	}

	if !onInit.Load().(bool) {
		t.Error("Initialization hook not executed")
	}

	if !onExecution.Load().(bool) {
		t.Error("Main task function not executed")
	}

	if !onStop.Load().(bool) {
		t.Error("After execution function not executed")
	}
}

func TestAsyncTaskPanics(t *testing.T) {
	// Panicking execution
	task1 := NewAsyncTask(
		"testTask1",
		func(l logging.LoggerInterface) error { panic("panic task1") },
		1,
		func(l logging.LoggerInterface) error { return nil },
		func(l logging.LoggerInterface) {},
		logging.NewLogger(nil),
	)

	task1.Start()

	// ---------------------------------

	// Panicking onInit()
	task2 := NewAsyncTask(
		"testTask1",
		func(l logging.LoggerInterface) error { return nil },
		1,
		func(l logging.LoggerInterface) error { panic("panic task2") },
		func(l logging.LoggerInterface) {},
		logging.NewLogger(nil),
	)

	task2.Start()

	// ---------------------------------

	// Panicking onStop()
	task3 := NewAsyncTask(
		"testTask1",
		func(l logging.LoggerInterface) error { return nil },
		1,
		func(l logging.LoggerInterface) error { return nil },
		func(l logging.LoggerInterface) { panic("panic task3") },
		logging.NewLogger(nil),
	)

	task3.Start()
	time.Sleep(time.Second * 2)
	task3.Stop()

	time.Sleep(time.Second * 2)

	if task1.IsRunning() {
		t.Error("Task1 is running and should be stopped")
	}
	if task2.IsRunning() {
		t.Error("Task2 is running and should be stopped")
	}
	if task3.IsRunning() {
		t.Error("Task3 is running and should be stopped")
	}
}

func TestAsyncTaskErrors(t *testing.T) {
	// Error in execution skips one iteration
	res := atomic.Value{}
	res.Store(0)
	task1 := NewAsyncTask(
		"testTask1",
		func(l logging.LoggerInterface) error { res.Store(res.Load().(int) + 1); return errors.New("") },
		1,
		func(l logging.LoggerInterface) error { return nil },
		func(l logging.LoggerInterface) {},
		logging.NewLogger(nil),
	)

	task1.Start()
	time.Sleep(time.Second * 3)
	task1.Stop()

	if res.Load().(int) < 2 {
		t.Error("Task should have executed at least two times")
	}

	res.Store(0)
	task2 := NewAsyncTask(
		"testTask1",
		func(l logging.LoggerInterface) error { res.Store(res.Load().(int) + 1); return nil },
		1,
		func(l logging.LoggerInterface) error { return errors.New("") },
		func(l logging.LoggerInterface) {},
		logging.NewLogger(nil),
	)

	task2.Start()
	time.Sleep(2 * time.Second)
	if res.Load().(int) != 0 {
		t.Error("Task should have never executed if there was an error when calling onInit()")
	}
}

func TestAsyncTaskWakeUp(t *testing.T) {
	res := atomic.Value{}
	res.Store(0)
	task1 := NewAsyncTask(
		"testTask1",
		func(l logging.LoggerInterface) error { res.Store(res.Load().(int) + 1); return nil },
		20,
		func(l logging.LoggerInterface) error { return nil },
		func(l logging.LoggerInterface) {},
		logging.NewLogger(nil),
	)

	task1.Start()
	time.Sleep(time.Second * 1)
	_ = task1.WakeUp()
	time.Sleep(time.Second * 1)
	_ = task1.WakeUp()
	time.Sleep(time.Second * 1)
	_ = task1.WakeUp()
	_ = task1.Stop()

	time.Sleep(time.Second * 2)

	if res.Load().(int) != 4 {
		t.Errorf("Task shuld have executed 4 times. It ran %d times", res)
	}
}
