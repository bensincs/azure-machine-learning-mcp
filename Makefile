# Azure ML MCP Server Makefile

.PHONY: build test test-verbose clean lint format deps help run

# Build the application
build:
	@echo "Building Azure ML MCP Server..."
	go build -o bin/mcp-server ./cmd/mcp-server

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run specific tests
test-helpers:
	@echo "Running helper tests..."
	go test -v ./internal/helpers/

test-azure:
	@echo "Running Azure client tests..."
	go test -v ./internal/azure/

test-tools:
	@echo "Running tools tests..."
	go test -v ./internal/tools/

test-server:
	@echo "Running server tests..."
	go test -v ./internal/server/

# Lint the code
lint:
	@echo "Running linter..."
	golangci-lint run

# Format the code
format:
	@echo "Formatting code..."
	go fmt ./...

# Run the server
run:
	@echo "Running Azure ML MCP Server..."
	go run ./cmd/mcp-server

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run all checks (format, lint, test)
check: format lint test

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o bin/mcp-server-linux-amd64 ./cmd/mcp-server
	GOOS=darwin GOARCH=amd64 go build -o bin/mcp-server-darwin-amd64 ./cmd/mcp-server
	GOOS=darwin GOARCH=arm64 go build -o bin/mcp-server-darwin-arm64 ./cmd/mcp-server
	GOOS=windows GOARCH=amd64 go build -o bin/mcp-server-windows-amd64.exe ./cmd/mcp-server

# Benchmark tests
bench:
	@echo "Running benchmarks..."
	go test -bench=. ./...

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  deps           - Install dependencies"
	@echo "  test           - Run tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-helpers   - Run helper tests only"
	@echo "  test-azure     - Run Azure client tests only"
	@echo "  test-tools     - Run tools tests only"
	@echo "  test-server    - Run server tests only"
	@echo "  lint           - Run linter"
	@echo "  format         - Format code"
	@echo "  run            - Run the server"
	@echo "  clean          - Clean build artifacts"
	@echo "  check          - Run format, lint, and test"
	@echo "  build-all      - Build for multiple platforms"
	@echo "  bench          - Run benchmarks"
	@echo "  help           - Show this help message"
