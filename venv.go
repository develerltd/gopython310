package gopython

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// InitializeWithVenv initializes the Python interpreter with virtual environment support
func (py *PureGoPython) InitializeWithVenv(config VirtualEnvConfig) error {
	if py.pyInitialize == nil {
		return errors.New("Python functions not registered")
	}

	// Validate and configure virtual environment before initialization
	if err := py.configureVirtualEnvironment(config); err != nil {
		return fmt.Errorf("virtual environment configuration failed: %v", err)
	}

	// Initialize Python interpreter
	py.pyInitialize()

	// Configure virtual environment paths after initialization
	if err := py.addSiteDirectories(config); err != nil {
		return fmt.Errorf("failed to configure virtual environment paths: %v", err)
	}

	return nil
}

// configureVirtualEnvironment validates the virtual environment exists
func (py *PureGoPython) configureVirtualEnvironment(config VirtualEnvConfig) error {
	if config.VenvPath == "" {
		return errors.New("virtual environment path cannot be empty")
	}

	// Check if virtual environment exists
	if _, err := os.Stat(config.VenvPath); os.IsNotExist(err) {
		return fmt.Errorf("virtual environment does not exist: %s", config.VenvPath)
	}

	// Validate that it looks like a proper venv
	venvLibDir := filepath.Join(config.VenvPath, "lib")
	if _, err := os.Stat(venvLibDir); os.IsNotExist(err) {
		return fmt.Errorf("invalid virtual environment: missing lib directory in %s", config.VenvPath)
	}

	// All path configuration will be done after initialization using site.addsitedir()
	// This avoids the Unicode encoding issues with Py_SetPath()
	return nil
}

// addSiteDirectories adds additional site directories after initialization
func (py *PureGoPython) addSiteDirectories(config VirtualEnvConfig) error {
	if len(config.SitePaths) == 0 && config.VenvPath == "" {
		return nil
	}

	// Import required modules
	siteCode := "import sys\nimport os\n"

	// Configure virtual environment properly
	if config.VenvPath != "" {
		venvLibDir := filepath.Join(config.VenvPath, "lib")
		var venvSitePackages string
		
		if entries, err := os.ReadDir(venvLibDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() && (entry.Name() == "python3.10" || entry.Name()[:6] == "python") {
					sitePackages := filepath.Join(venvLibDir, entry.Name(), "site-packages")
					if _, err := os.Stat(sitePackages); err == nil {
						venvSitePackages = sitePackages
						break
					}
				}
			}
		}
		
		if venvSitePackages != "" {
			// Set VIRTUAL_ENV environment variable for proper venv detection
			siteCode += fmt.Sprintf("os.environ['VIRTUAL_ENV'] = r'%s'\n", config.VenvPath)
			
			// Clean sys.path to only include essential paths
			siteCode += fmt.Sprintf("venv_site_packages = r'%s'\n", venvSitePackages)
			siteCode += `
# Save essential Python paths (stdlib only)
essential_paths = []
for path in sys.path:
    # Keep only essential Python standard library paths
    if (path.endswith('python310.zip') or 
        path.endswith('python3.10') or 
        path.endswith('lib-dynload') or
        path == ''):  # Empty string is current directory
        essential_paths.append(path)

# Replace sys.path with clean virtual environment setup
sys.path = [venv_site_packages] + essential_paths
`
			
			// Optionally add system site packages if SystemSite is True
			if config.SystemSite {
				siteCode += `
# Add system site packages as fallback (SystemSite=True)
import site
try:
    system_site_packages = site.getsitepackages()
    for path in system_site_packages:
        if path not in sys.path:
            sys.path.append(path)
except:
    pass  # Ignore if getsitepackages() fails
`
			}
		}
	}

	// Add custom site paths to the beginning as well
	for _, path := range config.SitePaths {
		siteCode += fmt.Sprintf("custom_path = r'%s'\n", path)
		siteCode += "if custom_path not in sys.path:\n"
		siteCode += "    sys.path.insert(0, custom_path)\n"
	}

	// Execute the site configuration
	return py.withGIL(func() error {
		cCode := stringToCString(siteCode)
		result := py.pyRunSimpleString(cCode)
		if result != 0 {
			return fmt.Errorf("failed to configure site directories")
		}
		return nil
	})
}