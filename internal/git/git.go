// Package git provides Git integration utilities.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Git provides Git operations.
type Git struct {
	workDir string
}

// New creates a new Git instance.
func New(workDir string) *Git {
	if workDir == "" {
		workDir = "."
	}
	return &Git{workDir: workDir}
}

// IsRepo checks if the current directory is a Git repository.
func (g *Git) IsRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = g.workDir
	return cmd.Run() == nil
}

// IsClean checks if the working directory has no uncommitted changes.
func (g *Git) IsClean() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}
	return len(bytes.TrimSpace(output)) == 0, nil
}

// Add stages files for commit.
func (g *Git) Add(paths ...string) error {
	args := append([]string{"add"}, paths...)
	cmd := exec.Command("git", args...)
	cmd.Dir = g.workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stage files: %w\n%s", err, output)
	}
	return nil
}

// Commit creates a new commit with the specified message.
func (g *Git) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = g.workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit: %w\n%s", err, output)
	}
	return nil
}

// Tag creates a new annotated tag.
func (g *Git) Tag(name, message string) error {
	var cmd *exec.Cmd
	if message != "" {
		cmd = exec.Command("git", "tag", "-a", name, "-m", message)
	} else {
		cmd = exec.Command("git", "tag", name)
	}
	cmd.Dir = g.workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create tag: %w\n%s", err, output)
	}
	return nil
}

// GetCurrentBranch returns the current branch name.
func (g *Git) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetLatestTag returns the latest tag.
func (g *Git) GetLatestTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", nil // No tags found is not an error
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCommitHash returns the current commit hash.
func (g *Git) GetCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get commit hash: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetShortCommitHash returns the short commit hash.
func (g *Git) GetShortCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get short commit hash: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
