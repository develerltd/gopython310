package main

import (
	"fmt"
	"log"
	"os"

	"github.com/develerltd/gopython310"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run test_generics.go <path-to-libpython3.10.so>")
	}

	libpythonPath := os.Args[1]

	// Create and initialize Python runtime
	py, err := gopython.NewPureGoPython(libpythonPath)
	if err != nil {
		log.Fatalf("Failed to create Python runtime: %v", err)
	}

	if err := py.Initialize(); err != nil {
		log.Fatalf("Failed to initialize Python: %v", err)
	}
	defer py.Finalize()

	// Define test functions in Python
	code := `
def add_numbers(data):
    return data['a'] + data['b']

def get_greeting(name):
    return f"Hello, {name}!"

def multiply_list(numbers):
    return [x * 2 for x in numbers]
`
	if err := py.RunString(code); err != nil {
		log.Fatalf("Error defining Python functions: %v", err)
	}

	// Test 1: Function that takes a map and returns a number
	fmt.Println("=== Test 1: Map input, int64 output ===")
	input1 := map[string]interface{}{"a": 10, "b": 25}
	result1, err := gopython.CallPyFunction[map[string]interface{}, int64](py, "__main__", "add_numbers", input1)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("add_numbers(%v) = %d\n", input1, result1)
	}

	// Test 2: Function that takes a string and returns a string
	fmt.Println("\n=== Test 2: String input, string output ===")
	input2 := "Go Developer"
	result2, err := gopython.CallPyFunction[string, string](py, "__main__", "get_greeting", input2)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("get_greeting(%q) = %q\n", input2, result2)
	}

	// Test 3: Function that takes a slice and returns a slice
	fmt.Println("\n=== Test 3: Slice input, slice output ===")
	input3 := []interface{}{1, 2, 3, 4, 5}
	result3, err := gopython.CallPyFunction[[]interface{}, []interface{}](py, "__main__", "multiply_list", input3)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("multiply_list(%v) = %v\n", input3, result3)
	}

	// Test 4: Built-in function with float
	fmt.Println("\n=== Test 4: Built-in math function ===")
	input4 := 16.0
	result4, err := gopython.CallPyFunction[float64, float64](py, "math", "sqrt", input4)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("math.sqrt(%f) = %f\n", input4, result4)
	}

	fmt.Println("\nAll tests completed!")
}