package tasks

import (
	"bytes"
	"html/template"
	"log"
	"strings"
	"sync"
)

// Storage - thread-safe type to store and compute tasks
type Storage struct {
	sync.Mutex
	Scheduler *Scheduler     // Schedules workers
	IDs       []uint         // For sorted printing
	Tasks     map[uint]*Task // For id-access
	cid       uint           // Current free id
}

// NewStorage - return initialized Storage
func NewStorage() *Storage {
	return &Storage{
		Scheduler: NewScheduler(),
		Tasks:     make(map[uint]*Task),
	}
}

// AddTask - adding new task
func (s *Storage) AddTask(t *Task) uint {
	s.Lock()
	id := s.cid
	s.cid++
	// TODO: infinite cid generation
	s.IDs = append(s.IDs, id)
	s.Tasks[id] = t
	s.Scheduler.pipe <- t
	s.Unlock()
	return id
}

// CancelTask - canceling task by id
func (s *Storage) CancelTask(id uint) {
	s.Lock()
	if val, ok := s.Tasks[id]; ok {
		if val.State == sWaiting || val.State == sRunning {
			val.cancel()
		} else {
			log.Printf("Can not cancel task: %s (ID: %d)", val.Name, id)
		}
	} else {
		log.Printf("Task with ID: %d not found", id)
	}
	s.Unlock()
}

// DeleteTask - deleting canceled task by id
func (s *Storage) DeleteTask(id uint) {
	s.Lock()
	if val, ok := s.Tasks[id]; ok {
		if val.State == sCanceled {
			delete(s.Tasks, id)
		} else {
			log.Printf("Can not delete task: %s (ID: %d)", val.Name, id)
		}
	} else {
		log.Printf("Task with ID: %d not found", id)
	}
	s.Unlock()
}

// PrintAll - printing task list to template
func (s *Storage) PrintAll() template.HTML {
	var buffer bytes.Buffer
	buffer.WriteString("<ul>")
	s.Lock()
	for _, id := range s.IDs {
		s.Tasks[id].Lock()
		output := "<br>" + strings.Replace(s.Tasks[id].Stdout.String(), "\n", "<br>", -1)
		buffer.WriteString("<li>" + s.Tasks[id].Name + " - " + s.Tasks[id].PrintInfo() + output + "</li>")
		s.Tasks[id].Unlock()
	}
	s.Unlock()
	buffer.WriteString("<ul>")
	return template.HTML(buffer.String())
}
