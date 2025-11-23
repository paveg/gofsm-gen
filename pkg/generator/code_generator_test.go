package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/gofsm-gen/pkg/model"
)

// createOrderStateMachine creates a realistic order state machine model for testing
func createOrderStateMachine(t *testing.T) *model.FSMModel {
	t.Helper()

	fsm, err := model.NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)
	fsm.Package = "orders"

	// Add states
	pending, _ := model.NewState("pending")
	pending.EntryAction = "logEntry"
	pending.ExitAction = "logExit"
	fsm.AddState(pending)

	approved, _ := model.NewState("approved")
	fsm.AddState(approved)

	rejected, _ := model.NewState("rejected")
	fsm.AddState(rejected)

	shipped, _ := model.NewState("shipped")
	shipped.EntryAction = "notifyCustomer"
	fsm.AddState(shipped)

	// Add events
	approve, _ := model.NewEvent("approve")
	fsm.AddEvent(approve)

	reject, _ := model.NewEvent("reject")
	fsm.AddEvent(reject)

	ship, _ := model.NewEvent("ship")
	fsm.AddEvent(ship)

	// Add transitions
	t1, _ := model.NewTransition("pending", "approved", "approve")
	t1.Guard = "hasPayment"
	t1.Action = "chargeCard"
	fsm.AddTransition(t1)

	t2, _ := model.NewTransition("pending", "rejected", "reject")
	t2.Action = "sendRejectionEmail"
	fsm.AddTransition(t2)

	t3, _ := model.NewTransition("approved", "shipped", "ship")
	t3.Action = "notifyShipping"
	fsm.AddTransition(t3)

	return fsm
}

func TestCodeGenerator_Generate_OrderStateMachine(t *testing.T) {
	fsm := createOrderStateMachine(t)

	gen, err := NewCodeGenerator()
	require.NoError(t, err, "Failed to create code generator")

	code, err := gen.Generate(fsm)
	require.NoError(t, err, "Failed to generate code")
	require.NotEmpty(t, code, "Generated code should not be empty")

	codeStr := string(code)

	// Verify package declaration
	assert.Contains(t, codeStr, "package orders", "Should have correct package declaration")

	// Verify state enum generation
	assert.Contains(t, codeStr, "type OrderStateMachineState int", "Should define state type")
	assert.Contains(t, codeStr, "OrderStateMachineStatePending", "Should define pending state constant")
	assert.Contains(t, codeStr, "OrderStateMachineStateApproved", "Should define approved state constant")
	assert.Contains(t, codeStr, "OrderStateMachineStateRejected", "Should define rejected state constant")
	assert.Contains(t, codeStr, "OrderStateMachineStateShipped", "Should define shipped state constant")

	// Verify event enum generation
	assert.Contains(t, codeStr, "type OrderStateMachineEvent int", "Should define event type")
	assert.Contains(t, codeStr, "OrderStateMachineEventApprove", "Should define approve event constant")
	assert.Contains(t, codeStr, "OrderStateMachineEventReject", "Should define reject event constant")
	assert.Contains(t, codeStr, "OrderStateMachineEventShip", "Should define ship event constant")

	// Verify exhaustive annotations
	assert.Contains(t, codeStr, "//exhaustive:enforce", "Should include exhaustive annotations")

	// Verify guard functions are generated
	assert.Contains(t, codeStr, "HasPayment func(ctx context.Context, c *OrderStateMachineContext) bool",
		"Should define hasPayment guard function")

	// Verify action functions are generated
	assert.Contains(t, codeStr, "ChargeCard func(ctx context.Context, from, to OrderStateMachineState, c *OrderStateMachineContext) error",
		"Should define chargeCard action function")
	assert.Contains(t, codeStr, "NotifyShipping func(ctx context.Context, from, to OrderStateMachineState, c *OrderStateMachineContext) error",
		"Should define notifyShipping action function")

	// Verify entry/exit actions
	assert.Contains(t, codeStr, "LogEntry func(ctx context.Context, c *OrderStateMachineContext) error",
		"Should define logEntry entry action")
	assert.Contains(t, codeStr, "NotifyCustomer func(ctx context.Context, c *OrderStateMachineContext) error",
		"Should define notifyCustomer entry action")

	// Verify constructor function
	assert.Contains(t, codeStr, "func NewOrderStateMachine(", "Should define constructor function")

	// Verify core API methods
	assert.Contains(t, codeStr, "func (sm *OrderStateMachine) State()", "Should define State method")
	assert.Contains(t, codeStr, "func (sm *OrderStateMachine) Transition(", "Should define Transition method")
	assert.Contains(t, codeStr, "func (sm *OrderStateMachine) PermittedEvents()", "Should define PermittedEvents method")
	assert.Contains(t, codeStr, "func (sm *OrderStateMachine) CanTransition(", "Should define CanTransition method")

	// Verify initial state is set correctly
	assert.Contains(t, codeStr, "currentState: OrderStateMachineStatePending", "Should set initial state to pending")
}

func TestCodeGenerator_Generate_SimpleDoorLock(t *testing.T) {
	// Test with a simpler state machine without guards/actions
	fsm, err := model.NewFSMModel("DoorLock", "locked")
	require.NoError(t, err)
	fsm.Package = "security"

	locked, _ := model.NewState("locked")
	fsm.AddState(locked)

	unlocked, _ := model.NewState("unlocked")
	fsm.AddState(unlocked)

	lockEvent, _ := model.NewEvent("lock")
	fsm.AddEvent(lockEvent)

	unlockEvent, _ := model.NewEvent("unlock")
	fsm.AddEvent(unlockEvent)

	t1, _ := model.NewTransition("locked", "unlocked", "unlock")
	fsm.AddTransition(t1)

	t2, _ := model.NewTransition("unlocked", "locked", "lock")
	fsm.AddTransition(t2)

	gen, err := NewCodeGenerator()
	require.NoError(t, err)

	code, err := gen.Generate(fsm)
	require.NoError(t, err)
	require.NotEmpty(t, code)

	codeStr := string(code)

	// Verify basic structure
	assert.Contains(t, codeStr, "package security")
	assert.Contains(t, codeStr, "type DoorLockState int")
	assert.Contains(t, codeStr, "DoorLockStateLocked")
	assert.Contains(t, codeStr, "DoorLockStateUnlocked")
	assert.Contains(t, codeStr, "type DoorLockEvent int")
	assert.Contains(t, codeStr, "DoorLockEventLock")
	assert.Contains(t, codeStr, "DoorLockEventUnlock")
}

func TestCodeGenerator_Generate_NilModel(t *testing.T) {
	gen, err := NewCodeGenerator()
	require.NoError(t, err)

	_, err = gen.Generate(nil)
	assert.Error(t, err, "Should return error for nil model")
	assert.Contains(t, err.Error(), "model cannot be nil")
}

func TestCodeGenerator_Generate_DefaultPackage(t *testing.T) {
	// Test that package defaults to "main" if not specified
	fsm, err := model.NewFSMModel("TestMachine", "idle")
	require.NoError(t, err)
	// Don't set Package, should default to "main"

	idle, _ := model.NewState("idle")
	fsm.AddState(idle)

	dummyEvent, _ := model.NewEvent("dummy")
	fsm.AddEvent(dummyEvent)

	gen, err := NewCodeGenerator()
	require.NoError(t, err)

	code, err := gen.Generate(fsm)
	require.NoError(t, err)

	codeStr := string(code)
	assert.Contains(t, codeStr, "package main", "Should default to main package")
}

func TestCodeGenerator_GenerateTo(t *testing.T) {
	fsm, err := model.NewFSMModel("TestMachine", "start")
	require.NoError(t, err)
	fsm.Package = "test"

	start, _ := model.NewState("start")
	fsm.AddState(start)

	end, _ := model.NewState("end")
	fsm.AddState(end)

	proceed, _ := model.NewEvent("proceed")
	fsm.AddEvent(proceed)

	t1, _ := model.NewTransition("start", "end", "proceed")
	fsm.AddTransition(t1)

	gen, err := NewCodeGenerator()
	require.NoError(t, err)

	var buf strings.Builder
	err = gen.GenerateTo(fsm, &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "package test")
	assert.Contains(t, output, "type TestMachineState int")
}

func TestTemplateFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function string
		input    string
		expected string
	}{
		{
			name:     "title case - simple",
			function: "title",
			input:    "pending",
			expected: "Pending",
		},
		{
			name:     "title case - snake_case",
			function: "title",
			input:    "order_approved",
			expected: "OrderApproved",
		},
		{
			name:     "title case - kebab-case",
			function: "title",
			input:    "user-logged-in",
			expected: "UserLoggedIn",
		},
		{
			name:     "camelCase - simple",
			function: "camelCase",
			input:    "has_payment",
			expected: "hasPayment",
		},
		{
			name:     "snakeCase - PascalCase",
			function: "snakeCase",
			input:    "OrderApproved",
			expected: "order_approved",
		},
	}

	funcs := TemplateFuncs()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, ok := funcs[tt.function]
			require.True(t, ok, "Function %s should exist", tt.function)

			funcCall := fn.(func(string) string)
			result := funcCall(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
