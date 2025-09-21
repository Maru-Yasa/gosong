package executor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Maru-Yasa/gosong/internal/common"
	"github.com/Maru-Yasa/gosong/internal/config"
)

type Executor interface {
	RunRaw(cmd string) (string, error)
	Run(cmd string, cwd string) error
	GetName() string
}

func GetLastIDFromHost(exec Executor, appPath string) (int, error) {
	script := fmt.Sprintf(`
		if [ ! -d %s/releases ]; then
			echo 0
			exit 0
		fi

		last_id=$(find %s/releases -maxdepth 1 -type d -printf "%%f\n" \
			| grep -E '^[0-9]+$' \
			| sort -n \
			| tail -1)

		if [ -z "$last_id" ]; then
			echo 0
		else
			echo "$last_id"
		fi
	`, appPath, appPath)

	outputFromExec, err := exec.RunRaw(script)
	if err != nil {
		return 0, err
	}

	id, convErr := strconv.Atoi(strings.TrimSpace(outputFromExec))
	if convErr != nil {
		return 0, convErr
	}

	return id, nil
}

func NewExecutorFromConfig(name string, cfg *config.RemoteHost) (Executor, error) {
	switch cfg.Type {
	case common.ExecutorSSH:
		sshExec, err := newSSHExecutor(name, cfg)
		if err != nil {
			return nil, err
		}
		return sshExec, nil
	case common.ExecutorLocal:
		return newLocalExecutor(name)
	default:
		return nil, fmt.Errorf("unknown executor type: %s", cfg.Type)
	}
}
