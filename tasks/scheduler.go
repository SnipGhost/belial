package tasks

import (
	"log"
	"time"
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
	t.State = sRunning
	log.Println("Task:", t.Name, "running")
	count := 5
	for i := 0; i < count; i++ {
		select {
		case <-t.ctx.Done():
			log.Println("Task:", t.Name, "finished by ctx")
			t.State = sCanceled
			return
		default:
		}
		// Emulating work ...
		// TODO: write real algorithm
		time.Sleep(1 * time.Second)
		log.Println("Task:", t.Name, "status", i, "/", count)
	}
	log.Println("Task:", t.Name, "done")
	t.State = sReady
}

// Run - waiting messages from pipe and run workers
func (s *Scheduler) Run() {
	log.Println("Running scheduler ...")
	for task := range s.pipe {
		go s.runWorker(task)
	}
}
