PHONY : clean test build install help

all: help

clean:
	@echo "Cleaning up..."
	@go clean -testcache

test:
	@echo "Running tests..."
	@go test ./...

build:
	@echo "Building the project..."
	@go build -o gomakase cmd/gomakase/main.go

install:
	@echo "Installing the project..."
	@go install ./cmd/gomakase

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  clean - Clean up"
	@echo "  test - Run tests"
	@echo "  build - Build the project (not implemented)"
	@echo "  install - Install the project (not implemented)"


