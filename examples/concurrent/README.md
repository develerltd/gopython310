# Concurrent Example

This example demonstrates the thread safety and concurrency features of the Go-Python library (Phase 4).

## Features Demonstrated

- **Thread Safety**: Multiple goroutines calling Python functions safely
- **Concurrent Access**: Testing different functions from multiple goroutines  
- **High Load**: Stress testing with many concurrent workers
- **Error Handling**: Tracking success/failure rates under load

## Running

```bash
go run main.go /path/to/libpython3.10.so
```

## Test Scenarios

### Test 1: Concurrent Different Functions
- 3 goroutines calling `worker_info()` 
- 3 goroutines calling `factorial()`
- Validates concurrent access to different functions

### Test 2: Concurrent Same Function  
- 5 goroutines calling `concurrent_counter()` with different arguments
- Tests thread safety when multiple goroutines call the same function
- Validates that results are correct (no race conditions)

### Test 3: High Concurrency
- 20 goroutines making multiple function calls each
- Mix of built-in (`math.sqrt`) and custom functions
- Reports success/failure rates and timing

## Expected Output

```
Test 1: Concurrent calls to different functions
Worker 0 result: {worker_id: 0, pid: 12345, thread_count: 1, time: 1234567890.123}
...

Test 2: Concurrent calls to same function (counter test)  
Counter worker 0: start=0, result=50, expected=50
...

Test 3: High concurrency test (20 goroutines)
High concurrency test completed in 123ms
Successful calls: 60, Failed calls: 0
Total calls: 60
✅ All concurrent tests passed! Library is thread-safe.
```

## What This Proves

- **No Race Conditions**: All function calls complete successfully
- **Correct Results**: Return values match expected calculations
- **Memory Safety**: No crashes or corruption under load
- **Thread Safety**: Multiple goroutines can safely use the library

This demonstrates that the library is production-ready for concurrent applications.