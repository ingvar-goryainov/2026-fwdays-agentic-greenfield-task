package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/ingvar-goryainov/agents-lint/internal/reporter"
	"github.com/spf13/cobra"
)

const defaultAgentsMDPath = "./AGENTS.md"

// runScan validates the AGENTS.md file at path and writes reporter-formatted
// output to w. It returns the process exit code (FR-CLI-02) and never calls
// os.Exit itself, keeping it fully unit-testable; the Cobra RunE wrapper
// owns the actual process exit.
func runScan(w io.Writer, path string) (exitCode int, err error) {
	findings, runErr := rules.Run(path)
	if runErr != nil {
		return 1, fmt.Errorf("scan %s: %w", path, runErr)
	}

	ruleCount := len(rules.DefaultRules()) + 1 // +1 for S001, run outside the Rule interface
	if writeErr := reporter.WriteText(w, findings, ruleCount, reporter.Options{}); writeErr != nil {
		return 1, writeErr
	}

	return exitCodeForFindings(findings), nil
}

// exitCodeForFindings implements FR-CLI-02: exit 1 if any finding is
// SeverityError, exit 0 otherwise (including when there are only
// SeverityWarning findings). Split out from runScan so the severity
// decision can be unit-tested directly against synthetic findings, since
// today's rule set (S001-S005) has no SeverityWarning rule to exercise
// this path through a real file.
func exitCodeForFindings(findings []lint.Finding) int {
	for _, f := range findings {
		if f.Severity == lint.SeverityError {
			return 1
		}
	}
	return 0
}

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Validate an AGENTS.md file against the schema",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := defaultAgentsMDPath
		if len(args) > 0 {
			path = args[0]
		}

		exitCode, err := runScan(cmd.OutOrStdout(), path)
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
	rootCmd.AddCommand(scanCmd)
}
