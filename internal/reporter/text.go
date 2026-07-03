package reporter

import (
	"fmt"
	"io"
	"os"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
)

// Options controls optional rendering behavior for WriteText.
type Options struct {
	// Color enables ANSI severity coloring. Ignored (treated as false)
	// when the NO_COLOR environment variable is set to a non-empty value.
	Color bool
}

const (
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
	ansiReset  = "\x1b[0m"
)

// WriteText renders findings as human-readable text to w, following
// FR-OUT-01 (per-line format), FR-OUT-02 (summary line), and FR-OUT-04
// (success line). ruleCount is the total number of rules evaluated to
// produce findings; it is only used to render the success line's "(N
// rules passed)" when findings is empty.
func WriteText(w io.Writer, findings []lint.Finding, ruleCount int, opts Options) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintf(w, "✓ AGENTS.md is valid (%d rules passed)\n", ruleCount)
		return err
	}

	colorEnabled := opts.Color && os.Getenv("NO_COLOR") == ""

	var errCount, warnCount int
	for _, f := range findings {
		severity := string(f.Severity)
		token := severity
		if color := severityColor(f.Severity); colorEnabled && color != "" {
			token = color + severity + ansiReset
		}
		if _, err := fmt.Fprintf(w, "%s  %s  %s:%d — %s\n", token, f.RuleID, f.File, f.Line, f.Message); err != nil {
			return err
		}

		switch f.Severity {
		case lint.SeverityError:
			errCount++
		case lint.SeverityWarning:
			warnCount++
		}
	}

	_, err := fmt.Fprintf(w, "%d error(s), %d warning(s)\n", errCount, warnCount)
	return err
}

// severityColor returns the ANSI color code for a severity, or "" for an
// unrecognized value (rendered uncolored rather than erroring).
func severityColor(s lint.Severity) string {
	switch s {
	case lint.SeverityError:
		return ansiRed
	case lint.SeverityWarning:
		return ansiYellow
	default:
		return ""
	}
}
