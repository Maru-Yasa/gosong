package main

import (
	"github.com/Maru-Yasa/gosong/cmd"
	"github.com/Maru-Yasa/gosong/pkg/logger"
)

func main() {
	logger.SetDefaultLogger(
		logger.NewMultiLogger(
			logger.NewConsoleLogger(),
			logger.NewFileLogger(),
		),
	)

	cmd.Execute()
}
