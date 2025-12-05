// Package files provides file handlers for different file types.
package files

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// DotNetHandler handles .csproj and AssemblyInfo.cs files.
type DotNetHandler struct{}

var (
	csprojVersionRegex   = regexp.MustCompile(`<Version>([^<]+)</Version>`)
	assemblyVersionRegex = regexp.MustCompile(`\[assembly:\s*AssemblyVersion\("([^"]+)"\)\]`)
)

// Type returns the handler type name.
func (h *DotNetHandler) Type() string {
	return "dotnet"
}

// Read reads the current version from a .NET project file.
func (h *DotNetHandler) Read(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}

	// Try .csproj format first
	matches := csprojVersionRegex.FindSubmatch(content)
	if matches != nil {
		return string(matches[1]), nil
	}

	// Try AssemblyInfo.cs format
	matches = assemblyVersionRegex.FindSubmatch(content)
	if matches != nil {
		return string(matches[1]), nil
	}

	return "", fmt.Errorf("version not found in %s", path)
}

// Write writes the new version to a .NET project file.
func (h *DotNetHandler) Write(path, oldVersion, newVersion string, dryRun bool) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	contentStr := string(content)
	var newContent string

	// Determine file type and replace appropriately
	if strings.HasSuffix(path, ".csproj") || csprojVersionRegex.MatchString(contentStr) {
		oldPattern := fmt.Sprintf("<Version>%s</Version>", oldVersion)
		if !strings.Contains(contentStr, oldPattern) {
			return fmt.Errorf("old version %s not found in %s", oldVersion, path)
		}
		newContent = strings.Replace(contentStr, oldPattern, fmt.Sprintf("<Version>%s</Version>", newVersion), 1)
	} else {
		oldPattern := fmt.Sprintf(`[assembly: AssemblyVersion("%s")]`, oldVersion)
		if !strings.Contains(contentStr, oldPattern) {
			return fmt.Errorf("old version %s not found in %s", oldVersion, path)
		}
		newContent = strings.Replace(contentStr, oldPattern, fmt.Sprintf(`[assembly: AssemblyVersion("%s")]`, newVersion), 1)
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
