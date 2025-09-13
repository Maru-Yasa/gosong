package cmd

import (
	"fmt"
	"os"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
)

func main() {
	cfg, err := config.Load("docs/example_v2.yaml")
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	exec := &executor.SimpleExecutor{}
	if err := executor.RunTaskSimple(cfg, exec, "deploy"); err != nil {
		fmt.Println("Execution error:", err)
		os.Exit(1)
	}
}
