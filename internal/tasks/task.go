package tasks

import "fmt"

type Step struct {
	Cd   string `yaml:"cd,omitempty"`
	Run  string `yaml:"run,omitempty"`
	Task string `yaml:"task,omitempty"`
}

type Task struct {
	Description string `yaml:"description,omitempty"`
	Steps       []Step
}

var BuiltInTasks = map[string]Task{}

func RegisterTask(name string, task Task) {
	BuiltInTasks[name] = task
}

func FindTask(name string, uTasks map[string]Task) (Task, error) {
	if t, ok := uTasks[name]; ok {
		return t, nil
	}
	if t, ok := BuiltInTasks[name]; ok {
		return t, nil
	}
	return Task{}, fmt.Errorf("task '%s' not found", name)
}
