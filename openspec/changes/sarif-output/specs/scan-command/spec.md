## MODIFIED Requirements

### Requirement: Scan output uses the resolved format's reporter
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
