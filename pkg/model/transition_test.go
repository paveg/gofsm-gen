package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransition_NewTransition(t *testing.T) {
	tests := []struct {
		name      string
		from      string
		to        string
		event     string
		wantErr   bool
	}{
		{
			name:    "valid transition",
			from:    "pending",
			to:      "approved",
			event:   "approve",
			wantErr: false,
		},
		{
			name:    "valid self-transition",
			from:    "pending",
			to:      "pending",
			event:   "refresh",
			wantErr: false,
		},
		{
			name:    "empty from state",
			from:    "",
			to:      "approved",
			event:   "approve",
			wantErr: true,
		},
		{
			name:    "empty to state",
			from:    "pending",
			to:      "",
			event:   "approve",
			wantErr: true,
		},
		{
			name:    "empty event",
			from:    "pending",
			to:      "approved",
			event:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transition, err := NewTransition(tt.from, tt.to, tt.event)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, transition)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transition)
				assert.Equal(t, tt.from, transition.From)
				assert.Equal(t, tt.to, transition.To)
				assert.Equal(t, tt.event, transition.Event)
			}
		})
	}
}

func TestTransition_WithGuard(t *testing.T) {
	transition, err := NewTransition("pending", "approved", "approve")
	assert.NoError(t, err)

	tests := []struct {
		name      string
		guardName string
		wantErr   bool
	}{
		{
			name:      "valid guard",
			guardName: "hasPayment",
			wantErr:   false,
		},
		{
			name:      "empty guard",
			guardName: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := transition.WithGuard(tt.guardName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.guardName, transition.Guard)
			}
		})
	}
}

func TestTransition_WithAction(t *testing.T) {
	transition, err := NewTransition("pending", "approved", "approve")
	assert.NoError(t, err)

	tests := []struct {
		name       string
		actionName string
		wantErr    bool
	}{
		{
			name:       "valid action",
			actionName: "chargeCard",
			wantErr:    false,
		},
		{
			name:       "empty action",
			actionName: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := transition.WithAction(tt.actionName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.actionName, transition.Action)
			}
		})
	}
}

func TestTransition_Validate(t *testing.T) {
	tests := []struct {
		name       string
		transition *Transition
		wantErr    bool
	}{
		{
			name: "valid basic transition",
			transition: &Transition{
				From:  "pending",
				To:    "approved",
				Event: "approve",
			},
			wantErr: false,
		},
		{
			name: "valid transition with guard",
			transition: &Transition{
				From:  "pending",
				To:    "approved",
				Event: "approve",
				Guard: "hasPayment",
			},
			wantErr: false,
		},
		{
			name: "valid transition with action",
			transition: &Transition{
				From:   "pending",
				To:     "approved",
				Event:  "approve",
				Action: "chargeCard",
			},
			wantErr: false,
		},
		{
			name: "valid transition with guard and action",
			transition: &Transition{
				From:   "pending",
				To:     "approved",
				Event:  "approve",
				Guard:  "hasPayment",
				Action: "chargeCard",
			},
			wantErr: false,
		},
		{
			name: "invalid transition with empty from",
			transition: &Transition{
				From:  "",
				To:    "approved",
				Event: "approve",
			},
			wantErr: true,
		},
		{
			name: "invalid transition with empty to",
			transition: &Transition{
				From:  "pending",
				To:    "",
				Event: "approve",
			},
			wantErr: true,
		},
		{
			name: "invalid transition with empty event",
			transition: &Transition{
				From:  "pending",
				To:    "approved",
				Event: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.transition.Validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTransition_IsSelfTransition(t *testing.T) {
	tests := []struct {
		name       string
		transition *Transition
		want       bool
	}{
		{
			name: "self transition",
			transition: &Transition{
				From:  "pending",
				To:    "pending",
				Event: "refresh",
			},
			want: true,
		},
		{
			name: "non-self transition",
			transition: &Transition{
				From:  "pending",
				To:    "approved",
				Event: "approve",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.transition.IsSelfTransition()
			assert.Equal(t, tt.want, got)
		})
	}
}
