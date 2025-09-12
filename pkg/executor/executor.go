package executor

import "fmt"

type ExecutorType string

const (
	ExecutorSSH   ExecutorType = "ssh"
	ExecutorLocal ExecutorType = "local"
)

type Executor interface {
	Run(cmd string) (string, error)
}

func ExecuteCommand(exec Executor, cmd string) {
	out, err := exec.Run(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out)
	}
}

func (e ExecutorType) IsValid() bool {
	switch e {
	case ExecutorSSH, ExecutorLocal:
		return true
	}
	return false
}
