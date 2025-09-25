package gopython

import (
	"errors"
	"fmt"
	"os"

	"github.com/ebitengine/purego"
)

// NewPureGoPython creates a new Python runtime instance
func NewPureGoPython(libpythonPath string) (*PureGoPython, error) {
	// Validate library path for current platform
	if err := ValidateLibraryPath(libpythonPath); err != nil {
		return nil, fmt.Errorf("invalid library path: %v", err)
	}

	// Load the Python library
	libHandle, err := purego.Dlopen(libpythonPath, purego.RTLD_NOW)
	if err != nil {
		return nil, fmt.Errorf("failed to load libpython from %s: %v", libpythonPath, err)
	}

	py := &PureGoPython{
		libHandle: libHandle,
	}

	// Register all Python functions
	if err := py.registerPythonFunctions(); err != nil {
		return nil, fmt.Errorf("failed to register Python functions: %v", err)
	}

	// Validate that critical functions are registered
	if err := py.validateFunctionRegistration(); err != nil {
		return nil, fmt.Errorf("function registration validation failed: %v", err)
	}

	return py, nil
}


// Initialize initializes the Python interpreter with default system configuration
func (py *PureGoPython) Initialize() error {
	if py.pyInitialize == nil {
		return errors.New("Python functions not registered")
	}

	py.pyInitialize()
	return nil
}

// IsInitialized returns true if the Python interpreter is initialized
func (py *PureGoPython) IsInitialized() bool {
	if py.pyIsInitialized == nil {
		return false
	}
	return py.pyIsInitialized() != 0
}

// Finalize shuts down the Python interpreter
func (py *PureGoPython) Finalize() error {
	if py.pyFinalizeEx == nil {
		return errors.New("Python functions not registered")
	}

	if !py.IsInitialized() {
		return errors.New("Python interpreter is not initialized")
	}

	// Try to clean up any remaining Python objects and threads
	py.withGIL(func() error {
		cleanupCode := `
import gc
import threading
import sys

# Force garbage collection
gc.collect()

# Try to join any remaining threads (except main thread)
main_thread = threading.main_thread()
for thread in threading.enumerate():
    if thread != main_thread and thread.is_alive():
        try:
            if hasattr(thread, 'join'):
                thread.join(timeout=0.1)  # Short timeout
        except:
            pass  # Ignore join errors

# Clear any remaining modules and references
if hasattr(sys, 'modules'):
    modules_to_clear = [name for name in sys.modules.keys() 
                       if not name.startswith('__') and name not in ['sys', 'builtins']]
    for name in modules_to_clear:
        try:
            del sys.modules[name]
        except:
            pass

# Final garbage collection
gc.collect()
`
		cCode := stringToCString(cleanupCode)
		py.pyRunSimpleString(cCode)
		return nil
	})

	result := py.pyFinalizeEx()
	if result < 0 {
		return fmt.Errorf("Python interpreter finalization failed with code: %d", result)
	}

	return nil
}

// RunString executes Python code from a string
func (py *PureGoPython) RunString(code string) error {
	if !py.IsInitialized() {
		return errors.New("Python interpreter is not initialized")
	}

	return py.withGIL(func() error {
		cCode := stringToCString(code)
		result := py.pyRunSimpleString(cCode)
		if result != 0 {
			return py.getPythonError()
		}
		return nil
	})
}

// RunFile executes Python code from a file
func (py *PureGoPython) RunFile(filename string) error {
	if !py.IsInitialized() {
		return errors.New("Python interpreter is not initialized")
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	// Read file content and execute as string (simpler and more reliable)
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	return py.RunString(string(content))
}

// CallFunction calls a Python function with the given arguments
func (py *PureGoPython) CallFunction(module, function string, args ...interface{}) (interface{}, error) {
	if !py.IsInitialized() {
		return nil, errors.New("Python interpreter is not initialized")
	}

	return py.withGILReturn(func() (interface{}, error) {
		return py.callFunctionUnsafe(module, function, args...)
	})
}

// callFunctionUnsafe performs the actual function call without GIL management
func (py *PureGoPython) callFunctionUnsafe(module, function string, args ...interface{}) (interface{}, error) {
	// Import the module
	moduleNameObj, err := py.goToPython(module)
	if err != nil {
		return nil, fmt.Errorf("failed to convert module name: %v", err)
	}
	defer py.safeDecRef(uintptr(moduleNameObj))

	moduleObj := py.pyImportImport(uintptr(moduleNameObj))
	if moduleObj == 0 {
		return nil, fmt.Errorf("failed to import module '%s': %v", module, py.getPythonError())
	}
	defer py.safeDecRef(moduleObj)

	// Get the function from the module
	functionNameObj, err := py.goToPython(function)
	if err != nil {
		return nil, fmt.Errorf("failed to convert function name: %v", err)
	}
	defer py.safeDecRef(uintptr(functionNameObj))

	functionObj := py.pyObjectGetAttr(moduleObj, uintptr(functionNameObj))
	if functionObj == 0 {
		return nil, fmt.Errorf("function '%s' not found in module '%s'", function, module)
	}
	defer py.safeDecRef(functionObj)

	// Build argument tuple
	argTuple, err := py.buildArgumentTuple(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to build arguments: %v", err)
	}
	defer py.safeDecRef(uintptr(argTuple))

	// Call the function
	resultObj := py.pyObjectCallObject(functionObj, uintptr(argTuple))
	if resultObj == 0 {
		return nil, fmt.Errorf("function call failed: %v", py.getPythonError())
	}
	defer py.safeDecRef(resultObj)

	// Convert result to Go
	return py.pythonToGo(PyObject(resultObj))
}

// CallPyFunction calls a Python function with type-safe generics for request and response types
func CallPyFunction[TRequest, TResponse any](py *PureGoPython, module, function string, request TRequest) (TResponse, error) {
	var zero TResponse

	if !py.IsInitialized() {
		return zero, errors.New("Python interpreter is not initialized")
	}

	// Call the underlying CallFunction with the request
	result, err := py.CallFunction(module, function, request)
	if err != nil {
		return zero, err
	}

	// Try to convert the result to the expected response type
	response, ok := result.(TResponse)
	if !ok {
		return zero, fmt.Errorf("failed to convert result to %T: got %T", zero, result)
	}

	return response, nil
}

// getPythonError extracts Python error information
func (py *PureGoPython) getPythonError() error {
	if py.pyErrOccurred() == 0 {
		return errors.New("unknown Python error")
	}

	var ptype, pvalue, ptraceback uintptr
	py.pyErrFetch(&ptype, &pvalue, &ptraceback)

	// Clear the error state
	py.pyErrClear()

	if pvalue == 0 {
		return errors.New("Python error occurred but no error message available")
	}

	// Convert error to string
	errorStr := py.pyObjectStr(pvalue)
	if errorStr == 0 {
		py.safeDecRef(ptype)
		py.safeDecRef(pvalue)
		py.safeDecRef(ptraceback)
		return errors.New("Python error occurred but failed to get error string")
	}

	cStr := py.pyUnicodeAsUTF8(errorStr)
	errorMessage := "Python error"
	if cStr != nil {
		errorMessage = cStringToGoString(cStr)
	}

	// Clean up error objects
	py.safeDecRef(ptype)
	py.safeDecRef(pvalue)
	py.safeDecRef(ptraceback)
	py.safeDecRef(errorStr)

	return fmt.Errorf("Python error: %s", errorMessage)
}