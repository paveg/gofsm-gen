# YAML Definition Reference

Complete reference for defining state machines using YAML syntax.

## Table of Contents

- [File Structure](#file-structure)
- [Machine Configuration](#machine-configuration)
- [States](#states)
- [Events](#events)
- [Transitions](#transitions)
- [Guards](#guards)
- [Actions](#actions)
- [Context Types](#context-types)
- [Options](#options)
- [Complete Examples](#complete-examples)

## File Structure

A YAML state machine definition consists of four main sections:

```yaml
machine:
  # Machine configuration

states:
  # State definitions

events:
  # Event definitions

transitions:
  # Transition definitions
```

## Machine Configuration

The `machine` section defines basic properties of the state machine.

### Syntax

```yaml
machine:
  name: <string>          # Required: Name of the state machine
  initial: <string>       # Required: Initial state
  description: <string>   # Optional: Documentation
  context: <string>       # Optional: Context type name
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Name of the generated state machine struct. Must be PascalCase. |
| `initial` | string | Yes | Name of the initial state. Must exist in states list. |
| `description` | string | No | Human-readable description for documentation. |
| `context` | string | No | Custom context type name. Defaults to `{Name}Context`. |

### Example

```yaml
machine:
  name: OrderStateMachine
  initial: pending
  description: "Manages the lifecycle of customer orders"
  context: OrderContext
```

## States

The `states` section defines all possible states in the machine.

### Simple Syntax

```yaml
states:
  - name: pending
  - name: approved
  - name: rejected
```

### Extended Syntax

```yaml
states:
  - name: <string>          # Required: State name
    description: <string>   # Optional: Documentation
    entry: <string>         # Optional: Entry action name
    exit: <string>          # Optional: Exit action name
    metadata: <map>         # Optional: Custom metadata
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | State identifier. Must be lowercase with underscores. |
| `description` | string | No | Human-readable description. |
| `entry` | string | No | Action to execute when entering this state. |
| `exit` | string | No | Action to execute when leaving this state. |
| `metadata` | map | No | Custom key-value data for code generation. |

### Example

```yaml
states:
  - name: pending
    description: "Order is awaiting approval"
    entry: logPending

  - name: processing
    description: "Order is being processed"
    entry: startProcessing
    exit: stopProcessing
    metadata:
      color: yellow
      timeout: 300

  - name: completed
    description: "Order has been fulfilled"
    entry: logCompletion
```

### State Naming Rules

- Use lowercase with underscores: `pending`, `in_progress`, `completed`
- Be descriptive but concise: `processing_payment` not `processing`
- Avoid reserved Go keywords: `type`, `func`, `interface`
- No spaces or special characters except underscore

## Events

The `events` section defines all possible triggers for state transitions.

### Simple Syntax

```yaml
events:
  - approve
  - reject
  - ship
```

### Extended Syntax

```yaml
events:
  - name: <string>          # Required: Event name
    description: <string>   # Optional: Documentation
    metadata: <map>         # Optional: Custom metadata
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Event identifier. Must be lowercase with underscores. |
| `description` | string | No | Human-readable description. |
| `metadata` | map | No | Custom key-value data for code generation. |

### Example

```yaml
events:
  - name: approve
    description: "Approve the order for processing"

  - name: reject
    description: "Reject the order"
    metadata:
      requires_reason: true

  - name: ship
    description: "Mark order as shipped"
    metadata:
      requires_tracking: true
```

### Event Naming Rules

- Use lowercase with underscores: `approve`, `send_email`, `timeout_occurred`
- Use imperative verbs: `submit`, `cancel`, `retry`
- Be specific: `payment_succeeded` not just `success`
- Avoid reserved Go keywords

## Transitions

The `transitions` section defines how states change in response to events.

### Syntax

```yaml
transitions:
  - from: <string>          # Required: Source state
    to: <string>            # Required: Target state
    on: <string>            # Required: Triggering event
    guard: <string>         # Optional: Guard function name
    action: <string>        # Optional: Action function name
    description: <string>   # Optional: Documentation
    metadata: <map>         # Optional: Custom metadata
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `from` | string | Yes | Source state name. Must exist in states list. |
| `to` | string | Yes | Target state name. Must exist in states list. |
| `on` | string | Yes | Event that triggers this transition. Must exist in events list. |
| `guard` | string | No | Name of guard function to check before transitioning. |
| `action` | string | No | Name of action function to execute during transition. |
| `description` | string | No | Human-readable description. |
| `metadata` | map | No | Custom key-value data for code generation. |

### Example

```yaml
transitions:
  - from: pending
    to: approved
    on: approve
    guard: hasPaymentMethod
    action: chargeCard
    description: "Approve order and charge payment"

  - from: approved
    to: shipped
    on: ship
    action: notifyShipping
    description: "Send order to shipping"

  - from: pending
    to: rejected
    on: reject
    action: notifyRejection
```

### Self-Transitions

A state can transition to itself:

```yaml
transitions:
  - from: processing
    to: processing
    on: retry
    guard: canRetry
    action: incrementRetryCount
```

### Multiple Transitions on Same Event

Multiple transitions can use the same event from the same state if they have different guards:

```yaml
transitions:
  # High priority customer path
  - from: pending
    to: express_processing
    on: submit
    guard: isHighPriority

  # Regular customer path
  - from: pending
    to: regular_processing
    on: submit
    guard: isRegularCustomer
```

**Warning**: Ensure guards are mutually exclusive to avoid non-determinism.

## Guards

Guards are predicate functions that control whether a transition can occur.

### Definition in YAML

```yaml
transitions:
  - from: pending
    to: approved
    on: approve
    guard: hasPaymentMethod
```

### Implementation in Go

Guards are defined as functions in the `Guards` struct:

```go
type OrderGuards struct {
    HasPaymentMethod func(ctx context.Context, c *OrderContext) bool
}
```

### Function Signature

```go
func(ctx context.Context, c *ContextType) bool
```

- **Parameters**:
  - `ctx`: Go context for cancellation and values
  - `c`: Pointer to context struct with state machine data
- **Returns**: `true` if transition is allowed, `false` otherwise

### Example Implementation

```go
guards := OrderGuards{
    HasPaymentMethod: func(ctx context.Context, c *OrderContext) bool {
        return c.PaymentMethod != "" && c.PaymentMethod != "none"
    },

    HasInventory: func(ctx context.Context, c *OrderContext) bool {
        return c.InventoryCount > 0
    },

    IsAuthorized: func(ctx context.Context, c *OrderContext) bool {
        return c.UserRole == "admin" || c.UserRole == "manager"
    },
}
```

### Best Practices

- **Pure Functions**: Guards should not modify state or have side effects
- **Fast Execution**: Guards should execute quickly (< 1ms)
- **No I/O**: Avoid database queries or API calls in guards
- **Clear Logic**: Each guard should check one condition
- **Testable**: Guards should be easy to unit test

## Actions

Actions are functions executed during state transitions or when entering/exiting states.

### Transition Actions

Defined in the `transitions` section:

```yaml
transitions:
  - from: pending
    to: approved
    on: approve
    action: chargeCard
```

### Entry/Exit Actions

Defined in the `states` section:

```yaml
states:
  - name: processing
    entry: startTimer
    exit: stopTimer
```

### Implementation in Go

Actions are defined in the `Actions` struct:

```go
type OrderActions struct {
    // Transition actions
    ChargeCard func(ctx context.Context, from, to OrderState, c *OrderContext) error

    // Entry/exit actions
    StartTimer func(ctx context.Context, state OrderState, c *OrderContext) error
    StopTimer  func(ctx context.Context, state OrderState, c *OrderContext) error
}
```

### Function Signatures

**Transition Action**:
```go
func(ctx context.Context, from, to StateType, c *ContextType) error
```

**Entry/Exit Action**:
```go
func(ctx context.Context, state StateType, c *ContextType) error
```

### Example Implementation

```go
actions := OrderActions{
    ChargeCard: func(ctx context.Context, from, to OrderState, c *OrderContext) error {
        log.Printf("Charging card for order %s", c.OrderID)
        if err := paymentService.Charge(c.PaymentMethod, c.Amount); err != nil {
            return fmt.Errorf("payment failed: %w", err)
        }
        c.ChargedAt = time.Now()
        return nil
    },

    StartTimer: func(ctx context.Context, state OrderState, c *OrderContext) error {
        c.ProcessingStarted = time.Now()
        return nil
    },

    StopTimer: func(ctx context.Context, state OrderState, c *OrderContext) error {
        duration := time.Since(c.ProcessingStarted)
        log.Printf("Processing took %v", duration)
        return nil
    },
}
```

### Error Handling

- Actions can return errors to abort the transition
- If an action returns an error, the state does not change
- Entry/exit actions can also return errors

```go
action: func(ctx context.Context, from, to State, c *Context) error {
    if err := someOperation(); err != nil {
        // Transition will be aborted
        return fmt.Errorf("operation failed: %w", err)
    }
    return nil
}
```

## Context Types

Context holds data that flows through the state machine.

### Default Context

If no `context` is specified in machine config, a default context is generated:

```go
type OrderStateMachineContext struct {
    // Add your fields here
}
```

### Custom Context

Specify in YAML:

```yaml
machine:
  name: OrderStateMachine
  context: OrderContext
```

Then define in your code:

```go
type OrderContext struct {
    OrderID       string
    CustomerID    string
    Amount        float64
    PaymentMethod string
    InventoryCount int
    ProcessingStarted time.Time
    ChargedAt     time.Time
}
```

### Context Usage

```go
sm := NewOrderStateMachine(guards, actions)

context := &OrderContext{
    OrderID:       "ORD-123",
    CustomerID:    "CUST-456",
    Amount:        99.99,
    PaymentMethod: "credit_card",
}

err := sm.TransitionWithContext(ctx, OrderEventApprove, context)
```

## Options

Configuration options for code generation.

### YAML Options

```yaml
machine:
  name: OrderStateMachine
  initial: pending

options:
  validation: true           # Enable runtime validation
  logging: true              # Generate logging code
  metrics: true              # Generate metrics collection
  zero_allocation: false     # Optimize for zero allocations
  concurrency_safe: true     # Add mutex protection
```

### Option Descriptions

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `validation` | bool | true | Enable runtime validation of transitions |
| `logging` | bool | false | Generate structured logging code |
| `metrics` | bool | false | Generate metrics collection hooks |
| `zero_allocation` | bool | false | Optimize for zero heap allocations |
| `concurrency_safe` | bool | false | Add mutex protection for concurrent access |

## Complete Examples

### Example 1: Simple Door Lock

```yaml
machine:
  name: DoorLock
  initial: locked
  description: "Simple door lock state machine"

states:
  - name: locked
  - name: unlocked

events:
  - unlock
  - lock

transitions:
  - from: locked
    to: unlocked
    on: unlock

  - from: unlocked
    to: locked
    on: lock
```

### Example 2: Order Processing with Guards and Actions

```yaml
machine:
  name: OrderStateMachine
  initial: pending
  context: OrderContext

states:
  - name: pending
    description: "Order is waiting for approval"

  - name: approved
    description: "Order has been approved"
    entry: notifyApproval

  - name: processing
    description: "Order is being processed"
    entry: startProcessing
    exit: stopProcessing

  - name: shipped
    description: "Order has been shipped"
    entry: notifyShipping

  - name: completed
    description: "Order is complete"

  - name: cancelled
    description: "Order was cancelled"
    entry: notifyCancellation

events:
  - approve
  - reject
  - process
  - ship
  - complete
  - cancel

transitions:
  - from: pending
    to: approved
    on: approve
    guard: hasPaymentMethod
    action: chargeCard

  - from: pending
    to: cancelled
    on: reject
    action: refundDeposit

  - from: approved
    to: processing
    on: process
    guard: hasInventory

  - from: processing
    to: shipped
    on: ship
    action: updateTracking

  - from: shipped
    to: completed
    on: complete

  - from: pending
    to: cancelled
    on: cancel

  - from: approved
    to: cancelled
    on: cancel
    action: refundPayment

options:
  validation: true
  logging: true
  concurrency_safe: true
```

### Example 3: Connection Management

```yaml
machine:
  name: Connection
  initial: disconnected

states:
  - name: disconnected
    entry: cleanupResources

  - name: connecting
    entry: startConnection
    exit: cancelConnection

  - name: connected
    entry: setupHeartbeat
    exit: stopHeartbeat

  - name: reconnecting
    entry: startReconnection

  - name: error
    entry: logError

events:
  - connect
  - connected
  - disconnect
  - connection_lost
  - reconnect
  - error

transitions:
  - from: disconnected
    to: connecting
    on: connect

  - from: connecting
    to: connected
    on: connected

  - from: connecting
    to: error
    on: error

  - from: connected
    to: disconnected
    on: disconnect

  - from: connected
    to: reconnecting
    on: connection_lost

  - from: reconnecting
    to: connected
    on: connected

  - from: reconnecting
    to: error
    on: error
    guard: maxRetriesExceeded

  - from: error
    to: connecting
    on: reconnect

options:
  validation: true
  logging: true
  concurrency_safe: true
```

## Validation Rules

The YAML parser validates:

1. **Required Fields**: `machine.name`, `machine.initial`, states, events, transitions
2. **State References**: All states referenced in transitions must be defined
3. **Event References**: All events referenced in transitions must be defined
4. **Initial State**: Must exist in states list
5. **Guard Names**: Must be valid Go identifiers
6. **Action Names**: Must be valid Go identifiers
7. **Uniqueness**: No duplicate state or event names
8. **Reachability**: All states should be reachable from initial state (warning)
9. **Determinism**: No conflicting unguarded transitions (warning)

## Next Steps

- See [Usage Guide](usage.md) for how to use these definitions
- Read [API Documentation](api.md) for generated code details
- Explore [Examples](../examples/) for complete working examples
