package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("default workdir", func(t *testing.T) {
		g := New("")
		if g.workDir != "." {
			t.Errorf("Expected workDir='.', got %q", g.workDir)
		}
	})

	t.Run("custom workdir", func(t *testing.T) {
		g := New("/custom/path")
		if g.workDir != "/custom/path" {
			t.Errorf("Expected workDir='/custom/path', got %q", g.workDir)
		}
	})
}

// setupTestRepo creates a temporary git repository for testing
func setupTestRepo(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git user (required for commits)
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git name: %v", err)
	}

	return tmpDir
}

func TestGit_IsRepo(t *testing.T) {
	t.Run("valid repo", func(t *testing.T) {
		tmpDir := setupTestRepo(t)
		g := New(tmpDir)

		if !g.IsRepo() {
			t.Error("Expected IsRepo() to return true for git repo")
		}
	})

	t.Run("not a repo", func(t *testing.T) {
		tmpDir := t.TempDir()
		g := New(tmpDir)

		if g.IsRepo() {
			t.Error("Expected IsRepo() to return false for non-repo")
		}
	})
}

func TestGit_IsClean(t *testing.T) {
	tmpDir := setupTestRepo(t)
	g := New(tmpDir)

	t.Run("clean repo", func(t *testing.T) {
		clean, err := g.IsClean()
		if err != nil {
			t.Fatalf("IsClean() error: %v", err)
		}
		if !clean {
			t.Error("Expected clean repo to return true")
		}
	})

	t.Run("dirty repo", func(t *testing.T) {
		// Create untracked file
		if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		clean, err := g.IsClean()
		if err != nil {
			t.Fatalf("IsClean() error: %v", err)
		}
		if clean {
			t.Error("Expected dirty repo to return false")
		}
	})
}

func TestGit_Add(t *testing.T) {
	tmpDir := setupTestRepo(t)
	g := New(tmpDir)

	// Create a file to add
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err := g.Add("test.txt")
	if err != nil {
		t.Errorf("Add() error: %v", err)
	}
}

func TestGit_Commit(t *testing.T) {
	tmpDir := setupTestRepo(t)
	g := New(tmpDir)

	// Create and stage a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := g.Add("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	err := g.Commit("Test commit message")
	if err != nil {
		t.Errorf("Commit() error: %v", err)
	}
}

func TestGit_Tag(t *testing.T) {
	tmpDir := setupTestRepo(t)
	g := New(tmpDir)

	// Create initial commit (required for tagging)
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := g.Add("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}
	if err := g.Commit("Initial commit"); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	t.Run("simple tag", func(t *testing.T) {
		err := g.Tag("v1.0.0", "")
		if err != nil {
			t.Errorf("Tag() error: %v", err)
		}
	})

	t.Run("annotated tag", func(t *testing.T) {
		err := g.Tag("v1.0.1", "Release v1.0.1")
		if err != nil {
			t.Errorf("Tag() with message error: %v", err)
		}
	})
}

func TestGit_GetCurrentBranch(t *testing.T) {
	tmpDir := setupTestRepo(t)
	g := New(tmpDir)

	// Create initial commit (required for branch)
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := g.Add("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}
	if err := g.Commit("Initial commit"); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	branch, err := g.GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch() error: %v", err)
	}

	// Default branch could be "master" or "main"
	if branch != "master" && branch != "main" {
		t.Errorf("Expected branch 'master' or 'main', got %q", branch)
	}
}

func TestGit_GetLatestTag(t *testing.T) {
	tmpDir := setupTestRepo(t)
	g := New(tmpDir)

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := g.Add("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}
	if err := g.Commit("Initial commit"); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	t.Run("no tags", func(t *testing.T) {
		tag, err := g.GetLatestTag()
		if err != nil {
			t.Errorf("GetLatestTag() error: %v", err)
		}
		if tag != "" {
			t.Errorf("Expected empty tag, got %q", tag)
		}
	})

	t.Run("with tag", func(t *testing.T) {
		if err := g.Tag("v1.0.0", ""); err != nil {
			t.Fatalf("Failed to create tag: %v", err)
		}

		tag, err := g.GetLatestTag()
		if err != nil {
			t.Errorf("GetLatestTag() error: %v", err)
		}
		if tag != "v1.0.0" {
			t.Errorf("Expected tag 'v1.0.0', got %q", tag)
		}
	})
}

func TestGit_GetCommitHash(t *testing.T) {
	tmpDir := setupTestRepo(t)
	g := New(tmpDir)

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := g.Add("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}
	if err := g.Commit("Initial commit"); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	hash, err := g.GetCommitHash()
	if err != nil {
		t.Fatalf("GetCommitHash() error: %v", err)
	}

	// SHA-1 hash is 40 characters
	if len(hash) != 40 {
		t.Errorf("Expected 40-char hash, got %d chars: %q", len(hash), hash)
	}
}

func TestGit_GetShortCommitHash(t *testing.T) {
	tmpDir := setupTestRepo(t)
	g := New(tmpDir)

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := g.Add("test.txt"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}
	if err := g.Commit("Initial commit"); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	hash, err := g.GetShortCommitHash()
	if err != nil {
		t.Fatalf("GetShortCommitHash() error: %v", err)
	}

	// Short hash is typically 7 characters
	if len(hash) < 7 {
		t.Errorf("Expected at least 7-char short hash, got %d chars: %q", len(hash), hash)
	}
}
