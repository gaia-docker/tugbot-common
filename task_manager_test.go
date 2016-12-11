package common

import (
	"github.com/stretchr/testify/assert"

	"errors"
	"sync"
	"testing"
	"time"
)

func TestRunningTaskReturnsError(t *testing.T) {
	ok := false
	assert.Error(t, recurring(
		Task{
			Name: "test",
			Job: func() error {
				ok = true

				return errors.New("expected :)")
			},
			Interval: 0},
		nil))
	assert.True(t, ok)
}

func TestStopRunningTask(t *testing.T) {
	quit := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	ok := false
	go func() {
		defer wg.Done()
		recurring(Task{
			Name: "test",
			Job: func() error {
				ok = true

				return nil
			},
			Interval: time.Second * 10,
		}, quit)
	}()
	quit <- true
	wg.Wait()
	assert.True(t, ok)
}

func TestTaskManagerRunTasks(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	manager := NewTaskManager()
	ok1, ok2 := false, false
	manager.RunNewTask(Task{
		Name: "t1",
		Job: func() error {
			defer wg.Done()
			ok1 = true

			return nil
		},
		Interval: time.Second * 10,
	})
	manager.RunNewTask(Task{
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
