package sb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir(t *testing.T) {
	// 1. Setup: Create temporary source and destination directories
	srcDir, err := os.MkdirTemp("", "test_src_*")
	if err != nil {
		t.Fatalf("Failed to create temp src dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	destDir, err := os.MkdirTemp("", "test_dest_*")
	if err != nil {
		t.Fatalf("Failed to create temp dest dir: %v", err)
	}
	defer os.RemoveAll(destDir)

	// 2. Create a file structure in the source directory
	// src/file1.txt
	file1Content := []byte("hello world")
	file1Path := filepath.Join(srcDir, "file1.txt")
	if err := os.WriteFile(file1Path, file1Content, 0644); err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}

	// src/subdir/file2.txt
	subDir := filepath.Join(srcDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	file2Content := []byte("another file")
	file2Path := filepath.Join(subDir, "file2.txt")
	// Use a specific permission to test if it's preserved
	if err := os.WriteFile(file2Path, file2Content, 0700); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	// 3. Run CopyDir
	if err := CopyDir(srcDir, destDir); err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	// 4. Assert: Check if destination directory is correct
	// Check file1.txt
	destFile1Path := filepath.Join(destDir, "file1.txt")
	content, err := os.ReadFile(destFile1Path)
	if err != nil {
		t.Fatalf("Failed to read dest file1: %v", err)
	}
	if string(content) != string(file1Content) {
		t.Errorf("File1 content mismatch: got %q, want %q", string(content), string(file1Content))
	}
	srcInfo, _ := os.Stat(file1Path)
	destInfo, _ := os.Stat(destFile1Path)
	if srcInfo.Mode() != destInfo.Mode() {
		t.Errorf("File1 permission mismatch: got %v, want %v", destInfo.Mode(), srcInfo.Mode())
	}

	// Check subdir/file2.txt
	destFile2Path := filepath.Join(destDir, "subdir", "file2.txt")
	content, err = os.ReadFile(destFile2Path)
	if err != nil {
		t.Fatalf("Failed to read dest file2: %v", err)
	}
	if string(content) != string(file2Content) {
		t.Errorf("File2 content mismatch: got %q, want %q", string(content), string(file2Content))
	}

	srcInfo, _ = os.Stat(file2Path)
	destInfo, _ = os.Stat(destFile2Path)
	if srcInfo.Mode() != destInfo.Mode() {
		t.Errorf("File2 permission mismatch: got %v, want %v", destInfo.Mode(), srcInfo.Mode())
	}
}
