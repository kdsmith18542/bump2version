package hooks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("default workdir", func(t *testing.T) {
		r := New("")
		if r.workDir != "." {
			t.Errorf("Expected workDir='.', got %q", r.workDir)
		}
	})

	t.Run("custom workdir", func(t *testing.T) {
		r := New("/custom/path")
		if r.workDir != "/custom/path" {
			t.Errorf("Expected workDir='/custom/path', got %q", r.workDir)
		}
	})
}

func TestHookContext(t *testing.T) {
	ctx := &HookContext{
		OldVersion: "1.0.0",
		NewVersion: "1.0.1",
		Part:       "patch",
		ConfigPath: ".bumpx.toml",
	}

	if ctx.OldVersion != "1.0.0" {
		t.Errorf("Expected OldVersion=1.0.0, got %s", ctx.OldVersion)
	}
	if ctx.NewVersion != "1.0.1" {
		t.Errorf("Expected NewVersion=1.0.1, got %s", ctx.NewVersion)
	}
}

func TestRunner_Run_EmptyCommand(t *testing.T) {
	r := New("")
	ctx := &HookContext{}

	// Empty command should return nil
	err := r.Run("", ctx)
	if err != nil {
		t.Errorf("Expected nil error for empty command, got %v", err)
	}
}

func TestRunner_Run_SimpleCommand(t *testing.T) {
	r := New("")
	ctx := &HookContext{
		OldVersion: "1.0.0",
		NewVersion: "1.0.1",
		Part:       "patch",
		ConfigPath: ".bumpx.toml",
	}

	// Simple echo command should succeed
	err := r.Run("echo test", ctx)
	if err != nil {
		t.Errorf("Expected nil error for echo command, got %v", err)
	}
}

func TestRunner_Run_WithEnvVars(t *testing.T) {
	// Create a temp file to capture env vars
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "env.txt")

	r := New(tmpDir)
	ctx := &HookContext{
		OldVersion: "1.0.0",
		NewVersion: "2.0.0",
		Part:       "major",
		ConfigPath: ".bumpx.toml",
	}

	// Write env vars to file
	cmd := "echo $BUMPX_OLD_VERSION $BUMPX_NEW_VERSION $BUMPX_PART > " + outFile
	err := r.Run(cmd, ctx)
	if err != nil {
		t.Fatalf("Hook failed: %v", err)
	}

	// Read and verify
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "1.0.0 2.0.0 major\n"
	if string(data) != expected {
		t.Errorf("Expected env output %q, got %q", expected, string(data))
	}
}

func TestRunner_Run_FailingCommand(t *testing.T) {
	r := New("")
	ctx := &HookContext{}

	// Command that fails
	err := r.Run("exit 1", ctx)
	if err == nil {
		t.Error("Expected error for failing command")
	}
}

func TestRunner_RunPreBump(t *testing.T) {
	r := New("")
	ctx := &HookContext{}

	err := r.RunPreBump("echo pre-bump", ctx)
	if err != nil {
		t.Errorf("RunPreBump failed: %v", err)
	}
}

func TestRunner_RunPostBump(t *testing.T) {
	r := New("")
	ctx := &HookContext{}

	err := r.RunPostBump("echo post-bump", ctx)
	if err != nil {
		t.Errorf("RunPostBump failed: %v", err)
	}
}

func TestRunner_RunPreBump_Failure(t *testing.T) {
	r := New("")
	ctx := &HookContext{}

	err := r.RunPreBump("exit 1", ctx)
	if err == nil {
		t.Error("Expected error for failing pre_bump hook")
	}
}

func TestRunner_RunPostBump_Failure(t *testing.T) {
	r := New("")
	ctx := &HookContext{}

	err := r.RunPostBump("exit 1", ctx)
	if err == nil {
		t.Error("Expected error for failing post_bump hook")
	}
}
