## Purpose

The `agents-lint scan` CLI command, wiring `internal/lint/rules.Run` and
`internal/reporter.WriteText` together to validate an AGENTS.md file and
report findings to the user with an exit code reflecting error-severity
results.

## Requirements

### Requirement: FR-CLI-01 — `scan [path]` with default target
`agents-lint scan [path]` SHALL validate the AGENTS.md file at `path` when given, and SHALL default `path` to `./AGENTS.md` when no positional argument is provided.

#### Scenario: No path argument uses the default
- **WHEN** `agents-lint scan` is run with no arguments in a directory containing `./AGENTS.md`
- **THEN** the command validates `./AGENTS.md`

#### Scenario: Explicit path argument is used as-is
- **WHEN** `agents-lint scan path/to/other.md` is run
- **THEN** the command validates `path/to/other.md` instead of the default

### Requirement: FR-CLI-02 — Exit code reflects error-severity findings
`agents-lint scan` SHALL exit with status `0` when the findings produced for the target file contain no `SeverityError` finding, and SHALL exit with status `1` when at least one `SeverityError` finding is present, regardless of how many warning-severity findings exist.

#### Scenario: Clean file exits 0
- **WHEN** `agents-lint scan` is run against an AGENTS.md file that produces zero findings
- **THEN** the process exits with status `0`

#### Scenario: Error-severity finding exits 1
- **WHEN** `agents-lint scan` is run against an AGENTS.md file that produces at least one `SeverityError` finding
- **THEN** the process exits with status `1`

#### Scenario: Warning-only findings still exit 0
- **WHEN** `agents-lint scan` is run against an AGENTS.md file that produces only `SeverityWarning` findings and no `SeverityError` findings
- **THEN** the process exits with status `0`

### Requirement: Scan output uses the text reporter
`agents-lint scan` SHALL render its findings using the reporter selected by the resolved `--format` value (FR-CLI-05): `text` (the default) SHALL use `internal/reporter.WriteText`, passing the total number of rules evaluated (`len(rules.DefaultRules()) + 1`) as the `ruleCount` argument, so the command's stdout output matches the FR-OUT-01/02/04 format; `sarif` SHALL use `internal/reporter.WriteSARIF` (FR-OUT-03). The choice of format SHALL NOT change which findings are produced or the FR-CLI-02 exit-code decision.

#### Scenario: Findings rendered in text reporter format by default
- **WHEN** `agents-lint scan` runs against a file with findings and no `--format` flag is given
- **THEN** stdout contains the per-finding lines and summary line produced by `reporter.WriteText` for those findings

#### Scenario: Clean file prints the text success line by default
- **WHEN** `agents-lint scan` runs against a file with zero findings and no `--format` flag is given
- **THEN** stdout contains the single success line `✓ AGENTS.md is valid (N rules passed)`, where `N` equals `len(rules.DefaultRules()) + 1`

#### Scenario: Findings rendered as SARIF when requested
- **WHEN** `agents-lint scan --format sarif` runs against a file with findings
- **THEN** stdout contains the `reporter.WriteSARIF` output for those findings instead of text-formatted output

#### Scenario: Exit code is unaffected by the output format
- **WHEN** `agents-lint scan --format sarif` runs against an AGENTS.md file that produces at least one `SeverityError` finding
- **THEN** the process exits with status `1`, exactly as it would with the default text format

### Requirement: Unexpected scan errors are reported and exit non-zero
When `internal/lint/rules.Run` returns a non-nil `error` (a failure unrelated to the file simply not existing, e.g. a read-permission error), `agents-lint scan` SHALL write a message describing the failure to stderr and SHALL exit with a non-zero status, without printing any reporter-formatted finding output.

#### Scenario: Unreadable file reports an error and exits non-zero
- **WHEN** `agents-lint scan` is run against a path that exists but cannot be read (e.g. permission denied)
- **THEN** the process writes an error message to stderr, exits with a non-zero status, and does not print `reporter.WriteText` output
