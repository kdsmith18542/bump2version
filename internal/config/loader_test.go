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

[files]
paths = ["package.json"]

[handlers.cargo]
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

	// Check files config
	if len(cfg.Files.Paths) != 1 || cfg.Files.Paths[0] != "package.json" {
		t.Errorf("Files.Paths = %v, want [package.json]", cfg.Files.Paths)
	}

	// Check handler config
	cargoConfig, ok := cfg.Handlers["cargo"]
	if !ok {
		t.Error("Missing 'cargo' handler config")
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

func TestFindConfigFile(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "bumpx-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "sub", "nested")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub dir: %v", err)
	}

	// Create config file in the root temp dir
	configPath := filepath.Join(tmpDir, DefaultConfigFileName)
	if err := os.WriteFile(configPath, []byte("[project]\nname = \"test\"\ncurrent_version = \"1.0.0\"\n"), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Save current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	// Change to subdirectory
	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// FindConfigFile should find the config in parent directory
	found, err := FindConfigFile()
	if err != nil {
		t.Fatalf("FindConfigFile() error = %v", err)
	}

	if found != configPath {
		t.Errorf("FindConfigFile() = %v, want %v", found, configPath)
	}
}

func TestFindConfigFileNotFound(t *testing.T) {
	// Create a temporary directory without a config file
	tmpDir, err := os.MkdirTemp("", "bumpx-test-noconfig")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	// Change to temporary directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// FindConfigFile should return an error
	_, err = FindConfigFile()
	if err == nil {
		t.Error("Expected error when config file not found")
	}
}

func TestNewLoader_DefaultPath(t *testing.T) {
	loader := NewLoader("")
	if loader.configPath != DefaultConfigFileName {
		t.Errorf("NewLoader(\"\").configPath = %v, want %v", loader.configPath, DefaultConfigFileName)
	}
}

func TestLoader_ConfigPath(t *testing.T) {
	loader := NewLoader("/custom/path/.bumpx.toml")
	if loader.ConfigPath() != "/custom/path/.bumpx.toml" {
		t.Errorf("ConfigPath() = %v, want /custom/path/.bumpx.toml", loader.ConfigPath())
	}
}
