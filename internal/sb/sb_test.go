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

func TestDeleteFromDestDir(t *testing.T) {
	// 1. Setup: Create temporary source and destination directories for the test run
	baseSrcDir, err := os.MkdirTemp("", "test_src_*")
	if err != nil {
		t.Fatalf("Failed to create temp src dir: %v", err)
	}
	defer os.RemoveAll(baseSrcDir)

	baseDestDir, err := os.MkdirTemp("", "test_dest_*")
	if err != nil {
		t.Fatalf("Failed to create temp dest dir: %v", err)
	}
	defer os.RemoveAll(baseDestDir)

	// --- Create a common source structure ---
	// src/file_to_keep.txt
	err = os.WriteFile(filepath.Join(baseSrcDir, "file_to_keep.txt"), []byte("keep"), 0644)
	if err != nil {
		t.Fatalf("Failed to create src/file_to_keep.txt: %v", err)
	}
	// src/subdir/another_kept_file.txt
	srcSubDir := filepath.Join(baseSrcDir, "subdir")
	if err := os.Mkdir(srcSubDir, 0755); err != nil {
		t.Fatalf("Failed to create src/subdir: %v", err)
	}
	err = os.WriteFile(filepath.Join(srcSubDir, "another_kept_file.txt"), []byte("keep"), 0644)
	if err != nil {
		t.Fatalf("Failed to create src/subdir/another_kept_file.txt: %v", err)
	}

	// --- Create a destination structure with extra files and dirs ---
	// dest/file_to_keep.txt (should be kept)
	err = os.WriteFile(filepath.Join(baseDestDir, "file_to_keep.txt"), []byte("keep"), 0644)
	if err != nil {
		t.Fatalf("Failed to create dest/file_to_keep.txt: %v", err)
	}
	// dest/file_to_delete.txt (should be deleted)
	fileToDeletePath := filepath.Join(baseDestDir, "file_to_delete.txt")
	err = os.WriteFile(fileToDeletePath, []byte("delete"), 0644)
	if err != nil {
		t.Fatalf("Failed to create dest/file_to_delete.txt: %v", err)
	}

	// dest/subdir/another_kept_file.txt (should be kept)
	destSubDir := filepath.Join(baseDestDir, "subdir")
	if err := os.Mkdir(destSubDir, 0755); err != nil {
		t.Fatalf("Failed to create dest/subdir: %v", err)
	}
	err = os.WriteFile(filepath.Join(destSubDir, "another_kept_file.txt"), []byte("keep"), 0644)
	if err != nil {
		t.Fatalf("Failed to create dest/subdir/another_kept_file.txt: %v", err)
	}

	// dest/dir_to_empty_and_delete/file.txt (file and then dir should be deleted)
	dirToBecomeEmpty := filepath.Join(baseDestDir, "dir_to_empty_and_delete")
	if err := os.Mkdir(dirToBecomeEmpty, 0755); err != nil {
		t.Fatalf("Failed to create dest/dir_to_empty_and_delete: %v", err)
	}
	err = os.WriteFile(filepath.Join(dirToBecomeEmpty, "file.txt"), []byte("delete"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file in dir_to_empty_and_delete: %v", err)
	}

	// dest/empty_dir_to_delete (should be deleted)
	emptyDirToDelete := filepath.Join(baseDestDir, "empty_dir_to_delete")
	if err := os.Mkdir(emptyDirToDelete, 0755); err != nil {
		t.Fatalf("Failed to create dest/empty_dir_to_delete: %v", err)
	}

	// 2. Run DeleteFromDestDir
	if err := DeleteFromDestDir(baseSrcDir, baseDestDir); err != nil {
		t.Fatalf("DeleteFromDestDir failed: %v", err)
	}

	// 3. Assert: Check if destination directory is correct
	t.Run("should delete files not in source", func(t *testing.T) {
		if _, err := os.Stat(fileToDeletePath); !os.IsNotExist(err) {
			t.Errorf("Expected file '%s' to be deleted, but it still exists", fileToDeletePath)
		}
	})

	t.Run("should keep files that are in source", func(t *testing.T) {
		keptFilePath := filepath.Join(baseDestDir, "file_to_keep.txt")
		if _, err := os.Stat(keptFilePath); err != nil {
			t.Errorf("Expected file '%s' to be kept, but it was deleted or is inaccessible: %v", keptFilePath, err)
		}
	})

	t.Run("should delete directories that become empty", func(t *testing.T) {
		if _, err := os.Stat(dirToBecomeEmpty); !os.IsNotExist(err) {
			t.Errorf("Expected directory '%s' to be deleted after becoming empty, but it still exists", dirToBecomeEmpty)
		}
	})

	t.Run("should delete directories that were already empty", func(t *testing.T) {
		if _, err := os.Stat(emptyDirToDelete); !os.IsNotExist(err) {
			t.Errorf("Expected empty directory '%s' to be deleted, but it still exists", emptyDirToDelete)
		}
	})

	t.Run("should not delete non-empty directories", func(t *testing.T) {
		if _, err := os.Stat(destSubDir); err != nil {
			t.Errorf("Expected directory '%s' to be kept, but it was deleted or is inaccessible: %v", destSubDir, err)
		}
	})
}
