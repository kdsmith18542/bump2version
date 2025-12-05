package core

import (
	"testing"

	"github.com/kdsmith18542/bumpx/internal/config"
)

func TestShow(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		wantVersion string
		wantErr     bool
	}{
		{
			name: "returns current version from config",
			config: &config.Config{
				Project: config.ProjectConfig{
					CurrentVersion: "1.2.3",
					VersionScheme:  "semver",
				},
			},
			wantVersion: "1.2.3",
			wantErr:     false,
		},
		{
			name: "returns different version",
			config: &config.Config{
				Project: config.ProjectConfig{
					CurrentVersion: "2.0.0-beta.1",
					VersionScheme:  "semver",
				},
			},
			wantVersion: "2.0.0-beta.1",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Show(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Show() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != nil && result.Version != tt.wantVersion {
				t.Errorf("Show() version = %v, want %v", result.Version, tt.wantVersion)
			}
		})
	}
}

func TestShow_NoVersionConfigured(t *testing.T) {
	cfg := &config.Config{
		Project: config.ProjectConfig{
			CurrentVersion: "",
			VersionScheme:  "semver",
		},
		Files:    config.FilesConfig{Paths: []string{}},
		Handlers: make(map[string]config.HandlerConfig),
	}

	_, err := Show(cfg)
	if err == nil {
		t.Error("Expected error when no version configured and no files to discover from")
	}
}

func TestNewBumper(t *testing.T) {
	cfg := &config.Config{
		Project: config.ProjectConfig{
			CurrentVersion: "1.0.0",
			VersionScheme:  "semver",
		},
	}

	bumper := NewBumper(cfg, ".bumpx.toml", false)
	if bumper == nil {
		t.Error("NewBumper() returned nil")
	}
	if bumper.config != cfg {
		t.Error("NewBumper() config not set correctly")
	}
	if bumper.configPath != ".bumpx.toml" {
		t.Errorf("NewBumper() configPath = %v, want .bumpx.toml", bumper.configPath)
	}
}

func TestBumper_replaceVersionPlaceholders(t *testing.T) {
	cfg := &config.Config{}
	bumper := NewBumper(cfg, "", false)

	tests := []struct {
		name       string
		template   string
		oldVersion string
		newVersion string
		want       string
	}{
		{
			name:       "replace version",
			template:   "Release {version}",
			oldVersion: "1.0.0",
			newVersion: "1.0.1",
			want:       "Release 1.0.1",
		},
		{
			name:       "replace new_version",
			template:   "Bump to {new_version}",
			oldVersion: "1.0.0",
			newVersion: "2.0.0",
			want:       "Bump to 2.0.0",
		},
		{
			name:       "replace old_version",
			template:   "Changed from {old_version}",
			oldVersion: "1.0.0",
			newVersion: "2.0.0",
			want:       "Changed from 1.0.0",
		},
		{
			name:       "replace multiple placeholders",
			template:   "Bump {old_version} -> {new_version}",
			oldVersion: "1.0.0",
			newVersion: "1.0.1",
			want:       "Bump 1.0.0 -> 1.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bumper.replaceVersionPlaceholders(tt.template, tt.oldVersion, tt.newVersion)
			if got != tt.want {
				t.Errorf("replaceVersionPlaceholders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHandlerForConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.HandlerConfig
		wantErr bool
	}{
		{
			name:    "npm handler",
			cfg:     config.HandlerConfig{Type: "npm"},
			wantErr: false,
		},
		{
			name:    "cargo handler",
			cfg:     config.HandlerConfig{Type: "cargo"},
			wantErr: false,
		},
		{
			name: "regex handler",
			cfg: config.HandlerConfig{
				Type:        "regex",
				Pattern:     "version = \"(.+)\"",
				Replacement: "version = \"{version}\"",
			},
			wantErr: false,
		},
		{
			name:    "invalid handler",
			cfg:     config.HandlerConfig{Type: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getHandlerForConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHandlerForConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
