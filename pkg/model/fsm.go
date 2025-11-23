package model

import "fmt"

// FSMModel represents the complete finite state machine model
type FSMModel struct {
	// Name is the name of the state machine
	Name string

	// Initial is the initial state
	Initial string

	// States is a map of state name to State
	States map[string]*State

	// Events is a map of event name to Event
	Events map[string]*Event

	// Transitions is a list of all transitions
	Transitions []*Transition

	// Package is the Go package name for generated code
	Package string

	// Description is an optional human-readable description
	Description string
}

// NewFSMModel creates a new FSMModel with the given name and initial state
func NewFSMModel(name, initial string) (*FSMModel, error) {
	if name == "" {
		return nil, fmt.Errorf("machine name cannot be empty")
	}

	if initial == "" {
		return nil, fmt.Errorf("initial state cannot be empty")
	}

	return &FSMModel{
		Name:        name,
		Initial:     initial,
		States:      make(map[string]*State),
		Events:      make(map[string]*Event),
		Transitions: make([]*Transition, 0),
	}, nil
}

// AddState adds a state to the FSM
func (f *FSMModel) AddState(state *State) error {
	if state == nil {
		return fmt.Errorf("cannot add nil state")
	}

	if state.Name == "" {
		return fmt.Errorf("state name cannot be empty")
	}

	if _, exists := f.States[state.Name]; exists {
		return fmt.Errorf("state %q already exists", state.Name)
	}

	f.States[state.Name] = state
	return nil
}

// AddEvent adds an event to the FSM
func (f *FSMModel) AddEvent(event *Event) error {
	if event == nil {
		return fmt.Errorf("cannot add nil event")
	}

	if event.Name == "" {
		return fmt.Errorf("event name cannot be empty")
	}

	if _, exists := f.Events[event.Name]; exists {
		return fmt.Errorf("event %q already exists", event.Name)
	}

	f.Events[event.Name] = event
	return nil
}

// AddTransition adds a transition to the FSM
func (f *FSMModel) AddTransition(transition *Transition) error {
	if transition == nil {
		return fmt.Errorf("cannot add nil transition")
	}

	// Validate that the from state exists
	if _, exists := f.States[transition.From]; !exists {
		return fmt.Errorf("from state %q is not defined", transition.From)
	}

	// Validate that the to state exists
	if _, exists := f.States[transition.To]; !exists {
		return fmt.Errorf("to state %q is not defined", transition.To)
	}

	// Validate that the event exists
	if _, exists := f.Events[transition.Event]; !exists {
		return fmt.Errorf("event %q is not defined", transition.Event)
	}

	f.Transitions = append(f.Transitions, transition)
	return nil
}

// Validate checks if the FSM model is valid
func (f *FSMModel) Validate() error {
	// Check that initial state is defined
	if _, exists := f.States[f.Initial]; !exists {
		return fmt.Errorf("initial state %q is not defined", f.Initial)
	}

	// Check that there is at least one state
	if len(f.States) == 0 {
		return fmt.Errorf("FSM must have at least one state")
	}

	// Check that there is at least one event
	if len(f.Events) == 0 {
		return fmt.Errorf("FSM must have at least one event")
	}

	// Validate all states
	for _, state := range f.States {
		if err := state.Validate(); err != nil {
			return fmt.Errorf("invalid state: %w", err)
		}
	}

	// Validate all events
	for _, event := range f.Events {
		if err := event.Validate(); err != nil {
			return fmt.Errorf("invalid event: %w", err)
		}
	}

	// Validate all transitions
	for _, transition := range f.Transitions {
		if err := transition.Validate(); err != nil {
			return fmt.Errorf("invalid transition: %w", err)
		}
	}

	return nil
}

// GetState returns the state with the given name, or nil if not found
func (f *FSMModel) GetState(name string) *State {
	return f.States[name]
}

// GetEvent returns the event with the given name, or nil if not found
func (f *FSMModel) GetEvent(name string) *Event {
	return f.Events[name]
}

// GetTransitionsFrom returns all transitions from the given state
func (f *FSMModel) GetTransitionsFrom(stateName string) []*Transition {
	transitions := make([]*Transition, 0)
	for _, t := range f.Transitions {
		if t.From == stateName {
			transitions = append(transitions, t)
		}
	}
	return transitions
}

// GetTransitionsTo returns all transitions to the given state
func (f *FSMModel) GetTransitionsTo(stateName string) []*Transition {
	transitions := make([]*Transition, 0)
	for _, t := range f.Transitions {
		if t.To == stateName {
			transitions = append(transitions, t)
		}
	}
	return transitions
}

// GetStateNames returns all state names (for template compatibility)
func (f *FSMModel) GetStateNames() []string {
	names := make([]string, 0, len(f.States))
	for name := range f.States {
		names = append(names, name)
	}
	return names
}

// GetEventNames returns all event names (for template compatibility)
func (f *FSMModel) GetEventNames() []string {
	names := make([]string, 0, len(f.Events))
	for name := range f.Events {
		names = append(names, name)
	}
	return names
}

// GetStatesSlice returns states as a slice (for template compatibility)
func (f *FSMModel) GetStatesSlice() []*State {
	states := make([]*State, 0, len(f.States))
	for _, state := range f.States {
		states = append(states, state)
	}
	return states
}

// GetEventsSlice returns events as a slice (for template compatibility)
func (f *FSMModel) GetEventsSlice() []*Event {
	events := make([]*Event, 0, len(f.Events))
	for _, event := range f.Events {
		events = append(events, event)
	}
	return events
}
