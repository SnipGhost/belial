package tasks

// Task - test struct
type Task struct {
	ID     int    // ID
	Name   string // Task title (?)
	Active bool   // Ready or not
}

// PrintInfo - for template generation
func (t *Task) PrintInfo() string {
	if !t.Active {
		return "Ready"
	}
	return "Active"
}
