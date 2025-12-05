// Package files provides file handlers for different file types.
package files

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// RegexHandler handles files using custom regex patterns.
type RegexHandler struct {
	pattern     *regexp.Regexp
	replacement string
}

// NewRegexHandlerWithPatterns creates a new RegexHandler with the specified pattern and replacement.
func NewRegexHandlerWithPatterns(pattern, replacement string) (*RegexHandler, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return &RegexHandler{
		pattern:     re,
		replacement: replacement,
	}, nil
}

// Type returns the handler type name.
func (h *RegexHandler) Type() string {
	return "regex"
}

// Read reads the current version from a file using the regex pattern.
func (h *RegexHandler) Read(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}

	matches := h.pattern.FindSubmatch(content)
	if matches == nil {
		return "", fmt.Errorf("pattern not found in %s", path)
	}

	// Look for a named group called "version"
	for i, name := range h.pattern.SubexpNames() {
		if name == "version" && i < len(matches) {
			return string(matches[i]), nil
		}
	}

	// Fall back to first capturing group
	if len(matches) > 1 {
		return string(matches[1]), nil
	}

	return "", fmt.Errorf("version not captured in pattern for %s", path)
}

// Write writes the new version to a file using the regex pattern.
func (h *RegexHandler) Write(path, oldVersion, newVersion string, dryRun bool) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	contentStr := string(content)

	// Check if the pattern matches
	if !h.pattern.MatchString(contentStr) {
		return fmt.Errorf("pattern not found in %s", path)
	}

	// Replace {version} placeholder in replacement string
	replacementStr := strings.Replace(h.replacement, "{version}", newVersion, -1)

	// Perform the replacement
	newContent := h.pattern.ReplaceAllString(contentStr, replacementStr)

	if contentStr == newContent {
		return fmt.Errorf("no changes made to %s", path)
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
