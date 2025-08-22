package gopython

import (
	"fmt"
	"unsafe"
)

// goToPython converts Go values to Python objects
func (py *PureGoPython) goToPython(value interface{}) (PyObject, error) {
	if value == nil {
		return 0, nil // Python None
	}

	switch v := value.(type) {
	case string:
		cStr := stringToCString(v)
		pyStr := py.pyUnicodeFromString(cStr)
		if pyStr == 0 {
			return 0, fmt.Errorf("failed to create Python string")
		}
		return PyObject(pyStr), nil

	case int:
		pyInt := py.pyLongFromLong(int64(v))
		if pyInt == 0 {
			return 0, fmt.Errorf("failed to create Python int")
		}
		return PyObject(pyInt), nil

	case int64:
		pyInt := py.pyLongFromLong(v)
		if pyInt == 0 {
			return 0, fmt.Errorf("failed to create Python int")
		}
		return PyObject(pyInt), nil

	case float64:
		pyFloat := py.pyFloatFromDouble(v)
		if pyFloat == 0 {
			return 0, fmt.Errorf("failed to create Python float")
		}
		return PyObject(pyFloat), nil

	case bool:
		var pyBool uintptr
		if v {
			pyBool = py.pyBoolFromLong(1)
		} else {
			pyBool = py.pyBoolFromLong(0)
		}
		if pyBool == 0 {
			return 0, fmt.Errorf("failed to create Python bool")
		}
		return PyObject(pyBool), nil

	case []interface{}:
		return py.sliceToPythonList(v)

	case map[string]interface{}:
		return py.mapToPythonDict(v)

	default:
		return 0, fmt.Errorf("unsupported Go type: %T", value)
	}
}

// sliceToPythonList converts a Go slice to a Python list
func (py *PureGoPython) sliceToPythonList(slice []interface{}) (PyObject, error) {
	pyList := py.pyListNew(len(slice))
	if pyList == 0 {
		return 0, fmt.Errorf("failed to create Python list")
	}

	for i, item := range slice {
		pyItem, err := py.goToPython(item)
		if err != nil {
			py.safeDecRef(pyList)
			return 0, fmt.Errorf("failed to convert slice item %d: %v", i, err)
		}

		// PyList_SetItem steals the reference, so we don't need to decref pyItem
		if py.pyListSetItem(pyList, i, uintptr(pyItem)) != 0 {
			py.safeDecRef(pyList)
			return 0, fmt.Errorf("failed to set list item %d", i)
		}
	}

	return PyObject(pyList), nil
}

// mapToPythonDict converts a Go map to a Python dictionary
func (py *PureGoPython) mapToPythonDict(m map[string]interface{}) (PyObject, error) {
	pyDict := py.pyDictNew()
	if pyDict == 0 {
		return 0, fmt.Errorf("failed to create Python dict")
	}

	for key, value := range m {
		pyValue, err := py.goToPython(value)
		if err != nil {
			py.safeDecRef(pyDict)
			return 0, fmt.Errorf("failed to convert dict value for key '%s': %v", key, err)
		}

		cKey := stringToCString(key)
		if py.pyDictSetItemString(pyDict, cKey, uintptr(pyValue)) != 0 {
			py.safeDecRef(pyDict)
			py.safeDecRef(uintptr(pyValue))
			return 0, fmt.Errorf("failed to set dict item for key '%s'", key)
		}

		// PyDict_SetItemString doesn't steal the reference, so we need to decref
		py.safeDecRef(uintptr(pyValue))
	}

	return PyObject(pyDict), nil
}

// pythonToGo converts Python objects to Go values
func (py *PureGoPython) pythonToGo(obj PyObject) (interface{}, error) {
	if py.isNone(obj) {
		return nil, nil
	}

	// Check string first
	if py.isString(obj) {
		cStr := py.pyUnicodeAsUTF8(uintptr(obj))
		if cStr == nil {
			return nil, fmt.Errorf("failed to convert Python string to UTF-8")
		}
		return cStringToGoString(cStr), nil
	}

	// Check bool first (since bool is a subclass of int in Python)
	if py.isBool(obj) {
		return py.pyLongAsLong(uintptr(obj)) != 0, nil
	}

	// Check integer
	if py.isInt(obj) {
		return py.pyLongAsLong(uintptr(obj)), nil
	}

	// Check float
	if py.isFloat(obj) {
		return py.pyFloatAsDouble(uintptr(obj)), nil
	}

	// Check list
	if py.isList(obj) {
		return py.pythonListToSlice(obj)
	}

	// Check dict
	if py.isDict(obj) {
		return py.pythonDictToMap(obj)
	}

	typeName := py.getTypeName(obj)
	return nil, fmt.Errorf("unsupported Python type: %s", typeName)
}

// pythonListToSlice converts a Python list to a Go slice
func (py *PureGoPython) pythonListToSlice(obj PyObject) ([]interface{}, error) {
	size := py.pyListSize(uintptr(obj))
	result := make([]interface{}, size)
	
	for i := 0; i < size; i++ {
		item := py.pyListGetItem(uintptr(obj), i)
		val, err := py.pythonToGo(PyObject(item))
		if err != nil {
			return nil, fmt.Errorf("failed to convert list item %d: %v", i, err)
		}
		result[i] = val
	}
	
	return result, nil
}

// pythonDictToMap converts a Python dictionary to a Go map
func (py *PureGoPython) pythonDictToMap(obj PyObject) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	keys := py.pyDictKeys(uintptr(obj))
	if keys == 0 {
		return nil, fmt.Errorf("failed to get dict keys")
	}
	defer py.safeDecRef(keys)

	size := py.pyListSize(keys)
	for i := 0; i < size; i++ {
		keyObj := py.pyListGetItem(keys, i)
		if !py.isString(PyObject(keyObj)) {
			continue // Skip non-string keys
		}

		cKey := py.pyUnicodeAsUTF8(keyObj)
		if cKey == nil {
			continue
		}

		// Convert key to Go string
		key := ""
		for j := 0; ; j++ {
			b := (*byte)(unsafe.Add(unsafe.Pointer(cKey), j))
			if *b == 0 {
				break
			}
			key += string(*b)
		}

		valObj := py.pyDictGetItemString(uintptr(obj), cKey)
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

// buildArgumentTuple converts Go arguments to a Python tuple for function calls
func (py *PureGoPython) buildArgumentTuple(args ...interface{}) (PyObject, error) {
	argTuple := py.pyTupleNew(len(args))
	if argTuple == 0 {
		return 0, fmt.Errorf("failed to create argument tuple")
	}

	for i, arg := range args {
		pyArg, err := py.goToPython(arg)
		if err != nil {
			py.safeDecRef(argTuple)
			return 0, fmt.Errorf("failed to convert argument %d: %v", i, err)
		}

		// PyTuple_SetItem steals the reference
		if py.pyTupleSetItem(argTuple, i, uintptr(pyArg)) != 0 {
			py.safeDecRef(argTuple)
			return 0, fmt.Errorf("failed to set tuple item %d", i)
		}
	}

	return PyObject(argTuple), nil
}