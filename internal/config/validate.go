// Package config provides configuration loading and validation.
package config

import (
	"fmt"
	"os"
)

// ValidationError represents a configuration validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult contains the results of validating a configuration.
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []string
}

// Validate validates the configuration and returns a ValidationResult.
func Validate(config *Config, strict bool) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []string{},
	}

	// Validate project configuration
	if config.Project.CurrentVersion == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "project.current_version",
			Message: "current_version is required",
		})
		result.Valid = false
	}

	// Validate version scheme
	scheme := config.Project.VersionScheme
	if scheme != "" && scheme != "semver" {
		if len(scheme) < 7 || (scheme[:7] != "calver:" && scheme[:7] != "custom:") {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "project.version_scheme",
				Message: "version_scheme must be 'semver', 'calver:<format>', or 'custom:<template>'",
			})
			result.Valid = false
		}
	}

	// Validate file configurations
	for name, fileConfig := range config.Files {
		// Validate type
		validTypes := map[string]bool{
			"cargo":     true,
			"npm":       true,
			"gomod":     true,
			"pyproject": true,
			"dotnet":    true,
			"regex":     true,
		}
		if !validTypes[fileConfig.Type] {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("files.%s.type", name),
				Message: fmt.Sprintf("invalid file type: %s (valid: cargo, npm, gomod, pyproject, dotnet, regex)", fileConfig.Type),
			})
			result.Valid = false
		}

		// Validate regex type has pattern
		if fileConfig.Type == "regex" && fileConfig.Pattern == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("files.%s.pattern", name),
				Message: "pattern is required for regex type",
			})
			result.Valid = false
		}

		// Check if paths exist
		for _, path := range fileConfig.Paths {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				if strict {
					result.Errors = append(result.Errors, ValidationError{
						Field:   fmt.Sprintf("files.%s.paths", name),
						Message: fmt.Sprintf("file not found: %s", path),
					})
					result.Valid = false
				} else {
					result.Warnings = append(result.Warnings,
						fmt.Sprintf("files.%s: file not found: %s", name, path))
				}
			}
		}
	}

	return result
}
