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
	t.Lock()
	t.State = sRunning
	log.Println("Task:", t.Name, "running")
	t.Unlock()
	count := 8
	for i := 0; i < count; i++ {
		select {
		case <-t.ctx.Done():
			t.Lock()
			log.Println("Task:", t.Name, "finished by ctx")
			t.State = sCanceled
			t.Unlock()
			return
		default:
		}
		// Emulating work ...
		// TODO: write real algorithm
		time.Sleep(1 * time.Second)
		t.Lock()
		log.Println("Task:", t.Name, "progress:", i, "/", count)
		t.Unlock()
	}
	t.Lock()
	log.Println("Task:", t.Name, "done")
	t.State = sReady
	t.Unlock()
}

// Run - waiting messages from pipe and run workers
func (s *Scheduler) Run() {
	log.Println("Running scheduler ...")
	for task := range s.pipe {
		go s.runWorker(task)
	}
}
