package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/ingvar-goryainov/agents-lint/internal/config"
	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/ingvar-goryainov/agents-lint/internal/reporter"
	"github.com/spf13/cobra"
)

const defaultAgentsMDPath = "./AGENTS.md"

// formatText and formatSARIF are the only values validateFormat accepts
// for --format (FR-CLI-05).
const (
	formatText  = "text"
	formatSARIF = "sarif"
)

// validateFormat implements FR-CLI-05's validation: format must be
// formatText or formatSARIF. It runs before any file I/O so an invalid
// --format value fails fast with a clear, actionable error.
func validateFormat(format string) error {
	if format != formatText && format != formatSARIF {
		return fmt.Errorf("unknown format %q: must be %q or %q", format, formatText, formatSARIF)
	}
	return nil
}

// runScan validates the AGENTS.md file resolved from argPath/configPath and
// writes reporter-formatted output to w. It returns the process exit code
// (FR-CLI-02) and never calls os.Exit itself, keeping it fully
// unit-testable; the Cobra RunE wrapper owns the actual process exit.
//
// argPath is the raw positional argument ("" if omitted), configPath is
// the --config flag value ("" if omitted), and format is the resolved
// --format value (formatText or formatSARIF); runScan resolves the
// effective config (FR-CFG-01/03, FR-CLI-04) and the effective scan path
// (FR-CFG-02) itself, since the precedence between them can't be decided
// before the config is loaded.
func runScan(w io.Writer, argPath, configPath, format string) (exitCode int, err error) {
	cfg, cfgErr := config.Resolve(configPath, rules.KnownRuleIDs())
	if cfgErr != nil {
		return 1, cfgErr
	}

	path := resolvePath(argPath, cfg.Path)

	findings, runErr := rules.Run(path)
	if runErr != nil {
		return 1, fmt.Errorf("scan %s: %w", path, runErr)
	}
	findings = cfg.Apply(findings)

	var writeErr error
	if format == formatSARIF {
		writeErr = reporter.WriteSARIF(w, findings)
	} else {
		ruleCount := effectiveRuleCount(cfg, rules.KnownRuleIDs())
		writeErr = reporter.WriteText(w, findings, ruleCount, reporter.Options{})
	}
	if writeErr != nil {
		return 1, writeErr
	}

	return exitCodeForFindings(findings), nil
}

// resolvePath implements FR-CFG-02's path precedence: an explicit
// positional argument wins, then the config file's path:, then the
// FR-CLI-01 default.
func resolvePath(argPath, cfgPath string) string {
	switch {
	case argPath != "":
		return argPath
	case cfgPath != "":
		return cfgPath
	default:
		return defaultAgentsMDPath
	}
}

// effectiveRuleCount is len(knownIDs) minus one for every rule ID that cfg
// explicitly disables, so the reporter's success line reflects the
// reduced rule set (FR-CFG-02).
func effectiveRuleCount(cfg *config.Config, knownIDs []string) int {
	count := len(knownIDs)
	for _, id := range knownIDs {
		if rc, ok := cfg.Rules[id]; ok && rc.Enabled != nil && !*rc.Enabled {
			count--
		}
	}
	return count
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
		if err := validateFormat(formatFlag); err != nil {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
			return err
		}

		argPath := ""
		if len(args) > 0 {
			argPath = args[0]
		}

		exitCode, err := runScan(cmd.OutOrStdout(), argPath, configFlag, formatFlag)
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
