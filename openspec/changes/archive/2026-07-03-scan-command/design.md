## Context

Two capabilities already exist and are unwired: `internal/lint/rules.Run(path
string) ([]lint.Finding, error)` (schema rules, fixed S001→S002→S003→S004→S005
order) and `internal/reporter.WriteText(w io.Writer, findings []lint.Finding,
ruleCount int, opts Options) error` (FR-OUT-01/02/04 text rendering). The CLI
side (`cmd/agents-lint/cmd/root.go`) is a bare Cobra root with no
subcommands. This change's only job is connecting the three: parse args →
call `Run` → call `WriteText` → set the process exit code per FR-CLI-02.

The one design wrinkle is exit-code control: FR-CLI-02 requires exit `0`/`1`
based on finding severity, which is a normal *result*, not a Go `error` —
but Cobra's `RunE` return value only cleanly maps to "error or not," and
`main.go` currently does `if err := cmd.Execute(); err != nil { os.Exit(1)
}`, which conflates "the command failed to run" with "the command ran fine
and found lint errors."

## Goals / Non-Goals

**Goals:**
- `agents-lint scan [path]` wired end-to-end: argument parsing, `rules.Run`,
  `reporter.WriteText`, exit code.
- FR-CLI-01: `path` defaults to `./AGENTS.md`; an explicit argument is
  passed through unchanged.
- FR-CLI-02: exit `0` when no `SeverityError` finding is present, exit `1`
  when at least one is — decided from `[]lint.Finding` returned by `Run`,
  not by whether `Run`'s own `error` return is nil.
- A genuine `Run` error (e.g. read-permission failure — see
  `TestRun_UnreadableFilePropagatesError` in
  `internal/lint/rules/registry_test.go`) is reported on stderr and also
  exits non-zero, distinguishable in code from the "found lint issues" path
  even though both currently exit `1`.
- NFR-PERF-01 verified: a benchmark/timed test demonstrates a ≤500-line
  fixture scans in < 100ms.

**Non-Goals:**
- No `--format`, `--config`, `--version`, or `--color` flags — `scan` takes
  only an optional positional `path` argument in this change.
- No config-file loading (`.agents-lint.yaml`) — that's
  `configuration-support`, sequenced after this change per the capability
  plan.
- No codebase-awareness rules (FR-C001/C002) in the scan — `rules.Run`
  already fixes the rule set to S001–S005; layer 2 is
  `codebase-awareness-rules`.

## Decisions

**Exit-code decision lives in a pure, testable function; `os.Exit` is an
untested one-line wrapper.** `cmd/agents-lint/cmd/scan.go` splits the work
into:
```go
// runScan does the actual work and returns the process exit code plus any
// unexpected (non-lint) error. It never calls os.Exit, so it's fully
// unit-testable.
func runScan(w io.Writer, path string) (exitCode int, err error) {
    findings, runErr := rules.Run(path)
    if runErr != nil {
        return 1, fmt.Errorf("scan %s: %w", path, runErr)
    }

    ruleCount := len(rules.DefaultRules()) + 1 // +1 for S001
    if writeErr := reporter.WriteText(w, findings, ruleCount, reporter.Options{}); writeErr != nil {
        return 1, writeErr
    }

    for _, f := range findings {
        if f.Severity == lint.SeverityError {
            return 1, nil
        }
    }
    return 0, nil
}
```
The Cobra command's `RunE` is a thin wrapper: call `runScan`, print `err` to
stderr and return it if non-nil (so Cobra's own error path still fires for
genuinely unexpected failures), otherwise call `os.Exit(exitCode)` directly
when `exitCode != 0`. Rationale: this keeps 100% of the exit-code decision
logic (the part FR-CLI-02 actually specifies) reachable by ordinary table
tests against `runScan`'s return values, and confines the untestable
`os.Exit` call to a single trivial line consistent with how `main.go`
already isn't unit-tested. Alternative considered: return a Cobra-idiomatic
sentinel `error` (e.g. `errLintFailed`) from `RunE` and have `main.go`
special-case it to choose the exit code — rejected because it forces
`main.go` to import lint-specific sentinels, coupling the generic
entrypoint to one subcommand's semantics, for no benefit over returning an
`int` directly from the testable helper.

**`ruleCount` is `len(rules.DefaultRules()) + 1`.** This was flagged as an
open decision in `text-reporting`'s design (`WriteText`'s `ruleCount`
parameter is caller-supplied, not inferred). `rules.DefaultRules()` returns
the four `Rule` interface implementations (S002–S005); S001 runs outside
that interface (see `schema-validation-rules`'s design), so the true count
of rules evaluated by `Run` is `len(rules.DefaultRules()) + 1` — computed
inline in `runScan` rather than adding a new exported function to
`internal/lint/rules`, since it's a one-line derivation from an already
-exported function and doesn't justify a new API surface.

**Default path is the literal string `"./AGENTS.md"`.** FR-CLI-01 specifies
this exact default. `cobra.Command.Args: cobra.MaximumNArgs(1)` allows zero
or one positional argument; when zero, `runScan` receives `"./AGENTS.md"`
directly (not resolved to an absolute path — `os.Stat` and `rules.Run`
already handle relative paths correctly, and resolving it further would
diverge from what a user typed without any spec requiring it).

**Output goes to `cmd.OutOrStdout()`, not raw `os.Stdout`.** Cobra commands
expose `OutOrStdout()`/`OutOrStderr()` so tests can inject a buffer via
`cmd.SetOut(...)` without redirecting the real process's file descriptors.
`runScan` itself takes a plain `io.Writer` (matching `reporter.WriteText`'s
signature) so it doesn't need to know about Cobra at all — only the
`RunE` closure resolves `cmd.OutOrStdout()` and passes it in.

**NFR-PERF-01 verification: a `testing.B` benchmark, not a hard-coded
duration assertion in a regular test.** A ~500-line synthetic AGENTS.md
fixture (`testdata/scan/perf/large.md`) is scanned via `rules.Run` inside
`BenchmarkRun_LargeFile` in `internal/lint/rules` (this capability is where
NFR-PERF-01 first becomes exercisable end-to-end, even though the benchmark
itself lives next to the code it measures). Rationale: asserting a hard
wall-clock threshold inside `go test` is flaky across machines/CI runners;
a benchmark reports `ns/op` for a human (or a future CI perf-budget job) to
check against the 100ms target without failing the suite on slower
hardware. Alternative considered: a regular test with `time.Since` and a
hard `assert.Less(t, elapsed, 100*time.Millisecond)` — rejected as exactly
the kind of flaky, hardware-dependent assertion the Go benchmark mechanism
exists to avoid.

## Risks / Trade-offs

- **[Risk]** Calling `os.Exit` directly inside `RunE` bypasses any Cobra
  cleanup/deferred behavior that might be added later (there is none
  today). → **Mitigation**: the exit call is the last statement in a
  single-purpose command's `RunE`, after all writes have completed and been
  error-checked; if a future change needs cleanup hooks, this is the one
  place to add them.
- **[Risk]** A genuine `rules.Run` error and a "found error-severity
  findings" result currently both exit `1`, which FR-CLI-02 doesn't
  distinguish between (it only specifies 0-vs-1 on error severity). →
  **Mitigation**: `runScan`'s two return paths are distinguished internally
  (one returns a non-nil `error`, the other doesn't) even though the
  observable exit code is the same today; this leaves room to give
  unexpected failures a distinct exit code later without changing
  `runScan`'s tests.
- **[Trade-off]** No `--path`/`-p` flag, only a positional argument — matches
  FR-CLI-01's literal `scan [path]` syntax; a flag form isn't specified and
  isn't added speculatively.

## Open Questions

(none — the `ruleCount` question raised in `text-reporting`'s design is
resolved above)

## Amendments

- **2026-07-03, during apply:** the "Exit-code decision lives in a pure,
  testable function" decision originally had `runScan` compute the
  `SeverityError` check inline. While writing task 1.3's table tests it
  became clear that no rule in the current S001-S005 set ever produces a
  `SeverityWarning` finding, so there is no real fixture file that can
  exercise `runScan`'s "warning-only findings still exit 0" behavior
  end-to-end. Fixed by extracting that one loop into
  `exitCodeForFindings(findings []lint.Finding) int`, called by `runScan`,
  and unit-testing it directly against synthetic `[]lint.Finding` values
  (including a `SeverityWarning`-only case) in
  `TestExitCodeForFindings`. `runScan`'s own tests still use real fixture
  files for the clean/error/unreadable cases. This doesn't change any
  externally observable behavior — it's a decomposition, not a semantics
  change — so the spec (FR-CLI-02) is unaffected.
