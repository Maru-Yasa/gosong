BINARY_NAME=gosong
BUILD_DIR=build
VERSION?=$(shell git describe --tags --abbrev=0)

.PHONY: install build build-all clean release

install:
	go mod tidy

build:
	rm -rf $(BUILD_DIR)
	@echo "Building for current platform..."
	mkdir -p $(BUILD_DIR)/current
	go build -o $(BUILD_DIR)/current/$(BINARY_NAME) main.go

build-all:
	rm -rf $(BUILD_DIR)
	@echo "Building for all major platforms..."
	mkdir -p $(BUILD_DIR)
	GOOS=linux   GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64   main.go
	GOOS=linux   GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64   main.go
	GOOS=darwin  GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64  main.go
	GOOS=darwin  GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64  main.go
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go

clean:
	rm -rf $(BUILD_DIR)

release: build-all
	@echo "Packaging release for version $(VERSION)..."
	cd $(BUILD_DIR) && \
	zip $(BINARY_NAME)-linux-amd64.zip   $(BINARY_NAME)-linux-amd64 && \
	zip $(BINARY_NAME)-linux-arm64.zip   $(BINARY_NAME)-linux-arm64 && \
	zip $(BINARY_NAME)-darwin-amd64.zip  $(BINARY_NAME)-darwin-amd64 && \
	zip $(BINARY_NAME)-darwin-arm64.zip  $(BINARY_NAME)-darwin-arm64 && \
	zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Release artifacts ready in $(BUILD_DIR)/"
