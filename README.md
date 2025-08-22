# Go-Python Library

A Go library that enables calling Python 3.10 functions using the Ebitengine purego interface to bind CPython C API without Cgo.

## Features

- **No Cgo Required**: Uses purego to bind CPython C API directly
- **Thread-Safe**: Full concurrency support with proper GIL management
- **Type Conversion**: Automatic Go ↔ Python type conversion for basic and complex types
- **Virtual Environment Support**: Works with Python virtual environments
- **Error Handling**: Proper Python exception handling and Go error conversion

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

### Linux
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

### macOS
```bash
# Homebrew Python
go run examples/basic/main.go /opt/homebrew/lib/libpython3.10.dylib

# Or with pyenv
go run examples/basic/main.go ~/.pyenv/versions/3.10.15/lib/libpython3.10.dylib

# Virtual environment example (macOS)
go run examples/venv/main.go /opt/homebrew/lib/libpython3.10.dylib /path/to/your/venv
```


### Finding Python Libraries

**Linux:**
```bash
# Ubuntu/Debian
find /usr -name "libpython3.10.so*" 2>/dev/null

# Or check with pkg-config
pkg-config --libs python3.10
```

**macOS:**
```bash
# Homebrew
ls /opt/homebrew/lib/libpython3.10.dylib
ls /usr/local/lib/libpython3.10.dylib

# pyenv
ls ~/.pyenv/versions/*/lib/libpython3.10.dylib

# System Python (if available)
ls /Library/Frameworks/Python.framework/Versions/3.10/lib/libpython3.10.dylib
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

## Thread Safety

The library is fully thread-safe for concurrent access from multiple goroutines:

- **GIL Management**: Proper PyGILState_Ensure/Release for thread safety
- **Mutex Protection**: Go mutex serializes access to Python interpreter
- **Safe Reference Counting**: Protected memory management across threads
- **Concurrent Function Calls**: Multiple goroutines can safely call Python functions