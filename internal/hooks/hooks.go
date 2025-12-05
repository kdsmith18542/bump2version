// Package hooks provides external command execution for hooks.
package hooks

import (
	"fmt"
	"os"
	"os/exec"
)

// Runner executes hook commands.
type Runner struct {
	workDir string
}

// New creates a new hook runner.
func New(workDir string) *Runner {
	if workDir == "" {
		workDir = "."
	}
	return &Runner{workDir: workDir}
}

// HookContext contains information passed to hooks via environment variables.
type HookContext struct {
	OldVersion string
	NewVersion string
	Part       string
	ConfigPath string
}

// Run executes a hook command.
// Note: Hook commands are executed via shell and come from the config file.
// Ensure the config file is trusted as commands are executed without additional validation.
// Hook commands should be used for trusted local scripts only.
func (r *Runner) Run(command string, ctx *HookContext) error {
	if command == "" {
		return nil
	}

	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = r.workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment variables
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("BUMPX_OLD_VERSION=%s", ctx.OldVersion),
		fmt.Sprintf("BUMPX_NEW_VERSION=%s", ctx.NewVersion),
		fmt.Sprintf("BUMPX_PART=%s", ctx.Part),
		fmt.Sprintf("BUMPX_CONFIG_PATH=%s", ctx.ConfigPath),
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("hook command failed: %w", err)
	}

	return nil
}

// RunPreBump executes the pre-bump hook.
func (r *Runner) RunPreBump(command string, ctx *HookContext) error {
	if err := r.Run(command, ctx); err != nil {
		return fmt.Errorf("pre_bump hook failed: %w", err)
	}
	return nil
}

// RunPostBump executes the post-bump hook.
func (r *Runner) RunPostBump(command string, ctx *HookContext) error {
	if err := r.Run(command, ctx); err != nil {
		return fmt.Errorf("post_bump hook failed: %w", err)
	}
	return nil
}
