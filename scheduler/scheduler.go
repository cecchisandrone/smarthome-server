package scheduler

import (
	"github.com/jasonlvhit/gocron"
)

type TaskFunc func()

type SchedulerManager struct {
}

func (s *SchedulerManager) Start() {
	go runTaskScheduler()
}

func (s *SchedulerManager) ScheduleExecution(n uint64, task TaskFunc) {
	gocron.Every(n).Seconds().Do(task)
}

func runTaskScheduler() {
	<-gocron.Start()
}
