package common

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"errors"
	"sync"
	"testing"
	"time"
)

func TestRunningTaskReturnsError(t *testing.T) {
	ok := false
	assert.Error(t, recurring(nil,
		Task{
			Name: "test",
			Job: func(interface{}) error {
				ok = true

				return errors.New("expected :)")
			},
			Interval: 0},
	))
	assert.True(t, ok)
}

func TestStopRunningTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	ok := false
	go func() {
		defer wg.Done()
		recurring(ctx, Task{
			Name: "test",
			Job: func(param interface{}) error {
				ok = true
				assert.Equal(t, "check me", param)

				return nil
			},
			JobParam: "check me",
			Interval: time.Second * 10,
		})
	}()
	cancel()
	wg.Wait()
	assert.True(t, ok)
}

func TestTaskManagerRunTasks(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	manager := NewTaskManager()
	ok1, ok2 := false, false
	assert.True(t, manager.RunNewRecurringTask(Task{
		ID:   "t1-id",
		Name: "t1",
		Job: func(interface{}) error {
			defer wg.Done()
			ok1 = true

			return nil
		},
		Interval: time.Second * 10,
	}))
	assert.True(t, manager.RunNewRecurringTask(Task{
		ID:   "t2-id",
		Name: "t2",
		Job: func(interface{}) error {
			defer wg.Done()
			ok2 = true

			return nil
		},
		Interval: time.Second * 10,
	}))
	wg.Wait()
	manager.StopAllTasks()
	assert.True(t, ok1)
	assert.True(t, ok2)
}

func TestTaskManagerRunTasks_TaskAlreadyExist(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	manager := NewTaskManager()
	ok1 := false
	assert.True(t, manager.RunNewRecurringTask(Task{
		ID:   "t1-id",
		Name: "t1",
		Job: func(interface{}) error {
			defer wg.Done()
			ok1 = true

			return nil
		},
		Interval: time.Second * 10,
	}))
	assert.False(t, manager.RunNewRecurringTask(Task{
		ID: "t1-id",
	}))
	wg.Wait()
	manager.StopAllTasks()
	assert.True(t, ok1)
}
