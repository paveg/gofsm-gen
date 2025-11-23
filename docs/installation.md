# Installation Guide

This guide covers all methods of installing and setting up gofsm-gen.

## Prerequisites

- **Go 1.21 or higher**: gofsm-gen requires Go 1.21+ for generics and other modern features
- **Git**: For cloning the repository (optional)

Check your Go version:

```bash
go version
```

## Installation Methods

### Method 1: Install via go install (Recommended)

The easiest way to install gofsm-gen is using `go install`:

```bash
go install github.com/yourusername/gofsm-gen/cmd/gofsm-gen@latest
```

This will install the `gofsm-gen` binary to your `$GOPATH/bin` directory.

Verify the installation:

```bash
gofsm-gen -version
```

### Method 2: Install from Source

Clone the repository and build from source:

```bash
# Clone the repository
git clone https://github.com/yourusername/gofsm-gen.git
cd gofsm-gen

# Build the CLI tool
go build -o bin/gofsm-gen ./cmd/gofsm-gen

# Optional: Install to $GOPATH/bin
go install ./cmd/gofsm-gen
```

### Method 3: Download Pre-built Binaries

Download pre-built binaries from the [releases page](https://github.com/yourusername/gofsm-gen/releases):

1. Navigate to the latest release
2. Download the binary for your platform:
   - `gofsm-gen-linux-amd64`
   - `gofsm-gen-darwin-amd64` (macOS Intel)
   - `gofsm-gen-darwin-arm64` (macOS Apple Silicon)
   - `gofsm-gen-windows-amd64.exe`

3. Make the binary executable (Linux/macOS):

```bash
chmod +x gofsm-gen-*
```

4. Move to a directory in your PATH:

```bash
sudo mv gofsm-gen-* /usr/local/bin/gofsm-gen
```

### Method 4: Docker

Run gofsm-gen using Docker:

```bash
# Pull the image
docker pull ghcr.io/yourusername/gofsm-gen:latest

# Run code generation
docker run --rm -v $(pwd):/workspace \
  ghcr.io/yourusername/gofsm-gen:latest \
  -spec=/workspace/fsm.yaml \
  -out=/workspace/fsm.gen.go
```

Create an alias for easier usage:

```bash
alias gofsm-gen='docker run --rm -v $(pwd):/workspace ghcr.io/yourusername/gofsm-gen:latest'
```

## Setting Up Your Project

### 1. Initialize Go Module

If you haven't already, initialize a Go module in your project:

```bash
go mod init github.com/yourusername/yourproject
```

### 2. Create FSM Definition Directory

Create a directory to store your state machine definitions:

```bash
mkdir -p statemachines
```

### 3. Install Static Analyzer (Optional but Recommended)

To enable exhaustiveness checking, install the `exhaustive` static analyzer:

```bash
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

Verify installation:

```bash
exhaustive -version
```

### 4. Install Additional Tools (Optional)

For enhanced development experience:

**Graphviz** (for diagram visualization):

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Windows
choco install graphviz
```

**Staticcheck** (for additional static analysis):

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
```

## IDE Integration

### VSCode

1. Install the Go extension:
   - Open Extensions (Cmd+Shift+X / Ctrl+Shift+X)
   - Search for "Go" by Go Team at Google
   - Click Install

2. Configure settings for generated files:

Add to `.vscode/settings.json`:

```json
{
  "go.buildTags": "",
  "go.lintTool": "staticcheck",
  "go.lintFlags": [],
  "files.watcherExclude": {
    "**/*.gen.go": true
  },
  "go.formatTool": "gofmt",
  "editor.formatOnSave": true
}
```

3. Add a task for code generation:

Create `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Generate FSM",
      "type": "shell",
      "command": "gofsm-gen",
      "args": [
        "-spec=${file}",
        "-out=${fileDirname}/${fileBasenameNoExtension}.gen.go"
      ],
      "problemMatcher": [],
      "group": {
        "kind": "build",
        "isDefault": true
      }
    }
  ]
}
```

### GoLand / IntelliJ IDEA

1. Configure File Watcher for automatic generation:
   - Go to Settings → Tools → File Watchers
   - Click "+" to add new watcher
   - Configure:
     - **File type**: YAML
     - **Scope**: Project Files
     - **Program**: `gofsm-gen`
     - **Arguments**: `-spec=$FilePath$ -out=$FileNameWithoutExtension$.gen.go`
     - **Working directory**: `$FileDir$`

2. Mark generated files:
   - Right-click on `*.gen.go` files
   - Select "Mark as" → "Generated Sources Root"

### Vim / Neovim

Add to your `.vimrc` or `init.vim`:

```vim
" Auto-generate FSM on save
autocmd BufWritePost *.fsm.yaml !gofsm-gen -spec=% -out=%:r.gen.go
```

## Verifying Installation

Run these commands to verify everything is set up correctly:

```bash
# Check gofsm-gen version
gofsm-gen -version

# Check Go version
go version

# Check exhaustive analyzer (if installed)
exhaustive -version

# Check staticcheck (if installed)
staticcheck -version
```

## Updating gofsm-gen

### Update via go install

```bash
go install github.com/yourusername/gofsm-gen/cmd/gofsm-gen@latest
```

### Update from source

```bash
cd gofsm-gen
git pull origin main
go install ./cmd/gofsm-gen
```

### Check for updates

```bash
gofsm-gen -check-update
```

## Troubleshooting

### gofsm-gen: command not found

Ensure `$GOPATH/bin` is in your PATH:

```bash
# Add to ~/.bashrc, ~/.zshrc, or equivalent
export PATH=$PATH:$(go env GOPATH)/bin
```

Then reload your shell:

```bash
source ~/.bashrc  # or ~/.zshrc
```

### Permission denied

If you get permission errors on Linux/macOS:

```bash
chmod +x $(which gofsm-gen)
```

### Old version still running

Clear Go's build cache:

```bash
go clean -cache
go install github.com/yourusername/gofsm-gen/cmd/gofsm-gen@latest
```

### Module errors

Update your Go modules:

```bash
go get -u github.com/yourusername/gofsm-gen
go mod tidy
```

## Next Steps

- Read the [Basic Usage Guide](usage.md) to start creating state machines
- Explore [examples/](../examples/) for complete examples
- Review the [YAML Definition Reference](yaml-reference.md) for detailed syntax

## Uninstalling

To remove gofsm-gen:

```bash
# Remove binary
rm $(which gofsm-gen)

# Clean Go cache
go clean -cache -modcache
```
