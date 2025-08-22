# Go-Python Library Implementation Plan

## Overview
Create a Go library that enables calling Python 3.10 functions using the Ebitengine purego interface to bind CPython C API without Cgo.

## Research Findings

### Purego Capabilities
- Cross-platform C function calling without Cgo
- Dynamic library loading via `Dlopen()`
- Symbol registration with `RegisterLibFunc()` and `RegisterFunc()`
- Type conversion between Go and C types
- Supports Linux, macOS, Windows, FreeBSD on amd64/arm64

### CPython C API (Python 3.10)
- **Initialization**: `Py_Initialize()`, `Py_SetProgramName()`
- **Execution**: `PyRun_SimpleString()`, `PyRun_SimpleFile()`
- **Objects**: `PyObject_CallObject()`, `Py_DECREF()`, `Py_XDECREF()`
- **Cleanup**: `Py_FinalizeEx()`

## Architecture Design

### Core Components

1. **Python Runtime Manager**
   - Initialize/finalize Python interpreter
   - Manage interpreter state
   - Handle multiple interpreter instances

2. **Function Call Interface**
   - Execute Python strings/files
   - Call Python functions with arguments
   - Handle return values and exceptions

3. **Type Conversion Layer**
   - Go ↔ Python type mapping
   - Handle complex types (lists, dicts, objects)
   - Memory management for cross-language data

4. **Library Loader**
   - Dynamically load libpython3.10.so
   - Register CPython API functions via purego
   - Platform-specific library detection

### Proposed Interface

```go
// Core interface
type PyRuntime interface {
    Initialize() error
    Finalize() error
    RunString(code string) (interface{}, error)
    RunFile(filename string) (interface{}, error)
    CallFunction(module, function string, args ...interface{}) (interface{}, error)
}

// Implementation using purego
type PureGoPython struct {
    libHandle uintptr
    // CPython API function pointers
    pyInitialize func()
    pyFinalize func() int
    pyRunSimpleString func(string) int
    // ... other API functions
}
```

## Implementation Plan

### Phase 1: Foundation
**Goal**: Establish basic Python interpreter lifecycle management

- [ ] **Project Setup**
  - [ ] Initialize Go module structure
  - [ ] Add purego dependency
  - [ ] Create basic package layout

- [ ] **Library Detection**
  - [ ] Implement platform-specific libpython3.10.so detection
  - [ ] Handle common Python installation paths
  - [ ] Add fallback search mechanisms

- [ ] **Core API Binding**
  - [ ] Load libpython3.10.so using purego.Dlopen()
  - [ ] Register essential CPython functions via purego.RegisterLibFunc()
  - [ ] Create function pointer wrappers for type safety

- [ ] **Basic Runtime**
  - [ ] Implement Py_Initialize() wrapper
  - [ ] Implement Py_FinalizeEx() wrapper
  - [ ] Add basic error handling for initialization failures
  - [ ] Create PureGoPython struct with lifecycle methods

**Deliverable**: Can initialize and cleanup Python interpreter

### Phase 2: Execution Layer
**Goal**: Execute Python code from Go

- [ ] **String Execution**
  - [ ] Implement PyRun_SimpleString() wrapper
  - [ ] Add Go string to C string conversion
  - [ ] Handle execution success/failure return codes

- [ ] **File Execution**
  - [ ] Implement PyRun_SimpleFile() wrapper
  - [ ] Add file handling and validation
  - [ ] Handle file not found and permission errors

- [ ] **Error Handling**
  - [ ] Capture Python exceptions using PyErr_Occurred()
  - [ ] Convert Python errors to Go errors
  - [ ] Add error message extraction from Python

**Deliverable**: Can run Python strings and files, get basic error feedback

### Phase 3: Advanced Features
**Goal**: Function calling and data exchange

- [ ] **Function Calls**
  - [ ] Implement PyImport_Import() for module loading
  - [ ] Add PyObject_GetAttrString() for function lookup
  - [ ] Implement PyObject_CallObject() wrapper
  - [ ] Create argument tuple building

- [ ] **Type Conversion**
  - [ ] Go → Python converters (string, int, float, bool)
  - [ ] Python → Go converters (basic types)
  - [ ] Handle complex types (slices → lists, maps → dicts)
  - [ ] Add type checking functions (PyUnicode_Check, PyLong_Check, etc.)

- [ ] **Memory Management**
  - [ ] Implement Py_DECREF() and Py_XDECREF() wrappers
  - [ ] Add reference counting to conversion functions
  - [ ] Create cleanup mechanisms for Go-created Python objects

**Deliverable**: Can call Python functions with arguments and get typed return values

### Phase 4: Production Ready
**Goal**: Robust, safe, and performant library

- [ ] **Concurrent Safety**
  - [ ] Research Python GIL interaction with goroutines
  - [ ] Implement thread-safe interpreter access
  - [ ] Add synchronization for multi-goroutine usage
  - [ ] Handle interpreter state isolation

- [ ] **Performance Optimization**
  - [ ] Profile function call overhead
  - [ ] Optimize type conversion paths
  - [ ] Cache frequently used Python objects
  - [ ] Minimize memory allocations

- [ ] **Testing & Documentation**
  - [ ] Unit tests for all phases
  - [ ] Integration tests with real Python scripts
  - [ ] Benchmarks comparing to Cgo alternatives
  - [ ] API documentation and examples
  - [ ] Platform compatibility testing

**Deliverable**: Production-ready library with full feature set

## Key Technical Challenges

1. **Library Location**: Platform-specific paths for libpython3.10.so
2. **GIL Management**: Python's Global Interpreter Lock with Go goroutines
3. **Memory Safety**: Reference counting across language boundaries
4. **Type Marshaling**: Complex Python objects to Go types
5. **Error Propagation**: Python exceptions to Go errors

## Benefits of Purego Approach

- **Cross-compilation**: Build on any platform targeting any platform
- **No Cgo**: Faster builds, smaller binaries
- **Dynamic Loading**: Runtime Python version detection
- **Pure Go**: Maintains Go toolchain benefits

## Next Steps

1. Set up Go module structure
2. Implement library detection logic
3. Create basic purego bindings for core CPython functions
4. Build minimal runtime initialization/cleanup
5. Add string execution capability
6. Expand with advanced features iteratively