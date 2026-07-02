// Package cmd holds the Cobra command tree for the agents-lint binary.
// Subcommands (scan, init, ...) register themselves here in later changes.
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agents-lint",
	Short: "agents-lint validates AGENTS.md files against a formal schema",
}

// Execute runs the root command and returns any error encountered.
func Execute() error {
	return rootCmd.Execute()
}
