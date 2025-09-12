package common

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

type Step struct {
	Command string `yaml:"command"`
}

type Task struct {
	Name  string
	Steps []Step `yaml:"steps"`
}
