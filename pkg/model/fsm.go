package model

// FSMModel represents the complete state machine definition
type FSMModel struct {
	Name    string
	Package string
	Initial string
	States  []State
	Events  []Event
	Transitions []Transition
}

// State represents a state in the state machine
type State struct {
	Name  string
	Entry string // Entry action function name
	Exit  string // Exit action function name
}

// Event represents an event that can trigger transitions
type Event struct {
	Name string
}

// Transition represents a state transition
type Transition struct {
	From   string
	To     string
	On     string
	Guard  string // Optional guard function name
	Action string // Optional action function name
}

// GetStateNames returns all state names
func (m *FSMModel) GetStateNames() []string {
	names := make([]string, len(m.States))
	for i, s := range m.States {
		names[i] = s.Name
	}
	return names
}

// GetEventNames returns all event names
func (m *FSMModel) GetEventNames() []string {
	names := make([]string, len(m.Events))
	for i, e := range m.Events {
		names[i] = e.Name
	}
	return names
}

// GetTransitionsFrom returns all transitions from a given state
func (m *FSMModel) GetTransitionsFrom(state string) []Transition {
	var transitions []Transition
	for _, t := range m.Transitions {
		if t.From == state {
			transitions = append(transitions, t)
		}
	}
	return transitions
}
