## ADDED Requirements

### Requirement: FR-OUT-03 — SARIF v2.1.0 output
`internal/reporter.WriteSARIF` SHALL render a `[]lint.Finding` slice as a JSON document conforming to the SARIF v2.1.0 schema: a top-level `$schema` and `version: "2.1.0"`, exactly one entry in `runs`, that run's `tool.driver.name` equal to `"agents-lint"`, and one entry in `results` for each input `Finding`, in the order the findings were supplied.

#### Scenario: Single finding renders one SARIF result
- **WHEN** `WriteSARIF` is called with one `Finding{RuleID: "S002", Severity: SeverityError, File: "AGENTS.md", Line: 3, Message: "missing ## Agents section"}`
- **THEN** the output is a JSON document with `version: "2.1.0"`, `runs[0].tool.driver.name: "agents-lint"`, and exactly one entry in `runs[0].results` whose `ruleId` is `"S002"` and whose `message.text` is `"missing ## Agents section"`

#### Scenario: Multiple findings preserve input order
- **WHEN** `WriteSARIF` is called with findings in the order `[S002, S004, S005]`
- **THEN** `runs[0].results` contains the corresponding SARIF results in that same order

#### Scenario: No findings produces an empty results array
- **WHEN** `WriteSARIF` is called with an empty `findings` slice
- **THEN** the output is a valid SARIF v2.1.0 document whose `runs[0].results` is an empty array

### Requirement: FR-OUT-03 — Severity mapped to SARIF level
Each SARIF result's `level` SHALL be derived from the source `Finding.Severity`: `SeverityError` SHALL map to `"error"` and `SeverityWarning` SHALL map to `"warning"`.

#### Scenario: Error severity maps to error level
- **WHEN** `WriteSARIF` is called with a `Finding` whose `Severity` is `SeverityError`
- **THEN** the corresponding SARIF result's `level` is `"error"`

#### Scenario: Warning severity maps to warning level
- **WHEN** `WriteSARIF` is called with a `Finding` whose `Severity` is `SeverityWarning`
- **THEN** the corresponding SARIF result's `level` is `"warning"`

### Requirement: FR-OUT-03 — Result location from File and Line
Each SARIF result SHALL include a `locations` entry whose `physicalLocation.artifactLocation.uri` is the source `Finding.File` verbatim. When `Finding.Line` is greater than zero, the `physicalLocation` SHALL also include a `region.startLine` equal to `Finding.Line`. When `Finding.Line` is zero, the `physicalLocation` SHALL omit `region` entirely rather than emitting `region.startLine: 0`.

#### Scenario: Finding with a line number includes a region
- **WHEN** `WriteSARIF` is called with a `Finding{File: "AGENTS.md", Line: 12, ...}`
- **THEN** the corresponding result's `physicalLocation.artifactLocation.uri` is `"AGENTS.md"` and `physicalLocation.region.startLine` is `12`

#### Scenario: Finding with no line number omits the region
- **WHEN** `WriteSARIF` is called with a `Finding{File: "AGENTS.md", Line: 0, ...}` (e.g. FR-S001's file-not-found finding)
- **THEN** the corresponding result's `physicalLocation.artifactLocation.uri` is `"AGENTS.md"` and `physicalLocation` has no `region` field

### Requirement: FR-OUT-03 — Rule catalog reflects distinct rule IDs present
`runs[0].tool.driver.rules` SHALL contain exactly one entry per distinct `RuleID` present in the input `findings`, each with `id` equal to that `RuleID`, in first-seen order, with no duplicates and no entries for rule IDs absent from `findings`.

#### Scenario: Repeated rule ID appears once in the catalog
- **WHEN** `WriteSARIF` is called with findings `[{RuleID: "S002"}, {RuleID: "S004"}, {RuleID: "S002"}]`
- **THEN** `runs[0].tool.driver.rules` contains exactly two entries, with `id` values `"S002"` and `"S004"`, in that order

#### Scenario: Empty findings produces an empty rule catalog
- **WHEN** `WriteSARIF` is called with an empty `findings` slice
- **THEN** `runs[0].tool.driver.rules` is an empty array

### Requirement: FR-CLI-05 — `--format` flag selects the output format
`agents-lint --format <value>` SHALL accept `text` or `sarif` and SHALL default to `text` when the flag is omitted. A value other than `text` or `sarif` SHALL cause the command to fail with a clear, actionable error identifying the invalid value and the accepted values, written to stderr, exiting non-zero, before any file is scanned.

#### Scenario: Default format is text
- **WHEN** `agents-lint scan` is run with no `--format` flag
- **THEN** the command produces `reporter.WriteText`-formatted output

#### Scenario: Explicit sarif format selects the SARIF reporter
- **WHEN** `agents-lint scan --format sarif` is run
- **THEN** the command produces `reporter.WriteSARIF`-formatted output instead of text output

#### Scenario: Unknown format value fails clearly and scans nothing
- **WHEN** `agents-lint scan --format xml` is run
- **THEN** the command writes an error to stderr naming `"xml"` and the accepted values (`text`, `sarif`), exits with a non-zero status, and does not perform the scan
