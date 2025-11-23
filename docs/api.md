# API Documentation

Complete reference for the generated code API.

## Table of Contents

- [Generated Types](#generated-types)
- [State Machine API](#state-machine-api)
- [Guard Functions](#guard-functions)
- [Action Functions](#action-functions)
- [Context Management](#context-management)
- [Error Types](#error-types)
- [Options and Configuration](#options-and-configuration)

## Generated Types

### State Enum

For a state machine named `Order`, states are generated as:

```go
type OrderState int

const (
    OrderStatePending OrderState = iota
    OrderStateApproved
    OrderStateRejected
    OrderStateShipped
)

// String returns the string representation of the state
func (s OrderState) String() string

// IsValid checks if the state value is valid
func (s OrderState) IsValid() bool
```

**Usage**:
```go
state := OrderStatePending
fmt.Println(state) // Output: "pending"

if state.IsValid() {
    // State is valid
}
```

### Event Enum

For a state machine named `Order`, events are generated as:

```go
type OrderEvent int

const (
    OrderEventApprove OrderEvent = iota
    OrderEventReject
    OrderEventShip
)

// String returns the string representation of the event
func (e OrderEvent) String() string

// IsValid checks if the event value is valid
func (e OrderEvent) IsValid() bool
```

**Usage**:
```go
event := OrderEventApprove
fmt.Println(event) // Output: "approve"

if event.IsValid() {
    // Event is valid
}
```

### Context Type

Default context structure:

```go
type OrderContext struct {
    // Add your custom fields here
}
```

You should extend this with your application-specific data:

```go
type OrderContext struct {
    OrderID       string
    CustomerID    string
    Amount        float64
    PaymentMethod string
    CreatedAt     time.Time
}
```

## State Machine API

### Constructor

```go
func NewOrderStateMachine(
    guards OrderGuards,
    actions OrderActions,
    opts ...OrderStateMachineOption,
) *OrderStateMachine
```

**Parameters**:
- `guards`: Struct containing guard functions
- `actions`: Struct containing action functions
- `opts`: Optional configuration options

**Returns**: Pointer to initialized state machine

**Example**:
```go
sm := NewOrderStateMachine(
    OrderGuards{},
    OrderActions{},
    WithLogger(logger),
    WithValidationMode(true),
)
```

### State Machine Methods

#### State()

```go
func (sm *OrderStateMachine) State() OrderState
```

Returns the current state of the machine.

**Example**:
```go
currentState := sm.State()
fmt.Printf("Current state: %s\n", currentState)
```

#### Transition()

```go
func (sm *OrderStateMachine) Transition(
    ctx context.Context,
    event OrderEvent,
) error
```

Triggers a state transition.

**Parameters**:
- `ctx`: Go context for cancellation and timeout
- `event`: Event to trigger

**Returns**: Error if transition fails

**Example**:
```go
ctx := context.Background()
if err := sm.Transition(ctx, OrderEventApprove); err != nil {
    log.Printf("Transition failed: %v", err)
}
```

#### TransitionWithContext()

```go
func (sm *OrderStateMachine) TransitionWithContext(
    ctx context.Context,
    event OrderEvent,
    context *OrderContext,
) error
```

Triggers a transition with custom context data.

**Parameters**:
- `ctx`: Go context
- `event`: Event to trigger
- `context`: State machine context with application data

**Returns**: Error if transition fails

**Example**:
```go
orderCtx := &OrderContext{
    OrderID: "ORD-123",
    Amount:  99.99,
}

err := sm.TransitionWithContext(ctx, OrderEventApprove, orderCtx)
```

#### PermittedEvents()

```go
func (sm *OrderStateMachine) PermittedEvents() []OrderEvent
```

Returns all events that are valid for the current state.

**Returns**: Slice of permitted events

**Example**:
```go
events := sm.PermittedEvents()
for _, event := range events {
    fmt.Printf("Can trigger: %s\n", event)
}
```

#### IsPermitted()

```go
func (sm *OrderStateMachine) IsPermitted(event OrderEvent) bool
```

Checks if a specific event is allowed in the current state.

**Parameters**:
- `event`: Event to check

**Returns**: `true` if event is permitted

**Example**:
```go
if sm.IsPermitted(OrderEventApprove) {
    // Event is allowed
    sm.Transition(ctx, OrderEventApprove)
}
```

#### CanTransitionTo()

```go
func (sm *OrderStateMachine) CanTransitionTo(
    state OrderState,
) bool
```

Checks if a transition to the given state is possible from current state.

**Parameters**:
- `state`: Target state to check

**Returns**: `true` if transition is possible

**Example**:
```go
if sm.CanTransitionTo(OrderStateApproved) {
    fmt.Println("Can transition to approved")
}
```

#### Reset()

```go
func (sm *OrderStateMachine) Reset()
```

Resets the state machine to its initial state.

**Example**:
```go
sm.Reset()
fmt.Printf("State after reset: %s\n", sm.State())
```

## Guard Functions

### Guard Struct

```go
type OrderGuards struct {
    HasPaymentMethod func(ctx context.Context, c *OrderContext) bool
    HasInventory     func(ctx context.Context, c *OrderContext) bool
    IsAuthorized     func(ctx context.Context, c *OrderContext) bool
}
```

### Guard Function Signature

```go
type GuardFunc func(ctx context.Context, c *ContextType) bool
```

**Parameters**:
- `ctx`: Go context for cancellation
- `c`: Pointer to context struct

**Returns**: `true` to allow transition, `false` to reject

### Implementation Example

```go
guards := OrderGuards{
    HasPaymentMethod: func(ctx context.Context, c *OrderContext) bool {
        // Check if payment method is set
        return c.PaymentMethod != "" && c.PaymentMethod != "none"
    },

    HasInventory: func(ctx context.Context, c *OrderContext) bool {
        // Check inventory availability
        return c.InventoryCount > 0
    },

    IsAuthorized: func(ctx context.Context, c *OrderContext) bool {
        // Check user authorization
        user, ok := ctx.Value("user").(User)
        if !ok {
            return false
        }
        return user.Role == "admin" || user.Role == "manager"
    },
}
```

### Guard Best Practices

1. **Pure Functions**: No side effects or state modifications
2. **Fast Execution**: Should complete in microseconds
3. **No I/O**: Avoid database calls or network requests
4. **Clear Logic**: One guard per condition
5. **Context Usage**: Use `ctx` for timeout and cancellation

## Action Functions

### Action Struct

```go
type OrderActions struct {
    // Transition actions
    ChargeCard      func(ctx context.Context, from, to OrderState, c *OrderContext) error
    NotifyShipping  func(ctx context.Context, from, to OrderState, c *OrderContext) error

    // Entry/exit actions
    LogEntry        func(ctx context.Context, state OrderState, c *OrderContext) error
    LogExit         func(ctx context.Context, state OrderState, c *OrderContext) error
}
```

### Action Function Signatures

**Transition Action**:
```go
type TransitionActionFunc func(
    ctx context.Context,
    from StateType,
    to StateType,
    c *ContextType,
) error
```

**Entry/Exit Action**:
```go
type StateActionFunc func(
    ctx context.Context,
    state StateType,
    c *ContextType,
) error
```

### Implementation Example

```go
actions := OrderActions{
    ChargeCard: func(ctx context.Context, from, to OrderState, c *OrderContext) error {
        log.Printf("Charging card for order %s", c.OrderID)

        if err := paymentService.Charge(ctx, c.PaymentMethod, c.Amount); err != nil {
            return fmt.Errorf("payment failed: %w", err)
        }

        c.ChargedAt = time.Now()
        return nil
    },

    NotifyShipping: func(ctx context.Context, from, to OrderState, c *OrderContext) error {
        return shippingService.NotifyNewOrder(ctx, c.OrderID)
    },

    LogEntry: func(ctx context.Context, state OrderState, c *OrderContext) error {
        log.Printf("Entering state %s for order %s", state, c.OrderID)
        return nil
    },

    LogExit: func(ctx context.Context, state OrderState, c *OrderContext) error {
        log.Printf("Exiting state %s for order %s", state, c.OrderID)
        return nil
    },
}
```

### Action Error Handling

- Actions can return errors to abort the transition
- If an action returns an error, the state does NOT change
- The error is propagated to the caller

```go
action := func(ctx context.Context, from, to State, c *Context) error {
    if err := criticalOperation(); err != nil {
        // Transition will be aborted, state unchanged
        return fmt.Errorf("critical operation failed: %w", err)
    }
    return nil
}
```

### Action Execution Order

For a transition from `StateA` to `StateB` triggered by `Event`:

1. Check guards (if any)
2. Execute `StateA` exit action (if defined)
3. Execute transition action (if defined)
4. Change state from `StateA` to `StateB`
5. Execute `StateB` entry action (if defined)

If any step returns an error, the process stops and state remains unchanged.

## Context Management

### Passing Context Data

```go
// Create context with data
orderCtx := &OrderContext{
    OrderID:       "ORD-123",
    CustomerID:    "CUST-456",
    Amount:        99.99,
    PaymentMethod: "credit_card",
}

// Pass to transition
err := sm.TransitionWithContext(ctx, OrderEventApprove, orderCtx)

// Context is available in guards and actions
```

### Context in Guards

```go
guards := OrderGuards{
    HasPaymentMethod: func(ctx context.Context, c *OrderContext) bool {
        // Access context fields
        return c.PaymentMethod != ""
    },
}
```

### Context in Actions

```go
actions := OrderActions{
    ChargeCard: func(ctx context.Context, from, to OrderState, c *OrderContext) error {
        // Read from context
        amount := c.Amount

        // Modify context
        c.ChargedAt = time.Now()

        return nil
    },
}
```

### Using Go Context

The Go `context.Context` parameter supports:

**Cancellation**:
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// If cancelled, guards/actions should return quickly
err := sm.TransitionWithContext(ctx, event, orderCtx)
```

**Timeout**:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := sm.TransitionWithContext(ctx, event, orderCtx)
```

**Values**:
```go
ctx := context.WithValue(context.Background(), "user", currentUser)

guards := OrderGuards{
    IsAuthorized: func(ctx context.Context, c *OrderContext) bool {
        user, ok := ctx.Value("user").(User)
        return ok && user.IsAuthorized()
    },
}
```

## Error Types

### Standard Errors

```go
var (
    // ErrInvalidTransition is returned when the event is not valid for current state
    ErrInvalidTransition = errors.New("invalid state transition")

    // ErrGuardRejected is returned when a guard function returns false
    ErrGuardRejected = errors.New("guard rejected transition")

    // ErrInvalidState is returned when state value is not valid
    ErrInvalidState = errors.New("invalid state")

    // ErrInvalidEvent is returned when event value is not valid
    ErrInvalidEvent = errors.New("invalid event")
)
```

### Error Checking

```go
err := sm.Transition(ctx, OrderEventApprove)
if err != nil {
    switch {
    case errors.Is(err, ErrInvalidTransition):
        // Event not allowed for current state
        log.Printf("Cannot approve from state %s", sm.State())

    case errors.Is(err, ErrGuardRejected):
        // Guard condition not met
        log.Println("Guard condition failed")

    default:
        // Action returned an error
        log.Printf("Action failed: %v", err)
    }
}
```

### Custom Errors from Actions

Actions can return custom errors:

```go
var ErrPaymentFailed = errors.New("payment processing failed")

actions := OrderActions{
    ChargeCard: func(ctx context.Context, from, to OrderState, c *OrderContext) error {
        if err := processPayment(c); err != nil {
            return fmt.Errorf("%w: %v", ErrPaymentFailed, err)
        }
        return nil
    },
}

// Check for custom error
err := sm.Transition(ctx, OrderEventApprove)
if errors.Is(err, ErrPaymentFailed) {
    // Handle payment failure
}
```

## Options and Configuration

### Option Functions

```go
type OrderStateMachineOption func(*OrderStateMachine)

// WithLogger sets a custom logger
func WithLogger(logger Logger) OrderStateMachineOption

// WithValidationMode enables strict validation
func WithValidationMode(enabled bool) OrderStateMachineOption

// WithMetrics enables metrics collection
func WithMetrics(collector MetricsCollector) OrderStateMachineOption

// WithConcurrencySafe enables mutex protection
func WithConcurrencySafe() OrderStateMachineOption
```

### Using Options

```go
sm := NewOrderStateMachine(
    guards,
    actions,
    WithLogger(logger),
    WithValidationMode(true),
    WithMetrics(metricsCollector),
    WithConcurrencySafe(),
)
```

### Logger Interface

```go
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
}
```

**Example**:
```go
type MyLogger struct{}

func (l *MyLogger) Info(msg string, fields ...Field) {
    log.Printf("INFO: %s %v", msg, fields)
}

// Use with state machine
sm := NewOrderStateMachine(
    guards,
    actions,
    WithLogger(&MyLogger{}),
)
```

### Metrics Interface

```go
type MetricsCollector interface {
    RecordTransition(from, to State, event Event, duration time.Duration)
    RecordGuardEvaluation(guard string, result bool, duration time.Duration)
    RecordActionExecution(action string, err error, duration time.Duration)
}
```

**Example**:
```go
type PrometheusCollector struct {
    transitionCounter *prometheus.CounterVec
    durationHistogram *prometheus.HistogramVec
}

func (c *PrometheusCollector) RecordTransition(from, to State, event Event, duration time.Duration) {
    c.transitionCounter.WithLabelValues(
        from.String(),
        to.String(),
        event.String(),
    ).Inc()

    c.durationHistogram.WithLabelValues(
        from.String(),
        to.String(),
    ).Observe(duration.Seconds())
}
```

## Thread Safety

By default, state machines are **not thread-safe**. For concurrent access:

### Option 1: Use WithConcurrencySafe()

```go
sm := NewOrderStateMachine(
    guards,
    actions,
    WithConcurrencySafe(),
)

// Now safe for concurrent access
go sm.Transition(ctx, OrderEventApprove)
go sm.Transition(ctx, OrderEventReject)
```

### Option 2: External Synchronization

```go
var mu sync.Mutex

mu.Lock()
err := sm.Transition(ctx, event)
mu.Unlock()
```

### Option 3: One State Machine per Goroutine

```go
// Each order has its own state machine
func handleOrder(orderID string) {
    sm := NewOrderStateMachine(guards, actions)
    // Use sm only in this goroutine
}
```

## Next Steps

- See [Usage Guide](usage.md) for practical examples
- Read [YAML Reference](yaml-reference.md) for definition syntax
- Explore [Examples](../examples/) for complete applications
