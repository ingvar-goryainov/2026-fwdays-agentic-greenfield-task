## 1. SARIF writer (FR-OUT-03)

- [x] 1.1 Add `internal/reporter/sarif.go` with the SARIF v2.1.0 structs (`sarifLog`, `sarifRun`, `sarifTool`, `sarifDriver`, `sarifReportingDescriptor`, `sarifResult`, `sarifMessage`, `sarifLocation`, `sarifPhysicalLocation`, `sarifArtifactLocation`, `sarifRegion`) and `WriteSARIF(w io.Writer, findings []lint.Finding) error`
- [x] 1.2 Implement severity → SARIF `level` mapping (`SeverityError` → `"error"`, `SeverityWarning` → `"warning"`)
- [x] 1.3 Implement per-result `physicalLocation`: `artifactLocation.uri` from `Finding.File`; include `region.startLine` only when `Finding.Line > 0`
- [x] 1.4 Implement `tool.driver.rules` de-duplication in first-seen order from the distinct `RuleID`s in `findings`
- [x] 1.5 Add `internal/reporter/sarif_test.go` covering: single finding, findings-order preservation, empty findings (empty `results` and empty `rules`), error/warning level mapping, `Line == 0` omitting `region`, and repeated `RuleID` collapsing to one rule-catalog entry
- [x] 1.6 Run `go test ./internal/reporter/... -run TestWriteSARIF` (or the actual test names) until green

## 2. `--format` flag (FR-CLI-05)

- [x] 2.1 Add a persistent `--format` string flag to `cmd/agents-lint/cmd/root.go` (default `"text"`), alongside the existing `--config` flag
- [x] 2.2 Add format validation (accepts only `"text"`/`"sarif"`) that runs before `runScan` does any file I/O, producing a clear error naming the invalid value and the accepted values on stderr with a non-zero exit
- [x] 2.3 Add/extend command-level tests in `cmd/agents-lint/cmd` for: default format, `--format sarif` accepted, unknown `--format` value rejected with no scan performed

## 3. Wire format into `scan` (scan-command delta)

- [x] 3.1 Update `runScan` in `cmd/agents-lint/cmd/scan.go` to accept the resolved format and branch between `reporter.WriteText` and `reporter.WriteSARIF`
- [x] 3.2 Update `scanCmd`'s `RunE` to pass the `--format` flag value through to `runScan`
- [x] 3.3 Add/extend tests in `cmd/agents-lint/cmd/scan_test.go` for: default text output unchanged, `--format sarif` producing SARIF output, exit code identical across formats for the same findings

## 4. Verify

- [x] 4.1 Run full suite: `go test ./...`
- [x] 4.2 Run `go test -cover ./internal/...` and confirm `internal/reporter` stays ≥ 90% line coverage
- [x] 4.3 Run `golangci-lint run` and fix any findings
- [x] 4.4 Manually run `go run ./cmd/agents-lint scan --format sarif` against a fixture AGENTS.md and confirm the output is valid JSON matching the SARIF v2.1.0 shape
