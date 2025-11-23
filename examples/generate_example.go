package main

import (
	"fmt"
	"os"

	"github.com/paveg/gofsm-gen/pkg/model"
)

func main() {
	// Create a sample FSM model
	// Note: In the final implementation, this would be parsed from YAML
	// and used with a code generator. This is a placeholder example.
	fsm, err := model.NewFSMModel("OrderStateMachine", "pending")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating FSM model: %v\n", err)
		os.Exit(1)
	}

	// Add states
	if err := fsm.AddState(&model.State{Name: "pending", EntryAction: "logEntry", ExitAction: "logExit"}); err != nil {
		fmt.Fprintf(os.Stderr, "Error adding state: %v\n", err)
		os.Exit(1)
	}
	fsm.AddState(&model.State{Name: "approved"})
	fsm.AddState(&model.State{Name: "rejected"})
	fsm.AddState(&model.State{Name: "shipped", EntryAction: "notifyCustomer"})

	// Add events
	fsm.AddEvent(&model.Event{Name: "approve"})
	fsm.AddEvent(&model.Event{Name: "reject"})
	fsm.AddEvent(&model.Event{Name: "ship"})

	// Add transitions
	fsm.AddTransition(&model.Transition{
		From:   "pending",
		To:     "approved",
		Event:  "approve",
		Guard:  "hasPayment",
		Action: "chargeCard",
	})
	fsm.AddTransition(&model.Transition{
		From:   "pending",
		To:     "rejected",
		Event:  "reject",
		Action: "sendRejectionEmail",
	})
	fsm.AddTransition(&model.Transition{
		From:   "approved",
		To:     "shipped",
		Event:  "ship",
		Action: "notifyShipping",
	})

	// Validate the FSM
	if err := fsm.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "FSM validation failed: %v\n", err)
		os.Exit(1)
	}

	// Print FSM summary
	fmt.Printf("FSM Model: %s\n", fsm.Name)
	fmt.Printf("Initial State: %s\n", fsm.Initial)
	fmt.Printf("States: %d\n", len(fsm.States))
	fmt.Printf("Events: %d\n", len(fsm.Events))
	fmt.Printf("Transitions: %d\n", len(fsm.Transitions))

	// Note: Code generation will be implemented in Phase 1
	fmt.Println("\nCode generation will be available in Phase 1.")
}
