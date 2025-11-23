# Basic Usage Guide

This guide walks you through creating, generating, and using state machines with gofsm-gen.

## Table of Contents

- [Creating Your First State Machine](#creating-your-first-state-machine)
- [Generating Code](#generating-code)
- [Using Generated Code](#using-generated-code)
- [Adding Guards and Actions](#adding-guards-and-actions)
- [Handling Errors](#handling-errors)
- [Testing State Machines](#testing-state-machines)
- [Common Patterns](#common-patterns)

## Creating Your First State Machine

### Step 1: Define the State Machine

Create a YAML file `door.yaml`:

```yaml
machine:
  name: DoorStateMachine
  initial: closed

states:
  - name: closed
  - name: open
  - name: locked

events:
  - open_door
  - close_door
  - lock
  - unlock

transitions:
  - from: closed
    to: open
    on: open_door

  - from: open
    to: closed
    on: close_door

  - from: closed
    to: locked
    on: lock

  - from: locked
    to: closed
    on: unlock
```

### Step 2: Generate the Code

```bash
gofsm-gen -spec=door.yaml -out=door_fsm.gen.go
```

This creates `door_fsm.gen.go` with:
- Type-safe state and event enums
- State machine struct
- Transition methods
- Validation logic

### Step 3: Use in Your Application

Create `main.go`:

```go
package main

import (
    "context"
    "fmt"
    "log"
)

func main() {
    // Create state machine with empty guards and actions
    sm := NewDoorStateMachine(
        DoorGuards{},
        DoorActions{},
    )

    ctx := context.Background()

    // Trigger transitions
    fmt.Printf("Current state: %s\n", sm.State())

    err := sm.Transition(ctx, DoorEventOpenDoor)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("After open: %s\n", sm.State())

    err = sm.Transition(ctx, DoorEventCloseDoor)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("After close: %s\n", sm.State())
}
```

Run it:

```bash
go run main.go door_fsm.gen.go
```

Output:
```
Current state: closed
After open: open
After close: closed
```

## Generating Code

### Basic Generation

```bash
gofsm-gen -spec=fsm.yaml -out=fsm.gen.go
```

### Generation Options

```bash
# Specify package name
gofsm-gen -spec=fsm.yaml -out=fsm.gen.go -package=myfsm

# Generate with tests
gofsm-gen -spec=fsm.yaml -out=fsm.gen.go -generate-tests

# Generate with mocks
gofsm-gen -spec=fsm.yaml -out=fsm.gen.go -generate-mocks

# Generate visualization
gofsm-gen -spec=fsm.yaml -visualize=mermaid -out=diagram.md

# Combine multiple options
gofsm-gen -spec=fsm.yaml -out=fsm.gen.go \
  -package=myfsm \
  -generate-tests \
  -generate-mocks \
  -visualize=mermaid
```

### Output Files

When using all generation options, you get:

```
fsm.gen.go           # Main state machine code
fsm_test.gen.go      # Generated unit tests
fsm_mock.gen.go      # Generated mocks for testing
diagram.md           # Mermaid visualization
```

## Using Generated Code

### Creating State Machines

```go
// Basic creation
sm := NewDoorStateMachine(DoorGuards{}, DoorActions{})

// With options
sm := NewDoorStateMachine(
    DoorGuards{},
    DoorActions{},
    WithLogger(logger),
    WithValidationMode(true),
)
```

### Checking State

```go
// Get current state
state := sm.State()
fmt.Printf("Current state: %s\n", state)

// Check specific state
if sm.State() == DoorStateClosed {
    fmt.Println("Door is closed")
}
```

### Triggering Transitions

```go
ctx := context.Background()

// Trigger a transition
err := sm.Transition(ctx, DoorEventOpenDoor)
if err != nil {
    // Handle invalid transition
    log.Printf("Cannot open door: %v", err)
}
```

### Checking Permitted Events

```go
// Get all events valid for current state
events := sm.PermittedEvents()

for _, event := range events {
    fmt.Printf("Can trigger: %s\n", event)
}

// Check if specific event is permitted
if sm.IsPermitted(DoorEventLock) {
    sm.Transition(ctx, DoorEventLock)
}
```

## Adding Guards and Actions

### Guards: Conditional Transitions

Guards control whether a transition can occur. Update `door.yaml`:

```yaml
transitions:
  - from: closed
    to: locked
    on: lock
    guard: hasKey
```

Implement the guard:

```go
type DoorContext struct {
    HasKey bool
    IsAuthorized bool
}

guards := DoorGuards{
    HasKey: func(ctx context.Context, c *DoorContext) bool {
        return c.HasKey
    },
}

sm := NewDoorStateMachine(guards, DoorActions{})

// Transition will only succeed if HasKey is true
context := &DoorContext{HasKey: true}
err := sm.TransitionWithContext(ctx, DoorEventLock, context)
```

### Actions: Side Effects

Actions execute code during transitions. Update `door.yaml`:

```yaml
transitions:
  - from: closed
    to: open
    on: open_door
    action: logDoorOpen
```

Implement the action:

```go
actions := DoorActions{
    LogDoorOpen: func(ctx context.Context, from, to DoorState, c *DoorContext) error {
        log.Printf("Door opened: %s -> %s", from, to)
        return sendNotification("Door opened")
    },
}

sm := NewDoorStateMachine(DoorGuards{}, actions)
```

### Entry and Exit Actions

Execute code when entering/leaving states:

```yaml
states:
  - name: open
    entry: startTimer
    exit: stopTimer
```

```go
actions := DoorActions{
    StartTimer: func(ctx context.Context, state DoorState, c *DoorContext) error {
        log.Println("Starting door open timer")
        return nil
    },
    StopTimer: func(ctx context.Context, state DoorState, c *DoorContext) error {
        log.Println("Stopping door open timer")
        return nil
    },
}
```

## Handling Errors

### Transition Errors

```go
err := sm.Transition(ctx, DoorEventLock)
if err != nil {
    switch {
    case errors.Is(err, ErrInvalidTransition):
        // Transition not allowed from current state
        log.Printf("Cannot lock door from state: %s", sm.State())

    case errors.Is(err, ErrGuardRejected):
        // Guard function returned false
        log.Println("Guard condition not met")

    default:
        // Action returned an error
        log.Printf("Action failed: %v", err)
    }
}
```

### Guard Errors

```go
guards := DoorGuards{
    HasKey: func(ctx context.Context, c *DoorContext) bool {
        // Guards return bool, not error
        // Check condition and return true/false
        if c == nil {
            return false
        }
        return c.HasKey
    },
}
```

### Action Errors

```go
actions := DoorActions{
    LogDoorOpen: func(ctx context.Context, from, to DoorState, c *DoorContext) error {
        // Actions can return errors
        if err := database.LogEvent("door_open"); err != nil {
            return fmt.Errorf("failed to log event: %w", err)
        }
        return nil
    },
}
```

## Testing State Machines

### Unit Testing

```go
func TestDoorStateMachine_OpenClose(t *testing.T) {
    sm := NewDoorStateMachine(DoorGuards{}, DoorActions{})
    ctx := context.Background()

    // Test initial state
    assert.Equal(t, DoorStateClosed, sm.State())

    // Test valid transition
    err := sm.Transition(ctx, DoorEventOpenDoor)
    assert.NoError(t, err)
    assert.Equal(t, DoorStateOpen, sm.State())

    // Test closing
    err = sm.Transition(ctx, DoorEventCloseDoor)
    assert.NoError(t, err)
    assert.Equal(t, DoorStateClosed, sm.State())
}

func TestDoorStateMachine_InvalidTransition(t *testing.T) {
    sm := NewDoorStateMachine(DoorGuards{}, DoorActions{})
    ctx := context.Background()

    // Cannot lock an open door
    sm.Transition(ctx, DoorEventOpenDoor)
    err := sm.Transition(ctx, DoorEventLock)

    assert.Error(t, err)
    assert.ErrorIs(t, err, ErrInvalidTransition)
    assert.Equal(t, DoorStateOpen, sm.State()) // State unchanged
}
```

### Testing with Guards

```go
func TestDoorStateMachine_GuardCondition(t *testing.T) {
    guards := DoorGuards{
        HasKey: func(ctx context.Context, c *DoorContext) bool {
            return c.HasKey
        },
    }

    sm := NewDoorStateMachine(guards, DoorActions{})
    ctx := context.Background()

    // Should succeed with key
    contextWithKey := &DoorContext{HasKey: true}
    err := sm.TransitionWithContext(ctx, DoorEventLock, contextWithKey)
    assert.NoError(t, err)

    // Should fail without key
    sm.Transition(ctx, DoorEventUnlock) // Back to closed
    contextNoKey := &DoorContext{HasKey: false}
    err = sm.TransitionWithContext(ctx, DoorEventLock, contextNoKey)
    assert.ErrorIs(t, err, ErrGuardRejected)
}
```

### Using Generated Tests

If you used `-generate-tests`, you'll get basic tests automatically:

```bash
gofsm-gen -spec=door.yaml -out=door_fsm.gen.go -generate-tests

# Run generated tests
go test -v
```

## Common Patterns

### Pattern 1: Request Workflow

```yaml
machine:
  name: RequestStateMachine
  initial: pending

states:
  - name: pending
  - name: approved
  - name: rejected
  - name: completed

events:
  - approve
  - reject
  - complete
  - retry

transitions:
  - from: pending
    to: approved
    on: approve
    guard: hasPermission
    action: notifyApprover

  - from: pending
    to: rejected
    on: reject
    action: notifyRequester

  - from: approved
    to: completed
    on: complete
    action: finalizeRequest

  - from: rejected
    to: pending
    on: retry
```

### Pattern 2: Connection State

```yaml
machine:
  name: ConnectionStateMachine
  initial: disconnected

states:
  - name: disconnected
    entry: cleanupResources
  - name: connecting
    entry: startConnection
  - name: connected
    entry: setupHeartbeat
    exit: stopHeartbeat
  - name: error
    entry: logError

events:
  - connect
  - disconnect
  - connection_established
  - connection_failed
  - retry

transitions:
  - from: disconnected
    to: connecting
    on: connect

  - from: connecting
    to: connected
    on: connection_established

  - from: connecting
    to: error
    on: connection_failed

  - from: error
    to: connecting
    on: retry

  - from: connected
    to: disconnected
    on: disconnect
```

### Pattern 3: Multi-step Process

```yaml
machine:
  name: DeploymentStateMachine
  initial: ready

states:
  - name: ready
  - name: building
  - name: testing
  - name: deploying
  - name: deployed
  - name: failed

events:
  - start_build
  - build_complete
  - test_complete
  - deploy_complete
  - fail
  - rollback

transitions:
  - from: ready
    to: building
    on: start_build

  - from: building
    to: testing
    on: build_complete

  - from: building
    to: failed
    on: fail

  - from: testing
    to: deploying
    on: test_complete

  - from: testing
    to: failed
    on: fail

  - from: deploying
    to: deployed
    on: deploy_complete

  - from: deploying
    to: failed
    on: fail

  - from: failed
    to: ready
    on: rollback
```

## Next Steps

- Learn about advanced YAML features in [YAML Definition Reference](yaml-reference.md)
- Explore the generated API in [API Documentation](api.md)
- See complete examples in [examples/](../examples/)
- Read about [Testing Strategies](testing.md)

## Tips and Best Practices

1. **Start Simple**: Begin with basic states and transitions, add complexity later
2. **Use Guards Wisely**: Guards should be pure functions without side effects
3. **Keep Actions Focused**: Each action should do one thing well
4. **Test Exhaustively**: Test all valid transitions and invalid ones
5. **Document Your FSM**: Add comments to your YAML explaining business logic
6. **Use Visualization**: Generate diagrams to communicate with non-technical stakeholders
7. **Version Your Specs**: Keep YAML files in version control
8. **Regenerate After Changes**: Always regenerate code after modifying YAML
