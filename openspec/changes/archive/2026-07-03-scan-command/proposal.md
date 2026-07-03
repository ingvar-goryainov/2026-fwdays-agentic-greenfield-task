## Why

`schema-validation-rules` produces `[]lint.Finding` and `text-reporting`
renders it as text, but there is still no way to invoke either from the
command line — `cmd/agents-lint/cmd` only has an empty root command. `scan`
is the first end-to-end usable slice of the product (per
`docs/openspec-capability-plan.md`'s build sequence): wiring the parser,
rule engine, and reporter together into `agents-lint scan [path]` is what
makes the "pass or fail like a compiler" pitch demoable for the first time.

## What Changes

- Add `agents-lint scan [path]` as a new Cobra subcommand in
  `cmd/agents-lint/cmd/scan.go`.
- Implement FR-CLI-01: `path` defaults to `./AGENTS.md` when omitted; when
  given, it's used as-is (no additional resolution beyond what
  `rules.Run`/`os.Stat` already do).
- Implement FR-CLI-02: the process exits `0` if no finding has
  `SeverityError`, and exits `1` if at least one does — independent of
  whether `Run` returns a Go `error` (a distinct exit path, since that's an
  unexpected failure, not a lint result).
- Wire `internal/lint/rules.Run(path)` → `internal/reporter.WriteText` so
  `scan`'s stdout output matches the FR-OUT-01/02/04 format already
  implemented and tested by `text-reporting`.
- Verify NFR-PERF-01 (full scan of a ≤500-line AGENTS.md completes in
  < 100ms) with a benchmark or timed test against a realistic fixture —
  this capability is the first point where the requirement is actually
  exercisable end-to-end.
- No new flags (`--format`, `--config`, `--version`) — those belong to
  `sarif-output`, `configuration-support`, and `cli-version-command`
  respectively.

## Capabilities

### New Capabilities
- `scan-command`: the `agents-lint scan [path]` CLI command — wires the
  existing rule engine and text reporter together with FR-CLI-01/02's
  path-resolution and exit-code semantics.

### Modified Capabilities
(none — this change is pure wiring; it does not alter the behavior
contracts already specified by `schema-validation-rules` or
`text-reporting`)

## Impact

- Affected code: new `cmd/agents-lint/cmd/scan.go` (registered on
  `rootCmd` in `cmd/agents-lint/cmd/root.go`); no changes to
  `internal/lint`, `internal/lint/rules`, or `internal/reporter`.
- Dependencies introduced: none — reuses `rules.Run` and
  `reporter.WriteText`, both already in the module.
- User-facing behavior: `agents-lint scan` becomes a real, runnable command
  for the first time. `agents-lint` with no subcommand still just prints
  Cobra's help/usage (unchanged).
