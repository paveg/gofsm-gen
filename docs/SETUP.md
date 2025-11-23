# Development Environment Setup Guide

This guide will help you set up your development environment for contributing to gofsm-gen.

## Prerequisites

### Required Tools

1. **Go 1.23 or later**
   ```bash
   # Check your Go version
   go version

   # Install Go (if needed)
   # macOS
   brew install go

   # Linux (Ubuntu/Debian)
   sudo apt update
   sudo apt install golang-go

   # Windows
   # Download from https://go.dev/dl/
   ```

2. **Git**
   ```bash
   # Check if Git is installed
   git --version

   # Install Git (if needed)
   # macOS
   brew install git

   # Linux (Ubuntu/Debian)
   sudo apt install git

   # Windows
   # Download from https://git-scm.com/download/win
   ```

### Recommended Tools

1. **golangci-lint** (for linting)
   ```bash
   # macOS/Linux
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

   # Windows
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. **staticcheck** (additional static analysis)
   ```bash
   go install honnef.co/go/tools/cmd/staticcheck@latest
   ```

3. **exhaustive** (for switch exhaustiveness checking)
   ```bash
   go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
   ```

4. **make** (optional, for using Makefile commands)
   ```bash
   # macOS
   # Already installed with Xcode Command Line Tools

   # Linux (Ubuntu/Debian)
   sudo apt install build-essential

   # Windows
   # Install via Chocolatey: choco install make
   ```

## Initial Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/gofsm-gen.git
cd gofsm-gen
```

### 2. Install Dependencies

```bash
# Download all dependencies
go mod download

# Verify dependencies
go mod verify
```

### 3. Build the Project

```bash
# Build the CLI tool
go build -o bin/gofsm-gen ./cmd/gofsm-gen

# Or use the shorthand
make build
```

### 4. Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Or use the Makefile
make test
make coverage
```

## Development Workflow

### Test-Driven Development (TDD)

This project follows **Test-Driven Development** as the primary methodology:

1. **Write a failing test first**
   ```bash
   # Create or edit test file
   vim pkg/parser/yaml_test.go

   # Run the specific test to see it fail
   go test ./pkg/parser -run TestYAMLParser_YourNewTest -v
   ```

2. **Implement the minimum code to pass**
   ```bash
   # Edit the implementation
   vim pkg/parser/yaml.go

   # Run the test again
   go test ./pkg/parser -run TestYAMLParser_YourNewTest -v
   ```

3. **Refactor while keeping tests green**
   ```bash
   # Make improvements
   # Run all tests to ensure nothing breaks
   go test ./...
   ```

### Code Quality Checks

Before committing, run these checks:

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run static analysis
go vet ./...
staticcheck ./...

# Check for exhaustive switch statements
exhaustive ./...

# Or use the Makefile
make lint
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem ./benchmarks/

# Run specific benchmark
go test -bench=BenchmarkStateTransition -benchmem ./benchmarks/

# Or use the Makefile
make bench
```

### Working with Examples

```bash
# Generate code from example
./bin/gofsm-gen -spec=examples/order/order.yaml -out=examples/order/fsm.gen.go

# Run example
go run examples/order/main.go
```

## IDE Setup

### VSCode

1. Install the Go extension:
   ```
   ext install golang.go
   ```

2. Recommended settings (`.vscode/settings.json`):
   ```json
   {
     "go.useLanguageServer": true,
     "go.lintTool": "golangci-lint",
     "go.lintOnSave": "package",
     "go.formatTool": "goimports",
     "go.testFlags": ["-v", "-race"],
     "go.coverOnSave": true,
     "editor.formatOnSave": true,
     "[go]": {
       "editor.defaultFormatter": "golang.go"
     }
   }
   ```

### GoLand / IntelliJ IDEA

1. Enable Go modules support:
   - Settings → Go → Go Modules → Enable Go modules integration

2. Configure file watchers:
   - Settings → Tools → File Watchers
   - Add watchers for gofmt and goimports

### Vim/Neovim

1. Install vim-go plugin:
   ```vim
   Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }
   ```

2. Add to your `.vimrc`:
   ```vim
   let g:go_fmt_command = "goimports"
   let g:go_auto_type_info = 1
   let g:go_def_mode='gopls'
   let g:go_info_mode='gopls'
   ```

## Common Development Tasks

### Creating a New Package

```bash
# Create package directory
mkdir -p pkg/mypackage

# Create package file
cat > pkg/mypackage/mypackage.go << 'EOF'
package mypackage

// Package documentation here
EOF

# Create test file
cat > pkg/mypackage/mypackage_test.go << 'EOF'
package mypackage

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
    // Write your test first!
    t.Skip("Not implemented yet")
}
EOF
```

### Adding a New Dependency

```bash
# Add dependency
go get github.com/example/package@v1.2.3

# Update go.mod and go.sum
go mod tidy

# Commit the changes
git add go.mod go.sum
git commit -m "deps: add github.com/example/package v1.2.3"
```

### Debugging

```bash
# Run with debug output
go run -gcflags="all=-N -l" ./cmd/gofsm-gen -spec=example.yaml

# Use delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/gofsm-gen -- -spec=example.yaml
```

### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./benchmarks/
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=. ./benchmarks/
go tool pprof mem.prof

# Generate profile visualization
go tool pprof -http=:8080 cpu.prof
```

## Troubleshooting

### Module-related Issues

```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
rm go.sum
go mod download
```

### Build Issues

```bash
# Clean build cache
go clean -cache

# Rebuild everything
go build -a -v ./...
```

### Test Failures

```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v ./pkg/parser -run TestYAMLParser_ParseStates

# Run tests with race detector
go test -race ./...
```

## Contributing

Before submitting a pull request:

1. Ensure all tests pass: `make test`
2. Ensure linting passes: `make lint`
3. Ensure code is formatted: `go fmt ./...`
4. Update documentation if needed
5. Follow TDD methodology: tests first, then implementation
6. Maintain >90% code coverage

## Getting Help

- Check the [documentation](docs/)
- Read the [CLAUDE.md](CLAUDE.md) file for project guidelines
- Open an issue on GitHub
- Review existing examples in `examples/`

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [golangci-lint Linters](https://golangci-lint.run/usage/linters/)
- [testify Documentation](https://pkg.go.dev/github.com/stretchr/testify)
