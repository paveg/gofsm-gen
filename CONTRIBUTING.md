# Contributing to gofsm-gen

Thank you for your interest in contributing to gofsm-gen! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Testing Guidelines](#testing-guidelines)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Project Structure](#project-structure)
- [Communication](#communication)

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors, regardless of experience level, background, or identity.

### Expected Behavior

- Be respectful and considerate
- Accept constructive criticism gracefully
- Focus on what is best for the project
- Show empathy towards other contributors

### Unacceptable Behavior

- Harassment or discriminatory language
- Trolling or insulting comments
- Publishing others' private information
- Other unprofessional conduct

## Getting Started

### Prerequisites

- **Go 1.21+**: Required for development
- **Git**: For version control
- **Make**: For running build tasks (optional)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

```bash
git clone https://github.com/YOUR_USERNAME/gofsm-gen.git
cd gofsm-gen
```

3. Add upstream remote:

```bash
git remote add upstream https://github.com/yourusername/gofsm-gen.git
```

4. Verify remotes:

```bash
git remote -v
```

### Install Dependencies

```bash
# Download Go dependencies
go mod download

# Install development tools
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

### Build the Project

```bash
# Build the CLI tool
go build -o bin/gofsm-gen ./cmd/gofsm-gen

# Verify it works
./bin/gofsm-gen -version
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./pkg/parser/...
```

## Development Workflow

We follow **Test-Driven Development (TDD)** as our primary methodology.

### TDD Workflow

1. **Write a failing test first**

```bash
# Create or edit test file
vim pkg/parser/yaml_test.go

# Write a test that defines the expected behavior
# Run it to ensure it fails
go test ./pkg/parser -run TestYAMLParser_NewFeature
```

2. **Write minimal code to make it pass**

```bash
# Implement the feature
vim pkg/parser/yaml.go

# Run the test again
go test ./pkg/parser -run TestYAMLParser_NewFeature
```

3. **Refactor while keeping tests green**

```bash
# Improve the implementation
# Run all tests to ensure nothing breaks
go test ./...
```

4. **Commit test and implementation together**

```bash
git add pkg/parser/yaml.go pkg/parser/yaml_test.go
git commit -m "feat(parser): add support for new feature"
```

### Feature Development Process

1. **Create a feature branch**

```bash
git checkout -b feature/your-feature-name
```

2. **Write tests first** (TDD approach)

3. **Implement the feature**

4. **Ensure all tests pass**

```bash
go test ./...
```

5. **Run static analysis**

```bash
go vet ./...
staticcheck ./...
```

6. **Update documentation** if needed

7. **Commit and push**

```bash
git add .
git commit -m "feat: description of feature"
git push origin feature/your-feature-name
```

8. **Open a pull request**

## Testing Guidelines

### Test Quality Standards

**IMPORTANT**: Write meaningful tests that verify actual behavior. Avoid anti-patterns.

#### DO NOT:

- ❌ Write tests with hardcoded magic values that don't represent real use cases
- ❌ Create tests that assert `result == result` or other tautologies
- ❌ Test implementation details instead of behavior
- ❌ Write brittle tests that break with any refactoring
- ❌ Use arbitrary test data that doesn't reflect domain knowledge

#### DO:

- ✅ Test real-world scenarios and edge cases
- ✅ Use domain-meaningful test data (e.g., realistic FSM definitions)
- ✅ Verify behavior, not implementation
- ✅ Write tests that document expected behavior clearly
- ✅ Use table-driven tests with descriptive test cases
- ✅ Test error conditions and boundary cases explicitly

### Test Organization

**Good Test Example**:

```go
func TestYAMLParser_ParseOrderStateMachine(t *testing.T) {
    yaml := `
machine:
  name: OrderStateMachine
  initial: pending
states:
  - name: pending
  - name: approved
  - name: shipped
events:
  - approve
  - ship
transitions:
  - from: pending
    to: approved
    on: approve
`
    parser := NewYAMLParser()
    model, err := parser.Parse(strings.NewReader(yaml))

    require.NoError(t, err)
    assert.Equal(t, "OrderStateMachine", model.Name)
    assert.Equal(t, "pending", model.Initial)
    assert.Len(t, model.States, 3)
}
```

**Bad Test Example**:

```go
// ❌ Avoid this
func TestParser(t *testing.T) {
    result := Parse("abc123")
    assert.Equal(t, "abc123", result) // Meaningless assertion
}
```

### Table-Driven Tests

Use table-driven tests for multiple scenarios:

```go
func TestStateTransition(t *testing.T) {
    tests := []struct {
        name         string
        initialState string
        event        string
        wantState    string
        wantErr      bool
    }{
        {
            name:         "order approval workflow",
            initialState: "pending",
            event:        "approve",
            wantState:    "approved",
            wantErr:      false,
        },
        {
            name:         "invalid transition rejected",
            initialState: "shipped",
            event:        "approve",
            wantState:    "shipped",
            wantErr:      true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            fsm := NewStateMachine(tt.initialState)
            err := fsm.Transition(tt.event)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantState, fsm.CurrentState())
            }
        })
    }
}
```

### Coverage Requirements

- **Minimum coverage**: 90% for new code
- **Focus areas**:
  - All exported functions and methods
  - Error handling paths
  - Edge cases and boundary conditions
  - Guard and action execution

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./pkg/parser -run TestYAMLParser_ParseStates

# Run with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. -benchmem ./benchmarks/
```

## Code Style

### Go Conventions

- Follow standard Go formatting: `gofmt`
- Use `golint` for linting
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use meaningful variable names
- Keep functions small and focused

### Formatting

Format code before committing:

```bash
# Format all Go files
go fmt ./...

# Or use gofmt directly
gofmt -w .
```

### Naming Conventions

**Packages**:
- Use lowercase, single-word names: `parser`, `generator`
- Avoid underscore or mixed caps

**Types**:
- Use PascalCase: `StateMachine`, `OrderEvent`
- Use descriptive names: `YAMLParser` not `YP`

**Functions**:
- Use PascalCase for exported: `ParseYAML`
- Use camelCase for unexported: `parseTransition`

**Variables**:
- Use camelCase: `currentState`, `eventName`
- Use short names for short scopes: `i`, `err`
- Use descriptive names for larger scopes

### Comments

**Package comments**:
```go
// Package parser provides YAML and DSL parsing for FSM definitions.
package parser
```

**Function comments**:
```go
// ParseYAML parses a YAML state machine definition from the given reader.
// It returns an FSMModel or an error if parsing fails.
func ParseYAML(r io.Reader) (*model.FSMModel, error) {
    // Implementation
}
```

**Inline comments**:
```go
// Check if initial state exists in states list
if !stateExists(model.Initial, model.States) {
    return nil, ErrInvalidInitialState
}
```

### Error Handling

**Define errors as sentinels**:
```go
var (
    ErrInvalidState = errors.New("invalid state")
    ErrInvalidEvent = errors.New("invalid event")
)
```

**Wrap errors with context**:
```go
if err := parseStates(data); err != nil {
    return nil, fmt.Errorf("parsing states: %w", err)
}
```

**Check errors explicitly**:
```go
// ✅ Good
if err != nil {
    return err
}

// ❌ Bad - don't ignore errors
parseStates(data) // Missing error check
```

### Static Analysis

Run before committing:

```bash
# Vet for common mistakes
go vet ./...

# Staticcheck for advanced analysis
staticcheck ./...

# Exhaustive switch checking
exhaustive ./...
```

## Submitting Changes

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

**Format**:
```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Build process or auxiliary tool changes

**Examples**:
```bash
feat(parser): add support for guard conditions

- Implement guard parsing in YAML parser
- Add validation for guard function names
- Update tests with guard examples

Closes #123
```

```bash
fix(generator): handle empty states list

Previously would panic on empty states.
Now returns proper validation error.

Fixes #456
```

### Pull Request Process

1. **Ensure tests pass**:
```bash
go test ./...
go vet ./...
staticcheck ./...
```

2. **Update documentation** if needed:
   - Update README.md for user-facing changes
   - Update docs/ for API changes
   - Add/update examples if applicable

3. **Create pull request**:
   - Use a descriptive title
   - Reference related issues
   - Describe changes clearly
   - Include test coverage info

4. **PR template**:
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Added new tests
- [ ] All tests pass
- [ ] Coverage >= 90%

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-reviewed the code
- [ ] Commented complex sections
- [ ] Updated documentation
- [ ] No breaking changes (or documented)

## Related Issues
Closes #123
```

5. **Code review**:
   - Address reviewer feedback
   - Keep commits clean and logical
   - Update PR description if needed

6. **Merge**:
   - Wait for approval from maintainers
   - Ensure CI checks pass
   - Squash commits if requested

## Project Structure

Understanding the project layout:

```
gofsm-gen/
├── cmd/
│   └── gofsm-gen/         # CLI entry point
│       └── main.go
├── pkg/
│   ├── parser/            # YAML/DSL parsers
│   │   ├── yaml.go
│   │   ├── yaml_test.go
│   │   └── ...
│   ├── model/             # FSM data model
│   │   ├── fsm.go
│   │   ├── fsm_test.go
│   │   └── ...
│   ├── generator/         # Code generators
│   │   ├── code_generator.go
│   │   ├── code_generator_test.go
│   │   └── ...
│   ├── analyzer/          # Static analysis
│   │   └── ...
│   ├── visualizer/        # Diagram generation
│   │   └── ...
│   └── runtime/           # Runtime support
│       └── ...
├── templates/             # Code generation templates
├── examples/              # Example FSM definitions
├── benchmarks/            # Performance benchmarks
├── docs/                  # Documentation
├── CLAUDE.md              # Project instructions
├── README.md              # Project overview
├── CONTRIBUTING.md        # This file
├── LICENSE                # License
└── go.mod                 # Go module definition
```

### Adding New Packages

1. Create directory under `pkg/`
2. Add `doc.go` with package comment
3. Implement functionality
4. Add comprehensive tests
5. Update documentation

## Communication

### Reporting Issues

Use GitHub Issues for:
- Bug reports
- Feature requests
- Documentation improvements

**Bug report template**:
```markdown
**Describe the bug**
Clear description of what the bug is.

**To Reproduce**
Steps to reproduce:
1. Define FSM with '...'
2. Run command '...'
3. See error

**Expected behavior**
What you expected to happen.

**Environment**
- OS: [e.g. macOS 13.0]
- Go version: [e.g. 1.21.0]
- gofsm-gen version: [e.g. v0.1.0]

**Additional context**
YAML definition, error output, etc.
```

**Feature request template**:
```markdown
**Is your feature request related to a problem?**
Clear description of the problem.

**Describe the solution you'd like**
What you want to happen.

**Describe alternatives you've considered**
Other approaches you've thought about.

**Additional context**
Examples, use cases, etc.
```

### Discussions

Use GitHub Discussions for:
- Questions about usage
- Design discussions
- Ideas and brainstorming
- General help

### Getting Help

- Check [documentation](docs/)
- Search [existing issues](https://github.com/yourusername/gofsm-gen/issues)
- Ask in [discussions](https://github.com/yourusername/gofsm-gen/discussions)
- Join our [Discord](#) (if applicable)

## Development Tips

### IDE Setup

**VSCode**:
```json
{
  "go.testFlags": ["-v"],
  "go.coverOnSave": true,
  "go.lintTool": "staticcheck",
  "editor.formatOnSave": true
}
```

**GoLand**:
- Enable Go Modules
- Set code style to gofmt
- Enable File Watchers for gofmt

### Debugging

**Print debugging**:
```go
import "log"

log.Printf("Debug: state=%v, event=%v", state, event)
```

**Delve debugger**:
```bash
# Install
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug test
dlv test ./pkg/parser -- -test.run TestYAMLParser_ParseStates
```

### Performance Profiling

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## Release Process

For maintainers:

1. Update version in code
2. Update CHANGELOG.md
3. Tag release: `git tag -a v0.1.0 -m "Release v0.1.0"`
4. Push tag: `git push origin v0.1.0`
5. GitHub Actions will build and publish release

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

Don't hesitate to ask questions! We're here to help:
- Open an issue with the "question" label
- Start a discussion on GitHub Discussions
- Reach out to maintainers

Thank you for contributing to gofsm-gen!
