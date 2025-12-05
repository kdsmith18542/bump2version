// Package files provides file handlers for different file types.
package files

import (
	"fmt"
	"os"
	"regexp"
)

// GoModHandler handles go.mod files.
type GoModHandler struct{}

var goModVersionRegex = regexp.MustCompile(`(?m)^// version: (.+)$`)

// Type returns the handler type name.
func (h *GoModHandler) Type() string {
	return "gomod"
}

// Read reads the current version from a go.mod file.
// Note: go.mod doesn't have a standard version field, so we use a comment.
// Format: // version: x.y.z
func (h *GoModHandler) Read(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}

	matches := goModVersionRegex.FindSubmatch(content)
	if matches == nil {
		return "", fmt.Errorf("version comment not found in %s (expected format: // version: x.y.z)", path)
	}

	return string(matches[1]), nil
}

// Write writes the new version to a go.mod file.
func (h *GoModHandler) Write(path, oldVersion, newVersion string, dryRun bool) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	oldPattern := fmt.Sprintf("// version: %s", oldVersion)
	newPattern := fmt.Sprintf("// version: %s", newVersion)

	contentStr := string(content)
	if !regexp.MustCompile(regexp.QuoteMeta(oldPattern)).MatchString(contentStr) {
		return fmt.Errorf("old version comment not found in %s (expected: %s)", path, oldPattern)
	}

	if dryRun {
		return nil
	}

	newContent := regexp.MustCompile(regexp.QuoteMeta(oldPattern)).ReplaceAllString(contentStr, newPattern)

	// Get original file permissions
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", path, err)
	}

	return os.WriteFile(path, []byte(newContent), info.Mode())
}
