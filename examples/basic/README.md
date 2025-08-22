# Basic Example

This example demonstrates the core functionality of the Go-Python library across all phases:

## Features Demonstrated

- **Phase 1**: Python interpreter initialization and cleanup
- **Phase 2**: Executing Python strings and files  
- **Phase 3**: Calling Python functions with type conversion

## Files

- `main.go` - Main example program
- `test.py` - Python file for testing file execution
- `mathtools.py` - Python module with functions for testing

## Running

```bash
go run main.go /path/to/libpython3.10.so
```

## What it does

1. Initializes the Python interpreter
2. Executes Python code from strings
3. Tests error handling with invalid Python code
4. Executes Python code from files
5. Creates Python functions in `__main__` module
6. Calls functions with different argument types
7. Demonstrates type conversion (strings, ints, lists, dicts)
8. Tests built-in module access (math.sqrt)

This example shows the complete workflow from basic interpreter management to advanced function calling with automatic type conversion.