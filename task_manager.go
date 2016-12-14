package common

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"time"
)

type Task struct {
	ID       string
	Name     string
	Job      func(interface{}) error
	JobParam interface{}
	Interval time.Duration
}

type TaskManager interface {
	RunNewRecurringTask(Task) bool
	StopAllTasks()
}

type taskManagerImpl struct {
	taskIdToCancel map[string]context.CancelFunc
}

func NewTaskManager() TaskManager {
	return &taskManagerImpl{taskIdToCancel: make(map[string]context.CancelFunc)}
}

func (manager *taskManagerImpl) RunNewRecurringTask(task Task) bool {
	ret := false
	if _, ok := manager.taskIdToCancel[task.ID]; !ok {
		ctx, cancel := context.WithCancel(context.Background())
		manager.taskIdToCancel[task.ID] = cancel
		go func(ctx context.Context) {
			recurring(ctx, task)
		}(ctx)
		ret = true
	}

	return ret
}

func (manager *taskManagerImpl) StopAllTasks() {
	for _, currTaskCancel := range manager.taskIdToCancel {
		currTaskCancel()
	}
}

func recurring(ctx context.Context, task Task) error {
	for {
		log.Debugf("Running task %s...", task.Name)
		if err := task.Job(task.JobParam); err != nil {
			log.Errorf("Task %s returned an error: %v", task.Name, err)

			return err
		}
		timer := time.NewTimer(task.Interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Debugf("Task %s stopped", task.Name)

			return nil
		case <-timer.C:
			log.Debugf("Timer expired. Interval: %v Task: %s", task.Interval, task.Name)
		}
	}

	return nil
}
