# Go-Python Library (Phase 1-3)

A Go library that enables calling Python 3.10 functions using the Ebitengine purego interface to bind CPython C API without Cgo.

## Phase 1: Foundation ✅

This phase implements basic Python interpreter lifecycle management:

- ✅ Initialize Go module structure
- ✅ Add purego dependency  
- ✅ Create basic package layout
- ✅ Load libpython using purego.Dlopen() with user-provided path
- ✅ Register essential CPython functions via purego.RegisterLibFunc()
- ✅ Create function pointer wrappers for type safety
- ✅ Implement Py_Initialize() wrapper
- ✅ Implement Py_FinalizeEx() wrapper
- ✅ Add basic error handling for initialization failures
- ✅ Create PureGoPython struct with lifecycle methods

**Deliverable**: Can initialize and cleanup Python interpreter ✅

## Phase 2: Execution Layer ✅

This phase adds Python code execution capabilities:

- ✅ Implement PyRun_SimpleString() wrapper
- ✅ Add Go string to C string conversion
- ✅ Handle execution success/failure return codes
- ✅ Implement PyRun_SimpleFile() wrapper
- ✅ Add file handling and validation
- ✅ Handle file not found and permission errors
- ✅ Capture Python exceptions using PyErr_Occurred()
- ✅ Convert Python errors to Go errors
- ✅ Add error message extraction from Python

**Deliverable**: Can run Python strings and files, get basic error feedback ✅

## Phase 3: Advanced Features ✅

This phase adds function calling and bidirectional type conversion:

- ✅ Implement PyImport_Import() for module loading
- ✅ Add PyObject_GetAttr() for function lookup
- ✅ Implement PyObject_CallObject() wrapper
- ✅ Create argument tuple building
- ✅ Go → Python converters (string, int, float, bool, slices, maps)
- ✅ Python → Go converters (all basic and complex types)
- ✅ Handle complex types (slices ↔ lists, maps ↔ dicts)
- ✅ Add type checking functions (PyUnicode_Check, PyLong_Check, etc.)
- ✅ Implement Py_DECREF() and Py_XDECREF() wrappers
- ✅ Add reference counting to conversion functions
- ✅ Create cleanup mechanisms for Go-created Python objects

**Deliverable**: Can call Python functions with arguments and get typed return values ✅

## Phase 4: Production Ready (Partial) ✅

This phase focuses on making the library production-ready with concurrency safety:

- ✅ Research Python GIL interaction with goroutines
- ✅ Implement thread-safe interpreter access using PyGILState management
- ✅ Add synchronization for multi-goroutine usage with mutex protection
- ✅ Add GIL state management functions (PyGILState_Ensure/Release)
- ✅ Create thread-safe wrappers for Python calls
- ✅ Test concurrent access scenarios
- 🔄 Handle interpreter state isolation (future work)
- 🔄 Performance optimization (future work)
- 🔄 Comprehensive testing suite (future work)

**Deliverable**: Thread-safe Python calls from multiple goroutines ✅

## Thread Safety

The library is now **fully thread-safe** for concurrent access from multiple goroutines:

- **GIL Management**: Proper PyGILState_Ensure/Release for thread safety
- **Mutex Protection**: Go mutex serializes access to Python interpreter
- **Safe Reference Counting**: Protected memory management across threads
- **Concurrent Function Calls**: Multiple goroutines can safely call Python functions

## Usage

```go
package main

import (
    "log"
    "gopython"
)

func main() {
    // Create Python runtime with library path
    py, err := gopython.NewPureGoPython("/usr/lib/x86_64-linux-gnu/libpython3.10.so")
    if err != nil {
        log.Fatal(err)
    }

    // Initialize Python interpreter
    if err := py.Initialize(); err != nil {
        log.Fatal(err)
    }
    defer py.Finalize()

    // Execute Python code from string
    code := `
print("Hello from Python!")
x = 2 + 3
print(f"2 + 3 = {x}")
`
    if err := py.RunString(code); err != nil {
        log.Printf("Error: %v", err)
    }

    // Execute Python code from file
    if err := py.RunFile("script.py"); err != nil {
        log.Printf("Error: %v", err)
    }

    // Call Python functions with type conversion
    result, err := py.CallFunction("math", "sqrt", 16.0)
    if err != nil {
        log.Printf("Error: %v", err)
    }
    fmt.Printf("math.sqrt(16.0) = %v\n", result)

    // Call custom functions with complex types
    data := map[string]interface{}{
        "numbers": []interface{}{1, 2, 3, 4, 5},
        "text": "Hello World",
    }
    result, err = py.CallFunction("mymodule", "process_data", data)
    if err != nil {
        log.Printf("Error: %v", err)
    }
    fmt.Printf("Result: %v\n", result)
}
```

## Running the Examples

```bash
# Find your libpython3.10.so location first
find /usr -name "libpython3.10.so*" 2>/dev/null

# Run the basic example with the library path
go run examples/basic/main.go /usr/lib/x86_64-linux-gnu/libpython3.10.so.1.0

# Run the concurrent safety test
go run examples/concurrent/main.go /usr/lib/x86_64-linux-gnu/libpython3.10.so.1.0

# Run the virtual environment example
go run examples/venv/main.go /usr/lib/x86_64-linux-gnu/libpython3.10.so.1.0 /path/to/your/venv
```

## API Reference

### `NewPureGoPython(libpythonPath string) (*PureGoPython, error)`
Creates a new Python runtime instance by loading the specified libpython library.

### `Initialize() error`
Initializes the Python interpreter. Must be called before any Python operations.

### `Finalize() error` 
Shuts down the Python interpreter and cleans up resources.

### `Initialize() error`
Initializes the Python interpreter with default system configuration.

### `InitializeWithVenv(config VirtualEnvConfig) error`
Initializes the Python interpreter with virtual environment support.

```go
config := gopython.VirtualEnvConfig{
    VenvPath:   "/path/to/venv",     // Virtual environment directory
    SystemSite: true,                // Include system packages
    SitePaths:  []string{},          // Additional package directories
    PythonHome: "",                  // Python installation directory (optional)
}
err := py.InitializeWithVenv(config)
```

### `IsInitialized() bool`
Returns true if the Python interpreter is currently initialized.

### `RunString(code string) error`
Executes Python code from a string. Returns error if execution fails.

### `RunFile(filename string) error`
Executes Python code from a file. Validates file existence and handles errors.

### `CallFunction(module, function string, args ...interface{}) (interface{}, error)`
Calls a Python function with automatic type conversion for arguments and return values.

**Supported Types:**
- **Go → Python**: `string`, `int`, `int64`, `float64`, `bool`, `[]interface{}`, `map[string]interface{}`
- **Python → Go**: `str`, `int`, `float`, `bool`, `list`, `dict`

## Type Conversion Examples

```go
// Basic types
result, _ := py.CallFunction("builtins", "len", "hello")           // string → str
result, _ := py.CallFunction("builtins", "abs", -42)               // int → int  
result, _ := py.CallFunction("math", "sqrt", 16.0)                 // float64 → float
result, _ := py.CallFunction("builtins", "bool", true)             // bool → bool

// Complex types
list := []interface{}{1, 2, 3}
result, _ := py.CallFunction("builtins", "sum", list)              // slice → list

data := map[string]interface{}{"key": "value"}
result, _ := py.CallFunction("json", "dumps", data)                // map → dict
```

## Concurrent Usage Example

```go
// The library is thread-safe - multiple goroutines can safely call Python
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(workerID int) {
        defer wg.Done()
        
        // Safe concurrent access from multiple goroutines
        result, err := py.CallFunction("math", "sqrt", float64(workerID*workerID))
        if err != nil {
            fmt.Printf("Worker %d error: %v\n", workerID, err)
            return
        }
        fmt.Printf("Worker %d: sqrt(%d) = %v\n", workerID, workerID*workerID, result)
    }(i)
}

wg.Wait()
```

## Testing Concurrency

Run the concurrent test suite:

```bash
go run examples/concurrent/main.go /path/to/libpython3.10.so
```

## Limitations

**Important**: Some Python features don't work in embedded environments:

- ❌ **`multiprocessing` module** - Cannot spawn new Python processes
- ⚠️ **`subprocess` operations** - May behave unexpectedly  
- ✅ **`threading` module** - Works fine as alternative
- ✅ **`concurrent.futures.ThreadPoolExecutor`** - Recommended for parallelism
- ✅ **Most libraries** - NumPy, SciPy, Pandas, etc. work perfectly

**Workaround**: Use Go goroutines for parallelism with our thread-safe library!

See [LIMITATIONS.md](LIMITATIONS.md) for detailed information and workarounds.

## Next Steps

Future enhancements could include:
- Sub-interpreter support for true isolation  
- Performance optimizations and benchmarking
- Comprehensive test suite with edge cases
- Production monitoring and metrics