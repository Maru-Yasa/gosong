package common

type Step struct {
	Cd   string `yaml:"cd,omitempty"`
	Run  string `yaml:"run,omitempty"`
	Task string `yaml:"task,omitempty"`
}

type Task struct {
	Description string `yaml:"description,omitempty"`
	Steps       []Step
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
