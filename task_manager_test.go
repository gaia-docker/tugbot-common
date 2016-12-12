package common

import (
	"github.com/stretchr/testify/assert"

	"errors"
	"golang.org/x/net/context"
	"sync"
	"testing"
	"time"
)

func TestRunningTaskReturnsError(t *testing.T) {
	ok := false
	assert.Error(t, recurring(nil,
		Task{
			Name: "test",
			Job: func() error {
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
			Job: func() error {
				ok = true

				return nil
			},
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
	manager.RunNewRecurringTask(Task{
		ID:   "t1-id",
		Name: "t1",
		Job: func() error {
			defer wg.Done()
			ok1 = true

			return nil
		},
		Interval: time.Second * 10,
	})
	manager.RunNewRecurringTask(Task{
		ID:   "t2-id",
		Name: "t2",
		Job: func() error {
			defer wg.Done()
			ok2 = true

			return nil
		},
		Interval: time.Second * 10,
	})
	wg.Wait()
	manager.StopTasks()
	assert.True(t, ok1)
	assert.True(t, ok2)
}
