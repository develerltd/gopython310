# RunString with Return Values Example

This example demonstrates how to use `RunString` to define Python functions that return values, then use `CallFunction` to retrieve those values with type safety.

## Features Demonstrated

- **RunString + CallFunction Pattern**: Define functions with `RunString`, call them with `CallFunction`
- **Complex Return Values**: Functions returning dicts, lists, and computed results
- **Type Safety**: Using generic `CallPyFunction[TRequest, TResponse]` for compile-time type checking
- **Real-world Use Cases**: Math calculations, data processing, text analysis

## Files

- `main.go` - Main example program demonstrating the pattern

## Running

```bash
go run main.go /path/to/libpython3.10.so
```

## What it does

1. **Math Calculations**: Defines a function that performs math operations and returns structured results
2. **Data Processing**: Creates a function that processes arrays and returns analysis
3. **Text Analysis**: Shows string manipulation and analysis with return values
4. **Type-Safe Calls**: Demonstrates using generics for compile-time type safety

## Key Pattern

```go
// 1. Define function with RunString
code := `
def my_function():
    # Complex Python logic here
    return {"result": calculated_value}
`
py.RunString(code)

// 2. Call function and get typed result
result, err := gopython.CallPyFunction[interface{}, map[string]interface{}](
    py, "__main__", "my_function", nil)
```

This approach combines the flexibility of Python code execution with the type safety and convenience of function calls.