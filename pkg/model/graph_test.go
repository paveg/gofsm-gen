package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateGraph_NewStateGraph(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	graph := NewStateGraph(fsm)
	assert.NotNil(t, graph)
	assert.Equal(t, fsm, graph.FSM)
}

func TestStateGraph_Build(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FSMModel
		wantErr bool
	}{
		{
			name: "build simple linear graph",
			setup: func() *FSMModel {
				fsm, _ := NewFSMModel("OrderStateMachine", "pending")
				fsm.AddState(&State{Name: "pending"})
				fsm.AddState(&State{Name: "approved"})
				fsm.AddState(&State{Name: "shipped"})
				fsm.AddEvent(&Event{Name: "approve"})
				fsm.AddEvent(&Event{Name: "ship"})
				fsm.AddTransition(&Transition{From: "pending", To: "approved", Event: "approve"})
				fsm.AddTransition(&Transition{From: "approved", To: "shipped", Event: "ship"})
				return fsm
			},
			wantErr: false,
		},
		{
			name: "build graph with branching",
			setup: func() *FSMModel {
				fsm, _ := NewFSMModel("OrderStateMachine", "pending")
				fsm.AddState(&State{Name: "pending"})
				fsm.AddState(&State{Name: "approved"})
				fsm.AddState(&State{Name: "rejected"})
				fsm.AddEvent(&Event{Name: "approve"})
				fsm.AddEvent(&Event{Name: "reject"})
				fsm.AddTransition(&Transition{From: "pending", To: "approved", Event: "approve"})
				fsm.AddTransition(&Transition{From: "pending", To: "rejected", Event: "reject"})
				return fsm
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := tt.setup()
			graph := NewStateGraph(fsm)

			err := graph.Build()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, graph.adjacencyList)
			}
		})
	}
}

func TestStateGraph_GetOutgoingTransitions(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	// Setup states
	fsm.AddState(&State{Name: "pending"})
	fsm.AddState(&State{Name: "approved"})
	fsm.AddState(&State{Name: "rejected"})

	// Setup events
	fsm.AddEvent(&Event{Name: "approve"})
	fsm.AddEvent(&Event{Name: "reject"})

	// Setup transitions
	t1 := &Transition{From: "pending", To: "approved", Event: "approve"}
	t2 := &Transition{From: "pending", To: "rejected", Event: "reject"}
	fsm.AddTransition(t1)
	fsm.AddTransition(t2)

	graph := NewStateGraph(fsm)
	graph.Build()

	tests := []struct {
		name      string
		state     string
		wantCount int
	}{
		{
			name:      "state with multiple outgoing transitions",
			state:     "pending",
			wantCount: 2,
		},
		{
			name:      "state with no outgoing transitions",
			state:     "approved",
			wantCount: 0,
		},
		{
			name:      "unknown state",
			state:     "unknown",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transitions := graph.GetOutgoingTransitions(tt.state)
			assert.Len(t, transitions, tt.wantCount)
		})
	}
}

func TestStateGraph_GetIncomingTransitions(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	// Setup states
	fsm.AddState(&State{Name: "pending"})
	fsm.AddState(&State{Name: "approved"})
	fsm.AddState(&State{Name: "rejected"})
	fsm.AddState(&State{Name: "shipped"})

	// Setup events
	fsm.AddEvent(&Event{Name: "approve"})
	fsm.AddEvent(&Event{Name: "reject"})
	fsm.AddEvent(&Event{Name: "ship"})

	// Setup transitions
	fsm.AddTransition(&Transition{From: "pending", To: "approved", Event: "approve"})
	fsm.AddTransition(&Transition{From: "pending", To: "rejected", Event: "reject"})
	fsm.AddTransition(&Transition{From: "approved", To: "shipped", Event: "ship"})

	graph := NewStateGraph(fsm)
	graph.Build()

	tests := []struct {
		name      string
		state     string
		wantCount int
	}{
		{
			name:      "state with no incoming transitions",
			state:     "pending",
			wantCount: 0,
		},
		{
			name:      "state with one incoming transition",
			state:     "approved",
			wantCount: 1,
		},
		{
			name:      "unknown state",
			state:     "unknown",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transitions := graph.GetIncomingTransitions(tt.state)
			assert.Len(t, transitions, tt.wantCount)
		})
	}
}

func TestStateGraph_IsReachable(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	// Setup states
	fsm.AddState(&State{Name: "pending"})
	fsm.AddState(&State{Name: "approved"})
	fsm.AddState(&State{Name: "shipped"})
	fsm.AddState(&State{Name: "orphan"})

	// Setup events
	fsm.AddEvent(&Event{Name: "approve"})
	fsm.AddEvent(&Event{Name: "ship"})

	// Setup transitions (orphan state is not connected)
	fsm.AddTransition(&Transition{From: "pending", To: "approved", Event: "approve"})
	fsm.AddTransition(&Transition{From: "approved", To: "shipped", Event: "ship"})

	graph := NewStateGraph(fsm)
	graph.Build()

	tests := []struct {
		name  string
		state string
		want  bool
	}{
		{
			name:  "initial state is reachable",
			state: "pending",
			want:  true,
		},
		{
			name:  "directly connected state is reachable",
			state: "approved",
			want:  true,
		},
		{
			name:  "indirectly connected state is reachable",
			state: "shipped",
			want:  true,
		},
		{
			name:  "orphan state is not reachable",
			state: "orphan",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := graph.IsReachable(tt.state)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStateGraph_GetUnreachableStates(t *testing.T) {
	fsm, err := NewFSMModel("OrderStateMachine", "pending")
	require.NoError(t, err)

	// Setup states
	fsm.AddState(&State{Name: "pending"})
	fsm.AddState(&State{Name: "approved"})
	fsm.AddState(&State{Name: "shipped"})
	fsm.AddState(&State{Name: "orphan1"})
	fsm.AddState(&State{Name: "orphan2"})

	// Setup events
	fsm.AddEvent(&Event{Name: "approve"})
	fsm.AddEvent(&Event{Name: "ship"})

	// Setup transitions (orphan states are not connected)
	fsm.AddTransition(&Transition{From: "pending", To: "approved", Event: "approve"})
	fsm.AddTransition(&Transition{From: "approved", To: "shipped", Event: "ship"})

	graph := NewStateGraph(fsm)
	graph.Build()

	unreachable := graph.GetUnreachableStates()
	assert.Len(t, unreachable, 2)
	assert.Contains(t, unreachable, "orphan1")
	assert.Contains(t, unreachable, "orphan2")
}

func TestStateGraph_HasCycles(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *FSMModel
		want  bool
	}{
		{
			name: "linear graph has no cycles",
			setup: func() *FSMModel {
				fsm, _ := NewFSMModel("OrderStateMachine", "pending")
				fsm.AddState(&State{Name: "pending"})
				fsm.AddState(&State{Name: "approved"})
				fsm.AddState(&State{Name: "shipped"})
				fsm.AddEvent(&Event{Name: "approve"})
				fsm.AddEvent(&Event{Name: "ship"})
				fsm.AddTransition(&Transition{From: "pending", To: "approved", Event: "approve"})
				fsm.AddTransition(&Transition{From: "approved", To: "shipped", Event: "ship"})
				return fsm
			},
			want: false,
		},
		{
			name: "graph with self-transition has cycle",
			setup: func() *FSMModel {
				fsm, _ := NewFSMModel("OrderStateMachine", "pending")
				fsm.AddState(&State{Name: "pending"})
				fsm.AddEvent(&Event{Name: "refresh"})
				fsm.AddTransition(&Transition{From: "pending", To: "pending", Event: "refresh"})
				return fsm
			},
			want: true,
		},
		{
			name: "graph with cycle between states",
			setup: func() *FSMModel {
				fsm, _ := NewFSMModel("OrderStateMachine", "pending")
				fsm.AddState(&State{Name: "pending"})
				fsm.AddState(&State{Name: "approved"})
				fsm.AddEvent(&Event{Name: "approve"})
				fsm.AddEvent(&Event{Name: "reject"})
				fsm.AddTransition(&Transition{From: "pending", To: "approved", Event: "approve"})
				fsm.AddTransition(&Transition{From: "approved", To: "pending", Event: "reject"})
				return fsm
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := tt.setup()
			graph := NewStateGraph(fsm)
			graph.Build()

			got := graph.HasCycles()
			assert.Equal(t, tt.want, got)
		})
	}
}
