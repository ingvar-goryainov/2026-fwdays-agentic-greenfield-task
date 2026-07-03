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

// configFlag holds the --config flag value (FR-CLI-04): a path that, when
// set, overrides the default .agents-lint.yaml lookup location.
var configFlag string

// Execute runs the root command and returns any error encountered.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "path to the config file (default: ./.agents-lint.yaml)")
}
