// Package gopython provides a Go interface to Python 3.10 using purego
package gopython

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ebitengine/purego"
)

// PureGoPython represents a Python runtime instance
type PureGoPython struct {
	libHandle uintptr

	// CPython API function pointers
	pyInitialize      func()
	pyFinalizeEx      func() int
	pyIsInitialized   func() int
	pyRunSimpleString func(*byte) int
	pyErrOccurred     func() uintptr
	pyErrPrint        func()
	pyErrClear        func()
}

// PyRuntime interface defines the core Python runtime operations
type PyRuntime interface {
	Initialize() error
	Finalize() error
	IsInitialized() bool
	RunString(code string) error
	RunFile(filename string) error
}

// NewPureGoPython creates a new PureGoPython instance with the given library path
func NewPureGoPython(libpythonPath string) (*PureGoPython, error) {
	if libpythonPath == "" {
		return nil, errors.New("libpython path cannot be empty")
	}

	py := &PureGoPython{}

	// Load the Python library
	lib, err := purego.Dlopen(libpythonPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, fmt.Errorf("failed to load libpython from %s: %v", libpythonPath, err)
	}
	py.libHandle = lib

	// Register CPython API functions
	if err := py.registerFunctions(); err != nil {
		return nil, fmt.Errorf("failed to register Python functions: %v", err)
	}

	return py, nil
}

// registerFunctions binds CPython API functions using purego
func (py *PureGoPython) registerFunctions() error {
	// Register Py_Initialize
	purego.RegisterLibFunc(&py.pyInitialize, py.libHandle, "Py_Initialize")
	if py.pyInitialize == nil {
		return errors.New("failed to register Py_Initialize")
	}

	// Register Py_FinalizeEx
	purego.RegisterLibFunc(&py.pyFinalizeEx, py.libHandle, "Py_FinalizeEx")
	if py.pyFinalizeEx == nil {
		return errors.New("failed to register Py_FinalizeEx")
	}

	// Register Py_IsInitialized
	purego.RegisterLibFunc(&py.pyIsInitialized, py.libHandle, "Py_IsInitialized")
	if py.pyIsInitialized == nil {
		return errors.New("failed to register Py_IsInitialized")
	}

	// Register PyRun_SimpleString
	purego.RegisterLibFunc(&py.pyRunSimpleString, py.libHandle, "PyRun_SimpleString")
	if py.pyRunSimpleString == nil {
		return errors.New("failed to register PyRun_SimpleString")
	}

	// Register PyErr_Occurred
	purego.RegisterLibFunc(&py.pyErrOccurred, py.libHandle, "PyErr_Occurred")
	if py.pyErrOccurred == nil {
		return errors.New("failed to register PyErr_Occurred")
	}

	// Register PyErr_Print
	purego.RegisterLibFunc(&py.pyErrPrint, py.libHandle, "PyErr_Print")
	if py.pyErrPrint == nil {
		return errors.New("failed to register PyErr_Print")
	}

	// Register PyErr_Clear
	purego.RegisterLibFunc(&py.pyErrClear, py.libHandle, "PyErr_Clear")
	if py.pyErrClear == nil {
		return errors.New("failed to register PyErr_Clear")
	}

	return nil
}

// Initialize initializes the Python interpreter
func (py *PureGoPython) Initialize() error {
	if py.pyInitialize == nil {
		return errors.New("Python functions not registered")
	}

	// Check if already initialized
	if py.IsInitialized() {
		return errors.New("Python interpreter is already initialized")
	}

	py.pyInitialize()

	// Verify initialization succeeded
	if !py.IsInitialized() {
		return errors.New("Python interpreter initialization failed")
	}

	return nil
}

// Finalize shuts down the Python interpreter
func (py *PureGoPython) Finalize() error {
	if py.pyFinalizeEx == nil {
		return errors.New("Python functions not registered")
	}

	if !py.IsInitialized() {
		return errors.New("Python interpreter is not initialized")
	}

	result := py.pyFinalizeEx()
	if result < 0 {
		return fmt.Errorf("Python interpreter finalization failed with code: %d", result)
	}

	return nil
}

// IsInitialized returns true if the Python interpreter is initialized
func (py *PureGoPython) IsInitialized() bool {
	if py.pyIsInitialized == nil {
		return false
	}
	return py.pyIsInitialized() != 0
}

// stringToCString converts a Go string to a C string
func stringToCString(s string) *byte {
	if s == "" {
		return nil
	}
	bytes := []byte(s + "\x00") // null terminate
	return &bytes[0]
}

// checkPythonError checks if a Python error occurred and returns it as a Go error
func (py *PureGoPython) checkPythonError() error {
	if py.pyErrOccurred == nil {
		return nil
	}

	// Check if an error occurred
	if py.pyErrOccurred() != 0 {
		// Print the error to stderr (Python's default behavior)
		if py.pyErrPrint != nil {
			py.pyErrPrint()
		}

		// Clear the error
		if py.pyErrClear != nil {
			py.pyErrClear()
		}

		return errors.New("Python execution error occurred (see stderr for details)")
	}

	return nil
}

// RunString executes Python code from a string
func (py *PureGoPython) RunString(code string) error {
	if !py.IsInitialized() {
		return errors.New("Python interpreter is not initialized")
	}

	if py.pyRunSimpleString == nil {
		return errors.New("PyRun_SimpleString not registered")
	}

	if code == "" {
		return errors.New("code string cannot be empty")
	}

	// Convert Go string to C string
	cCode := stringToCString(code)
	if cCode == nil {
		return errors.New("failed to convert code to C string")
	}

	// Execute the Python code
	result := py.pyRunSimpleString(cCode)

	// Check for Python errors first
	if err := py.checkPythonError(); err != nil {
		return err
	}

	// Check return code (0 = success, -1 = error)
	if result != 0 {
		return fmt.Errorf("PyRun_SimpleString failed with return code: %d", result)
	}

	return nil
}

// RunFile executes Python code from a file using Python's exec() function
func (py *PureGoPython) RunFile(filename string) error {
	if !py.IsInitialized() {
		return errors.New("Python interpreter is not initialized")
	}

	if filename == "" {
		return errors.New("filename cannot be empty")
	}

	// Check if file exists and is readable
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename)
	} else if err != nil {
		return fmt.Errorf("cannot access file %s: %v", filename, err)
	}

	// Convert to absolute path to avoid any path issues
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %v", filename, err)
	}

	// Use Python's exec() function to execute the file
	// This is safer than using PyRun_SimpleFile and avoids FILE* pointer issues
	execCode := fmt.Sprintf(`
try:
    with open(r'%s', 'r', encoding='utf-8') as f:
        exec(f.read(), globals())
except Exception as e:
    import traceback
    traceback.print_exc()
    raise e
`, absPath)

	// Execute using RunString with proper error handling
	if err := py.RunString(execCode); err != nil {
		return fmt.Errorf("error executing file %s: %v", filename, err)
	}

	return nil
}
