package main

import (
	"fmt"
	"os"
)

const version = "0.1.0-dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("gofsm-gen version %s\n", version)
		os.Exit(0)
	}

	fmt.Fprintln(os.Stderr, "gofsm-gen: FSM code generator for Go")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage: gofsm-gen [options]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "  -spec=FILE       FSM specification file (YAML)")
	fmt.Fprintln(os.Stderr, "  -out=FILE        Output file path")
	fmt.Fprintln(os.Stderr, "  -package=NAME    Package name for generated code")
	fmt.Fprintln(os.Stderr, "  --version        Show version information")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "This is a minimal placeholder implementation.")
	fmt.Fprintln(os.Stderr, "Full functionality will be added in upcoming phases.")
	os.Exit(1)
}
