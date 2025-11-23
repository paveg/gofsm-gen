# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

gofsm-gen is a code generation-based state machine library for Go that provides Rust-level exhaustiveness checking for state transitions. It generates type-safe state machine code from YAML or Go DSL definitions.

**Core Concept**: Provide compile-time safety through code generation + static analysis (using the `exhaustive` tool) to ensure all state transitions are handled, mimicking Rust's enum-based pattern matching.

## Architecture

### Layer Structure

```
Application Layer       <- User application code
Generated Code Layer    <- Type-safe generated code (*.gen.go)
Code Generation Layer   <- gofsm-gen CLI tool
Definition Layer        <- FSM definitions (YAML/DSL)
```

### Key Components

- **Definition Parser** (`pkg/parser/`): Parses YAML/HCL/Go DSL definitions
- **Model Builder** (`pkg/model/`): Builds internal FSM representation
- **Code Generator** (`pkg/generator/`): Generates Go code from FSM model
- **Static Analyzer** (`pkg/analyzer/`): Validates exhaustiveness and reachability
- **Runtime Validator** (`pkg/runtime/`): Provides runtime validation and logging
- **Visualizer** (`pkg/visualizer/`): Generates Mermaid/Graphviz diagrams

## Commands

### Build and Development

```bash
# Build the CLI tool
go build -o bin/gofsm-gen ./cmd/gofsm-gen

# Run tests
go test ./...

# Run tests with coverage
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem ./benchmarks/

# Run static analysis
go vet ./...
staticcheck ./...
```

### Code Generation

```bash
# Generate from YAML specification
gofsm-gen -spec=statemachine.yaml -out=fsm.gen.go

# Generate with tests and mocks
gofsm-gen -spec=fsm.yaml -out=fsm.gen.go -generate-tests -generate-mocks

# Generate from Go DSL
gofsm-gen -type=OrderStateMachine

# Generate with visualization
gofsm-gen -spec=fsm.yaml -out=fsm.gen.go -visualize=mermaid

# Type inference mode
gofsm-gen -infer -type=DoorLock

# Full options example
gofsm-gen \
  -spec=fsm.yaml \
  -out=fsm.gen.go \
  -package=myfsm \
  -generate-tests \
  -generate-mocks \
  -visualize=mermaid
```

## FSM Definition Format

### YAML Definition

```yaml
machine:
  name: OrderStateMachine
  initial: pending

states:
  - name: pending
    entry: logEntry
    exit: logExit
  - name: approved
  - name: rejected
  - name: shipped

events:
  - approve
  - reject
  - ship

transitions:
  - from: pending
    to: approved
    on: approve
    guard: hasPayment
    action: chargeCard
  - from: approved
    to: shipped
    on: ship
    action: notifyShipping
```

### Generated Code Usage

```go
// Define guards and actions
guards := OrderGuards{
    HasPayment: func(ctx context.Context, c *OrderContext) bool {
        return c.PaymentMethod != ""
    },
}

actions := OrderActions{
    ChargeCard: func(ctx context.Context, from, to OrderState, c *OrderContext) error {
        return processPayment(c)
    },
}

// Create state machine
sm := NewOrderStateMachine(guards, actions,
    WithLogger(logger),
    WithValidationMode(true),
)

// Trigger transitions
err := sm.Transition(ctx, OrderEventApprove)

// Check current state
state := sm.State()

// Get permitted events
events := sm.PermittedEvents()
```

## Package Structure

```
github.com/yourusername/gofsm-gen/
├── cmd/gofsm-gen/          # CLI entry point
├── pkg/
│   ├── parser/             # YAML/DSL/AST parsers
│   ├── model/              # Internal FSM data model
│   │   ├── fsm.go          # FSMModel, State, Event, Transition
│   │   └── graph.go        # StateGraph for analysis
│   ├── generator/          # Code generators
│   │   ├── code_generator.go
│   │   ├── test_generator.go
│   │   └── mock_generator.go
│   ├── analyzer/           # Static analysis
│   │   ├── exhaustive.go   # Exhaustiveness checking
│   │   ├── validator.go    # Model validation
│   │   └── graph.go        # Graph analysis
│   ├── visualizer/         # Diagram generation
│   │   ├── mermaid.go
│   │   └── graphviz.go
│   └── runtime/            # Runtime support
│       ├── logger.go
│       ├── validator.go
│       └── context.go
├── templates/              # Code generation templates
│   ├── state_machine.tmpl
│   ├── test.tmpl
│   └── mock.tmpl
├── examples/               # Example FSM definitions
└── benchmarks/             # Performance benchmarks
```

## Key Implementation Details

### Code Generation

The generator produces:

1. **Type-safe enums**: States and events as Go `int` constants with exhaustive switch enforcement
2. **Exhaustive annotations**: `//exhaustive:enforce` comments for static analysis
3. **Guard functions**: Optional predicates that control transition execution
4. **Action functions**: Code executed during transitions
5. **Entry/Exit actions**: Code executed when entering/leaving states

### Static Analysis Integration

The generated code includes `//exhaustive:enforce` annotations that the `exhaustive` static analyzer validates. This ensures:

- All states are handled in switch statements
- All events are handled for each state
- No unreachable code paths
- Compile-time safety similar to Rust's enum exhaustiveness

### Validation Checks

The analyzer performs:

- **Reachability**: All states reachable from initial state
- **Determinism**: No conflicting unguarded transitions
- **Completeness**: All referenced states/events are defined
- **Guard conflicts**: Warns if multiple guards could be true simultaneously

### Performance Goals

- State transition: < 50ns/transition
- Memory usage: < 1KB/instance
- Zero allocations in hot path (when ZeroAllocation option enabled)
- Code generation: < 1 second for 1000 states

## Development Phases

The project is planned in 4 phases:

**Phase 1**: YAML definitions + basic code generation + exhaustive integration
**Phase 2**: Guards/actions + Go DSL support
**Phase 3**: VSCode extension + enhanced tooling
**Phase 4**: Hierarchical state machines + history states

## Testing Strategy

### Test-Driven Development (TDD)

This project follows **Test-Driven Development** as the primary development methodology:

1. **Write tests first**: Before implementing any feature, write the test that defines the expected behavior
2. **Red-Green-Refactor cycle**:
   - Red: Write a failing test
   - Green: Write minimal code to make the test pass
   - Refactor: Improve the code while keeping tests green
3. **Tests as specification**: Tests serve as living documentation of expected behavior

### Test Quality Guidelines

**IMPORTANT**: Write only meaningful tests that verify actual behavior. Avoid the following anti-patterns:

❌ **DO NOT**:

- Write tests with hardcoded magic values that don't represent real use cases
- Create tests that simply assert `result == result` or other tautologies
- Test implementation details instead of behavior
- Write brittle tests that break with any refactoring
- Use arbitrary test data that doesn't reflect domain knowledge

✅ **DO**:

- Test real-world scenarios and edge cases
- Use domain-meaningful test data (e.g., realistic state machine definitions)
- Verify behavior, not implementation
- Write tests that document expected behavior clearly
- Use table-driven tests with descriptive test cases
- Test error conditions and boundary cases explicitly

### Test Types

- **Unit tests**: All packages require >90% coverage
  - Focus on single functions/methods with clear inputs and outputs
  - Use realistic test data representing actual FSM definitions
  - Test both happy paths and error conditions
- **Integration tests**: End-to-end generation and validation
  - Test complete workflows (parse → validate → generate)
  - Use real YAML/DSL examples from the examples/ directory
- **Benchmarks**: Performance regression tracking
  - Measure actual performance metrics against goals
  - Use realistic state machine sizes (10, 100, 1000 states)
- **Golden file tests**: Generated code comparison
  - Compare generated code against known-good examples
  - Update golden files only when intentionally changing output
- **Property-based tests**: FSM model validation using fuzzing
  - Use Go's native fuzzing support (go test -fuzz)
  - Verify invariants hold for arbitrary valid inputs

### Test Organization

```go
// Good: Descriptive, domain-meaningful test
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

// Bad: Hardcoded magic values without meaning
func TestParser(t *testing.T) {
    result := Parse("abc123")
    assert.Equal(t, "abc123", result) // Meaningless assertion
}
```

## Code Style

### General Guidelines

- Follow standard Go conventions (gofmt, golint)
- Use exhaustive switch statements (with `//exhaustive:enforce`)
- Document all exported APIs with godoc comments
- Keep functions small and focused
- Prefer composition over inheritance
- Use interfaces for extensibility

### Test-Driven Development Workflow

When implementing new features or fixing bugs:

1. **Start with a failing test**
   ```bash
   # Write the test first
   vim pkg/parser/yaml_test.go

   # Verify it fails for the right reason
   go test ./pkg/parser -run TestYAMLParser_ParseGuardConditions
   ```

2. **Implement minimal code to pass**
   ```bash
   # Write just enough code to make the test pass
   vim pkg/parser/yaml.go

   # Run the test again
   go test ./pkg/parser -run TestYAMLParser_ParseGuardConditions
   ```

3. **Refactor while keeping tests green**
   ```bash
   # Improve the implementation
   # Run all tests to ensure nothing breaks
   go test ./...
   ```

4. **Commit with test and implementation together**
   ```bash
   git add pkg/parser/yaml.go pkg/parser/yaml_test.go
   git commit -m "feat(parser): add guard condition parsing"
   ```

### Avoiding Hardcoded Values

**Bad**: Magic numbers and strings without context
```go
// ❌ Avoid this
func TestSomething(t *testing.T) {
    result := Process(123, "abc")
    assert.Equal(t, 456, result)  // What do these numbers mean?
}
```

**Good**: Use constants or variables with meaningful names
```go
// ✅ Do this
func TestStateTransition_ValidTransition(t *testing.T) {
    const (
        initialState = "pending"
        targetState  = "approved"
        triggerEvent = "approve"
    )

    fsm := NewStateMachine(initialState)
    err := fsm.Transition(triggerEvent)

    require.NoError(t, err)
    assert.Equal(t, targetState, fsm.CurrentState())
}
```

**Better**: Use table-driven tests with descriptive cases
```go
// ✅ Even better: table-driven with meaningful scenarios
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

### Code Review Checklist

Before submitting code, ensure:

- [ ] Tests were written **before** implementation (TDD)
- [ ] All tests use meaningful, domain-relevant test data
- [ ] No hardcoded magic values without explanation
- [ ] Test names clearly describe what is being tested
- [ ] Both success and failure cases are tested
- [ ] Code coverage is >90% for new code
- [ ] All tests pass: `go test ./...`
- [ ] Static analysis passes: `go vet ./... && staticcheck ./...`
