.DEFAULT_GOAL := build

.PHONY: fmt vet build 

fmt:
	@echo "Running go fmt..."
	@go fmt ./...

vet: fmt
	@echo "Running go vet..."
	@go vet ./...

build: vet
	@echo "Building the project..."
	@go build -ldflags "-s -w" -o bin/myapp ./main.go