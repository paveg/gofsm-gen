package model

// StateGraph represents a graph-based view of the FSM for analysis
type StateGraph struct {
	// FSM is the underlying FSM model
	FSM *FSMModel

	// adjacencyList maps state names to their outgoing transitions
	adjacencyList map[string][]*Transition

	// reverseAdjacencyList maps state names to their incoming transitions
	reverseAdjacencyList map[string][]*Transition

	// reachable tracks which states are reachable from the initial state
	reachable map[string]bool
}

// NewStateGraph creates a new StateGraph from an FSM model
func NewStateGraph(fsm *FSMModel) *StateGraph {
	return &StateGraph{
		FSM:                  fsm,
		adjacencyList:        make(map[string][]*Transition),
		reverseAdjacencyList: make(map[string][]*Transition),
		reachable:            make(map[string]bool),
	}
}

// Build constructs the graph structure from the FSM model
func (g *StateGraph) Build() error {
	// Initialize adjacency lists for all states
	for stateName := range g.FSM.States {
		g.adjacencyList[stateName] = make([]*Transition, 0)
		g.reverseAdjacencyList[stateName] = make([]*Transition, 0)
	}

	// Build adjacency lists from transitions
	for _, transition := range g.FSM.Transitions {
		g.adjacencyList[transition.From] = append(g.adjacencyList[transition.From], transition)
		g.reverseAdjacencyList[transition.To] = append(g.reverseAdjacencyList[transition.To], transition)
	}

	// Compute reachability using DFS
	g.computeReachability()

	return nil
}

// computeReachability computes which states are reachable from the initial state
func (g *StateGraph) computeReachability() {
	visited := make(map[string]bool)
	g.dfs(g.FSM.Initial, visited)
	g.reachable = visited
}

// dfs performs depth-first search to find all reachable states
func (g *StateGraph) dfs(state string, visited map[string]bool) {
	if visited[state] {
		return
	}

	visited[state] = true

	for _, transition := range g.adjacencyList[state] {
		g.dfs(transition.To, visited)
	}
}

// GetOutgoingTransitions returns all transitions leaving the given state
func (g *StateGraph) GetOutgoingTransitions(state string) []*Transition {
	if transitions, exists := g.adjacencyList[state]; exists {
		return transitions
	}
	return []*Transition{}
}

// GetIncomingTransitions returns all transitions entering the given state
func (g *StateGraph) GetIncomingTransitions(state string) []*Transition {
	if transitions, exists := g.reverseAdjacencyList[state]; exists {
		return transitions
	}
	return []*Transition{}
}

// IsReachable returns true if the state is reachable from the initial state
func (g *StateGraph) IsReachable(state string) bool {
	return g.reachable[state]
}

// GetUnreachableStates returns a list of states that are not reachable from the initial state
func (g *StateGraph) GetUnreachableStates() []string {
	unreachable := make([]string, 0)

	for stateName := range g.FSM.States {
		if !g.reachable[stateName] {
			unreachable = append(unreachable, stateName)
		}
	}

	return unreachable
}

// HasCycles returns true if the graph contains cycles
func (g *StateGraph) HasCycles() bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for stateName := range g.FSM.States {
		if !visited[stateName] {
			if g.hasCycleUtil(stateName, visited, recStack) {
				return true
			}
		}
	}

	return false
}

// hasCycleUtil is a utility function for cycle detection using DFS
func (g *StateGraph) hasCycleUtil(state string, visited, recStack map[string]bool) bool {
	visited[state] = true
	recStack[state] = true

	for _, transition := range g.adjacencyList[state] {
		if !visited[transition.To] {
			if g.hasCycleUtil(transition.To, visited, recStack) {
				return true
			}
		} else if recStack[transition.To] {
			return true
		}
	}

	recStack[state] = false
	return false
}
