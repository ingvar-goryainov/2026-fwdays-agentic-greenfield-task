## Why

`agents-lint scan` only renders findings as human-readable text
(`internal/reporter.WriteText`). Per `docs/openspec-capability-plan.md`,
`sarif-output` is the next capability after `codebase-awareness-rules`:
CI-native tooling (GitHub Code Scanning, VS Code's SARIF viewer) can't
consume text output, so there is no way to wire `agents-lint` into those
workflows without a machine-readable format. Both rule layers (schema and
codebase-awareness) now exist, so the SARIF writer can be built against the
full finding surface rather than a partial one.

## What Changes

- Add a new `internal/reporter/sarif.go` writer,
  `WriteSARIF(w io.Writer, findings []lint.Finding) error`, that renders
  `[]lint.Finding` as a SARIF v2.1.0 log (FR-OUT-03): one `run`, one `tool.driver`
  (name `agents-lint`), one `result` per finding (`ruleId`, `level`,
  `message.text`, and a `physicalLocation` with `artifactLocation.uri` and
  `region.startLine` derived from `File`/`Line`), and a `tool.driver.rules`
  entry for every distinct `RuleID` present in `findings`.
- Map `lint.Severity` to SARIF `level`: `SeverityError` → `"error"`,
  `SeverityWarning` → `"warning"`.
- Implement FR-CLI-05: add a persistent `--format [text|sarif]` flag on the
  root command, default `text`. An unrecognized value fails with a clear,
  actionable error before any scan runs.
- Wire the flag into `agents-lint scan`: when `--format sarif` is given, the
  command writes `WriteSARIF`'s output instead of `reporter.WriteText`'s;
  the FR-CLI-02 exit-code contract (0/1 based on `SeverityError` findings)
  is unchanged and applies identically regardless of format.

## Capabilities

### New Capabilities
- `sarif-output`: the SARIF v2.1.0 writer (FR-OUT-03) and the `--format`
  flag (FR-CLI-05) that selects it, plus `scan`'s wiring to choose between
  the text and SARIF reporters.

### Modified Capabilities
- `scan-command`: the existing "Scan output uses the text reporter"
  requirement becomes conditional on the resolved `--format` (default
  `text`, using `reporter.WriteText`; `sarif` using
  `reporter.WriteSARIF`). FR-CLI-01/02 (default path, exit-code) are
  unchanged in either format.

## Impact

- Affected code: new `internal/reporter/sarif.go` (+ tests);
  `cmd/agents-lint/cmd/root.go` (new persistent `--format` flag, validated
  before scan runs); `cmd/agents-lint/cmd/scan.go` (reporter selection
  based on the resolved format).
- Dependencies introduced: none — SARIF v2.1.0 is a JSON document; rendered
  with `encoding/json` from the standard library (no new SARIF-specific
  dependency).
- User-facing behavior: `agents-lint scan --format sarif` prints a SARIF
  v2.1.0 log instead of text output. Default behavior (`--format` omitted)
  is unchanged.
