// Package config provides configuration loading and validation.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// DefaultConfigFileName is the default name for the config file.
const DefaultConfigFileName = ".bumpx.toml"

// Loader handles loading configuration from files.
type Loader struct {
	configPath string
}

// NewLoader creates a new configuration loader.
func NewLoader(configPath string) *Loader {
	if configPath == "" {
		configPath = DefaultConfigFileName
	}
	return &Loader{configPath: configPath}
}

// Load loads the configuration from the config file.
func (l *Loader) Load() (*Config, error) {
	config := NewDefaultConfig()

	// Check if config file exists
	if _, err := os.Stat(l.configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", l.configPath)
	}

	// Read and parse the config file
	if _, err := toml.DecodeFile(l.configPath, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// ConfigPath returns the path to the config file.
func (l *Loader) ConfigPath() string {
	return l.configPath
}

// FindConfigFile searches for a config file in the current directory
// and parent directories.
func FindConfigFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		configPath := filepath.Join(dir, DefaultConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("config file not found: %s", DefaultConfigFileName)
}

// SaveConfig saves the configuration to a file.
func SaveConfig(config *Config, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
