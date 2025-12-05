// Package core provides functionality to show the current project version.
package core

import (
	"fmt"

	"github.com/kdsmith18542/bumpx/internal/config"
	"github.com/kdsmith18542/bumpx/internal/files"
	"github.com/kdsmith18542/bumpx/internal/output"
)

// Show returns the current project version.
func Show(cfg *config.Config) (*output.ShowResult, error) {
	version := cfg.Project.CurrentVersion

	// Discover version from files if not set in config
	if version == "" {
		discovered, err := discoverVersionFromFiles(cfg)
		if err != nil {
			return nil, fmt.Errorf("no current_version in config and failed to discover: %w", err)
		}
		version = discovered
	}

	return &output.ShowResult{
		Version: version,
	}, nil
}

// discoverVersionFromFiles attempts to read version from configured files.
func discoverVersionFromFiles(cfg *config.Config) (string, error) {
	// Try simple paths first
	for _, path := range cfg.Files.Paths {
		handler, err := files.DetectHandler(path)
		if err != nil {
			continue
		}
		ver, err := handler.Read(path)
		if err == nil && ver != "" {
			return ver, nil
		}
	}

	// Try handlers
	for _, handlerCfg := range cfg.Handlers {
		handler, err := getHandlerForConfig(handlerCfg)
		if err != nil {
			continue
		}
		for _, path := range handlerCfg.Paths {
			ver, err := handler.Read(path)
			if err == nil && ver != "" {
				return ver, nil
			}
		}
	}

	return "", fmt.Errorf("no version found in configured files")
}

// getHandlerForConfig returns the appropriate file handler for a handler config.
func getHandlerForConfig(cfg config.HandlerConfig) (files.Handler, error) {
	if cfg.Type == "regex" {
		return files.NewRegexHandler(cfg.Pattern, cfg.Replacement)
	}
	return files.GetHandler(cfg.Type)
}
