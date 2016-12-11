package common

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunningTaskReturnsError(t *testing.T) {
	assert.Error(t, Recurring(
		Task{
			Name:     "test",
			Job:      func() error { return errors.New("expected :)") },
			Interval: 0},
		nil))
}

func TestStopRunningTask(t *testing.T) {
	quit := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Recurring(Task{
			Name: "test",
			Job: func() error {
				return nil
			},
			Interval: time.Second * 10,
		}, quit)
	}()
	fmt.Println("helooo effi ")
	quit <- true
	wg.Wait()
}
