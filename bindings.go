package gopython

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

// registerPythonFunctions registers all CPython API functions with purego
func (py *PureGoPython) registerPythonFunctions() error {
	// Core interpreter functions
	purego.RegisterLibFunc(&py.pyInitialize, py.libHandle, "Py_Initialize")
	purego.RegisterLibFunc(&py.pyFinalizeEx, py.libHandle, "Py_FinalizeEx")
	purego.RegisterLibFunc(&py.pyIsInitialized, py.libHandle, "Py_IsInitialized")
	purego.RegisterLibFunc(&py.pySetProgramName, py.libHandle, "Py_SetProgramName")
	purego.RegisterLibFunc(&py.pySetPythonHome, py.libHandle, "Py_SetPythonHome")
	purego.RegisterLibFunc(&py.pySetPath, py.libHandle, "Py_SetPath")

	// Code execution functions
	purego.RegisterLibFunc(&py.pyRunSimpleString, py.libHandle, "PyRun_SimpleString")

	// Module and import functions
	purego.RegisterLibFunc(&py.pyImportImport, py.libHandle, "PyImport_Import")
	purego.RegisterLibFunc(&py.pyImportAddModule, py.libHandle, "PyImport_AddModule")
	purego.RegisterLibFunc(&py.pyModuleGetDict, py.libHandle, "PyModule_GetDict")
	purego.RegisterLibFunc(&py.pyDictGetItemString, py.libHandle, "PyDict_GetItemString")

	// Object attribute functions
	purego.RegisterLibFunc(&py.pyObjectGetAttr, py.libHandle, "PyObject_GetAttr")
	purego.RegisterLibFunc(&py.pyObjectCallObject, py.libHandle, "PyObject_CallObject")
	purego.RegisterLibFunc(&py.pyObjectType, py.libHandle, "PyObject_Type")
	purego.RegisterLibFunc(&py.pyObjectStr, py.libHandle, "PyObject_Str")
	purego.RegisterLibFunc(&py.pyObjectRepr, py.libHandle, "PyObject_Repr")

	// String/Unicode functions
	purego.RegisterLibFunc(&py.pyUnicodeFromString, py.libHandle, "PyUnicode_FromString")
	purego.RegisterLibFunc(&py.pyUnicodeAsUTF8, py.libHandle, "PyUnicode_AsUTF8")

	// Integer functions
	purego.RegisterLibFunc(&py.pyLongFromLong, py.libHandle, "PyLong_FromLong")
	purego.RegisterLibFunc(&py.pyLongAsLong, py.libHandle, "PyLong_AsLong")
	purego.RegisterLibFunc(&py.pyLongFromSize, py.libHandle, "PyLong_FromSize_t")
	purego.RegisterLibFunc(&py.pyBoolFromLong, py.libHandle, "PyBool_FromLong")

	// Float functions
	purego.RegisterLibFunc(&py.pyFloatFromDouble, py.libHandle, "PyFloat_FromDouble")
	purego.RegisterLibFunc(&py.pyFloatAsDouble, py.libHandle, "PyFloat_AsDouble")

	// List functions
	purego.RegisterLibFunc(&py.pyListNew, py.libHandle, "PyList_New")
	purego.RegisterLibFunc(&py.pyListSetItem, py.libHandle, "PyList_SetItem")
	purego.RegisterLibFunc(&py.pyListGetItem, py.libHandle, "PyList_GetItem")
	purego.RegisterLibFunc(&py.pyListSize, py.libHandle, "PyList_Size")

	// Dictionary functions
	purego.RegisterLibFunc(&py.pyDictNew, py.libHandle, "PyDict_New")
	purego.RegisterLibFunc(&py.pyDictSetItemString, py.libHandle, "PyDict_SetItemString")
	purego.RegisterLibFunc(&py.pyDictKeys, py.libHandle, "PyDict_Keys")

	// Tuple functions
	purego.RegisterLibFunc(&py.pyTupleNew, py.libHandle, "PyTuple_New")
	purego.RegisterLibFunc(&py.pyTupleSetItem, py.libHandle, "PyTuple_SetItem")
	purego.RegisterLibFunc(&py.pyTupleGetItem, py.libHandle, "PyTuple_GetItem")
	purego.RegisterLibFunc(&py.pyTupleSize, py.libHandle, "PyTuple_Size")

	// Type checking functions - Note: PyType_GetName only available in Python 3.11+
	// We'll use an alternative approach for Python 3.10 compatibility

	// Reference counting functions
	purego.RegisterLibFunc(&py.pyIncRef, py.libHandle, "Py_IncRef")
	purego.RegisterLibFunc(&py.pyDecRef, py.libHandle, "Py_DecRef")

	// Error handling functions
	purego.RegisterLibFunc(&py.pyErrOccurred, py.libHandle, "PyErr_Occurred")
	purego.RegisterLibFunc(&py.pyErrFetch, py.libHandle, "PyErr_Fetch")
	purego.RegisterLibFunc(&py.pyErrClear, py.libHandle, "PyErr_Clear")

	// GIL functions (for future use if needed)
	purego.RegisterLibFunc(&py.pyGILStateEnsure, py.libHandle, "PyGILState_Ensure")
	purego.RegisterLibFunc(&py.pyGILStateRelease, py.libHandle, "PyGILState_Release")

	return nil
}

// Type checking helper functions using runtime type inspection
// These replace the macro-based type checking that caused undefined symbol errors

// getTypeName returns the type name of a Python object using Python 3.10 compatible approach
func (py *PureGoPython) getTypeName(obj PyObject) string {
	if obj == 0 {
		return "NoneType"
	}

	typeObj := py.pyObjectType(uintptr(obj))
	if typeObj == 0 {
		return "unknown"
	}
	defer py.safeDecRef(typeObj)

	// Get the __name__ attribute from the type object (Python 3.10 compatible)
	// Create "__name__" string directly to avoid circular dependency
	nameAttrStr := stringToCString("__name__")
	nameAttrObj := py.pyUnicodeFromString(nameAttrStr)
	if nameAttrObj == 0 {
		return "unknown"
	}
	defer py.safeDecRef(nameAttrObj)

	nameObj := py.pyObjectGetAttr(typeObj, nameAttrObj)
	if nameObj == 0 {
		return "unknown"
	}
	defer py.safeDecRef(nameObj)

	// Convert to string - use direct Unicode conversion
	cStr := py.pyUnicodeAsUTF8(nameObj)
	if cStr != nil {
		return cStringToGoString(cStr)
	}

	return "unknown"
}

// isStringUnsafe checks if object is string without using getTypeName (to avoid circular dependency)
func (py *PureGoPython) isStringUnsafe(obj PyObject) bool {
	if obj == 0 {
		return false
	}
	// Try to convert to UTF-8 - if it succeeds, it's likely a string
	cStr := py.pyUnicodeAsUTF8(uintptr(obj))
	return cStr != nil
}

// isString checks if a Python object is a string
func (py *PureGoPython) isString(obj PyObject) bool {
	typeName := py.getTypeName(obj)
	return typeName == "str"
}

// isInt checks if a Python object is an integer
func (py *PureGoPython) isInt(obj PyObject) bool {
	typeName := py.getTypeName(obj)
	return typeName == "int"
}

// isBool checks if a Python object is a boolean
func (py *PureGoPython) isBool(obj PyObject) bool {
	typeName := py.getTypeName(obj)
	return typeName == "bool"
}

// isFloat checks if a Python object is a float
func (py *PureGoPython) isFloat(obj PyObject) bool {
	typeName := py.getTypeName(obj)
	return typeName == "float"
}

// isList checks if a Python object is a list
func (py *PureGoPython) isList(obj PyObject) bool {
	typeName := py.getTypeName(obj)
	return typeName == "list"
}

// isDict checks if a Python object is a dictionary
func (py *PureGoPython) isDict(obj PyObject) bool {
	typeName := py.getTypeName(obj)
	return typeName == "dict"
}

// isTuple checks if a Python object is a tuple
func (py *PureGoPython) isTuple(obj PyObject) bool {
	typeName := py.getTypeName(obj)
	return typeName == "tuple"
}

// isNone checks if a Python object is None
func (py *PureGoPython) isNone(obj PyObject) bool {
	return obj == 0 || py.getTypeName(obj) == "NoneType"
}

// safeDecRef safely decrements reference count, handling nil/zero pointers
func (py *PureGoPython) safeDecRef(obj uintptr) {
	if obj != 0 && py.pyDecRef != nil {
		py.pyDecRef(obj)
	}
}

// cStringToGoString converts a C string to a Go string
func cStringToGoString(ptr *byte) string {
	if ptr == nil {
		return ""
	}

	var result []byte
	for i := 0; ; i++ {
		b := (*byte)(unsafe.Add(unsafe.Pointer(ptr), i))
		if *b == 0 {
			break
		}
		result = append(result, *b)
	}
	return string(result)
}

// validateFunctionRegistration checks that all critical functions are registered
func (py *PureGoPython) validateFunctionRegistration() error {
	if py.pyInitialize == nil {
		return fmt.Errorf("failed to register Py_Initialize")
	}
	if py.pyFinalizeEx == nil {
		return fmt.Errorf("failed to register Py_FinalizeEx")
	}
	if py.pyRunSimpleString == nil {
		return fmt.Errorf("failed to register PyRun_SimpleString")
	}
	if py.pyImportImport == nil {
		return fmt.Errorf("failed to register PyImport_Import")
	}
	if py.pyObjectCallObject == nil {
		return fmt.Errorf("failed to register PyObject_CallObject")
	}
	return nil
}