// Package gopython provides a Go library for embedding Python 3.10 using the purego interface
// to bind CPython C API without Cgo. This allows calling Python functions from Go with
// bidirectional type conversion and thread safety.
//
// Key Features:
// - Pure Go implementation using ebitengine/purego
// - Thread-safe Python function calls from multiple goroutines  
// - Bidirectional type conversion between Go and Python
// - Virtual environment support
// - Comprehensive error handling
// - Memory management with reference counting
//
// Basic Usage:
//   py, err := gopython.NewPureGoPython("/path/to/libpython3.10.so")
//   if err != nil {
//       log.Fatal(err)
//   }
//   
//   if err := py.Initialize(); err != nil {
//       log.Fatal(err)
//   }
//   defer py.Finalize()
//   
//   result, err := py.CallFunction("math", "sqrt", 16.0)
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("sqrt(16) = %v\n", result)
//
// Virtual Environment Usage:
//   config := gopython.VirtualEnvConfig{
//       VenvPath:   "/path/to/venv",
//       SystemSite: true,
//   }
//   if err := py.InitializeWithVenv(config); err != nil {
//       log.Fatal(err)
//   }
//
// Thread Safety:
// All public methods are thread-safe and can be called from multiple goroutines
// concurrently. The library uses Go mutex-based protection rather than Python's
// GIL state management for better reliability in embedded contexts.
//
// Supported Type Conversions:
// Go → Python: string→str, int→int, float64→float, bool→bool, []interface{}→list, map[string]interface{}→dict
// Python → Go: str→string, int→int64, float→float64, bool→bool, list→[]interface{}, dict→map[string]interface{}
package gopython

// This file serves as the main public API interface.
// Implementation details are split across multiple files:
//
// - types.go: Type definitions and structures
// - bindings.go: CPython API function bindings via purego
// - conversion.go: Go ↔ Python type conversion functions  
// - interpreter.go: Python interpreter lifecycle management
// - venv.go: Virtual environment support
// - threading.go: Thread safety wrappers and concurrency utilities
//
// This modular approach improves code organization and maintainability
// while keeping the public API simple and focused.

// Public API Documentation:

// NewPureGoPython creates a new Python runtime instance by loading the specified
// libpython library. The libpythonPath should point to a valid Python 3.10 shared library.
//
// Example:
//   py, err := gopython.NewPureGoPython("/usr/lib/x86_64-linux-gnu/libpython3.10.so")
//
// The function loads the library, registers all CPython API functions, and validates
// that critical functions are available. Returns an error if the library cannot be
// loaded or if required functions are missing.

// Initialize initializes the Python interpreter with default system configuration.
// This must be called before any Python operations can be performed.
//
// Example:
//   if err := py.Initialize(); err != nil {
//       log.Fatal(err)
//   }
//   defer py.Finalize()
//
// For virtual environment support, use InitializeWithVenv instead.

// InitializeWithVenv initializes the Python interpreter with virtual environment
// support. This configures the Python path to prioritize virtual environment
// packages over system packages.
//
// Example:
//   config := gopython.VirtualEnvConfig{
//       VenvPath:   "/path/to/venv",      // Virtual environment directory
//       SystemSite: true,                 // Include system packages as fallback
//       SitePaths:  []string{},           // Additional package directories
//   }
//   if err := py.InitializeWithVenv(config); err != nil {
//       log.Fatal(err)
//   }

// Finalize shuts down the Python interpreter and cleans up resources.
// This should be called when the Python runtime is no longer needed,
// typically using defer after initialization.
//
// The function performs cleanup of Python objects and threads before
// calling Py_FinalizeEx. May take some time to complete if there are
// background threads or network connections to clean up.

// IsInitialized returns true if the Python interpreter is currently initialized
// and ready to execute Python code.

// RunString executes Python code from a string. Returns an error if the
// execution fails or if there are Python syntax/runtime errors.
//
// Example:
//   code := `
//   import math
//   result = math.sqrt(16)
//   print(f"sqrt(16) = {result}")
//   `
//   if err := py.RunString(code); err != nil {
//       log.Printf("Error: %v", err)
//   }

// RunFile executes Python code from a file. The file is validated for
// existence before execution. Returns an error if the file doesn't exist
// or if there are Python execution errors.
//
// Example:
//   if err := py.RunFile("script.py"); err != nil {
//       log.Printf("Error: %v", err)
//   }

// CallFunction calls a Python function with automatic type conversion for
// arguments and return values. The module parameter specifies the Python
// module name (use "__main__" for the main module, or specific module names
// like "math", "json", etc.).
//
// Example:
//   // Call built-in function
//   result, err := py.CallFunction("math", "sqrt", 16.0)
//   
//   // Call custom function defined in __main__
//   result, err := py.CallFunction("__main__", "my_function", "arg1", 42, true)
//   
//   // Call function with complex types
//   data := map[string]interface{}{
//       "numbers": []interface{}{1, 2, 3, 4, 5},
//       "text": "Hello World",
//   }
//   result, err := py.CallFunction("mymodule", "process_data", data)
//
// Supported argument types: string, int, int64, float64, bool, []interface{}, map[string]interface{}
// Supported return types: string, int64, float64, bool, []interface{}, map[string]interface{}, nil
//
// The function is thread-safe and can be called from multiple goroutines concurrently.

// Thread Safety:
// All public methods in this package are thread-safe and use Go mutex-based
// protection. Multiple goroutines can safely call Python functions concurrently
// without additional synchronization required by the caller.
//
// The library serializes all Python operations through a single mutex, which
// ensures thread safety but means that Python operations cannot run truly in
// parallel. For CPU-intensive workloads, consider using multiple Python
// interpreter instances or leveraging Go's concurrency for the parallel work.

// Error Handling:
// All functions return descriptive errors that include context about what
// operation failed. Python exceptions are captured and converted to Go errors
// with the original Python error message preserved.

// Memory Management:
// The library handles Python reference counting automatically. Users do not
// need to manually manage Python object lifetimes. However, it's important
// to call Finalize() when done to ensure proper cleanup of Python resources.

// Limitations:
// - The multiprocessing module doesn't work in embedded Python environments
// - Some subprocess operations may behave unexpectedly
// - Threading module and concurrent.futures.ThreadPoolExecutor work fine
// - Most third-party libraries (NumPy, SciPy, Pandas, etc.) work perfectly
//
// See LIMITATIONS.md for detailed information and workarounds.