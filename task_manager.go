package common

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

type Task struct {
	Name     string
	Job      func() error
	Interval time.Duration
}

func Recurring(task Task, quit chan bool) error {
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
