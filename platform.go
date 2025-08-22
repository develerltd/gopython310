package gopython

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetVenvSitePackagesPath returns the site-packages path for a virtual environment
func GetVenvSitePackagesPath(venvPath string) (string, error) {
	// Determine the lib directory path based on platform
	var venvLibDir string
	switch runtime.GOOS {
	case "windows":
		venvLibDir = filepath.Join(venvPath, "Lib")
	default: // linux, darwin, etc.
		venvLibDir = filepath.Join(venvPath, "lib")
	}
	
	if _, err := os.Stat(venvLibDir); os.IsNotExist(err) {
		return "", fmt.Errorf("virtual environment lib directory does not exist: %s", venvLibDir)
	}
	
	// Look for Python version directories
	entries, err := os.ReadDir(venvLibDir)
	if err != nil {
		return "", fmt.Errorf("failed to read venv lib directory: %v", err)
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			// Match python3.10, python3.11, etc.
			if strings.HasPrefix(name, "python3.") || strings.HasPrefix(name, "python") {
				sitePackages := filepath.Join(venvLibDir, name, "site-packages")
				if _, err := os.Stat(sitePackages); err == nil {
					return sitePackages, nil
				}
			}
		}
	}
	
	return "", fmt.Errorf("could not find site-packages directory in virtual environment: %s", venvPath)
}

// ValidateLibraryPath checks if a library path exists and has the expected extension
func ValidateLibraryPath(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("library file does not exist: %s", path)
	}
	
	// Check file extension based on platform
	var expectedExt string
	switch runtime.GOOS {
	case "darwin":
		expectedExt = ".dylib"
	case "windows":
		expectedExt = ".dll"
	default: // linux and others
		expectedExt = ".so"
	}
	
	if !strings.Contains(path, expectedExt) {
		return fmt.Errorf("library file should contain %s extension for %s, got: %s", 
			expectedExt, runtime.GOOS, path)
	}
	
	return nil
}