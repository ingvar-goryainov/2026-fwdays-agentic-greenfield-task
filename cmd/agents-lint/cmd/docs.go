package cmd

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/spf13/cobra"
)

// runDocs looks up ruleID in the rule registry and writes its full
// specification to w (FR-CLI-07): description, default severity, a valid
// example, an invalid example, and fix guidance. An unknown ruleID
// produces an actionable error rather than empty output (BC-OUTPUT-02);
// Cobra's ExactArgs(1) on docsCmd rejects a missing argument before this
// is ever called.
func runDocs(w io.Writer, ruleID string) error {
	doc, ok := rules.RuleDocs()[ruleID]
	if !ok {
		return fmt.Errorf(
			"unknown rule %q — known rule IDs: %s",
			ruleID, strings.Join(sortedRuleIDs(), ", "),
		)
	}

	_, err := fmt.Fprintf(w,
		"%s (%s)\n\n%s\n\nValid:\n%s\n\nInvalid:\n%s\n\nFix:\n%s\n",
		doc.ID, doc.Severity, doc.Description, indentBlock(doc.Valid), indentBlock(doc.Invalid), indentBlock(doc.FixGuidance),
	)
	return err
}

// indentBlock prefixes every line of s with two spaces, so multi-line
// examples (which may contain blank lines, e.g. an AGENTS.md snippet)
// render as a single visually distinct block instead of only the first
// line being indented.
func indentBlock(s string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, line := range lines {
		lines[i] = "  " + line
	}
	return strings.Join(lines, "\n")
}

// sortedRuleIDs returns rules.KnownRuleIDs() sorted, so the "known rule
// IDs" hint in an unknown-rule error is stable and easy to scan.
func sortedRuleIDs() []string {
	ids := append([]string(nil), rules.KnownRuleIDs()...)
	sort.Strings(ids)
	return ids
}

var docsCmd = &cobra.Command{
	Use:   "docs <rule-id>",
	Short: "Print the full specification for a rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runDocs(cmd.OutOrStdout(), args[0]); err != nil {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
