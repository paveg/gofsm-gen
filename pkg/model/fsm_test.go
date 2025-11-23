package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFSMModel_NewFSMModel(t *testing.T) {
	tests := []struct {
		name         string
		machineName  string
		initialState string
		wantErr      bool
	}{
		{
			name:         "valid order state machine",
			machineName:  "OrderStateMachine",
			initialState: "pending",
			wantErr:      false,
		},
		{
			name:         "valid door lock machine",
			machineName:  "DoorLock",
			initialState: "locked",
			wantErr:      false,
		},
		{
			name:         "empty machine name",
			machineName:  "",
			initialState: "pending",
			wantErr:      true,
		},
		{
			name:         "empty initial state",
			machineName:  "OrderStateMachine",
			initialState: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm, err := NewFSMModel(tt.machineName, tt.initialState)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, fsm)
			} else {
				require.NoError(t, err)
				require.NotNil(t, fsm)
				assert.Equal(t, tt.machineName, fsm.Name)
				assert.Equal(t, tt.initialState, fsm.Initial)
				assert.NotNil(t, fsm.States)
				assert.NotNil(t, fsm.Events)
				assert.NotNil(t, fsm.Transitions)
			}
		})
	}
}

func TestFSMModel_AddState(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	tests := []struct {
		name    string
		state   *State
		wantErr bool
	}{
		{
			name:    "add valid state",
			state:   &State{Name: "approved"},
			wantErr: false,
		},
		{
			name:    "add another valid state",
			state:   &State{Name: "shipped"},
			wantErr: false,
		},
		{
			name:    "add nil state",
			state:   nil,
			wantErr: true,
		},
		{
			name:    "add state with empty name",
			state:   &State{Name: ""},
			wantErr: true,
		},
		{
			name:    "add duplicate state",
			state:   &State{Name: "approved"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fsm.AddState(tt.state)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, fsm.States, tt.state.Name)
			}
		})
	}
}

func TestFSMModel_AddEvent(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	tests := []struct {
		name    string
		event   *Event
		wantErr bool
	}{
		{
			name:    "add valid event",
			event:   &Event{Name: "approve"},
			wantErr: false,
		},
		{
			name:    "add another valid event",
			event:   &Event{Name: "reject"},
			wantErr: false,
		},
		{
			name:    "add nil event",
			event:   nil,
			wantErr: true,
		},
		{
			name:    "add event with empty name",
			event:   &Event{Name: ""},
			wantErr: true,
		},
		{
			name:    "add duplicate event",
			event:   &Event{Name: "approve"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fsm.AddEvent(tt.event)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, fsm.Events, tt.event.Name)
			}
		})
	}
}

func TestFSMModel_AddTransition(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	// Setup states and events
	require.NoError(t, fsm.AddState(&State{Name: "pending"}))
	require.NoError(t, fsm.AddState(&State{Name: "approved"}))
	require.NoError(t, fsm.AddEvent(&Event{Name: "approve"}))

	tests := []struct {
		name       string
		transition *Transition
		wantErr    bool
	}{
		{
			name: "add valid transition",
			transition: &Transition{
				From:  "pending",
				To:    "approved",
				Event: "approve",
			},
			wantErr: false,
		},
		{
			name:       "add nil transition",
			transition: nil,
			wantErr:    true,
		},
		{
			name: "add transition with undefined from state",
			transition: &Transition{
				From:  "unknown",
				To:    "approved",
				Event: "approve",
			},
			wantErr: true,
		},
		{
			name: "add transition with undefined to state",
			transition: &Transition{
				From:  "pending",
				To:    "unknown",
				Event: "approve",
			},
			wantErr: true,
		},
		{
			name: "add transition with undefined event",
			transition: &Transition{
				From:  "pending",
				To:    "approved",
				Event: "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fsm.AddTransition(tt.transition)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, fsm.Transitions, tt.transition)
			}
		})
	}
}

func TestFSMModel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FSMModel
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid complete state machine",
			setup: func() *FSMModel {
				fsm, _ := NewFSMModel("OrderStateMachine", "pending")
				fsm.AddState(&State{Name: "pending"})
				fsm.AddState(&State{Name: "approved"})
				fsm.AddEvent(&Event{Name: "approve"})
				fsm.AddTransition(&Transition{
					From:  "pending",
					To:    "approved",
					Event: "approve",
				})
				return fsm
			},
			wantErr: false,
		},
		{
			name: "initial state not defined",
			setup: func() *FSMModel {
				fsm, _ := NewFSMModel("OrderStateMachine", "pending")
				fsm.AddState(&State{Name: "approved"})
				return fsm
			},
			wantErr: true,
			errMsg:  "initial state",
		},
		{
			name: "no states defined",
			setup: func() *FSMModel {
				fsm, _ := NewFSMModel("OrderStateMachine", "pending")
				return fsm
			},
			wantErr: true,
		},
		{
			name: "no events defined",
			setup: func() *FSMModel {
				fsm, _ := NewFSMModel("OrderStateMachine", "pending")
				fsm.AddState(&State{Name: "pending"})
				return fsm
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := tt.setup()
			err := fsm.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFSMModel_GetState(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	pendingState := &State{Name: "pending"}
	approvedState := &State{Name: "approved"}
	fsm.AddState(pendingState)
	fsm.AddState(approvedState)

	tests := []struct {
		name      string
		stateName string
		want      *State
	}{
		{
			name:      "get existing state",
			stateName: "pending",
			want:      pendingState,
		},
		{
			name:      "get another existing state",
			stateName: "approved",
			want:      approvedState,
		},
		{
			name:      "get non-existent state",
			stateName: "unknown",
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fsm.GetState(tt.stateName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFSMModel_GetEvent(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	approveEvent := &Event{Name: "approve"}
	rejectEvent := &Event{Name: "reject"}
	fsm.AddEvent(approveEvent)
	fsm.AddEvent(rejectEvent)

	tests := []struct {
		name      string
		eventName string
		want      *Event
	}{
		{
			name:      "get existing event",
			eventName: "approve",
			want:      approveEvent,
		},
		{
			name:      "get another existing event",
			eventName: "reject",
			want:      rejectEvent,
		},
		{
			name:      "get non-existent event",
			eventName: "unknown",
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fsm.GetEvent(tt.eventName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFSMModel_GetTransitionsFrom(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	// Setup
	fsm.AddState(&State{Name: "pending"})
	fsm.AddState(&State{Name: "approved"})
	fsm.AddState(&State{Name: "rejected"})
	fsm.AddEvent(&Event{Name: "approve"})
	fsm.AddEvent(&Event{Name: "reject"})

	t1 := &Transition{From: "pending", To: "approved", Event: "approve"}
	t2 := &Transition{From: "pending", To: "rejected", Event: "reject"}
	fsm.AddTransition(t1)
	fsm.AddTransition(t2)

	tests := []struct {
		name      string
		stateName string
		wantCount int
	}{
		{
			name:      "state with multiple outgoing transitions",
			stateName: "pending",
			wantCount: 2,
		},
		{
			name:      "state with no outgoing transitions",
			stateName: "approved",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transitions := fsm.GetTransitionsFrom(tt.stateName)
			assert.Len(t, transitions, tt.wantCount)
		})
	}
}

func TestFSMModel_GetTransitionsTo(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	// Setup
	fsm.AddState(&State{Name: "pending"})
	fsm.AddState(&State{Name: "approved"})
	fsm.AddState(&State{Name: "rejected"})
	fsm.AddEvent(&Event{Name: "approve"})
	fsm.AddEvent(&Event{Name: "reject"})

	t1 := &Transition{From: "pending", To: "approved", Event: "approve"}
	t2 := &Transition{From: "pending", To: "rejected", Event: "reject"}
	fsm.AddTransition(t1)
	fsm.AddTransition(t2)

	tests := []struct {
		name      string
		stateName string
		wantCount int
	}{
		{
			name:      "state with no incoming transitions",
			stateName: "pending",
			wantCount: 0,
		},
		{
			name:      "state with one incoming transition",
			stateName: "approved",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transitions := fsm.GetTransitionsTo(tt.stateName)
			assert.Len(t, transitions, tt.wantCount)
		})
	}
}
