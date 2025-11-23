package main

import (
	"fmt"
	"os"

	"github.com/yourusername/gofsm-gen/pkg/generator"
	"github.com/yourusername/gofsm-gen/pkg/model"
)

func main() {
	// Create a sample FSM model (in real usage, this would be parsed from YAML)
	fsm := &model.FSMModel{
		Name:    "OrderStateMachine",
		Package: "orders",
		Initial: "pending",
		States: []model.State{
			{Name: "pending", Entry: "logEntry", Exit: "logExit"},
			{Name: "approved", Entry: "", Exit: ""},
			{Name: "rejected", Entry: "", Exit: ""},
			{Name: "shipped", Entry: "notifyCustomer", Exit: ""},
		},
		Events: []model.Event{
			{Name: "approve"},
			{Name: "reject"},
			{Name: "ship"},
		},
		Transitions: []model.Transition{
			{From: "pending", To: "approved", On: "approve", Guard: "hasPayment", Action: "chargeCard"},
			{From: "pending", To: "rejected", On: "reject", Guard: "", Action: "sendRejectionEmail"},
			{From: "approved", To: "shipped", On: "ship", Guard: "", Action: "notifyShipping"},
		},
	}

	// Create code generator
	gen, err := generator.NewCodeGenerator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating generator: %v\n", err)
		os.Exit(1)
	}

	// Generate code
	code, err := gen.Generate(fsm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	// Print generated code
	fmt.Println(string(code))
}
