package common

type ExecutorType string

const (
	ExecutorSSH   ExecutorType = "ssh"
	ExecutorLocal ExecutorType = "local"
)

type Step struct {
	Cd   string `yaml:"cd,omitempty"`
	Run  string `yaml:"run,omitempty"`
	Task string `yaml:"task,omitempty"`
}

type UTask struct {
	Description string `yaml:"description,omitempty"`
	Steps       []Step
}

func (e ExecutorType) IsValid() bool {
	switch e {
	case ExecutorSSH, ExecutorLocal:
		return true
	}
	return false
}
