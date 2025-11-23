package model

import (
	"fmt"
	"regexp"
)

// State represents a single state in the finite state machine
type State struct {
	// Name is the unique identifier for this state
	Name string

	// EntryAction is the optional action to execute when entering this state
	EntryAction string

	// ExitAction is the optional action to execute when leaving this state
	ExitAction string

	// Description is an optional human-readable description
	Description string
}

// validNamePattern matches valid Go identifiers (letters, digits, underscores)
var validNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// NewState creates a new State with the given name
func NewState(name string) (*State, error) {
	if name == "" {
		return nil, fmt.Errorf("state name cannot be empty")
	}

	if !validNamePattern.MatchString(name) {
		return nil, fmt.Errorf("state name %q contains invalid characters (use only letters, digits, and underscores)", name)
	}

	return &State{
		Name: name,
	}, nil
}

// WithEntryAction sets the entry action for this state
func (s *State) WithEntryAction(actionName string) error {
	if actionName == "" {
		return fmt.Errorf("entry action name cannot be empty")
	}

	s.EntryAction = actionName
	return nil
}

// WithExitAction sets the exit action for this state
func (s *State) WithExitAction(actionName string) error {
	if actionName == "" {
		return fmt.Errorf("exit action name cannot be empty")
	}

	s.ExitAction = actionName
	return nil
}

// Validate checks if the state is valid
func (s *State) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("state name cannot be empty")
	}

	if !validNamePattern.MatchString(s.Name) {
		return fmt.Errorf("state name %q contains invalid characters (use only letters, digits, and underscores)", s.Name)
	}

	return nil
}
