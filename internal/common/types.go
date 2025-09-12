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
