package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "bumpx-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath = filepath.Join(tmpDir, "bumpx")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		panic("failed to build bumpx: " + err.Error())
	}

	os.Exit(m.Run())
}

func runBumpx(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return string(output), exitCode
}

func setupTestProject(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	config := "[project]\nname = \"test-project\"\ncurrent_version = \"1.0.0\"\nversion_scheme = \"semver\"\n\n[files]\npaths = [\"package.json\"]\n\n[git]\ncommit = false\ntag = false\n"
	if err := os.WriteFile(filepath.Join(tmpDir, ".bumpx.toml"), []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	packageJSON := "{\n  \"name\": \"test-project\",\n  \"version\": \"1.0.0\"\n}\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(packageJSON), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	return tmpDir
}

func TestCLI_Version(t *testing.T) {
	output, exitCode := runBumpx(t, ".", "--version")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(output, "bumpx version") {
		t.Errorf("Expected version output, got: %s", output)
	}
}

func TestCLI_Show(t *testing.T) {
	tmpDir := setupTestProject(t)
	output, exitCode := runBumpx(t, tmpDir, "show")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d: %s", exitCode, output)
	}
	if !strings.Contains(output, "1.0.0") {
		t.Errorf("Expected version 1.0.0, got: %s", output)
	}
}

func TestCLI_Validate(t *testing.T) {
	tmpDir := setupTestProject(t)
	output, exitCode := runBumpx(t, tmpDir, "validate")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d: %s", exitCode, output)
	}
	if !strings.Contains(output, "valid") {
		t.Errorf("Expected validation success, got: %s", output)
	}
}

func TestCLI_Init(t *testing.T) {
	tmpDir := t.TempDir()
	output, exitCode := runBumpx(t, tmpDir, "init", "--yes")
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d: %s", exitCode, output)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".bumpx.toml")); os.IsNotExist(err) {
		t.Error("Expected .bumpx.toml to be created")
	}
}

func TestCLI_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	_, exitCode := runBumpx(t, tmpDir, "show")
	if exitCode == 0 {
		t.Error("Expected non-zero exit code when no config exists")
	}
}
