// Package files provides file handlers for different file types.
package files

import (
	"encoding/json"
	"fmt"
	"os"
)

// NPMHandler handles package.json files.
type NPMHandler struct{}

// Type returns the handler type name.
func (h *NPMHandler) Type() string {
	return "npm"
}

// Read reads the current version from a package.json file.
func (h *NPMHandler) Read(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", path, err)
	}

	version, ok := pkg["version"].(string)
	if !ok {
		return "", fmt.Errorf("version not found or not a string in %s", path)
	}

	return version, nil
}

// Write writes the new version to a package.json file.
func (h *NPMHandler) Write(path, oldVersion, newVersion string, dryRun bool) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	currentVersion, ok := pkg["version"].(string)
	if !ok || currentVersion != oldVersion {
		return fmt.Errorf("expected version %s but found %v in %s", oldVersion, pkg["version"], path)
	}

	if dryRun {
		return nil
	}

	pkg["version"] = newVersion

	newContent, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Add trailing newline
	newContent = append(newContent, '\n')

	// Get original file permissions
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", path, err)
	}

	return os.WriteFile(path, newContent, info.Mode())
}
