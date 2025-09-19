/*
Copyright © 2025 Eric Ekholm <eric.ekholm@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ekholme/syncbuddy/internal/sb"
	"github.com/spf13/cobra"
)

var (
	sourceDir string
	destDir   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "syncbuddy",
	Short: "Sync files from a source directory to a destination directory",
	Long: `syncbuddy is a tool for copying all files from a source directory
to a destination directory.`,
	// Using RunE allows us to return an error from the command, which Cobra will
	// print to stderr. This is more idiomatic than handling errors with os.Exit.
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate that the source directory exists and is a directory.
		srcInfo, err := os.Stat(sourceDir)
		if err != nil {
			return fmt.Errorf("error accessing source '%s': %w", sourceDir, err)
		}
		if !srcInfo.IsDir() {
			return fmt.Errorf("source '%s' is not a directory", sourceDir)
		}

		fmt.Printf("Syncing from %s to %s...\n", sourceDir, destDir)
		if err := sb.CopyDir(sourceDir, destDir); err != nil {
			return fmt.Errorf("synchronization failed: %w", err)
		}
		fmt.Println("Synchronization complete!")
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&sourceDir, "source", "s", "", "Source directory to sync from (required)")
	rootCmd.Flags().StringVarP(&destDir, "destination", "d", "", "Destination directory to sync to (required)")

	//mark the above flags as required
	rootCmd.MarkFlagRequired("source")
	rootCmd.MarkFlagRequired("destination")
}
