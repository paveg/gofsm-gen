package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/gofsm-gen/pkg/model"
)

func TestCodeGenerator_Generate_OrderStateMachine(t *testing.T) {
	// Create a realistic order state machine model
	fsm := &model.FSMModel{
		Name:    "OrderStateMachine",
		Package: "orders",
		Initial: "pending",
		States: []model.State{
			{Name: "pending", Entry: "logEntry", Exit: "logExit"},
			{Name: "approved", Entry: "", Exit: ""},
			{Name: "rejected", Entry: "", Exit: ""},
			{Name: "shipped", Entry: "notifyCustomer", Exit: ""},
		},
		Events: []model.Event{
			{Name: "approve"},
			{Name: "reject"},
			{Name: "ship"},
		},
		Transitions: []model.Transition{
			{From: "pending", To: "approved", On: "approve", Guard: "hasPayment", Action: "chargeCard"},
			{From: "pending", To: "rejected", On: "reject", Guard: "", Action: "sendRejectionEmail"},
			{From: "approved", To: "shipped", On: "ship", Guard: "", Action: "notifyShipping"},
		},
	}

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
	fsm := &model.FSMModel{
		Name:    "DoorLock",
		Package: "security",
		Initial: "locked",
		States: []model.State{
			{Name: "locked", Entry: "", Exit: ""},
			{Name: "unlocked", Entry: "", Exit: ""},
		},
		Events: []model.Event{
			{Name: "lock"},
			{Name: "unlock"},
		},
		Transitions: []model.Transition{
			{From: "locked", To: "unlocked", On: "unlock", Guard: "", Action: ""},
			{From: "unlocked", To: "locked", On: "lock", Guard: "", Action: ""},
		},
	}

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
	fsm := &model.FSMModel{
		Name:    "TestMachine",
		Package: "", // Empty package
		Initial: "idle",
		States: []model.State{
			{Name: "idle", Entry: "", Exit: ""},
		},
		Events: []model.Event{},
		Transitions: []model.Transition{},
	}

	gen, err := NewCodeGenerator()
	require.NoError(t, err)

	code, err := gen.Generate(fsm)
	require.NoError(t, err)

	codeStr := string(code)
	assert.Contains(t, codeStr, "package main", "Should default to main package")
}

func TestCodeGenerator_GenerateTo(t *testing.T) {
	fsm := &model.FSMModel{
		Name:    "TestMachine",
		Package: "test",
		Initial: "start",
		States: []model.State{
			{Name: "start", Entry: "", Exit: ""},
			{Name: "end", Entry: "", Exit: ""},
		},
		Events: []model.Event{
			{Name: "proceed"},
		},
		Transitions: []model.Transition{
			{From: "start", To: "end", On: "proceed", Guard: "", Action: ""},
		},
	}

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
