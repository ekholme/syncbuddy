package sb

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	sourceDir string
	destDir   string
)

// CopyDir recursively copies a directory from src to dst.
// It creates directories and copies files, preserving file permissions.
func CopyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//determine the relative path of the current file/dir from the source directory
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		//construct the full destination path
		destPath := filepath.Join(dst, relPath)

		//check if src file is a directory & make the directory if it is
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		//if src file is a file, copy it
		if info.Mode().IsRegular() {
			sourceFile, err := os.Open(path)
			if err != nil {
				return err
			}

			defer sourceFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}

			defer destFile.Close()

			if _, err = io.Copy(destFile, sourceFile); err != nil {
				return err
			}

			// Set the destination file's permissions to match the source file's.
			return os.Chmod(destPath, info.Mode())
		}

		return nil
	})
}

// create a func that returns the copy command
func NewCopyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy a directory from src to dst",
		RunE: func(cmd *cobra.Command, args []string) error {
			srcInfo, err := os.Stat(sourceDir)
			if err != nil {
				return fmt.Errorf("error accessing source '%s': %w", sourceDir, err)
			}

			if srcInfo.IsDir() {
				return fmt.Errorf("source '%s' is not a directory", sourceDir)
			}

			fmt.Printf("Copying from %s to %s...\n", sourceDir, destDir)
			if err := CopyDir(sourceDir, destDir); err != nil {
				return fmt.Errorf("copying failed: %w", err)
			}
			fmt.Println("Copying complete!")
			return nil
		},
	}
	cmd.Flags().StringVarP(&sourceDir, "source", "s", "", "Source directory to sync from (required)")
	cmd.Flags().StringVarP(&destDir, "destination", "d", "", "Destination directory to sync to (required)")

	//mark the above flags as required
	cmd.MarkFlagRequired("source")
	cmd.MarkFlagRequired("destination")
	return cmd
}

// deletes any files in dest not in src, including empty directories
func DeleteFromDestDir(src string, dest string) error {
	//part 1: handle files -------------
	srcFiles, err := getAllFiles(src)
	if err != nil {
		return fmt.Errorf("error getting files from src: %w", err)
	}

	destFiles, err := getAllFiles(dest)
	if err != nil {
		return fmt.Errorf("error getting files from dest: %w", err)

	}

	for pathInDest := range destFiles {
		if _, ok := srcFiles[pathInDest]; !ok {
			fullPath := filepath.Join(dest, pathInDest)
			fmt.Printf("Deleting %s\n", fullPath)
			err := os.Remove(fullPath)
			if err != nil {
				return fmt.Errorf("error deleting %s: %w", fullPath, err)
			}
		}
	}

	//part 2: handle empty dirs ------------
	destDirs, err := getAllDirs(dest)
	if err != nil {
		return fmt.Errorf("error getting dirs from dest: %w", err)
	}

	//sort in place to ensure we process subdirectories before parent dirs
	sort.Slice(destDirs, func(i, j int) bool {
		return len(destDirs[i]) > len(destDirs[j])
	})

	// Iterate through the sorted directories and delete any that are empty.
	for _, dirPath := range destDirs {
		fullPathToDelete := filepath.Join(dest, dirPath)
		// Use os.ReadDir to check if the directory is empty.
		entries, err := os.ReadDir(fullPathToDelete)
		if err != nil {
			// If we can't read it, it's likely already gone. Continue to the next one.
			continue
		}

		if len(entries) == 0 {
			fmt.Printf("Deleting empty directory: %s\n", fullPathToDelete)
			if err := os.Remove(fullPathToDelete); err != nil {
				// We don't want to fail the entire sync if we can't remove an empty directory.
				fmt.Printf("Warning: failed to delete directory %s: %v\n", fullPathToDelete, err)
			}
		}
	}

	return nil
}

// wrapping the delete functionality in a cobra command
func NewDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete files in destination not in source",
		RunE: func(cmd *cobra.Command, args []string) error {
			srcInfo, err := os.Stat(sourceDir)
			if err != nil {
				return fmt.Errorf("error accessing source '%s': %w", sourceDir, err)
			}

			if srcInfo.IsDir() {
				return fmt.Errorf("source '%s' is not a directory", sourceDir)
			}

			fmt.Printf("Deleting files from %s not in %s...\n", destDir, sourceDir)
			if err := DeleteFromDestDir(sourceDir, destDir); err != nil {
				return fmt.Errorf("deleting failed: %w", err)
			}
			fmt.Println("Deletion complete!")
			return nil

		},
	}
	cmd.Flags().StringVarP(&sourceDir, "source", "s", "", "Source directory to sync from (required)")
	cmd.Flags().StringVarP(&destDir, "destination", "d", "", "Destination directory to sync to (required)")

	//mark the above flags as required
	cmd.MarkFlagRequired("source")
	cmd.MarkFlagRequired("destination")
	return cmd

}

// helper to get all files in a directory
func getAllFiles(root string) (map[string]bool, error) {
	files := make(map[string]bool)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil //skip directories for now
		}
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		//replace \ with / for cross-platform consistency.
		//although this won't really matter for me right now
		relPath = strings.ReplaceAll(relPath, "\\", "/")
		files[relPath] = true
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, err
}

// helper function to walk a directory tree and return a slice of all directory paths relative to the root
func getAllDirs(root string) ([]string, error) {
	var dirs []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if relPath != "." {
				dirs = append(dirs, relPath)
			}
		}
		return nil
	})
	return dirs, err
}
