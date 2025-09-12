package common

type StepType string

const (
	StepTypeCommand StepType = "command"
	StepTypeTask    StepType = "task"
)

type Step struct {
	Type    StepType `yaml:"type"`
	Command string   `yaml:"command,omitempty"`
	Task    string   `yaml:"task,omitempty"`
}

type Task struct {
	Description string `yaml:"description,omitempty"`
	Steps       []Step `yaml:"steps"`
}

type ExecutorType string

const (
	ExecutorSSH   ExecutorType = "ssh"
	ExecutorLocal ExecutorType = "local"
)

func (e ExecutorType) IsValid() bool {
	switch e {
	case ExecutorSSH, ExecutorLocal:
		return true
	}
	return false
}
