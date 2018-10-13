package tasks

import (
	"context"
	"time"
)

type statetype uint8

const (
	sWaiting  statetype = iota // State before starting the worker
	sRunning                   // State after starting the worker
	sReady                     // State after successful worker completion
	sCanceled                  // State after context canceling
)

// Task - test struct
type Task struct {
	Name   string             // Task title
	State  statetype          // Ready or not
	ctx    context.Context    // Context
	cancel context.CancelFunc // Function to canceling context
}

// NewTask - return initialized Task
func NewTask(name string, timeout uint) *Task {
	workTime := time.Duration(timeout) * time.Millisecond
	ctx, finish := context.WithTimeout(context.Background(), workTime)
	return &Task{
		ctx:    ctx,
		cancel: finish,
		Name:   name,
		State:  sWaiting,
	}
}

// PrintInfo - for template generation
func (t *Task) PrintInfo() string {
	switch t.State {
	case sWaiting:
		return "Scheduling"
	case sRunning:
		return "Running"
	case sReady:
		return "Ready"
	case sCanceled:
		return "Canceled"
	default:
		return "Unknown"
	}
}
