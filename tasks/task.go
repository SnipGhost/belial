package tasks

// Task - test struct
type Task struct {
	ID     int
	Name   string
	Active bool
}

// PrintInfo - for template generation
func (t *Task) PrintInfo() string {
	if !t.Active {
		return "Ready"
	}
	return "Active"
}
