/*
Copyright © 2025 Eric Ekholm <eric.ekholm@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/ekholme/syncbuddy/internal/sb"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{Use: "syncbuddy"}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(sb.NewCopyCommand())
	rootCmd.AddCommand(sb.NewDeleteCommand())
	rootCmd.AddCommand(sb.NewSyncCommand())
}
