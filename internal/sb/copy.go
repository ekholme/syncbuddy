package sb

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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
