## Purpose

Human-readable text rendering of `[]lint.Finding`, decoupled from any CLI
command so both `scan` and future consumers (e.g. `sarif-output` reusing the
same `[]Finding` input) can share it.

## Requirements

### Requirement: FR-OUT-01 — Per-finding text line
For each `Finding` in the input slice, `reporter.WriteText` SHALL write exactly one line to the output in the form `severity  rule-id  file:line — message` (two spaces between `severity`, `rule-id`, and `file:line`; ` — ` before `message`), using the `Finding`'s `Severity`, `RuleID`, `File`, `Line`, and `Message` fields verbatim, in the order the findings were supplied.

#### Scenario: Single error finding
- **WHEN** `WriteText` is called with one `Finding{RuleID: "S002", Severity: SeverityError, File: "AGENTS.md", Line: 3, Message: "missing ## Agents section"}`
- **THEN** the output contains the line `error  S002  AGENTS.md:3 — missing ## Agents section`

#### Scenario: Multiple findings preserve input order
- **WHEN** `WriteText` is called with findings in the order `[S002, S004, S005]`
- **THEN** the per-finding lines appear in the output in that same order, with no reordering or grouping by severity

### Requirement: FR-OUT-02 — Summary line
After all per-finding lines, `reporter.WriteText` SHALL write a summary line in the literal form `N error(s), M warning(s)`, where `N` is the count of findings with `Severity == SeverityError` and `M` is the count of findings with `Severity == SeverityWarning`, whenever `findings` is non-empty.

#### Scenario: Mixed severities
- **WHEN** `WriteText` is called with 2 findings of `SeverityError` and 1 finding of `SeverityWarning`
- **THEN** the output ends with the line `2 error(s), 1 warning(s)`

#### Scenario: Only warnings
- **WHEN** `WriteText` is called with 0 error-severity findings and 3 warning-severity findings
- **THEN** the output ends with the line `0 error(s), 3 warning(s)`

### Requirement: FR-OUT-04 — Success line on no findings
When `findings` is empty, `reporter.WriteText` SHALL write a single line, `✓ AGENTS.md is valid (N rules passed)`, where `N` is the `ruleCount` argument, and SHALL NOT write any per-finding lines or a summary line.

#### Scenario: No findings produces success line only
- **WHEN** `WriteText` is called with an empty `findings` slice and `ruleCount = 5`
- **THEN** the entire output is the single line `✓ AGENTS.md is valid (5 rules passed)`

### Requirement: BC-OUTPUT-01 — Readable output, optional color, NO_COLOR respected
`reporter.WriteText` SHALL produce output that is readable without a pager (no ANSI codes when color is disabled) and SHALL support optional ANSI severity coloring via `Options.Color`; when the `NO_COLOR` environment variable is set to any non-empty value, `reporter.WriteText` SHALL NOT emit ANSI color codes regardless of `Options.Color`.

#### Scenario: Color disabled by default
- **WHEN** `WriteText` is called with `Options{Color: false}`
- **THEN** the output contains no ANSI escape sequences

#### Scenario: Color requested and NO_COLOR unset
- **WHEN** `WriteText` is called with `Options{Color: true}` and the `NO_COLOR` environment variable is unset
- **THEN** the severity token in each per-finding line is wrapped in ANSI color codes

#### Scenario: NO_COLOR overrides Options.Color
- **WHEN** `WriteText` is called with `Options{Color: true}` and the `NO_COLOR` environment variable is set to a non-empty value
- **THEN** the output contains no ANSI escape sequences

### Requirement: BC-OUTPUT-02 — Messages preserved verbatim
`reporter.WriteText` SHALL write each `Finding.Message` unmodified — it SHALL NOT truncate, reword, or drop message text — so that actionable messages authored by rules reach the output intact.

#### Scenario: Long message is not truncated
- **WHEN** `WriteText` is called with a `Finding` whose `Message` is longer than a typical terminal width (e.g. 200 characters)
- **THEN** the full, untruncated message text appears in the corresponding output line
