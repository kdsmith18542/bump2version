// Package config provides configuration loading and validation.
package config

import (
	"github.com/kdsmith18542/bumpx/internal/version"
)

// Config represents the main configuration for bumpx.
type Config struct {
	Project  ProjectConfig            `toml:"project"`
	Files    FilesConfig              `toml:"files"`
	Handlers map[string]HandlerConfig `toml:"handlers,omitempty"` // optional detailed handlers
	Git      GitConfig                `toml:"git"`
	Hooks    HooksConfig              `toml:"hooks"`
}

// FilesConfig represents the files section per spec.
type FilesConfig struct {
	Paths []string `toml:"paths"`
}

// ProjectConfig represents the project-level configuration.
type ProjectConfig struct {
	Name           string `toml:"name"`
	CurrentVersion string `toml:"current_version"`
	VersionScheme  string `toml:"version_scheme"` // "semver" | "calver:YYYY.MM.DD" | "custom:<template>"
}

// HandlerConfig represents the configuration for a specific file handler.
type HandlerConfig struct {
	Paths       []string `toml:"paths"`
	Type        string   `toml:"type"`                  // "cargo", "npm", "gomod", "pyproject", "dotnet", "regex"
	Pattern     string   `toml:"pattern,omitempty"`     // For regex type
	Replacement string   `toml:"replacement,omitempty"` // For regex type
}

// GitConfig represents the git integration configuration.
type GitConfig struct {
	Enable        bool   `toml:"enable"`
	Commit        bool   `toml:"commit"`
	CommitMessage string `toml:"commit_message"`
	Tag           bool   `toml:"tag"`
	TagName       string `toml:"tag_name"`
	TagMessage    string `toml:"tag_message"`
	AllowDirty    bool   `toml:"allow_dirty"`
}

// HooksConfig represents the hooks configuration.
type HooksConfig struct {
	PreBump  string `toml:"pre_bump"`
	PostBump string `toml:"post_bump"`
}

// GetVersionScheme returns the parsed version scheme.
func (c *Config) GetVersionScheme() version.Scheme {
	switch c.Project.VersionScheme {
	case "semver", "":
		return version.SchemeSemVer
	default:
		if len(c.Project.VersionScheme) > 7 && c.Project.VersionScheme[:7] == "calver:" {
			return version.SchemeCalVer
		}
		return version.SchemeSemVer
	}
}

// NewDefaultConfig returns a new configuration with default values.
func NewDefaultConfig() *Config {
	return &Config{
		Project: ProjectConfig{
			VersionScheme: "semver",
		},
		Files: FilesConfig{
			Paths: []string{},
		},
		Handlers: make(map[string]HandlerConfig),
		Git: GitConfig{
			Enable:        false,
			Commit:        false, // default per spec: no side-effects
			CommitMessage: "chore(release): {version}",
			Tag:           false, // default per spec: no side-effects
			TagName:       "v{version}",
			TagMessage:    "Release {version}",
			AllowDirty:    false,
		},
		Hooks: HooksConfig{},
	}
}
