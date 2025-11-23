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

- **Unit tests**: All packages require >90% coverage
- **Integration tests**: End-to-end generation and validation
- **Benchmarks**: Performance regression tracking
- **Golden file tests**: Generated code comparison
- **Property-based tests**: FSM model validation using fuzzing

## Code Style

- Follow standard Go conventions (gofmt, golint)
- Use exhaustive switch statements (with `//exhaustive:enforce`)
- Document all exported APIs with godoc comments
- Keep functions small and focused
- Prefer composition over inheritance
- Use interfaces for extensibility
