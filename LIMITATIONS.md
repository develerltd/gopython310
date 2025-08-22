# Library Limitations and Workarounds

## Overview

While this Go-Python library provides robust functionality for most use cases, there are some inherent limitations when embedding Python that users should be aware of.

## Major Limitations

### 1. Python Multiprocessing Module

**Problem**: Python's `multiprocessing` module will not work correctly in this embedded environment.

**Why**: 
- `multiprocessing` expects to spawn new Python interpreter processes
- It tries to execute the current binary (your Go program) instead of `python.exe`
- `sys.argv` is not properly set in embedded environments
- Process spawning mechanisms assume a standard Python installation

**Example of what WON'T work**:
```python
import multiprocessing

def worker(x):
    return x * x

# This will fail in embedded Python!
with multiprocessing.Pool() as pool:
    results = pool.map(worker, [1, 2, 3, 4])
```

**Error you might see**:
```
RuntimeError: An attempt has been made to start a new process before the current process has finished its bootstrapping phase.
```

### 2. Subprocess Module Limitations

**Problem**: Some subprocess operations may behave unexpectedly.

**Why**:
- The embedded environment may not have access to all system executables
- PATH and environment variables might not be fully inherited
- Shell commands may execute in unexpected contexts

### 3. Threading Module Interactions

**Problem**: Python's `threading` module works, but with caveats.

**Considerations**:
- Python threads will still be subject to the GIL
- Our library's mutex serializes all Python operations anyway
- Mixing Python threading with Go goroutines requires careful consideration

## Workarounds and Alternatives

### Instead of Multiprocessing: Use Go Goroutines

**Problem**: Need parallel processing
**Solution**: Use Go's concurrency with our thread-safe library

```go
// Instead of Python multiprocessing, use Go goroutines
var wg sync.WaitGroup
results := make(chan int, 4)

for _, x := range []int{1, 2, 3, 4} {
    wg.Add(1)
    go func(value int) {
        defer wg.Done()
        
        // Call Python function from Go goroutine (thread-safe!)
        result, err := py.CallFunction("__main__", "worker", value)
        if err == nil {
            if intResult, ok := result.(int64); ok {
                results <- int(intResult)
            }
        }
    }(x)
}

wg.Wait()
close(results)

// Collect results
for result := range results {
    fmt.Println("Result:", result)
}
```

### Multiprocessing Alternatives in Python

**Option 1: Use `concurrent.futures` with ThreadPoolExecutor**
```python
from concurrent.futures import ThreadPoolExecutor

def worker(x):
    return x * x

# This works in embedded Python
with ThreadPoolExecutor(max_workers=4) as executor:
    results = list(executor.map(worker, [1, 2, 3, 4]))
```

**Option 2: Use `threading` module directly**
```python
import threading
import queue

def worker(x, result_queue):
    result_queue.put(x * x)

# This works in embedded Python
result_queue = queue.Queue()
threads = []

for x in [1, 2, 3, 4]:
    t = threading.Thread(target=worker, args=(x, result_queue))
    t.start()
    threads.append(t)

for t in threads:
    t.join()

results = []
while not result_queue.empty():
    results.append(result_queue.get())
```

### For CPU-Intensive Tasks: Hybrid Approach

**Pattern**: Use Go for parallelism, Python for computation

```go
// Process data in parallel using Go goroutines
func ProcessDataParallel(py *gopython.PureGoPython, data []int) []int {
    type result struct {
        index int
        value int
    }
    
    resultChan := make(chan result, len(data))
    var wg sync.WaitGroup
    
    // Process each item in a separate goroutine
    for i, item := range data {
        wg.Add(1)
        go func(index, value int) {
            defer wg.Done()
            
            // Call expensive Python computation
            pyResult, err := py.CallFunction("math_module", "complex_calculation", value)
            if err == nil {
                if intResult, ok := pyResult.(int64); ok {
                    resultChan <- result{index: index, value: int(intResult)}
                }
            }
        }(i, item)
    }
    
    wg.Wait()
    close(resultChan)
    
    // Collect results in order
    results := make([]int, len(data))
    for res := range resultChan {
        results[res.index] = res.value
    }
    
    return results
}
```

## Best Practices

### 1. Design for Single-Process Parallelism
- Use Go goroutines for parallelism
- Use Python for domain-specific computation
- Leverage our library's thread safety

### 2. Avoid Process-Based Python Modules
- Don't use `multiprocessing`
- Be cautious with `subprocess` 
- Prefer `threading` or `concurrent.futures.ThreadPoolExecutor`

### 3. Test Carefully
- Test any new Python libraries in embedded environment
- Some libraries may have hidden multiprocessing dependencies
- Consider fallback strategies for unsupported operations

### 4. Communication Patterns
```go
// Good: Go orchestrates, Python computes
func ProcessWorkflow(py *gopython.PureGoPython) {
    // Use Go for orchestration
    for i := 0; i < 100; i++ {
        go func(taskID int) {
            // Call Python for domain-specific work
            result, _ := py.CallFunction("analytics", "process_task", taskID)
            // Use Go for result handling
            saveResult(result)
        }(i)
    }
}
```

## Testing for Compatibility

When evaluating Python libraries for use with this embedded environment, test:

```python
# Test script to check library compatibility
import sys
print("Testing library in embedded environment...")

try:
    import your_library
    print("✓ Import successful")
    
    # Test basic functionality
    result = your_library.some_function()
    print("✓ Basic functionality works")
    
    # Test for multiprocessing usage (will fail in embedded env)
    import multiprocessing
    if hasattr(your_library, 'parallel_function'):
        print("⚠ Library may use multiprocessing - test carefully")
    
except Exception as e:
    print(f"✗ Error: {e}")
```

## Summary

While these limitations exist, they don't significantly impact most use cases. The library excels at:

- ✅ Mathematical computations
- ✅ Data processing and analysis  
- ✅ Machine learning inference
- ✅ Scientific computing
- ✅ String processing and manipulation
- ✅ File I/O and parsing
- ✅ Most third-party libraries (NumPy, SciPy, Pandas, etc.)

The key is understanding these constraints and designing your architecture accordingly, leveraging Go's excellent concurrency for parallelism while using Python for domain-specific computation.