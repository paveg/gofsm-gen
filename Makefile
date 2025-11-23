.PHONY: help build test test-verbose test-race coverage lint fmt vet clean bench install run examples

# Default target
help:
	@echo "gofsm-gen Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build         - Build the CLI tool"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-race     - Run tests with race detector"
	@echo "  coverage      - Generate and view test coverage report"
	@echo "  lint          - Run golangci-lint"
	@echo "  fmt           - Format code with gofmt and goimports"
	@echo "  vet           - Run go vet and staticcheck"
	@echo "  clean         - Remove build artifacts and caches"
	@echo "  bench         - Run benchmarks"
	@echo "  install       - Install the CLI tool"
	@echo "  run           - Run the CLI tool (use ARGS=... to pass arguments)"
	@echo "  examples      - Generate code for all examples"
	@echo "  all           - Run fmt, vet, lint, test, and build"

# Build the CLI tool
build:
	@echo "Building gofsm-gen..."
	@mkdir -p bin
	@go build -o bin/gofsm-gen ./cmd/gofsm-gen
	@echo "Build complete: bin/gofsm-gen"

# Run all tests
test:
	@echo "Running tests..."
	@go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	@go test -v ./...

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@go test -race ./...

# Generate and view coverage report
coverage:
	@echo "Generating coverage report..."
	@go test -cover -coverprofile=coverage.out ./...
	@echo "Coverage summary:"
	@go tool cover -func=coverage.out | tail -1
	@echo "Opening coverage report in browser..."
	@go tool cover -html=coverage.out

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --timeout=5m

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@which goimports > /dev/null && goimports -w . || echo "goimports not installed, skipping"

# Run static analysis
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "Running staticcheck..."
	@which staticcheck > /dev/null && staticcheck ./... || echo "staticcheck not installed, skipping"

# Clean build artifacts and caches
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf dist/
	@rm -f coverage.out coverage.html
	@rm -f *.prof
	@go clean -cache -testcache -modcache
	@echo "Clean complete"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./benchmarks/

# Install the CLI tool to GOPATH/bin
install:
	@echo "Installing gofsm-gen..."
	@go install ./cmd/gofsm-gen
	@echo "Installed to $(shell go env GOPATH)/bin/gofsm-gen"

# Run the CLI tool
run: build
	@./bin/gofsm-gen $(ARGS)

# Generate code for all examples
examples: build
	@echo "Generating code for examples..."
	@for dir in examples/*/; do \
		if [ -f $$dir/*.yaml ]; then \
			echo "Generating $$dir..."; \
			./bin/gofsm-gen -spec=$$dir/*.yaml -out=$$dir/fsm.gen.go; \
		fi \
	done
	@echo "Examples generated"

# Run all checks and build
all: fmt vet lint test build
	@echo "All checks passed!"

# Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify
	@echo "Dependencies ready"

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated"
