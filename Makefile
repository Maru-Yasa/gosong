# Makefile for gosong


BINARY_NAME=gosong
BUILD_DIR=build

.PHONY: install build clean

install:
	go mod tidy


build:
	@echo "Building for current platform..."
	mkdir -p $(BUILD_DIR)/current
	go build -o $(BUILD_DIR)/current/$(BINARY_NAME) main.go


build-all:
	@echo "Building for all major platforms..."
	mkdir -p $(BUILD_DIR)/linux-amd64 $(BUILD_DIR)/darwin-amd64 $(BUILD_DIR)/darwin-arm64 $(BUILD_DIR)/windows-amd64
	GOOS=linux   GOARCH=amd64 go build -o $(BUILD_DIR)/linux-amd64/$(BINARY_NAME)   main.go
	GOOS=darwin  GOARCH=amd64 go build -o $(BUILD_DIR)/darwin-amd64/$(BINARY_NAME)  main.go
	GOOS=darwin  GOARCH=arm64 go build -o $(BUILD_DIR)/darwin-arm64/$(BINARY_NAME)  main.go
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe main.go

	rm -rf $(BUILD_DIR)
