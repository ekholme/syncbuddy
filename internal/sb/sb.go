package sb

import (
	"io"
	"os"
	"path/filepath"
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
