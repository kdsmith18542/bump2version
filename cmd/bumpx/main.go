// bumpx is a CLI tool for version management and release chores.
package main

import (
	"fmt"
	"os"

	"github.com/kdsmith18542/bumpx/internal/config"
	"github.com/kdsmith18542/bumpx/internal/core"
	"github.com/kdsmith18542/bumpx/internal/output"
	"github.com/spf13/cobra"
)

var (
	// Version is set during build
	Version = "0.1.0"
)

// Exit codes per specification
const (
	ExitSuccess        = 0 // success
	ExitUsageError     = 1 // usage or validation error
	ExitOperationError = 2 // operation failed (file I/O, parse errors)
)

var (
	cfgFile      string
	outputFormat string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		// Cobra already handles usage errors with exit code 1
		os.Exit(ExitUsageError)
	}
}

var rootCmd = &cobra.Command{
	Use:     "bumpx",
	Short:   "A version management tool for polyglot repositories",
	Long:    `bumpx is a single static CLI binary for version management and release chores across polyglot repositories.`,
	Version: Version,
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current project version",
	Long:  `Show the current project version as understood from the config file.`,
	RunE:  runShow,
}

var bumpCmd = &cobra.Command{
	Use:   "bump <part>",
	Short: "Bump the version",
	Long: `Bump the version and update all configured files.

For SemVer: major, minor, patch, pre, build
For CalVer: date, build`,
	Args: cobra.ExactArgs(1),
	RunE: runBump,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new bumpx configuration",
	Long:  `Create a starter .bumpx.toml in the current directory.`,
	RunE:  runInit,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	Long:  `Validate .bumpx.toml and report any errors or warnings.`,
	RunE:  runValidate,
}

// Bump command flags
var (
	dryRun     bool
	noGit      bool
	gitCommit  bool
	gitTag     bool
	noTag      bool
	noCommit   bool
	allowDirty bool
	message    string
	tagName    string
)

// Init command flags
var (
	initYes    bool
	initScheme string
)

// Validate command flags
var (
	validateStrict bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is .bumpx.toml)")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "format", "human", "output format: human or json")

	// Show command
	rootCmd.AddCommand(showCmd)

	// Bump command
	bumpCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "don't write any files, just pretend")
	bumpCmd.Flags().BoolVar(&noGit, "no-git", false, "ignore git settings in config")
	bumpCmd.Flags().BoolVar(&gitCommit, "git-commit", false, "create a git commit (overrides config)")
	bumpCmd.Flags().BoolVar(&gitTag, "git-tag", false, "create a git tag (overrides config)")
	bumpCmd.Flags().BoolVar(&noTag, "no-tag", false, "do not create a tag")
	bumpCmd.Flags().BoolVar(&noCommit, "no-commit", false, "do not commit")
	bumpCmd.Flags().BoolVar(&allowDirty, "allow-dirty", false, "allow operation on dirty working directory")
	bumpCmd.Flags().StringVarP(&message, "message", "m", "", "override commit message")
	bumpCmd.Flags().StringVar(&tagName, "tag-name", "", "override tag name")
	rootCmd.AddCommand(bumpCmd)

	// Init command
	initCmd.Flags().BoolVarP(&initYes, "yes", "y", false, "non-interactive, accept defaults")
	initCmd.Flags().StringVar(&initScheme, "scheme", "semver", "version scheme (semver or calver:YYYY.MM.DD)")
	rootCmd.AddCommand(initCmd)

	// Validate command
	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "treat missing paths as errors")
	rootCmd.AddCommand(validateCmd)
}

func loadConfig() (*config.Config, string, error) {
	configPath := cfgFile
	if configPath == "" {
		var err error
		configPath, err = config.FindConfigFile()
		if err != nil {
			return nil, "", err
		}
	}

	loader := config.NewLoader(configPath)
	cfg, err := loader.Load()
	if err != nil {
		return nil, "", err
	}

	return cfg, configPath, nil
}

func runShow(_ *cobra.Command, _ []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	result, err := core.Show(cfg)
	if err != nil {
		return err
	}

	jsonMode := outputFormat == "json"
	printer := output.NewPrinter(jsonMode)
	if jsonMode {
		return printer.JSON(result)
	}

	printer.Println("%s", result.Version)
	return nil
}

func runBump(_ *cobra.Command, args []string) error {
	part := args[0]

	cfg, configPath, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitOperationError)
		return nil
	}

	jsonMode := outputFormat == "json"
	bumper := core.NewBumper(cfg, configPath, jsonMode)
	opts := core.BumpOptions{
		DryRun:     dryRun,
		NoGit:      noGit,
		GitCommit:  gitCommit,
		GitTag:     gitTag,
		NoTag:      noTag,
		NoCommit:   noCommit,
		AllowDirty: allowDirty,
		Message:    message,
		TagName:    tagName,
		JSON:       jsonMode,
	}

	result, err := bumper.Bump(part, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitOperationError)
		return nil
	}

	printer := output.NewPrinter(jsonMode)
	if jsonMode {
		return printer.JSON(result)
	}

	if dryRun {
		printer.Println("Dry run: would bump %s -> %s", result.OldVersion, result.NewVersion)
	} else {
		printer.Println("Bumped %s -> %s", result.OldVersion, result.NewVersion)
	}

	return nil
}

func runInit(_ *cobra.Command, _ []string) error {
	configPath := ".bumpx.toml"
	if cfgFile != "" {
		configPath = cfgFile
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists: %s", configPath)
	}

	cfg := config.NewDefaultConfig()
	cfg.Project.Name = "my-project"
	cfg.Project.CurrentVersion = "0.1.0"
	cfg.Project.VersionScheme = initScheme

	// Try to detect files in the current directory
	cfg.Files.Paths = detectProjectFiles()

	if err := config.SaveConfig(cfg, configPath); err != nil {
		return err
	}

	printer := output.NewPrinter(outputFormat == "json")
	printer.Println("Created %s", configPath)

	return nil
}

func runValidate(_ *cobra.Command, _ []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	result := config.Validate(cfg, validateStrict)

	jsonMode := outputFormat == "json"
	printer := output.NewPrinter(jsonMode)

	if jsonMode {
		jsonResult := output.ValidationResult{
			Valid:    result.Valid,
			Warnings: result.Warnings,
		}
		for _, e := range result.Errors {
			jsonResult.Errors = append(jsonResult.Errors, output.ValidationError{
				Field:   e.Field,
				Message: e.Message,
			})
		}
		return printer.JSON(jsonResult)
	}

	if result.Valid {
		printer.Println("Configuration is valid")
	} else {
		printer.Println("Configuration has errors:")
		for _, e := range result.Errors {
			printer.Println("  - %s: %s", e.Field, e.Message)
		}
	}

	if len(result.Warnings) > 0 {
		printer.Println("Warnings:")
		for _, w := range result.Warnings {
			printer.Println("  - %s", w)
		}
	}

	if !result.Valid {
		os.Exit(ExitOperationError)
	}

	return nil
}

func detectProjectFiles() []string {
	var detected []string

	// Common project files to detect
	filesToCheck := []string{
		"Cargo.toml",
		"package.json",
		"go.mod",
		"pyproject.toml",
	}

	for _, filename := range filesToCheck {
		if _, err := os.Stat(filename); err == nil {
			detected = append(detected, filename)
		}
	}

	return detected
}
