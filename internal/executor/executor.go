package executor

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/internal/common"
	"github.com/Maru-Yasa/gosong/internal/config"
)

type IExecutor interface {
	Run(cmd string) (string, error)
}

func NewExecutorFromConfig(cfg *config.RemoteHost) (IExecutor, error) {
	switch cfg.Type {
	case common.ExecutorSSH:
		sshExec, err := newSSHExecutor(cfg)
		if err != nil {
			return nil, err
		}
		return sshExec, nil
	case common.ExecutorLocal:
		return newLocalExecutor()
	default:
		return nil, fmt.Errorf("unknown executor type: %s", cfg.Type)
	}
}
