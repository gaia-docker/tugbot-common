package common

import (
	"errors"
	"testing"

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

//func TestStopRunningTask(t *testing.T) {
//	quit := make(chan bool)
//	var wg sync.WaitGroup
//	wg.Add(1)
//	go Recurring(Task{
//		Name: "test",
//		Job: func() error {
//			defer wg.Done()
//			quit <- true
//
//			return nil
//		},
//		Interval: time.Second * 10,
//	}, quit)
//	wg.Wait()
//}
