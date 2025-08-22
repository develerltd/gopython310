# Virtual Environment Example

This example demonstrates how to use the Go-Python library with Python virtual environments, allowing you to access packages installed in a specific venv.

## Prerequisites

1. **Create a virtual environment**:
   ```bash
   python3 -m venv /path/to/your/venv
   source /path/to/your/venv/bin/activate
   ```

2. **Install some packages in the venv**:
   ```bash
   pip install requests numpy pandas flask
   ```

3. **Deactivate the venv**:
   ```bash
   deactivate
   ```

## Features Demonstrated

- **Virtual Environment Detection**: Automatically finds venv site-packages
- **Package Access**: Import and use packages installed in the venv
- **System Fallback**: Optionally include system site-packages
- **Path Configuration**: Shows how Python paths are configured
- **Function Calls**: Call Python functions that use venv packages from Go

## Running

```bash
go run examples/venv/main.go /path/to/libpython3.10.so /path/to/your/venv
```

## Example Output

```
Testing virtual environment support
Libpython: /usr/lib/x86_64-linux-gnu/libpython3.10.so
Virtual environment: /home/user/myenv

Initializing Python with virtual environment...

=== Test 1: Python Path Configuration ===
Python executable: /usr/bin/python3
Python version: 3.10.12 (main, Nov 20 2023, 15:14:05) [GCC 11.4.0] on linux
Virtual environment check: True

Python path entries:
  0: /home/user/myenv/lib/python3.10/site-packages
  1: /usr/lib/python310.zip
  2: /usr/lib/python3.10
  3: /usr/lib/python3.10/lib-dynload

=== Test 2: Available Packages ===
Found 15 installed packages:
  - certifi
  - charset-normalizer
  - idna
  - numpy
  - pandas
  - pip
  - requests
  - setuptools
  - urllib3
  - wheel

=== Test 3: Package Import Tests ===
✓ requests: Available
  Version: 2.31.0
✓ numpy: Available  
  Version: 1.24.3
✓ pandas: Available
  Version: 2.0.3
✗ flask: Not available
...

=== Test 4: Functionality Test ===
  requests: HTTP 200
  numpy: mean of [1,2,3,4,5] = 3.0
  pandas: DataFrame shape = (3, 2)

=== Test 5: Go-to-Python Function Calls with Venv ===
Math test result: map[mean:5.5 std:3.0276503540974917 sum:55 using:numpy]

=== Virtual Environment Test Complete ===
✓ Successfully initialized Python with virtual environment
✓ Virtual environment packages are accessible  
✓ Go-to-Python function calls work with venv packages
```

## Configuration Options

```go
venvConfig := gopython.VirtualEnvConfig{
    VenvPath:   "/path/to/venv",           // Required: path to virtual environment
    SystemSite: true,                      // Optional: include system packages
    SitePaths:  []string{"/custom/path"},  // Optional: additional package directories
    PythonHome: "/usr",                    // Optional: Python installation directory
}
```

## How It Works

1. **Pre-initialization**: Configures Python paths using `Py_SetPath()` before `Py_Initialize()`
2. **Virtual Environment Detection**: Automatically finds `lib/python3.10/site-packages` in the venv
3. **Site Directory Addition**: Uses Python's `site.addsitedir()` to add additional paths
4. **Package Resolution**: Python imports resolve to venv packages first, then system packages

## Common Use Cases

### Data Science Workflow
```go
// Use a venv with scipy, numpy, pandas, matplotlib
venvConfig := gopython.VirtualEnvConfig{
    VenvPath:   "/path/to/datascience-venv",
    SystemSite: false, // Only venv packages
}

py.InitializeWithVenv(venvConfig)

// Now you can call functions that use these packages
result, _ := py.CallFunction("analytics", "process_dataframe", data)
```

### Web Development
```go
// Use a venv with flask, django, requests
venvConfig := gopython.VirtualEnvConfig{
    VenvPath:   "/path/to/web-venv", 
    SystemSite: true, // Include system packages as fallback
}

py.InitializeWithVenv(venvConfig)

// Call web-related Python functions
response, _ := py.CallFunction("web_utils", "make_api_request", url, params)
```

### Machine Learning
```go
// Use a venv with tensorflow, scikit-learn, etc.
venvConfig := gopython.VirtualEnvConfig{
    VenvPath:   "/path/to/ml-venv",
    SystemSite: false,
}

py.InitializeWithVenv(venvConfig)

// Call ML functions
prediction, _ := py.CallFunction("ml_model", "predict", features)
```

## Troubleshooting

### Virtual Environment Not Found
```
Error: virtual environment does not exist: /path/to/venv
```
- Check that the venv path is correct
- Ensure the venv was created with `python3 -m venv`

### Package Import Failures
```
✗ package_name: Not available
```
- Activate the venv and install the package: `pip install package_name`
- Check that the package is in the venv's site-packages directory

### Path Configuration Issues
- Enable system site packages if you need fallback: `SystemSite: true`
- Add custom paths if packages are in non-standard locations
- Check Python path output to verify configuration

This example shows that virtual environments work seamlessly with the embedded Python library!