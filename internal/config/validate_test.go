package config

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		config     *Config
		strict     bool
		wantValid  bool
		wantErrors int
	}{
		{
			name: "valid config",
			config: &Config{
				Project: ProjectConfig{
					CurrentVersion: "1.0.0",
					VersionScheme:  "semver",
				},
				Files: make(map[string]FileConfig),
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "missing current_version",
			config: &Config{
				Project: ProjectConfig{
					VersionScheme: "semver",
				},
				Files: make(map[string]FileConfig),
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "invalid version scheme",
			config: &Config{
				Project: ProjectConfig{
					CurrentVersion: "1.0.0",
					VersionScheme:  "invalid",
				},
				Files: make(map[string]FileConfig),
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "valid calver scheme",
			config: &Config{
				Project: ProjectConfig{
					CurrentVersion: "2024.01.15",
					VersionScheme:  "calver:YYYY.MM.DD",
				},
				Files: make(map[string]FileConfig),
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "invalid file type",
			config: &Config{
				Project: ProjectConfig{
					CurrentVersion: "1.0.0",
					VersionScheme:  "semver",
				},
				Files: map[string]FileConfig{
					"test": {
						Paths: []string{"test.txt"},
						Type:  "invalid",
					},
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "regex without pattern",
			config: &Config{
				Project: ProjectConfig{
					CurrentVersion: "1.0.0",
					VersionScheme:  "semver",
				},
				Files: map[string]FileConfig{
					"test": {
						Paths: []string{"test.txt"},
						Type:  "regex",
					},
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "valid regex config",
			config: &Config{
				Project: ProjectConfig{
					CurrentVersion: "1.0.0",
					VersionScheme:  "semver",
				},
				Files: map[string]FileConfig{
					"test": {
						Paths:       []string{"test.txt"},
						Type:        "regex",
						Pattern:     `version=(?P<version>.+)`,
						Replacement: "version={version}",
					},
				},
			},
			wantValid:  true,
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.config, tt.strict)
			if result.Valid != tt.wantValid {
				t.Errorf("Validate() valid = %v, want %v", result.Valid, tt.wantValid)
			}
			if len(result.Errors) != tt.wantErrors {
				t.Errorf("Validate() errors = %d, want %d", len(result.Errors), tt.wantErrors)
				for _, e := range result.Errors {
					t.Logf("  Error: %s: %s", e.Field, e.Message)
				}
			}
		})
	}
}
