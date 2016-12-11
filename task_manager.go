package common

import (
	log "github.com/Sirupsen/logrus"

	"time"
)

type Task struct {
	ID       string
	Name     string
	Job      func() error
	Interval time.Duration
}

type TaskManager interface {
	RunNewTask(Task)
	StopTasks()
}

type taskManagerImpl struct {
	taskIdToStopChannel map[string]chan bool
}

func NewTaskManager() TaskManager {
	return &taskManagerImpl{taskIdToStopChannel: make(map[string]chan bool)}
}

func (manager *taskManagerImpl) RunNewTask(task Task) {
	if _, ok := manager.taskIdToStopChannel[task.ID]; !ok {
		stop := make(chan bool)
		manager.taskIdToStopChannel[task.ID] = stop
		go func(quit chan bool) {
			recurring(task, quit)
		}(stop)
	}
}

func (manager *taskManagerImpl) StopTasks() {
	for _, currTask := range manager.taskIdToStopChannel {
		currTask <- true
	}
}

func recurring(task Task, quit chan bool) error {
	for {
		log.Debugf("Running task %s...", task.Name)
		if err := task.Job(); err != nil {
			log.Errorf("Task %s returned an error: %v", task.Name, err)

			return err
		}
		timer := time.NewTimer(task.Interval)
		select {
		case <-quit:
			timer.Stop()
			log.Debugf("Task %s stopped", task.Name)

			return nil
		case <-timer.C:
			log.Debugf("Timer expired. Interval: %v Task: %s", task.Interval, task.Name)
		}
	}

	return nil
}
