package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCargoHandler(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "bumpx-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test Cargo.toml file
	cargoPath := filepath.Join(tmpDir, "Cargo.toml")
	cargoContent := `[package]
name = "my-project"
version = "1.2.3"
edition = "2021"

[dependencies]
serde = "1.0"
`
	if err := os.WriteFile(cargoPath, []byte(cargoContent), 0644); err != nil {
		t.Fatalf("Failed to write Cargo.toml: %v", err)
	}

	handler := &CargoHandler{}

	// Test Read
	version, err := handler.Read(cargoPath)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if version != "1.2.3" {
		t.Errorf("Read() = %s, want 1.2.3", version)
	}

	// Test Write (dry run)
	if err := handler.Write(cargoPath, "1.2.3", "1.3.0", true); err != nil {
		t.Fatalf("Write(dryRun=true) error = %v", err)
	}

	// Verify file unchanged
	version, _ = handler.Read(cargoPath)
	if version != "1.2.3" {
		t.Error("Dry run modified the file")
	}

	// Test Write (actual)
	if err := handler.Write(cargoPath, "1.2.3", "1.3.0", false); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify file changed
	version, _ = handler.Read(cargoPath)
	if version != "1.3.0" {
		t.Errorf("After Write(), version = %s, want 1.3.0", version)
	}
}

func TestNPMHandler(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "bumpx-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test package.json file
	pkgPath := filepath.Join(tmpDir, "package.json")
	pkgContent := `{
  "name": "my-project",
  "version": "1.2.3",
  "description": "A test project"
}
`
	if err := os.WriteFile(pkgPath, []byte(pkgContent), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	handler := &NPMHandler{}

	// Test Read
	version, err := handler.Read(pkgPath)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if version != "1.2.3" {
		t.Errorf("Read() = %s, want 1.2.3", version)
	}

	// Test Write
	if err := handler.Write(pkgPath, "1.2.3", "1.3.0", false); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify file changed
	version, _ = handler.Read(pkgPath)
	if version != "1.3.0" {
		t.Errorf("After Write(), version = %s, want 1.3.0", version)
	}
}

func TestRegexHandler(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "bumpx-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test VERSION file
	versionPath := filepath.Join(tmpDir, "VERSION")
	versionContent := `version=1.2.3
`
	if err := os.WriteFile(versionPath, []byte(versionContent), 0644); err != nil {
		t.Fatalf("Failed to write VERSION: %v", err)
	}

	handler, err := NewRegexHandlerWithPatterns(`version=(?P<version>[0-9\.]+)`, "version={version}")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Test Read
	version, err := handler.Read(versionPath)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if version != "1.2.3" {
		t.Errorf("Read() = %s, want 1.2.3", version)
	}

	// Test Write
	if err := handler.Write(versionPath, "1.2.3", "1.3.0", false); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify file changed
	version, _ = handler.Read(versionPath)
	if version != "1.3.0" {
		t.Errorf("After Write(), version = %s, want 1.3.0", version)
	}
}

func TestDetectFileType(t *testing.T) {
	tests := []struct {
		path     string
		wantType string
	}{
		{"Cargo.toml", "cargo"},
		{"backend/Cargo.toml", "cargo"},
		{"package.json", "npm"},
		{"frontend/package.json", "npm"},
		{"go.mod", "gomod"},
		{"pyproject.toml", "pyproject"},
		{"setup.cfg", "pyproject"},
		{"MyProject.csproj", "dotnet"},
		{"AssemblyInfo.cs", "dotnet"},
		{"VERSION.txt", "regex"},
		{"unknown.file", "regex"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := DetectFileType(tt.path)
			if got != tt.wantType {
				t.Errorf("DetectFileType(%s) = %s, want %s", tt.path, got, tt.wantType)
			}
		})
	}
}
