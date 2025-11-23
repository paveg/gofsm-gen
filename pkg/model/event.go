package model

import "fmt"

// Event represents an event that can trigger state transitions
type Event struct {
	// Name is the unique identifier for this event
	Name string

	// Description is an optional human-readable description
	Description string
}

// NewEvent creates a new Event with the given name
func NewEvent(name string) (*Event, error) {
	if name == "" {
		return nil, fmt.Errorf("event name cannot be empty")
	}

	if !validNamePattern.MatchString(name) {
		return nil, fmt.Errorf("event name %q contains invalid characters (use only letters, digits, and underscores)", name)
	}

	return &Event{
		Name: name,
	}, nil
}

// Validate checks if the event is valid
func (e *Event) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("event name cannot be empty")
	}

	if !validNamePattern.MatchString(e.Name) {
		return fmt.Errorf("event name %q contains invalid characters (use only letters, digits, and underscores)", e.Name)
	}

	return nil
}
