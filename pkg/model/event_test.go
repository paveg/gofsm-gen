package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvent_NewEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventName string
		wantErr   bool
	}{
		{
			name:      "valid event name",
			eventName: "approve",
			wantErr:   false,
		},
		{
			name:      "valid event with underscore",
			eventName: "submit_order",
			wantErr:   false,
		},
		{
			name:      "empty event name",
			eventName: "",
			wantErr:   true,
		},
		{
			name:      "event name with spaces",
			eventName: "submit order",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := NewEvent(tt.eventName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, tt.eventName, event.Name)
			}
		})
	}
}

func TestEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   *Event
		wantErr bool
	}{
		{
			name: "valid event",
			event: &Event{
				Name: "approve",
			},
			wantErr: false,
		},
		{
			name: "valid event with metadata",
			event: &Event{
				Name:        "approve",
				Description: "Approve the order",
			},
			wantErr: false,
		},
		{
			name: "invalid event with empty name",
			event: &Event{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "invalid event name with spaces",
			event: &Event{
				Name: "submit order",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
