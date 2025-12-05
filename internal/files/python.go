// Package files provides file handlers for different file types.
package files

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// PyProjectHandler handles pyproject.toml files.
type PyProjectHandler struct{}

var pyprojectVersionRegex = regexp.MustCompile(`(?m)^version\s*=\s*"([^"]+)"`)

// Type returns the handler type name.
func (h *PyProjectHandler) Type() string {
	return "pyproject"
}

// Read reads the current version from a pyproject.toml file.
func (h *PyProjectHandler) Read(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}

	matches := pyprojectVersionRegex.FindSubmatch(content)
	if matches == nil {
		return "", fmt.Errorf("version not found in %s", path)
	}

	return string(matches[1]), nil
}

// Write writes the new version to a pyproject.toml file.
func (h *PyProjectHandler) Write(path, oldVersion, newVersion string, dryRun bool) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	// Replace the version
	newContent := pyprojectVersionRegex.ReplaceAllString(
		string(content),
		fmt.Sprintf(`version = "%s"`, newVersion),
	)

	if !strings.Contains(string(content), fmt.Sprintf(`version = "%s"`, oldVersion)) {
		return fmt.Errorf("old version %s not found in %s", oldVersion, path)
	}

	if dryRun {
		return nil
	}

	// Get original file permissions
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", path, err)
	}

	return os.WriteFile(path, []byte(newContent), info.Mode())
}
