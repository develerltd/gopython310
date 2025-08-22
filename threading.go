package gopython

// withGIL executes a function with GIL protection (thread-safe)
func (py *PureGoPython) withGIL(fn func() error) error {
	py.mu.Lock()
	defer py.mu.Unlock()
	return fn()
}

// withGILReturn executes a function with GIL protection and returns a value (thread-safe)
func (py *PureGoPython) withGILReturn(fn func() (interface{}, error)) (interface{}, error) {
	py.mu.Lock()
	defer py.mu.Unlock()
	return fn()
}

// Thread-safe wrapper functions for public API

// RunStringThreadSafe executes Python code from a string (thread-safe)
func (py *PureGoPython) RunStringThreadSafe(code string) error {
	return py.RunString(code) // Already thread-safe internally
}

// RunFileThreadSafe executes Python code from a file (thread-safe)
func (py *PureGoPython) RunFileThreadSafe(filename string) error {
	return py.RunFile(filename) // Already thread-safe internally
}

// CallFunctionThreadSafe calls a Python function with arguments (thread-safe)
func (py *PureGoPython) CallFunctionThreadSafe(module, function string, args ...interface{}) (interface{}, error) {
	return py.CallFunction(module, function, args...) // Already thread-safe internally
}

// IsInitializedThreadSafe checks if Python interpreter is initialized (thread-safe)
func (py *PureGoPython) IsInitializedThreadSafe() bool {
	py.mu.Lock()
	defer py.mu.Unlock()
	return py.IsInitialized()
}

// FinalizeThreadSafe shuts down the Python interpreter (thread-safe)
func (py *PureGoPython) FinalizeThreadSafe() error {
	py.mu.Lock()
	defer py.mu.Unlock()
	return py.Finalize()
}

// Note: The library uses Go mutex-based thread safety instead of Python's GIL state management
// This approach was chosen because:
// 1. PyGILState_Ensure/Release caused fatal errors in embedded Python
// 2. Go mutex provides simpler and more reliable thread safety
// 3. All Python operations are serialized through the mutex, preventing race conditions
// 4. This is compatible with Python's threading model when called from embedded contexts

// Future enhancement: If true parallel Python execution is needed, consider:
// - Multiple sub-interpreters (PyInterpreterState)
// - Per-thread Python interpreter instances
// - Advanced GIL management patterns

// Thread Safety Architecture:
// ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
// │   Goroutine 1   │    │   Goroutine 2   │    │   Goroutine N   │
// └─────────────────┘    └─────────────────┘    └─────────────────┘
//          │                       │                       │
//          ▼                       ▼                       ▼
// ┌─────────────────────────────────────────────────────────────────┐
// │                         Go Mutex Lock                          │
// └─────────────────────────────────────────────────────────────────┘
//                                  │
//                                  ▼
// ┌─────────────────────────────────────────────────────────────────┐
// │                   Single Python Interpreter                    │
// │                    (Thread-Safe Access)                       │
// └─────────────────────────────────────────────────────────────────┘