package main

import (
	"fmt"
	"log"
	"os"

	"gopython"
)

func main() {
	// Check if libpython path was provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run example/main.go <path-to-libpython3.10.so>")
	}

	libpythonPath := os.Args[1]
	fmt.Printf("Attempting to load libpython from: %s\n", libpythonPath)

	// Create a new Python runtime instance
	py, err := gopython.NewPureGoPython(libpythonPath)
	if err != nil {
		log.Fatalf("Failed to create Python runtime: %v", err)
	}

	// Initialize the Python interpreter
	fmt.Println("Initializing Python interpreter...")
	if err := py.Initialize(); err != nil {
		log.Fatalf("Failed to initialize Python: %v", err)
	}
	defer py.Finalize()

	fmt.Printf("Python interpreter initialized: %v\n", py.IsInitialized())

	// Phase 2: Execute Python strings
	fmt.Println("\n=== Phase 2: Execution Layer ===")

	// Test string execution
	fmt.Println("Executing Python string...")
	code := `
print("Hello from Python!")
x = 2 + 3
print(f"2 + 3 = {x}")
import sys
print(f"Python version: {sys.version}")
`
	if err := py.RunString(code); err != nil {
		log.Printf("Error executing string: %v", err)
	} else {
		fmt.Println("String execution successful!")
	}

	// Test error handling
	fmt.Println("\nTesting error handling...")
	errorCode := `
print("This will work")
undefined_variable  # This will cause an error
print("This won't execute")
`
	if err := py.RunString(errorCode); err != nil {
		fmt.Printf("Expected error caught: %v\n", err)
	}

	// Test file execution
	fmt.Println("\nTesting file execution...")
	if err := py.RunFile("example/test.py"); err != nil {
		fmt.Printf("File execution error: %v\n", err)
	}

	fmt.Println("\nPhase 2 implementation complete!")
}