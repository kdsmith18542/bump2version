// Package core provides the main orchestration logic for bumpx.
package core

import (
	"fmt"
	"strings"

	"github.com/kdsmith18542/bumpx/internal/config"
	"github.com/kdsmith18542/bumpx/internal/files"
	"github.com/kdsmith18542/bumpx/internal/git"
	"github.com/kdsmith18542/bumpx/internal/hooks"
	"github.com/kdsmith18542/bumpx/internal/output"
	"github.com/kdsmith18542/bumpx/internal/version"
)

// BumpOptions contains options for the bump operation.
type BumpOptions struct {
	DryRun    bool
	NoGit     bool
	GitCommit bool // explicit flag to enable commit (overrides config)
	GitTag    bool // explicit flag to enable tag (overrides config)
	NoTag     bool
	NoCommit  bool
	Message   string
	TagName   string
	JSON      bool
}

// Bumper handles version bumping operations.
type Bumper struct {
	config     *config.Config
	configPath string
	printer    *output.Printer
	git        *git.Git
	hooks      *hooks.Runner
}

// NewBumper creates a new Bumper instance.
func NewBumper(cfg *config.Config, configPath string, jsonMode bool) *Bumper {
	return &Bumper{
		config:     cfg,
		configPath: configPath,
		printer:    output.NewPrinter(jsonMode),
		git:        git.New(""),
		hooks:      hooks.New(""),
	}
}

// Bump performs the version bump operation.
func (b *Bumper) Bump(part string, opts BumpOptions) (*output.BumpResult, error) {
	// Parse current version (discover from files if not set)
	currentVersion := b.config.Project.CurrentVersion
	if currentVersion == "" {
		discovered, err := b.discoverVersion()
		if err != nil {
			return nil, fmt.Errorf("no current_version in config and failed to discover: %w", err)
		}
		currentVersion = discovered
	}
	scheme := b.config.GetVersionScheme()

	ver, err := version.Parse(currentVersion, scheme)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current version: %w", err)
	}

	// Bump version
	newVer, err := ver.Bump(part)
	if err != nil {
		return nil, fmt.Errorf("failed to bump version: %w", err)
	}

	newVersion := newVer.String()

	// Run pre-bump hook
	hookCtx := &hooks.HookContext{
		OldVersion: currentVersion,
		NewVersion: newVersion,
		Part:       part,
		ConfigPath: b.configPath,
	}

	if b.config.Hooks.PreBump != "" && !opts.DryRun {
		if err := b.hooks.RunPreBump(b.config.Hooks.PreBump, hookCtx); err != nil {
			return nil, err
		}
	}

	// Check git status
	if b.config.Git.Enable && !opts.NoGit && !b.config.Git.AllowDirty {
		clean, err := b.git.IsClean()
		if err != nil {
			return nil, fmt.Errorf("failed to check git status: %w", err)
		}
		if !clean {
			return nil, fmt.Errorf("working directory is dirty (use --allow-dirty to override)")
		}
	}

	result := &output.BumpResult{
		OldVersion:   currentVersion,
		NewVersion:   newVersion,
		Part:         part,
		ChangedFiles: []string{},
	}

	// Update files from simple paths list (auto-detect handler)
	for _, path := range b.config.Files.Paths {
		handler, err := files.DetectHandler(path)
		if err != nil {
			return nil, fmt.Errorf("failed to detect handler for %s: %w", path, err)
		}

		if opts.DryRun {
			b.printer.Println("Would update %s", path)
		} else {
			b.printer.Println("Updating %s", path)
		}

		if err := handler.Write(path, currentVersion, newVersion, opts.DryRun); err != nil {
			return nil, fmt.Errorf("failed to update %s: %w", path, err)
		}

		result.ChangedFiles = append(result.ChangedFiles, path)
	}

	// Update files from detailed handlers
	for name, handlerCfg := range b.config.Handlers {
		handler, err := b.getHandler(handlerCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to get handler for %s: %w", name, err)
		}

		for _, path := range handlerCfg.Paths {
			if opts.DryRun {
				b.printer.Println("Would update %s", path)
			} else {
				b.printer.Println("Updating %s", path)
			}

			if err := handler.Write(path, currentVersion, newVersion, opts.DryRun); err != nil {
				return nil, fmt.Errorf("failed to update %s: %w", path, err)
			}

			result.ChangedFiles = append(result.ChangedFiles, path)
		}
	}

	// Update config file
	if !opts.DryRun {
		b.config.Project.CurrentVersion = newVersion
		if err := config.SaveConfig(b.config, b.configPath); err != nil {
			return nil, fmt.Errorf("failed to update config file: %w", err)
		}
		result.ChangedFiles = append(result.ChangedFiles, b.configPath)
	}

	// Git operations
	// Run git ops if: (config enabled OR explicit flag) AND not disabled
	enableGit := (b.config.Git.Enable || opts.GitCommit || opts.GitTag) && !opts.NoGit
	if enableGit && !opts.DryRun {
		// Stage files
		if err := b.git.Add(result.ChangedFiles...); err != nil {
			return nil, fmt.Errorf("failed to stage files: %w", err)
		}

		// Commit: enabled by config or explicit flag, not disabled
		doCommit := (b.config.Git.Commit || opts.GitCommit) && !opts.NoCommit
		if doCommit {
			message := b.config.Git.CommitMessage
			if opts.Message != "" {
				message = opts.Message
			}
			message = b.replaceVersionPlaceholders(message, currentVersion, newVersion)

			if err := b.git.Commit(message); err != nil {
				return nil, fmt.Errorf("failed to commit: %w", err)
			}

			hash, _ := b.git.GetShortCommitHash()
			result.GitCommit = hash
		}

		// Tag: enabled by config or explicit flag, not disabled
		doTag := (b.config.Git.Tag || opts.GitTag) && !opts.NoTag
		if doTag {
			tagName := b.config.Git.TagName
			if opts.TagName != "" {
				tagName = opts.TagName
			}
			tagName = b.replaceVersionPlaceholders(tagName, currentVersion, newVersion)

			tagMessage := b.replaceVersionPlaceholders(b.config.Git.TagMessage, currentVersion, newVersion)

			if err := b.git.Tag(tagName, tagMessage); err != nil {
				return nil, fmt.Errorf("failed to create tag: %w", err)
			}

			result.GitTag = tagName
		}
	}

	// Run post-bump hook
	if b.config.Hooks.PostBump != "" && !opts.DryRun {
		if err := b.hooks.RunPostBump(b.config.Hooks.PostBump, hookCtx); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// getHandler returns the appropriate file handler.
func (b *Bumper) getHandler(cfg config.HandlerConfig) (files.Handler, error) {
	if cfg.Type == "regex" {
		return files.NewRegexHandler(cfg.Pattern, cfg.Replacement)
	}
	return files.GetHandler(cfg.Type)
}

// replaceVersionPlaceholders replaces {version} and other placeholders in templates.
func (b *Bumper) replaceVersionPlaceholders(template, oldVersion, newVersion string) string {
	result := template
	result = strings.Replace(result, "{version}", newVersion, -1)
	result = strings.Replace(result, "{new_version}", newVersion, -1)
	result = strings.Replace(result, "{old_version}", oldVersion, -1)
	result = strings.Replace(result, "{current_version}", oldVersion, -1)
	return result
}

// discoverVersion attempts to read version from configured files.
func (b *Bumper) discoverVersion() (string, error) {
	// Try simple paths first
	for _, path := range b.config.Files.Paths {
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
	for _, handlerCfg := range b.config.Handlers {
		handler, err := b.getHandler(handlerCfg)
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
