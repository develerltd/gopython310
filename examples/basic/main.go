package main

import (
	"fmt"
	"log"
	"os"

	"github.com/develerltd/gopython310"
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
	if err := py.RunFile("examples/basic/test.py"); err != nil {
		fmt.Printf("File execution error: %v\n", err)
	}

	// Phase 3: Function calling and type conversion
	fmt.Println("\n=== Phase 3: Function Calling ===")

	// Create a test module by running a string first
	fmt.Println("Creating test module...")
	moduleCode := `
def add_numbers(a, b):
    return a + b

def process_list(items):
    return [x * 2 for x in items]

def get_info():
    return {
        "name": "Python Module", 
        "version": "1.0",
        "features": ["functions", "lists", "dicts"]
    }

def greet(name):
    return f"Hello, {name}!"
`
	if err := py.RunString(moduleCode); err != nil {
		log.Printf("Error creating module: %v", err)
	}

	// Test function calling with different types
	fmt.Println("\nTesting function calls...")

	// Test simple function with numbers
	result, err := py.CallFunction("__main__", "add_numbers", 10, 25)
	if err != nil {
		fmt.Printf("Error calling add_numbers: %v\n", err)
	} else {
		fmt.Printf("add_numbers(10, 25) = %v (type: %T)\n", result, result)
	}

	// Test function with string argument
	result, err = py.CallFunction("__main__", "greet", "Go Developer")
	if err != nil {
		fmt.Printf("Error calling greet: %v\n", err)
	} else {
		fmt.Printf("greet(\"Go Developer\") = %v\n", result)
	}

	// Test function with list argument and return
	list := []interface{}{1, 2, 3, 4, 5}
	result, err = py.CallFunction("__main__", "process_list", list)
	if err != nil {
		fmt.Printf("Error calling process_list: %v\n", err)
	} else {
		fmt.Printf("process_list(%v) = %v\n", list, result)
	}

	// Test function returning dictionary
	result, err = py.CallFunction("__main__", "get_info")
	if err != nil {
		fmt.Printf("Error calling get_info: %v\n", err)
	} else {
		fmt.Printf("get_info() = %v (type: %T)\n", result, result)
	}

	// Test calling built-in modules
	fmt.Println("\nTesting built-in module calls...")
	result, err = py.CallFunction("math", "sqrt", 16.0)
	if err != nil {
		fmt.Printf("Error calling math.sqrt: %v\n", err)
	} else {
		fmt.Printf("math.sqrt(16.0) = %v\n", result)
	}

	fmt.Println("\nPhase 3 implementation complete!")

	// Test limitations and compatibility
	fmt.Println("\n=== Testing Limitations and Compatibility ===")
	fmt.Println("Running compatibility tests...")
	if err := py.RunFile("examples/basic/test_limitations.py"); err != nil {
		fmt.Printf("Limitations test error: %v\n", err)
	}
}