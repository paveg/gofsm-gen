package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_NewState(t *testing.T) {
	tests := []struct {
		name      string
		stateName string
		wantErr   bool
	}{
		{
			name:      "valid state name",
			stateName: "pending",
			wantErr:   false,
		},
		{
			name:      "valid state with underscore",
			stateName: "in_progress",
			wantErr:   false,
		},
		{
			name:      "empty state name",
			stateName: "",
			wantErr:   true,
		},
		{
			name:      "state name with spaces",
			stateName: "in progress",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, err := NewState(tt.stateName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, state)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, state)
				assert.Equal(t, tt.stateName, state.Name)
			}
		})
	}
}

func TestState_WithEntryAction(t *testing.T) {
	state, err := NewState("pending")
	assert.NoError(t, err)

	tests := []struct {
		name       string
		actionName string
		wantErr    bool
	}{
		{
			name:       "valid entry action",
			actionName: "logEntry",
			wantErr:    false,
		},
		{
			name:       "empty entry action",
			actionName: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := state.WithEntryAction(tt.actionName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.actionName, state.EntryAction)
			}
		})
	}
}

func TestState_WithExitAction(t *testing.T) {
	state, err := NewState("pending")
	assert.NoError(t, err)

	tests := []struct {
		name       string
		actionName string
		wantErr    bool
	}{
		{
			name:       "valid exit action",
			actionName: "logExit",
			wantErr:    false,
		},
		{
			name:       "empty exit action",
			actionName: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := state.WithExitAction(tt.actionName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.actionName, state.ExitAction)
			}
		})
	}
}

func TestState_Validate(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		wantErr bool
	}{
		{
			name: "valid state with no actions",
			state: &State{
				Name: "pending",
			},
			wantErr: false,
		},
		{
			name: "valid state with entry action",
			state: &State{
				Name:        "pending",
				EntryAction: "logEntry",
			},
			wantErr: false,
		},
		{
			name: "valid state with both actions",
			state: &State{
				Name:        "pending",
				EntryAction: "logEntry",
				ExitAction:  "logExit",
			},
			wantErr: false,
		},
		{
			name: "invalid state with empty name",
			state: &State{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "invalid state name with spaces",
			state: &State{
				Name: "in progress",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.state.Validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
