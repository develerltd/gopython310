package main

import (
	"fmt"
	"log"
	"os"

	"github.com/develerltd/gopython310"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run example_runstring_with_return.go <path-to-libpython3.10.so>")
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

	// Example 1: Define a function that does calculations and returns result
	fmt.Println("=== Example 1: Math calculations ===")
	code1 := `
def calculate():
    import math
    result = math.sqrt(16) + math.pow(2, 3)
    return {"sqrt_16": math.sqrt(16), "pow_2_3": math.pow(2, 3), "sum": result}
`
	if err := py.RunString(code1); err != nil {
		log.Printf("Error: %v", err)
	} else {
		result, err := py.CallFunction("__main__", "calculate")
		if err != nil {
			log.Printf("Error calling function: %v", err)
		} else {
			fmt.Printf("calculate() = %v\n", result)
		}
	}

	// Example 2: Define a function that processes data
	fmt.Println("\n=== Example 2: Data processing ===")
	code2 := `
def process_data():
    data = [1, 2, 3, 4, 5]
    processed = []
    for x in data:
        processed.append(x * x + 1)
    return {
        "original": data,
        "processed": processed,
        "sum_original": sum(data),
        "sum_processed": sum(processed)
    }
`
	if err := py.RunString(code2); err != nil {
		log.Printf("Error: %v", err)
	} else {
		result, err := py.CallFunction("__main__", "process_data")
		if err != nil {
			log.Printf("Error calling function: %v", err)
		} else {
			fmt.Printf("process_data() = %v\n", result)
		}
	}

	// Example 3: Define a function that uses external libraries
	fmt.Println("\n=== Example 3: Using external libraries ===")
	code3 := `
def analyze_text():
    text = "Hello World! This is a test string with multiple words."
    words = text.split()

    analysis = {
        "text": text,
        "word_count": len(words),
        "char_count": len(text),
        "words": words,
        "longest_word": max(words, key=len),
        "uppercase": text.upper(),
        "lowercase": text.lower()
    }
    return analysis
`
	if err := py.RunString(code3); err != nil {
		log.Printf("Error: %v", err)
	} else {
		result, err := py.CallFunction("__main__", "analyze_text")
		if err != nil {
			log.Printf("Error calling function: %v", err)
		} else {
			fmt.Printf("analyze_text() = %v\n", result)
		}
	}

	// Example 4: Using the generic version for type safety
	fmt.Println("\n=== Example 4: Type-safe generic version ===")
	code4 := `
def get_number():
    return 42
`
	if err := py.RunString(code4); err != nil {
		log.Printf("Error: %v", err)
	} else {
		// Using the generic CallPyFunction for type safety
		result, err := gopython.CallPyFunction[interface{}, int64](py, "__main__", "get_number", nil)
		if err != nil {
			log.Printf("Error calling function: %v", err)
		} else {
			fmt.Printf("get_number() = %d (type: %T)\n", result, result)
		}
	}

	fmt.Println("\nAll examples completed!")
}