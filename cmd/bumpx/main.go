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

var (
	cfgFile    string
	jsonOutput bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
	dryRun   bool
	noGit    bool
	noTag    bool
	noCommit bool
	message  string
	tagName  string
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
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output in JSON format")

	// Show command
	rootCmd.AddCommand(showCmd)

	// Bump command
	bumpCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "don't write any files, just pretend")
	bumpCmd.Flags().BoolVar(&noGit, "no-git", false, "ignore git settings in config")
	bumpCmd.Flags().BoolVar(&noTag, "no-tag", false, "do not create a tag")
	bumpCmd.Flags().BoolVar(&noCommit, "no-commit", false, "do not commit")
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
	cfg, configPath, err := loadConfig()
	if err != nil {
		return err
	}

	bumper := core.NewBumper(cfg, configPath, jsonOutput)
	result, err := bumper.Show()
	if err != nil {
		return err
	}

	printer := output.NewPrinter(jsonOutput)
	if jsonOutput {
		return printer.JSON(result)
	}

	printer.Println("%s", result.Version)
	return nil
}

func runBump(_ *cobra.Command, args []string) error {
	part := args[0]

	cfg, configPath, err := loadConfig()
	if err != nil {
		return err
	}

	bumper := core.NewBumper(cfg, configPath, jsonOutput)
	opts := core.BumpOptions{
		DryRun:   dryRun,
		NoGit:    noGit,
		NoTag:    noTag,
		NoCommit: noCommit,
		Message:  message,
		TagName:  tagName,
		JSON:     jsonOutput,
	}

	result, err := bumper.Bump(part, opts)
	if err != nil {
		return err
	}

	printer := output.NewPrinter(jsonOutput)
	if jsonOutput {
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
	detectedFiles := detectProjectFiles()
	for name, fileConfig := range detectedFiles {
		cfg.Files[name] = fileConfig
	}

	if err := config.SaveConfig(cfg, configPath); err != nil {
		return err
	}

	printer := output.NewPrinter(jsonOutput)
	printer.Println("Created %s", configPath)

	return nil
}

func runValidate(_ *cobra.Command, _ []string) error {
	cfg, _, err := loadConfig()
	if err != nil {
		return err
	}

	result := config.Validate(cfg, validateStrict)

	printer := output.NewPrinter(jsonOutput)

	if jsonOutput {
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
		os.Exit(2)
	}

	return nil
}

func detectProjectFiles() map[string]config.FileConfig {
	detected := make(map[string]config.FileConfig)

	// Common project files to detect
	filesToCheck := map[string]struct {
		name     string
		fileType string
	}{
		"Cargo.toml":     {"cargo", "cargo"},
		"package.json":   {"npm", "npm"},
		"go.mod":         {"gomod", "gomod"},
		"pyproject.toml": {"python", "pyproject"},
	}

	for filename, info := range filesToCheck {
		if _, err := os.Stat(filename); err == nil {
			detected[info.name] = config.FileConfig{
				Paths: []string{filename},
				Type:  info.fileType,
			}
		}
	}

	return detected
}
