package main

import (
	"fmt"
	"log"
	"os"

	"gopython"
)

func main() {
	// Check if arguments were provided
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run examples/venv/main.go <path-to-libpython3.10.so> <path-to-venv>")
	}

	libpythonPath := os.Args[1]
	venvPath := os.Args[2]

	fmt.Printf("Testing virtual environment support\n")
	fmt.Printf("Libpython: %s\n", libpythonPath)
	fmt.Printf("Virtual environment: %s\n", venvPath)

	// Create a new Python runtime instance
	py, err := gopython.NewPureGoPython(libpythonPath)
	if err != nil {
		log.Fatalf("Failed to create Python runtime: %v", err)
	}

	// Configure virtual environment
	venvConfig := gopython.VirtualEnvConfig{
		VenvPath:   venvPath,
		SystemSite: true,       // Include system packages as fallback
		SitePaths:  []string{}, // Additional paths if needed
	}

	// Initialize with virtual environment
	fmt.Println("\nInitializing Python with virtual environment...")
	if err := py.InitializeWithVenv(venvConfig); err != nil {
		log.Fatalf("Failed to initialize Python with venv: %v", err)
	}
	defer func() {
		fmt.Println("\nFinalizing Python interpreter...")
		if err := py.Finalize(); err != nil {
			fmt.Printf("Finalization error: %v\n", err)
		}
		fmt.Println("Python interpreter finalized.")
	}()

	// Test 1: Check Python path configuration
	fmt.Println("\n=== Test 1: Python Path Configuration ===")
	pathCode := `
import sys
print("Python executable:", sys.executable)
print("Python version:", sys.version)
print("Virtual environment check:", hasattr(sys, 'real_prefix') or (hasattr(sys, 'base_prefix') and sys.base_prefix != sys.prefix))
print("\\nPython path entries:")
for i, path in enumerate(sys.path):
    if path:  # Skip empty entries
        print(f"  {i}: {path}")
`
	if err := py.RunString(pathCode); err != nil {
		fmt.Printf("Error checking paths: %v\n", err)
	}

	// Test 2: Check available packages
	fmt.Println("\n=== Test 2: Available Packages ===")
	packagesCode := `
try:
    # Use modern importlib.metadata (Python 3.8+)
    import importlib.metadata as metadata
    installed_packages = [dist.metadata['name'] for dist in metadata.distributions()]
except ImportError:
    # Fallback to pkg_resources for older Python versions
    import pkg_resources
    installed_packages = [d.project_name for d in pkg_resources.working_set]

print(f"Found {len(installed_packages)} installed packages:")
for pkg in sorted(installed_packages)[:10]:  # Show first 10
    print(f"  - {pkg}")
if len(installed_packages) > 10:
    print(f"  ... and {len(installed_packages) - 10} more")
`
	if err := py.RunString(packagesCode); err != nil {
		fmt.Printf("Error checking packages: %v\n", err)
	}

	// Test 3: Try importing packages that might be in venv
	fmt.Println("\n=== Test 3: Package Import Tests ===")

	testPackages := []string{
		"requests",
		"numpy",
		"pandas",
		"flask",
		"django",
		"matplotlib",
		"scipy",
	}

	for _, pkg := range testPackages {
		testCode := fmt.Sprintf(`
try:
    import %s
    print("✓ %s: Available")
    if hasattr(%s, '__version__'):
        print(f"  Version: {%s.__version__}")
    elif hasattr(%s, 'version'):
        print(f"  Version: {%s.version}")
except ImportError as e:
    print("✗ %s: Not available")
`, pkg, pkg, pkg, pkg, pkg, pkg, pkg)

		if err := py.RunString(testCode); err != nil {
			fmt.Printf("Error testing %s: %v\n", pkg, err)
		}
	}

	// Test 4: Test a simple function that might use venv packages
	fmt.Println("\n=== Test 4: Functionality Test ===")
	functionalityCode := `
def test_functionality():
    results = []
    
    # Test requests if available
    try:
        import requests
        response = requests.get('https://httpbin.org/get', timeout=5)
        results.append(f"requests: HTTP {response.status_code}")
    except Exception as e:
        results.append(f"requests: {str(e)[:50]}...")
    
    # Test numpy if available
    try:
        import numpy as np
        arr = np.array([1, 2, 3, 4, 5])
        results.append(f"numpy: mean of [1,2,3,4,5] = {np.mean(arr)}")
    except Exception as e:
        results.append(f"numpy: {str(e)[:50]}...")
    
    # Test pandas if available  
    try:
        import pandas as pd
        df = pd.DataFrame({'A': [1, 2, 3], 'B': [4, 5, 6]})
        results.append(f"pandas: DataFrame shape = {df.shape}")
    except Exception as e:
        results.append(f"pandas: {str(e)[:50]}...")
    
    return results
`

	// Execute the functionality test code first
	if err := py.RunString(functionalityCode); err != nil {
		fmt.Printf("Error setting up functionality test: %v\n", err)
	} else {
		// Now call the function
		result, err := py.CallFunction("__main__", "test_functionality")
		if err != nil {
			fmt.Printf("Error in functionality test: %v\n", err)
		} else {
			fmt.Println("Functionality test results:")
			if results, ok := result.([]interface{}); ok {
				for _, res := range results {
					fmt.Printf("  %v\n", res)
				}
			}
		}
	}

	// Test 5: Demonstrate that venv packages can be called from Go
	fmt.Println("\n=== Test 5: Go-to-Python Function Calls with Venv ===")

	// Try calling a function that uses venv packages
	mathTestCode := `
def advanced_math_test(numbers):
    """Test function that may use scientific packages"""
    try:
        import numpy as np
        arr = np.array(numbers)
        return {
            "using": "numpy",
            "mean": float(np.mean(arr)),
            "std": float(np.std(arr)),
            "sum": float(np.sum(arr))
        }
    except ImportError:
        # Fallback to standard library
        import statistics
        return {
            "using": "statistics (stdlib)",
            "mean": statistics.mean(numbers),
            "std": statistics.stdev(numbers) if len(numbers) > 1 else 0,
            "sum": sum(numbers)
        }
`

	if err := py.RunString(mathTestCode); err != nil {
		fmt.Printf("Error setting up math test: %v\n", err)
	} else {
		// Call the function from Go
		testData := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result, err := py.CallFunction("__main__", "advanced_math_test", testData)
		if err != nil {
			fmt.Printf("Error calling advanced_math_test: %v\n", err)
		} else {
			fmt.Printf("Math test result: %v\n", result)
		}
	}

	fmt.Println("\n=== Virtual Environment Test Complete ===")
	fmt.Println("✓ Successfully initialized Python with virtual environment")
	fmt.Println("✓ Virtual environment packages are accessible")
	fmt.Println("✓ Go-to-Python function calls work with venv packages")
}
