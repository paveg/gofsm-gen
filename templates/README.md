# gofsm-gen Templates

This directory contains code generation templates for the gofsm-gen state machine library.

## Template Files

### state_machine.tmpl

The main template that generates type-safe state machine code from FSM model definitions.

#### Generated Code Structure

The template generates the following components:

1. **State Enum**
   - Type-safe integer constants for each state
   - Exhaustive switch enforcement annotations
   - String() method for debugging

2. **Event Enum**
   - Type-safe integer constants for each event
   - Exhaustive switch enforcement annotations
   - String() method for debugging

3. **Context Structure**
   - Custom context type for passing data through transitions
   - User-extensible for domain-specific fields

4. **Guard Functions**
   - Type-safe guard function interfaces
   - Predicates that control whether transitions can execute
   - Optional - only generated if guards are defined

5. **Action Functions**
   - Type-safe action function interfaces
   - Code executed during state transitions
   - Receives from/to state and context

6. **Entry/Exit Actions**
   - State-specific entry actions (executed when entering a state)
   - State-specific exit actions (executed when leaving a state)
   - Optional - only generated if defined

7. **State Machine Type**
   - Thread-safe implementation with mutex
   - Functional options for configuration
   - Logger interface for observability

8. **Core Methods**
   - `State()` - Get current state
   - `Context()` - Get context
   - `SetContext()` - Update context
   - `Transition()` - Trigger state transition
   - `PermittedEvents()` - Get valid events for current state
   - `CanTransition()` - Check if transition is possible

#### Template Functions

Custom template functions available for use:

- `title` - Convert to PascalCase (e.g., "pending" → "Pending", "order_approved" → "OrderApproved")
- `lower` - Convert to lowercase
- `upper` - Convert to uppercase
- `camelCase` - Convert to camelCase (e.g., "has_payment" → "hasPayment")
- `snakeCase` - Convert to snake_case (e.g., "OrderApproved" → "order_approved")

#### Model Methods Used

The template relies on these FSMModel methods:

- `GetStateNames()` - Returns all state names
- `GetEventNames()` - Returns all event names
- `GetTransitionsFrom(state)` - Returns transitions from a specific state

#### Exhaustive Checking

The generated code includes `//exhaustive:enforce` annotations that work with the `exhaustive` static analysis tool to ensure:

- All states are handled in switch statements
- All events are handled for each state
- No unreachable code paths
- Compile-time safety similar to Rust's enum exhaustiveness

Example generated switch with exhaustive checking:

```go
//exhaustive:enforce
switch currentState {
case OrderStateMachineStatePending:
    //exhaustive:enforce
    switch event {
    case OrderStateMachineEventApprove:
        // transition logic
    case OrderStateMachineEventReject:
        // transition logic
    default:
        return fmt.Errorf("invalid event")
    }
// ... other states
}
```

#### Thread Safety

All generated state machines are thread-safe:

- Uses `sync.RWMutex` for concurrent access
- Read methods use `RLock()`
- Write methods use `Lock()`
- Safe for use in concurrent goroutines

#### Usage Example

Given a YAML definition:

```yaml
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
    guard: hasPayment
    action: chargeCard
  - from: approved
    to: shipped
    on: ship
```

The generated code can be used as:

```go
// Define guards and actions
guards := OrderStateMachineGuards{
    HasPayment: func(ctx context.Context, c *OrderStateMachineContext) bool {
        return c.PaymentMethod != ""
    },
}

actions := OrderStateMachineActions{
    ChargeCard: func(ctx context.Context, from, to OrderStateMachineState, c *OrderStateMachineContext) error {
        return processPayment(c)
    },
}

// Create state machine
sm := NewOrderStateMachine(guards, actions,
    WithLogger(logger),
    WithValidationMode(true),
)

// Trigger transitions
err := sm.Transition(ctx, OrderStateMachineEventApprove)

// Check current state
state := sm.State()

// Get permitted events
events := sm.PermittedEvents()
```

## Template Development

### Testing Templates

Run the generator tests to validate template changes:

```bash
go test ./pkg/generator/
```

### Adding New Features

When extending the template:

1. Update the `FSMModel` in `pkg/model/fsm.go` if new fields are needed
2. Add corresponding template logic in `state_machine.tmpl`
3. Update template function helpers in `pkg/generator/template_funcs.go` if needed
4. Add tests in `pkg/generator/code_generator_test.go`
5. Update this documentation

### Template Syntax

The templates use Go's `text/template` syntax:

- `{{.Field}}` - Access field
- `{{range .Items}}...{{end}}` - Iterate over slice
- `{{if .Condition}}...{{end}}` - Conditional
- `{{- /* trim whitespace */ -}}` - Control whitespace
- `{{$var := .Value}}` - Assign variable
- `{{.Field | title}}` - Apply template function

## Future Templates

Planned additional templates:

- `test.tmpl` - Generate unit tests for state machines
- `mock.tmpl` - Generate mock implementations for testing
- `diagram.tmpl` - Generate Mermaid/Graphviz diagrams
- `serialization.tmpl` - Generate JSON/protobuf serialization code
