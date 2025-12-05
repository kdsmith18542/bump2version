package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoader(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "bumpx-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test config file
	configPath := filepath.Join(tmpDir, ".bumpx.toml")
	configContent := `
[project]
name = "test-project"
current_version = "1.2.3"
version_scheme = "semver"

[files.cargo]
paths = ["Cargo.toml"]
type = "cargo"

[git]
enable = true
commit = true
commit_message = "Release {version}"
tag = true
tag_name = "v{version}"
tag_message = "Version {version}"
allow_dirty = false

[hooks]
pre_bump = ""
post_bump = ""
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Test loading
	loader := NewLoader(configPath)
	cfg, err := loader.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify config values
	if cfg.Project.Name != "test-project" {
		t.Errorf("Project.Name = %s, want test-project", cfg.Project.Name)
	}
	if cfg.Project.CurrentVersion != "1.2.3" {
		t.Errorf("Project.CurrentVersion = %s, want 1.2.3", cfg.Project.CurrentVersion)
	}
	if cfg.Project.VersionScheme != "semver" {
		t.Errorf("Project.VersionScheme = %s, want semver", cfg.Project.VersionScheme)
	}

	// Check file config
	cargoConfig, ok := cfg.Files["cargo"]
	if !ok {
		t.Error("Missing 'cargo' file config")
	} else {
		if len(cargoConfig.Paths) != 1 || cargoConfig.Paths[0] != "Cargo.toml" {
			t.Errorf("cargo.Paths = %v, want [Cargo.toml]", cargoConfig.Paths)
		}
		if cargoConfig.Type != "cargo" {
			t.Errorf("cargo.Type = %s, want cargo", cargoConfig.Type)
		}
	}

	// Check git config
	if !cfg.Git.Enable {
		t.Error("Git.Enable = false, want true")
	}
	if cfg.Git.CommitMessage != "Release {version}" {
		t.Errorf("Git.CommitMessage = %s, want 'Release {version}'", cfg.Git.CommitMessage)
	}
}

func TestLoaderMissingFile(t *testing.T) {
	loader := NewLoader("/nonexistent/path/.bumpx.toml")
	_, err := loader.Load()
	if err == nil {
		t.Error("Expected error for missing config file")
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "bumpx-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, ".bumpx.toml")

	cfg := NewDefaultConfig()
	cfg.Project.Name = "saved-project"
	cfg.Project.CurrentVersion = "2.0.0"

	if err := SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Load and verify
	loader := NewLoader(configPath)
	loaded, err := loader.Load()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loaded.Project.Name != "saved-project" {
		t.Errorf("Project.Name = %s, want saved-project", loaded.Project.Name)
	}
	if loaded.Project.CurrentVersion != "2.0.0" {
		t.Errorf("Project.CurrentVersion = %s, want 2.0.0", loaded.Project.CurrentVersion)
	}
}
