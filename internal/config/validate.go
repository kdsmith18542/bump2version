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

	// Validate files.paths
	for _, path := range config.Files.Paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if strict {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "files.paths",
					Message: fmt.Sprintf("file not found: %s", path),
				})
				result.Valid = false
			} else {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("files.paths: file not found: %s", path))
			}
		}
	}

	// Validate handler configurations
	validTypes := map[string]bool{
		"cargo":     true,
		"npm":       true,
		"gomod":     true,
		"pyproject": true,
		"dotnet":    true,
		"regex":     true,
	}

	for name, handlerConfig := range config.Handlers {
		// Validate type
		if !validTypes[handlerConfig.Type] {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("handlers.%s.type", name),
				Message: fmt.Sprintf("invalid file type: %s (valid: cargo, npm, gomod, pyproject, dotnet, regex)", handlerConfig.Type),
			})
			result.Valid = false
		}

		// Validate regex type has pattern
		if handlerConfig.Type == "regex" && handlerConfig.Pattern == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("handlers.%s.pattern", name),
				Message: "pattern is required for regex type",
			})
			result.Valid = false
		}

		// Check if paths exist
		for _, path := range handlerConfig.Paths {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				if strict {
					result.Errors = append(result.Errors, ValidationError{
						Field:   fmt.Sprintf("handlers.%s.paths", name),
						Message: fmt.Sprintf("file not found: %s", path),
					})
					result.Valid = false
				} else {
					result.Warnings = append(result.Warnings,
						fmt.Sprintf("handlers.%s: file not found: %s", name, path))
				}
			}
		}
	}

	return result
}
