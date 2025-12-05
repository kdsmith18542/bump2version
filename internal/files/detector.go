// Package files provides file handlers for different file types.
package files

import (
	"fmt"
)

// Handler is the interface that file handlers must implement.
type Handler interface {
	// Type returns the handler type name.
	Type() string
	// Read reads the current version from the file.
	Read(path string) (string, error)
	// Write writes the new version to the file.
	Write(path, oldVersion, newVersion string, dryRun bool) error
}

// GetHandler returns the appropriate handler for a file type.
func GetHandler(fileType string) (Handler, error) {
	switch fileType {
	case "cargo":
		return &CargoHandler{}, nil
	case "npm":
		return &NPMHandler{}, nil
	case "gomod":
		return &GoModHandler{}, nil
	case "pyproject":
		return &PyProjectHandler{}, nil
	case "dotnet":
		return &DotNetHandler{}, nil
	case "regex":
		return nil, fmt.Errorf("regex handler requires pattern and replacement; use NewRegexHandlerWithPatterns")
	default:
		return nil, fmt.Errorf("unknown file type: %s", fileType)
	}
}

// NewRegexHandler creates a new regex handler with the specified pattern and replacement.
func NewRegexHandler(pattern, replacement string) (*RegexHandler, error) {
	return NewRegexHandlerWithPatterns(pattern, replacement)
}

// DetectFileType attempts to detect the file type based on the filename.
func DetectFileType(path string) string {
	switch {
	case matchFilename(path, "Cargo.toml"):
		return "cargo"
	case matchFilename(path, "package.json"):
		return "npm"
	case matchFilename(path, "go.mod"):
		return "gomod"
	case matchFilename(path, "pyproject.toml"):
		return "pyproject"
	case matchFilename(path, "setup.cfg"):
		return "pyproject"
	case matchFilename(path, ".csproj"):
		return "dotnet"
	case matchFilename(path, "AssemblyInfo.cs"):
		return "dotnet"
	default:
		return "regex"
	}
}

// matchFilename checks if the path ends with the given filename or suffix.
func matchFilename(path, filename string) bool {
	if len(path) >= len(filename) {
		return path[len(path)-len(filename):] == filename
	}
	return false
}

// DetectHandler auto-detects and returns the appropriate handler for a file path.
func DetectHandler(path string) (Handler, error) {
	fileType := DetectFileType(path)
	if fileType == "regex" {
		// For unknown file types, use a generic version pattern
		return NewRegexHandlerWithPatterns(`"version"\s*:\s*"([^"]+)"`, `"version": "{version}"`)
	}
	return GetHandler(fileType)
}
