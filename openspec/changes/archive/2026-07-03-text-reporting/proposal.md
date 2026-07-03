## Why

`schema-validation-rules` produces `[]lint.Finding`, but nothing yet turns
that into something a human (or a CI log) can read. `scan-command` — the
first end-to-end demoable slice of the product — needs a text renderer to
call before it can be wired up. Building it as its own capability, decoupled
from any CLI command, means `scan` and future consumers (e.g. `sarif-output`
later reusing the same `[]Finding` input) don't duplicate formatting logic.

## What Changes

- Add a new `internal/reporter` package with a `RenderText(findings
  []lint.Finding, w io.Writer)` function (exact signature decided in
  design.md) that formats `[]lint.Finding` as human-readable text.
- Implement the per-line format (FR-OUT-01): one line per finding —
  `severity  rule-id  file:line — message`.
- Implement the summary line (FR-OUT-02): `N error(s), M warning(s)` after
  all finding lines.
- Implement the success line (FR-OUT-04): when `findings` is empty, print a
  single `✓ AGENTS.md is valid (N rules passed)` line instead of the
  per-line/summary output.
- Respect `NO_COLOR` and keep output pager-free (BC-OUTPUT-01): color is
  optional decoration only, never required to read the output, and is
  suppressed when the `NO_COLOR` environment variable is set.
- Preserve each `Finding.Message` verbatim so rule-authored, actionable
  messages (BC-OUTPUT-02) reach the terminal unmodified — the reporter does
  not truncate, reword, or drop message text.
- Add golden-file tests under `testdata/reporter/golden/` per the testing
  convention in `AGENTS.md` (no findings, one error, one warning, mixed).
- No CLI wiring — `agents-lint scan` calling this reporter is
  `scan-command`'s job, not this change's.

## Capabilities

### New Capabilities
- `text-reporting`: renders `[]lint.Finding` to human-readable text
  (per-line format, summary line, success line), decoupled from any CLI
  command so both `scan` and future consumers can reuse it.

### Modified Capabilities
(none — this change only reads `lint.Finding`/`lint.Severity` from
`schema-validation-rules`; it does not change their definition)

## Impact

- Affected code: new `internal/reporter/` package only (no changes to
  `internal/lint` or `internal/lint/rules`).
- New test fixtures: `testdata/reporter/golden/*.txt`.
- Dependencies introduced: none — text formatting uses only the standard
  library.
- No user-facing behavior yet: there is still no `scan` command to invoke
  this reporter from — that's `scan-command`, which depends on this change.
