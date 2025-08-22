package gopython

import (
	"sync"
	"unsafe"
)

// PyObject represents a Python object pointer
type PyObject uintptr

// VirtualEnvConfig contains configuration for virtual environment initialization
type VirtualEnvConfig struct {
	VenvPath   string   // Path to virtual environment directory
	SystemSite bool     // Include system site packages as fallback
	SitePaths  []string // Additional site package directories
	PythonHome string   // Python installation directory (optional)
}

// PureGoPython represents a Python runtime instance with CPython API bindings
type PureGoPython struct {
	libHandle uintptr
	mu        sync.Mutex // Thread safety protection

	// Core interpreter functions
	pyInitialize     func()
	pyFinalizeEx     func() int
	pyIsInitialized  func() int
	pySetProgramName func(*uint16)
	pySetPythonHome  func(*uint16)
	pySetPath        func(*uint16)

	// Code execution functions
	pyRunSimpleString func(*byte) int
	pyRunSimpleFile   func(uintptr, *byte) int

	// Module and import functions
	pyImportImport      func(uintptr) uintptr
	pyImportAddModule   func(*byte) uintptr
	pyModuleGetDict     func(uintptr) uintptr
	pyDictGetItemString func(uintptr, *byte) uintptr

	// Object attribute functions
	pyObjectGetAttr     func(uintptr, uintptr) uintptr
	pyObjectCallObject  func(uintptr, uintptr) uintptr
	pyObjectType        func(uintptr) uintptr
	pyObjectStr         func(uintptr) uintptr
	pyObjectRepr        func(uintptr) uintptr
	pyObjectGetTypeName func(uintptr) *byte

	// String/Unicode functions
	pyUnicodeFromString func(*byte) uintptr
	pyUnicodeAsUTF8     func(uintptr) *byte

	// Integer functions
	pyLongFromLong  func(int64) uintptr
	pyLongAsLong    func(uintptr) int64
	pyLongFromSize  func(int) uintptr
	pyBoolFromLong  func(int64) uintptr

	// Float functions
	pyFloatFromDouble func(float64) uintptr
	pyFloatAsDouble   func(uintptr) float64

	// List functions
	pyListNew     func(int) uintptr
	pyListSetItem func(uintptr, int, uintptr) int
	pyListGetItem func(uintptr, int) uintptr
	pyListSize    func(uintptr) int

	// Dictionary functions
	pyDictNew           func() uintptr
	pyDictSetItemString func(uintptr, *byte, uintptr) int
	pyDictKeys          func(uintptr) uintptr

	// Tuple functions
	pyTupleNew     func(int) uintptr
	pyTupleSetItem func(uintptr, int, uintptr) int
	pyTupleGetItem func(uintptr, int) uintptr
	pyTupleSize    func(uintptr) int

	// Type checking functions (using runtime type inspection - Python 3.10 compatible)

	// Reference counting functions
	pyIncRef func(uintptr)
	pyDecRef func(uintptr)

	// Error handling functions
	pyErrOccurred func() uintptr
	pyErrFetch    func(*uintptr, *uintptr, *uintptr)
	pyErrClear    func()

	// File operations
	pyFileFromFd func(int, *byte, *byte, int, *byte, *byte, *byte, int) uintptr

	// GIL functions (for future use if needed)
	pyGILStateEnsure  func() int
	pyGILStateRelease func(int)
}

// stringToCString converts a Go string to a null-terminated C string
func stringToCString(s string) *byte {
	if len(s) == 0 {
		return (*byte)(unsafe.Pointer(&[]byte{0}[0]))
	}
	bytes := make([]byte, len(s)+1)
	copy(bytes, s)
	bytes[len(s)] = 0
	return (*byte)(unsafe.Pointer(&bytes[0]))
}

// uint16ToCWString converts a Go string to a null-terminated wide C string (UTF-16)
func uint16ToCWString(s string) *uint16 {
	runes := []rune(s)
	utf16 := make([]uint16, len(runes)+1)
	for i, r := range runes {
		utf16[i] = uint16(r)
	}
	utf16[len(runes)] = 0
	return (*uint16)(unsafe.Pointer(&utf16[0]))
}