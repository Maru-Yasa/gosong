package main

import (
	"log"

	"github.com/Maru-Yasa/gosong/cmd"
	"github.com/Maru-Yasa/gosong/pkg/logger"
)

func main() {
	if err := logger.Init(); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	cmd.Execute()
}
