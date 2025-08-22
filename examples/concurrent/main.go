package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gopython"
)

func main() {
	// Check if libpython path was provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run example/concurrency.go <path-to-libpython3.10.so>")
	}

	libpythonPath := os.Args[1]
	fmt.Printf("Testing concurrent access with libpython from: %s\n", libpythonPath)

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

	// Create some Python functions for testing
	setupCode := `
import time
import threading

def fibonacci(n):
    """Calculate fibonacci number (slow for testing)"""
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

def factorial(n):
    """Calculate factorial"""
    if n <= 1:
        return 1
    return n * factorial(n-1)

def worker_info(worker_id):
    """Return worker information with thread safety test"""
    import os
    return {
        "worker_id": worker_id,
        "pid": os.getpid(),
        "thread_count": threading.active_count(),
        "time": time.time()
    }

def concurrent_counter(start, increment):
    """Test concurrent operations"""
    result = start
    for i in range(10):
        result += increment
        time.sleep(0.001)  # Small delay to encourage race conditions if not thread-safe
    return result
`

	fmt.Println("Setting up Python test functions...")
	if err := py.RunString(setupCode); err != nil {
		log.Fatalf("Failed to setup Python functions: %v", err)
	}

	fmt.Println("\n=== Testing Concurrent Function Calls ===")

	// Test 1: Concurrent calls to different functions
	fmt.Println("Test 1: Concurrent calls to different functions")
	var wg sync.WaitGroup
	results := make(chan string, 6)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Call worker_info function
			result, err := py.CallFunction("__main__", "worker_info", workerID)
			if err != nil {
				results <- fmt.Sprintf("Worker %d error: %v", workerID, err)
				return
			}
			results <- fmt.Sprintf("Worker %d result: %v", workerID, result)
		}(i)
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Call factorial function
			result, err := py.CallFunction("__main__", "factorial", 5+workerID)
			if err != nil {
				results <- fmt.Sprintf("Factorial worker %d error: %v", workerID, err)
				return
			}
			results <- fmt.Sprintf("Factorial worker %d result: %v", workerID, result)
		}(i + 10)
	}

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Println(result)
	}

	// Test 2: Concurrent calls to the same function with different arguments
	fmt.Println("\nTest 2: Concurrent calls to same function (counter test)")
	var wg2 sync.WaitGroup
	counterResults := make(chan string, 5)

	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go func(workerID int) {
			defer wg2.Done()

			startValue := workerID * 100
			increment := 5

			result, err := py.CallFunction("__main__", "concurrent_counter", startValue, increment)
			if err != nil {
				counterResults <- fmt.Sprintf("Counter worker %d error: %v", workerID, err)
				return
			}

			expected := startValue + (10 * increment) // 10 iterations, each adding increment
			counterResults <- fmt.Sprintf("Counter worker %d: start=%d, result=%v, expected=%d",
				workerID, startValue, result, expected)
		}(i)
	}

	wg2.Wait()
	close(counterResults)

	for result := range counterResults {
		fmt.Println(result)
	}

	// Test 3: High concurrency test
	fmt.Println("\nTest 3: High concurrency test (20 goroutines)")
	var wg3 sync.WaitGroup
	successCount := 0
	errorCount := 0
	var mu sync.Mutex

	start := time.Now()
	for i := 0; i < 20; i++ {
		wg3.Add(1)
		go func(workerID int) {
			defer wg3.Done()

			// Mix of different function calls
			functions := []struct {
				module   string
				function string
				args     []interface{}
			}{
				{"math", "sqrt", []interface{}{float64(workerID + 1)}},
				{"__main__", "factorial", []interface{}{4}},
				{"__main__", "worker_info", []interface{}{workerID}},
			}

			for _, fn := range functions {
				_, err := py.CallFunction(fn.module, fn.function, fn.args...)
				mu.Lock()
				if err != nil {
					errorCount++
				} else {
					successCount++
				}
				mu.Unlock()
			}
		}(i)
	}

	wg3.Wait()
	duration := time.Since(start)

	fmt.Printf("High concurrency test completed in %v\n", duration)
	fmt.Printf("Successful calls: %d, Failed calls: %d\n", successCount, errorCount)
	fmt.Printf("Total calls: %d\n", successCount+errorCount)

	if errorCount == 0 {
		fmt.Println("✅ All concurrent tests passed! Library is thread-safe.")
	} else {
		fmt.Printf("⚠️  Some calls failed - this may indicate thread safety issues or other problems.\n")
	}

	fmt.Println("\n=== Concurrency Safety Test Complete ===")
}
