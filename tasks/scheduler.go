package tasks

import (
	"log"
)

// Scheduler - contains pipe
// TODO: make many priority pipes
type Scheduler struct {
	pipe chan *Task
}

// NewScheduler - return initialized scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		pipe: make(chan *Task, 10),
	}
}

func (s *Scheduler) runWorker(t *Task) {
	t.Lock()
	t.State = sRunning
	log.Println("Task:", t.Name, "running")
	t.Unlock()
	errcode := t.check()
	if errcode == 0 {
		t.Lock()
		log.Println("Task:", t.Name, "done")
		t.State = sReady
		t.Unlock()
	}
}

// Run - waiting messages from pipe and run workers
func (s *Scheduler) Run() {
	log.Println("Running scheduler ...")
	for task := range s.pipe {
		if task.State == sWaiting {
			go s.runWorker(task)
		}
	}
}
