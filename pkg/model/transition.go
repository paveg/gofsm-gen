package model

import "fmt"

// Transition represents a state transition in the finite state machine
type Transition struct {
	// From is the source state
	From string

	// To is the target state
	To string

	// Event is the event that triggers this transition
	Event string

	// Guard is an optional guard condition that must be true for the transition to occur
	Guard string

	// Action is an optional action to execute during the transition
	Action string

	// Description is an optional human-readable description
	Description string
}

// NewTransition creates a new Transition
func NewTransition(from, to, event string) (*Transition, error) {
	if from == "" {
		return nil, fmt.Errorf("from state cannot be empty")
	}

	if to == "" {
		return nil, fmt.Errorf("to state cannot be empty")
	}

	if event == "" {
		return nil, fmt.Errorf("event cannot be empty")
	}

	return &Transition{
		From:  from,
		To:    to,
		Event: event,
	}, nil
}

// WithGuard sets the guard condition for this transition
func (t *Transition) WithGuard(guardName string) error {
	if guardName == "" {
		return fmt.Errorf("guard name cannot be empty")
	}

	t.Guard = guardName
	return nil
}

// WithAction sets the action for this transition
func (t *Transition) WithAction(actionName string) error {
	if actionName == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	t.Action = actionName
	return nil
}

// Validate checks if the transition is valid
func (t *Transition) Validate() error {
	if t.From == "" {
		return fmt.Errorf("from state cannot be empty")
	}

	if t.To == "" {
		return fmt.Errorf("to state cannot be empty")
	}

	if t.Event == "" {
		return fmt.Errorf("event cannot be empty")
	}

	return nil
}

// IsSelfTransition returns true if this is a self-transition (from and to are the same state)
func (t *Transition) IsSelfTransition() bool {
	return t.From == t.To
}
