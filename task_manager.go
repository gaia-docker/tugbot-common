package common

import (
	log "github.com/Sirupsen/logrus"

	"time"
)

type Task struct {
	Name     string
	Job      func() error
	Interval time.Duration
}

type TaskManager interface {
	RunNewTask(Task)
	StopTasks()
}

type taskManagerImpl struct {
	tasks []chan bool
}

func NewTaskManager() TaskManager {
	return &taskManagerImpl{}
}

func (manager *taskManagerImpl) RunNewTask(task Task) {
	stop := make(chan bool)
	manager.tasks = append(manager.tasks, stop)
	go func(quit chan bool) {
		recurring(task, quit)
	}(stop)
}

func (manager *taskManagerImpl) StopTasks() {
	for _, currTask := range manager.tasks {
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
