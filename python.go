// Package gopython provides a Go interface to Python 3.10 using purego
package gopython

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

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
	
	// Function calling and objects
	pyImportImport     func(uintptr) uintptr
	pyObjectGetAttr    func(uintptr, uintptr) uintptr
	pyObjectCallObject func(uintptr, uintptr) uintptr
	pyTupleNew         func(int) uintptr
	pyTupleSetItem     func(uintptr, int, uintptr) int
	
	// Type conversion functions
	pyUnicodeFromString func(*byte) uintptr
	pyUnicodeAsUTF8     func(uintptr) *byte
	pyLongFromLong      func(int64) uintptr
	pyLongAsLong        func(uintptr) int64
	pyFloatFromDouble   func(float64) uintptr
	pyFloatAsDouble     func(uintptr) float64
	pyBoolFromLong      func(int) uintptr
	pyListNew           func(int) uintptr
	pyListSetItem       func(uintptr, int, uintptr) int
	pyListGetItem       func(uintptr, int) uintptr
	pyListSize          func(uintptr) int
	pyDictNew           func() uintptr
	pyDictSetItemString func(uintptr, *byte, uintptr) int
	pyDictGetItemString func(uintptr, *byte) uintptr
	pyDictKeys          func(uintptr) uintptr
	
	// Type checking using PyObject_Type and name comparison
	pyObjectType    func(uintptr) uintptr
	pyObjectRepr    func(uintptr) uintptr
	pyObjectGetAttrString func(uintptr, *byte) uintptr
	
	// Memory management
	pyDecRef func(uintptr)
	pyIncRef func(uintptr)
}

// PyObject represents a Python object pointer
type PyObject uintptr

// PyRuntime interface defines the core Python runtime operations
type PyRuntime interface {
	Initialize() error
	Finalize() error
	IsInitialized() bool
	RunString(code string) error
	RunFile(filename string) error
	CallFunction(module, function string, args ...interface{}) (interface{}, error)
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

	// Register function calling APIs
	purego.RegisterLibFunc(&py.pyImportImport, py.libHandle, "PyImport_Import")
	if py.pyImportImport == nil {
		return errors.New("failed to register PyImport_Import")
	}

	purego.RegisterLibFunc(&py.pyObjectGetAttr, py.libHandle, "PyObject_GetAttr")
	if py.pyObjectGetAttr == nil {
		return errors.New("failed to register PyObject_GetAttr")
	}

	purego.RegisterLibFunc(&py.pyObjectCallObject, py.libHandle, "PyObject_CallObject")
	if py.pyObjectCallObject == nil {
		return errors.New("failed to register PyObject_CallObject")
	}

	purego.RegisterLibFunc(&py.pyTupleNew, py.libHandle, "PyTuple_New")
	if py.pyTupleNew == nil {
		return errors.New("failed to register PyTuple_New")
	}

	purego.RegisterLibFunc(&py.pyTupleSetItem, py.libHandle, "PyTuple_SetItem")
	if py.pyTupleSetItem == nil {
		return errors.New("failed to register PyTuple_SetItem")
	}

	// Register type conversion APIs
	purego.RegisterLibFunc(&py.pyUnicodeFromString, py.libHandle, "PyUnicode_FromString")
	if py.pyUnicodeFromString == nil {
		return errors.New("failed to register PyUnicode_FromString")
	}

	purego.RegisterLibFunc(&py.pyUnicodeAsUTF8, py.libHandle, "PyUnicode_AsUTF8")
	if py.pyUnicodeAsUTF8 == nil {
		return errors.New("failed to register PyUnicode_AsUTF8")
	}

	purego.RegisterLibFunc(&py.pyLongFromLong, py.libHandle, "PyLong_FromLong")
	if py.pyLongFromLong == nil {
		return errors.New("failed to register PyLong_FromLong")
	}

	purego.RegisterLibFunc(&py.pyLongAsLong, py.libHandle, "PyLong_AsLong")
	if py.pyLongAsLong == nil {
		return errors.New("failed to register PyLong_AsLong")
	}

	purego.RegisterLibFunc(&py.pyFloatFromDouble, py.libHandle, "PyFloat_FromDouble")
	if py.pyFloatFromDouble == nil {
		return errors.New("failed to register PyFloat_FromDouble")
	}

	purego.RegisterLibFunc(&py.pyFloatAsDouble, py.libHandle, "PyFloat_AsDouble")
	if py.pyFloatAsDouble == nil {
		return errors.New("failed to register PyFloat_AsDouble")
	}

	purego.RegisterLibFunc(&py.pyBoolFromLong, py.libHandle, "PyBool_FromLong")
	if py.pyBoolFromLong == nil {
		return errors.New("failed to register PyBool_FromLong")
	}

	purego.RegisterLibFunc(&py.pyListNew, py.libHandle, "PyList_New")
	if py.pyListNew == nil {
		return errors.New("failed to register PyList_New")
	}

	purego.RegisterLibFunc(&py.pyListSetItem, py.libHandle, "PyList_SetItem")
	if py.pyListSetItem == nil {
		return errors.New("failed to register PyList_SetItem")
	}

	purego.RegisterLibFunc(&py.pyListGetItem, py.libHandle, "PyList_GetItem")
	if py.pyListGetItem == nil {
		return errors.New("failed to register PyList_GetItem")
	}

	purego.RegisterLibFunc(&py.pyListSize, py.libHandle, "PyList_Size")
	if py.pyListSize == nil {
		return errors.New("failed to register PyList_Size")
	}

	purego.RegisterLibFunc(&py.pyDictNew, py.libHandle, "PyDict_New")
	if py.pyDictNew == nil {
		return errors.New("failed to register PyDict_New")
	}

	purego.RegisterLibFunc(&py.pyDictSetItemString, py.libHandle, "PyDict_SetItemString")
	if py.pyDictSetItemString == nil {
		return errors.New("failed to register PyDict_SetItemString")
	}

	purego.RegisterLibFunc(&py.pyDictGetItemString, py.libHandle, "PyDict_GetItemString")
	if py.pyDictGetItemString == nil {
		return errors.New("failed to register PyDict_GetItemString")
	}

	purego.RegisterLibFunc(&py.pyDictKeys, py.libHandle, "PyDict_Keys")
	if py.pyDictKeys == nil {
		return errors.New("failed to register PyDict_Keys")
	}

	// Register type checking APIs
	purego.RegisterLibFunc(&py.pyObjectType, py.libHandle, "PyObject_Type")
	if py.pyObjectType == nil {
		return errors.New("failed to register PyObject_Type")
	}

	purego.RegisterLibFunc(&py.pyObjectRepr, py.libHandle, "PyObject_Repr")
	if py.pyObjectRepr == nil {
		return errors.New("failed to register PyObject_Repr")
	}

	purego.RegisterLibFunc(&py.pyObjectGetAttrString, py.libHandle, "PyObject_GetAttrString")
	if py.pyObjectGetAttrString == nil {
		return errors.New("failed to register PyObject_GetAttrString")
	}

	// Register memory management APIs
	purego.RegisterLibFunc(&py.pyDecRef, py.libHandle, "Py_DecRef")
	if py.pyDecRef == nil {
		return errors.New("failed to register Py_DecRef")
	}

	purego.RegisterLibFunc(&py.pyIncRef, py.libHandle, "Py_IncRef")
	if py.pyIncRef == nil {
		return errors.New("failed to register Py_IncRef")
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

// getTypeName returns the name of a Python object's type
func (py *PureGoPython) getTypeName(obj uintptr) string {
	if obj == 0 {
		return "NoneType"
	}
	
	typeObj := py.pyObjectType(obj)
	if typeObj == 0 {
		return "unknown"
	}
	defer py.safeDecRef(typeObj)
	
	nameObj := py.pyObjectGetAttrString(typeObj, stringToCString("__name__"))
	if nameObj == 0 {
		return "unknown"
	}
	defer py.safeDecRef(nameObj)
	
	cStr := py.pyUnicodeAsUTF8(nameObj)
	if cStr == nil {
		return "unknown"
	}
	
	// Convert C string to Go string safely
	name := ""
	for i := 0; ; i++ {
		b := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(cStr)) + uintptr(i)))
		if *b == 0 {
			break
		}
		name += string(*b)
	}
	return name
}

// Type checking helper functions
func (py *PureGoPython) isString(obj uintptr) bool {
	return py.getTypeName(obj) == "str"
}

func (py *PureGoPython) isInt(obj uintptr) bool {
	return py.getTypeName(obj) == "int"
}

func (py *PureGoPython) isFloat(obj uintptr) bool {
	return py.getTypeName(obj) == "float"
}

func (py *PureGoPython) isBool(obj uintptr) bool {
	return py.getTypeName(obj) == "bool"
}

func (py *PureGoPython) isList(obj uintptr) bool {
	return py.getTypeName(obj) == "list"
}

func (py *PureGoPython) isDict(obj uintptr) bool {
	return py.getTypeName(obj) == "dict"
}

// safeDecRef safely decrements reference count, checking for null pointer
func (py *PureGoPython) safeDecRef(obj uintptr) {
	if obj != 0 {
		py.pyDecRef(obj)
	}
}

// goToPython converts a Go value to a Python object
func (py *PureGoPython) goToPython(value interface{}) (PyObject, error) {
	if value == nil {
		// Return None (we'll use 0 as None for simplicity)
		return PyObject(0), nil
	}

	switch v := value.(type) {
	case string:
		cStr := stringToCString(v)
		pyObj := py.pyUnicodeFromString(cStr)
		if pyObj == 0 {
			return 0, fmt.Errorf("failed to convert string to Python object")
		}
		return PyObject(pyObj), nil

	case int:
		pyObj := py.pyLongFromLong(int64(v))
		if pyObj == 0 {
			return 0, fmt.Errorf("failed to convert int to Python object")
		}
		return PyObject(pyObj), nil

	case int64:
		pyObj := py.pyLongFromLong(v)
		if pyObj == 0 {
			return 0, fmt.Errorf("failed to convert int64 to Python object")
		}
		return PyObject(pyObj), nil

	case float64:
		pyObj := py.pyFloatFromDouble(v)
		if pyObj == 0 {
			return 0, fmt.Errorf("failed to convert float64 to Python object")
		}
		return PyObject(pyObj), nil

	case bool:
		var intVal int
		if v {
			intVal = 1
		} else {
			intVal = 0
		}
		pyObj := py.pyBoolFromLong(intVal)
		if pyObj == 0 {
			return 0, fmt.Errorf("failed to convert bool to Python object")
		}
		return PyObject(pyObj), nil

	case []interface{}:
		pyList := py.pyListNew(len(v))
		if pyList == 0 {
			return 0, fmt.Errorf("failed to create Python list")
		}

		for i, item := range v {
			pyItem, err := py.goToPython(item)
			if err != nil {
				py.safeDecRef(pyList)
				return 0, fmt.Errorf("failed to convert list item %d: %v", i, err)
			}
			if py.pyListSetItem(pyList, i, uintptr(pyItem)) != 0 {
				py.safeDecRef(pyList)
				py.safeDecRef(uintptr(pyItem))
				return 0, fmt.Errorf("failed to set list item %d", i)
			}
		}
		return PyObject(pyList), nil

	case map[string]interface{}:
		pyDict := py.pyDictNew()
		if pyDict == 0 {
			return 0, fmt.Errorf("failed to create Python dict")
		}

		for key, val := range v {
			pyVal, err := py.goToPython(val)
			if err != nil {
				py.safeDecRef(pyDict)
				return 0, fmt.Errorf("failed to convert dict value for key '%s': %v", key, err)
			}
			cKey := stringToCString(key)
			if py.pyDictSetItemString(pyDict, cKey, uintptr(pyVal)) != 0 {
				py.safeDecRef(pyDict)
				py.safeDecRef(uintptr(pyVal))
				return 0, fmt.Errorf("failed to set dict item '%s'", key)
			}
		}
		return PyObject(pyDict), nil

	default:
		return 0, fmt.Errorf("unsupported Go type: %T", value)
	}
}

// pythonToGo converts a Python object to a Go value
func (py *PureGoPython) pythonToGo(pyObj PyObject) (interface{}, error) {
	if pyObj == 0 {
		return nil, nil
	}

	obj := uintptr(pyObj)

	// Check string
	if py.isString(obj) {
		cStr := py.pyUnicodeAsUTF8(obj)
		if cStr == nil {
			return nil, fmt.Errorf("failed to convert Python string to C string")
		}
		
		// Convert C string to Go string safely
		str := ""
		for i := 0; ; i++ {
			b := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(cStr)) + uintptr(i)))
			if *b == 0 {
				break
			}
			str += string(*b)
		}
		return str, nil
	}

	// Check bool first (since bool is a subclass of int in Python)
	if py.isBool(obj) {
		return py.pyLongAsLong(obj) != 0, nil
	}

	// Check integer
	if py.isInt(obj) {
		return py.pyLongAsLong(obj), nil
	}

	// Check float
	if py.isFloat(obj) {
		return py.pyFloatAsDouble(obj), nil
	}

	// Check list
	if py.isList(obj) {
		size := py.pyListSize(obj)
		result := make([]interface{}, size)
		for i := 0; i < size; i++ {
			item := py.pyListGetItem(obj, i)
			val, err := py.pythonToGo(PyObject(item))
			if err != nil {
				return nil, fmt.Errorf("failed to convert list item %d: %v", i, err)
			}
			result[i] = val
		}
		return result, nil
	}

	// Check dict
	if py.isDict(obj) {
		result := make(map[string]interface{})
		keys := py.pyDictKeys(obj)
		if keys == 0 {
			return nil, fmt.Errorf("failed to get dict keys")
		}
		defer py.safeDecRef(keys)

		size := py.pyListSize(keys)
		for i := 0; i < size; i++ {
			keyObj := py.pyListGetItem(keys, i)
			if !py.isString(keyObj) {
				continue // Skip non-string keys
			}

			cKey := py.pyUnicodeAsUTF8(keyObj)
			if cKey == nil {
				continue
			}
			
			// Convert key to Go string
			key := ""
			for j := 0; ; j++ {
				b := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(cKey)) + uintptr(j)))
				if *b == 0 {
					break
				}
				key += string(*b)
			}

			valObj := py.pyDictGetItemString(obj, cKey)
			if valObj == 0 {
				continue
			}

			val, err := py.pythonToGo(PyObject(valObj))
			if err != nil {
				return nil, fmt.Errorf("failed to convert dict value for key '%s': %v", key, err)
			}
			result[key] = val
		}
		return result, nil
	}

	typeName := py.getTypeName(obj)
	return nil, fmt.Errorf("unsupported Python type: %s", typeName)
}

// CallFunction calls a Python function with the given arguments
func (py *PureGoPython) CallFunction(module, function string, args ...interface{}) (interface{}, error) {
	if !py.IsInitialized() {
		return nil, errors.New("Python interpreter is not initialized")
	}

	// Import the module
	moduleNameObj, err := py.goToPython(module)
	if err != nil {
		return nil, fmt.Errorf("failed to convert module name: %v", err)
	}
	defer py.safeDecRef(uintptr(moduleNameObj))

	moduleObj := py.pyImportImport(uintptr(moduleNameObj))
	if moduleObj == 0 {
		if err := py.checkPythonError(); err != nil {
			return nil, fmt.Errorf("failed to import module '%s': %v", module, err)
		}
		return nil, fmt.Errorf("failed to import module '%s'", module)
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
		if err := py.checkPythonError(); err != nil {
			return nil, fmt.Errorf("failed to get function '%s' from module '%s': %v", function, module, err)
		}
		return nil, fmt.Errorf("function '%s' not found in module '%s'", function, module)
	}
	defer py.safeDecRef(functionObj)

	// Convert arguments to Python objects
	argTuple := py.pyTupleNew(len(args))
	if argTuple == 0 {
		return nil, fmt.Errorf("failed to create argument tuple")
	}
	defer py.safeDecRef(argTuple)

	for i, arg := range args {
		pyArg, err := py.goToPython(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to convert argument %d: %v", i, err)
		}
		if py.pyTupleSetItem(argTuple, i, uintptr(pyArg)) != 0 {
			py.safeDecRef(uintptr(pyArg))
			return nil, fmt.Errorf("failed to set argument %d in tuple", i)
		}
		// Note: PyTuple_SetItem steals the reference, so we don't DecRef pyArg
	}

	// Call the function
	resultObj := py.pyObjectCallObject(functionObj, argTuple)
	if resultObj == 0 {
		if err := py.checkPythonError(); err != nil {
			return nil, fmt.Errorf("error calling function '%s.%s': %v", module, function, err)
		}
		return nil, fmt.Errorf("function call '%s.%s' returned NULL", module, function)
	}
	defer py.safeDecRef(resultObj)

	// Convert result back to Go
	result, err := py.pythonToGo(PyObject(resultObj))
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %v", err)
	}

	return result, nil
}
