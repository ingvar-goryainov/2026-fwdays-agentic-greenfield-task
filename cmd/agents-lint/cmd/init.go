package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/ingvar-goryainov/agents-lint/internal/lint/skeleton"
	"github.com/spf13/cobra"
)

// runInit writes the AGENTS.md skeleton (FR-CLI-03) to path and reports
// progress to w. It returns the process exit code and never calls os.Exit
// itself, keeping it fully unit-testable; the Cobra RunE wrapper owns the
// actual process exit. Following the runScan pattern in scan.go.
func runInit(w io.Writer, path string, force bool) (exitCode int, err error) {
	_, statErr := os.Stat(path)
	switch {
	case statErr == nil && !force:
		return 1, fmt.Errorf(
			"init %s: file already exists — pass --force to overwrite it",
			path,
		)
	case statErr != nil && !os.IsNotExist(statErr):
		return 1, fmt.Errorf("init %s: %w", path, statErr)
	}

	if err := os.WriteFile(path, skeleton.Content(), 0o644); err != nil {
		return 1, fmt.Errorf("init %s: %w", path, err)
	}

	if _, err := fmt.Fprintf(w, "created %s\n", path); err != nil {
		return 1, err
	}
	return 0, nil
}

var initForce bool

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Generate a minimal AGENTS.md skeleton that passes the schema rules",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := defaultAgentsMDPath
		if len(args) > 0 {
			path = args[0]
		}

		exitCode, err := runInit(cmd.OutOrStdout(), path, initForce)
		if err != nil {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
			return err
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}
		return nil
	},
}

func init() {
	initCmd.Flags().BoolVar(&initForce, "force", false, "overwrite an existing AGENTS.md")
	rootCmd.AddCommand(initCmd)
}
