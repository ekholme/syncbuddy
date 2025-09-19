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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Syncing from %s to %s...\n", sourceDir, destDir)
		if err := sb.CopyDir(sourceDir, destDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Println("Synchronization complete!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
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
