## 1. `runScan` core logic (FR-CLI-01, FR-CLI-02, reporter wiring)

- [x] 1.1 Implement `runScan(w io.Writer, path string) (exitCode int, err error)`
      in `cmd/agents-lint/cmd/scan.go`: calls `rules.Run(path)`, computes
      `ruleCount := len(rules.DefaultRules()) + 1`, calls
      `reporter.WriteText(w, findings, ruleCount, reporter.Options{})`,
      returns `exitCode = 1` if any finding has `Severity ==
      lint.SeverityError`, else `0`
- [x] 1.2 Implement the `rules.Run` error path: return `(1, fmt.Errorf("scan
      %s: %w", path, err))` without calling `reporter.WriteText`
- [x] 1.3 Table-driven tests in `cmd/agents-lint/cmd/scan_test.go` for
      `runScan`: clean file → `(0, nil)`, file with an error-severity
      finding → `(1, nil)`, file with only warning-severity findings →
      `(0, nil)`, unreadable file → `(1, non-nil error)`. The severity
      decision was extracted into `exitCodeForFindings` and tested directly
      with synthetic findings (`TestExitCodeForFindings`), since no rule in
      today's S001-S005 set emits `SeverityWarning` to exercise that path
      through a real file — see design.md amendment below.
- [x] 1.4 Test that `runScan`'s stdout output (captured via a
      `bytes.Buffer`) matches what `reporter.WriteText` would produce
      directly for the same findings/ruleCount (confirms wiring, not
      reporter formatting itself — that's already covered by
      `internal/reporter`'s own tests)

## 2. Cobra command registration

- [x] 2.1 Define `scanCmd` in `cmd/agents-lint/cmd/scan.go`: `Use: "scan
      [path]"`, `Args: cobra.MaximumNArgs(1)`, defaults `path` to
      `"./AGENTS.md"` when no argument is given (FR-CLI-01)
- [x] 2.2 `RunE` closure: call `runScan(cmd.OutOrStdout(), path)`; if `err !=
      nil`, print it to `cmd.ErrOrStderr()` and return it; else if
      `exitCode != 0`, call `os.Exit(exitCode)`; else return `nil`
- [x] 2.3 Register `scanCmd` on `rootCmd` via `rootCmd.AddCommand(scanCmd)`
      in `cmd/agents-lint/cmd/root.go` (or an `init()` in `scan.go`,
      matching whichever registration pattern `root.go` already uses)
- [x] 2.4 Manual check: `go run ./cmd/agents-lint scan` against this repo's
      own `AGENTS.md` prints the success line and exits `0`
      (`echo $?` after running)
- [x] 2.5 Manual check: `go run ./cmd/agents-lint scan
      testdata/S002/invalid.md` prints an error-severity finding line and
      exits `1`

## 3. NFR-PERF-01 verification

- [x] 3.1 Generate a ~500-line synthetic AGENTS.md fixture at
      `testdata/scan/perf/large.md` (multiple agent blocks, representative
      of real content, sized to the NFR-PERF-01 boundary) — 514 lines, 34
      agent blocks
- [x] 3.2 Add `BenchmarkRun_LargeFile` in
      `internal/lint/rules/registry_test.go` (or a new
      `registry_bench_test.go`) benchmarking `rules.Run` against the
      fixture
- [x] 3.3 Run `go test ./internal/lint/rules/... -bench=BenchmarkRun_LargeFile
      -benchtime=10x` and confirm `ns/op` is well under the 100ms
      (100,000,000 ns) NFR-PERF-01 budget; note the result in the PR/commit,
      not asserted as a hard test failure (see design.md rationale) — result:
      ~898,000 ns/op (~0.9ms), well under the 100ms budget

## 4. Full verification

- [x] 4.1 Run `go test ./cmd/... -run TestRunScan` until green
- [x] 4.2 Run `go test ./...` — confirm green
- [x] 4.3 Run `go test -cover ./cmd/...` — confirm `runScan` itself is at or
      above 90% line coverage (the `RunE`/`os.Exit` glue in `scanCmd` is
      thin wiring and is not required to be unit-tested per design.md) —
      `runScan` and `exitCodeForFindings` are both at 100%
- [x] 4.4 Run `golangci-lint run` — confirm clean (caught and fixed one
      `errcheck` issue: unchecked `fmt.Fprintln` return in the `RunE`
      closure)
- [x] 4.5 Confirm `cmd/agents-lint/cmd/scan.go` is the only new file
      touching CLI wiring — no changes to `internal/lint`,
      `internal/lint/rules`, or `internal/reporter` behavior
- [x] 4.6 Confirm no new entries were added to `go.mod`
