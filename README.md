# Go-Python Library (Phase 1-2)

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
}
```

## Running the Example

```bash
# Find your libpython3.10.so location first
find /usr -name "libpython3.10.so*" 2>/dev/null

# Run the example with the library path
go run example/main.go /usr/lib/x86_64-linux-gnu/libpython3.10.so.1.0
```

## API Reference

### `NewPureGoPython(libpythonPath string) (*PureGoPython, error)`
Creates a new Python runtime instance by loading the specified libpython library.

### `Initialize() error`
Initializes the Python interpreter. Must be called before any Python operations.

### `Finalize() error` 
Shuts down the Python interpreter and cleans up resources.

### `IsInitialized() bool`
Returns true if the Python interpreter is currently initialized.

### `RunString(code string) error`
Executes Python code from a string. Returns error if execution fails.

### `RunFile(filename string) error`
Executes Python code from a file. Validates file existence and handles errors.

## Next Steps

Phase 3 will add function calling and bidirectional type conversion capabilities.