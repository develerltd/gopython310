# Development Guide

This guide covers how to get started contributing to the Go-Python library.

## Overview

The Go-Python library enables calling Python 3.10 functions from Go using the Ebitengine purego interface to bind CPython C API without Cgo. It provides thread-safe Python function calls with bidirectional type conversion and virtual environment support.

## Architecture

The codebase is organized into focused modules:

```
├── python.go           # Main public API and documentation
├── types.go           # Type definitions and structures  
├── bindings.go        # CPython API function bindings via purego
├── conversion.go      # Go ↔ Python type conversion functions
├── interpreter.go     # Python interpreter lifecycle management
├── venv.go           # Virtual environment support
├── threading.go      # Thread safety wrappers
├── platform.go       # Cross-platform compatibility utilities
└── examples/         # Usage examples and tests
    ├── basic/        # Basic functionality demonstration
    ├── concurrent/   # Thread safety testing
    └── venv/         # Virtual environment usage
```

## Prerequisites

### Required Software
- **Go 1.21+** (for `unsafe.Add` support)
- **Python 3.10** development headers and shared library
- **Git** for version control

### Platform-Specific Setup

#### Linux (Ubuntu/Debian)
```bash
# Install Python development headers
sudo apt update
sudo apt install python3.10-dev

# Verify libpython3.10.so location
find /usr -name "libpython3.10.so*" 2>/dev/null
```

#### macOS
```bash
# Homebrew Python (recommended)
brew install python@3.10

# Verify libpython3.10.dylib location
ls /opt/homebrew/lib/libpython3.10.dylib  # Apple Silicon
ls /usr/local/lib/libpython3.10.dylib     # Intel

# Or pyenv
pyenv install 3.10.15
ls ~/.pyenv/versions/3.10.15/lib/libpython3.10.dylib
```

#### Windows
```bash
# Standard Python installation
# Verify python310.dll location
dir "C:\Python310\python310.dll"
```

## Getting Started

### 1. Clone and Setup
```bash
git clone <repository-url>
cd go-python
go mod tidy
```

### 2. Find Your Python Library
The library requires the path to your Python 3.10 shared library:

**Linux:**
```bash
find /usr -name "libpython3.10.so*" 2>/dev/null
# Common: /usr/lib/x86_64-linux-gnu/libpython3.10.so.1.0
```

**macOS:**
```bash
# Homebrew
ls /opt/homebrew/lib/libpython3.10.dylib

# pyenv  
ls ~/.pyenv/versions/*/lib/libpython3.10.dylib
```

### 3. Run Examples
```bash
# Basic functionality
go run examples/basic/main.go <path-to-libpython3.10>

# Thread safety testing
go run examples/concurrent/main.go <path-to-libpython3.10>

# Virtual environment support
go run examples/venv/main.go <path-to-libpython3.10> <path-to-venv>
```

### 4. Run Tests
```bash
# Build all components
go build -v

# Test examples build
go build -v examples/basic/main.go
go build -v examples/concurrent/main.go  
go build -v examples/venv/main.go
```

## Development Workflow

### Code Organization Principles

1. **Single Responsibility**: Each file has one clear purpose
2. **Clean Separation**: Public API in `python.go`, implementation in specialized files
3. **Platform Agnostic**: Core logic works across Linux, macOS, Windows
4. **Thread Safety**: All public methods use mutex protection

### Key Design Decisions

**Thread Safety Strategy:**
- Uses Go mutex instead of Python GIL state management
- Serializes all Python operations for reliability
- Avoids `PyGILState_Ensure/Release` due to embedded interpreter issues

**Memory Management:**
- Automatic Python reference counting via `safeDecRef()`
- Proper cleanup in defer statements
- No manual memory management required by users

**Error Handling:**
- Python exceptions converted to Go errors
- Descriptive error messages with context
- Graceful fallbacks for missing features

### Adding New Features

#### 1. Type Conversion
To support new Go ↔ Python type conversions:

1. Add the type case to `goToPython()` in `conversion.go`
2. Add the type case to `pythonToGo()` in `conversion.go`
3. Add tests in the basic example
4. Update documentation in `python.go`

#### 2. New CPython Functions
To bind new CPython API functions:

1. Add function pointer to `PureGoPython` struct in `types.go`
2. Register function in `registerPythonFunctions()` in `bindings.go`
3. Add to validation in `validateFunctionRegistration()` if critical
4. Create wrapper functions as needed

#### 3. Platform Support
For new platform-specific features:

1. Add platform detection to `platform.go`
2. Update path handling in virtual environment functions
3. Test on target platform
4. Update documentation with platform-specific examples

### Testing Strategy

#### Manual Testing
- Run all examples on your development platform
- Test with different Python library paths
- Verify virtual environment functionality
- Test error conditions and edge cases

#### Cross-Platform Testing
- Test on Linux, macOS if possible
- Verify library path detection works correctly
- Test virtual environment path handling
- Ensure examples work across platforms

#### Thread Safety Testing
- Run concurrent example with high load
- Monitor for race conditions or deadlocks
- Test multiple goroutines calling Python simultaneously
- Verify cleanup happens correctly

### Common Development Tasks

#### Adding a New Example
1. Create directory under `examples/`
2. Add `main.go` with clear demonstration
3. Include error handling and cleanup
4. Update README.md with usage instructions

#### Debugging Issues
1. **Library Loading**: Verify libpython path is correct for your platform
2. **Function Binding**: Check if CPython function exists in your Python version
3. **Type Conversion**: Add debug prints in conversion functions
4. **Thread Safety**: Use race detector: `go run -race examples/concurrent/main.go`

#### Performance Optimization
1. **Profiling**: Use Go's built-in profiler on the concurrent example
2. **Memory**: Monitor Python object reference counting
3. **Bottlenecks**: Identify serial operations that could be optimized

## Code Style Guidelines

### Go Conventions
- Follow standard Go formatting (`gofmt`)
- Use descriptive variable names
- Add comments for exported functions
- Handle errors explicitly
- Use defer for cleanup

### Documentation
- Keep `python.go` focused on API documentation
- Add implementation details to relevant specialized files
- Update README.md for user-facing changes
- Include examples for new features

### Error Messages
- Include context about what operation failed
- Preserve original Python error messages
- Use consistent error formatting
- Provide actionable error messages

## Common Issues and Solutions

### Library Loading Errors
**Symptom**: `failed to load libpython`
**Solution**: Verify the library path is correct and file exists

### Symbol Not Found
**Symptom**: `undefined symbol: PyXxx_Xxx`
**Solution**: Function may not exist in Python 3.10; check CPython documentation

### Thread Safety Issues  
**Symptom**: Crashes or deadlocks in concurrent usage
**Solution**: Ensure all Python operations go through thread-safe wrappers

### Virtual Environment Problems
**Symptom**: Wrong packages loaded or import errors
**Solution**: Check venv structure and site-packages path detection

### Memory Issues
**Symptom**: Memory leaks or crashes
**Solution**: Verify reference counting and cleanup in defer statements

## Contributing

### Before Submitting
1. Test on your platform thoroughly
2. Ensure all examples build and run
3. Add appropriate error handling
4. Update documentation if needed
5. Follow existing code patterns

### Pull Request Process
1. Create focused, single-purpose changes
2. Include clear commit messages
3. Test cross-platform if possible
4. Update relevant documentation
5. Ensure backwards compatibility

### Reporting Issues
1. Include platform and Python version information
2. Provide minimal reproduction case
3. Include error messages and stack traces
4. Test with basic example first

## Advanced Topics

### Multi-Version Python Support
The current library targets Python 3.10. Supporting multiple versions would require:
- Runtime version detection
- Conditional function binding
- Version-specific fallbacks
- Extensive testing matrix

### Performance Optimization
Current performance characteristics:
- Thread-safe but serialized Python operations
- No true parallel Python execution
- Memory overhead from type conversions
- Suitable for I/O bound and mixed workloads

### Integration Patterns
Common integration approaches:
- Wrap Python libraries for Go consumption
- Execute Python scripts from Go applications  
- Use Python for data processing in Go pipelines
- Leverage Python's ecosystem from Go services

---

This library provides a solid foundation for Go-Python interoperability. The modular architecture makes it straightforward to extend and maintain while ensuring reliability across platforms.