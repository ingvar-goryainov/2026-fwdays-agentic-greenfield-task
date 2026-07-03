## Context

`schema-validation-rules` (archived) established `internal/lint.Finding` /
`internal/lint.Severity` and `internal/lint/rules.Run(path string)
([]lint.Finding, error)`. Nothing yet turns a `[]lint.Finding` into text.
`scan-command` will be the first caller, but per the repo's library-first
architecture (`internal/` returns typed results, no `fmt.Print` — see
`docs/requirements.md`'s architecture note and `AGENTS.md`'s code
conventions), the renderer itself has to live in `internal/`, write to a
caller-supplied `io.Writer`, and stay ignorant of Cobra/CLI wiring.

FR-OUT-04's success line — `✓ AGENTS.md is valid (N rules passed)` — needs a
rule count, which a `[]lint.Finding` alone can't supply (an empty slice
doesn't say whether 0 rules ran or all of them passed). That forces the
renderer's signature to accept the total rule count as an explicit input
rather than inferring it from the findings.

## Goals / Non-Goals

**Goals:**
- One function, `reporter.WriteText`, that renders `[]lint.Finding` plus a
  rule count to an `io.Writer` in the FR-OUT-01/02/04 text format.
- Deterministic output: same inputs always produce the same bytes, so it can
  be golden-file tested (`testdata/reporter/golden/*.txt` per `AGENTS.md`'s
  testing conventions).
- Optional ANSI color that's off by default and forcibly disabled when
  `NO_COLOR` is set (BC-OUTPUT-01), implemented with the standard library
  only — no new dependency.
- Message text passed through unmodified (BC-OUTPUT-02 is a rule-authoring
  concern; the reporter's job is to not undermine it by truncating or
  reformatting `Finding.Message`).

**Non-Goals:**
- No CLI wiring, no `--format` flag, no TTY auto-detection — those belong to
  `scan-command` (text is the only format anyway; `--format sarif` is
  `sarif-output`, out of scope here).
- No sorting or de-duplication of findings — `rules.Run` already returns
  them in a fixed, deterministic rule order; the reporter renders what it's
  given, in order, without imposing a second ordering policy.
- No SARIF or JSON output — that's `sarif-output` (FR-OUT-03), a separate
  capability reusing the same `[]lint.Finding` input later.

## Decisions

**Package and signature.** New package `internal/reporter`, one function:
```go
package reporter

// Options controls optional rendering behavior.
type Options struct {
    // Color enables ANSI severity coloring. Ignored (treated as false)
    // when the NO_COLOR environment variable is set.
    Color bool
}

// WriteText renders findings as human-readable text to w, following
// FR-OUT-01 (per-line format), FR-OUT-02 (summary line), and FR-OUT-04
// (success line). ruleCount is the total number of rules evaluated to
// produce findings; it is only used to render the success line's "(N
// rules passed)" when findings is empty.
func WriteText(w io.Writer, findings []lint.Finding, ruleCount int, opts Options) error
```
Rationale: an `io.Writer`-based function (not a `string`-returning one)
lets callers stream to `os.Stdout` without an intermediate allocation, is
trivially testable with `bytes.Buffer` and golden files, and keeps
`internal/reporter` free of `fmt.Print*` per the library-first constraint —
the CLI layer owns deciding *where* output goes. Alternative considered:
`RenderText(findings []lint.Finding, ruleCount int, opts Options) string` —
rejected, it forces every caller (including tests) through a full string
allocation and slightly violates the "no `fmt.Print` in `internal/`" spirit
less directly than the writer form, but mainly loses the streaming option
for no benefit at this data size.

**`ruleCount` is caller-supplied, not inferred.** The reporter cannot
compute "N rules passed" from an empty `[]lint.Finding` alone — zero
findings and zero rules run look identical from inside the slice. Making
`ruleCount` an explicit parameter keeps the reporter a pure function of its
inputs instead of reaching for global rule-registry state (which would also
create an unwanted `internal/reporter` → `internal/lint/rules` dependency).
`scan-command`'s design will pass `len(rules.DefaultRules()) + 1` (the four
schema rules plus S001) as this argument — noted here so that later design
doesn't have to re-derive it. Alternative considered: bundle findings and
rule count into a `Result` struct (`type Result struct { Findings
[]lint.Finding; RuleCount int }`) — rejected as an unnecessary abstraction
for two values that are only ever used together in this one function call.

**Per-line format, literally.** FR-OUT-01's format string — `severity
rule-id  file:line — message` — is implemented as written: two literal
spaces between `severity`, `rule-id`, and `file:line`, then ` — ` (space,
em dash, space) before `message`. `severity` renders as the `Severity`
value verbatim (`"error"` / `"warning"`, already lowercase in
`internal/lint/finding.go`). No column alignment/padding beyond the fixed
two-space separator — alignment isn't specified by FR-OUT-01 and would add
complexity (computing max field width) for a cosmetic effect the
requirement doesn't ask for.

**Summary line uses the literal `(s)` suffix.** FR-OUT-02 specifies the
summary as `N error(s), M warning(s)` — that string is implemented exactly,
i.e. the output literally contains the substring `error(s)` and
`warning(s)` regardless of count (`"1 error(s), 0 warning(s)"`, not
grammatically pluralized `"1 error, 0 warnings"`). Rationale: this matches
the requirement's literal text with no branching logic, and mirrors a
common linter convention (npm, several Go linters) of an unconditional
`(s)` suffix. Alternative considered: grammatical pluralization (`1 error`
vs `2 errors`) — reads more naturally but isn't what FR-OUT-02 specifies
and adds a branch per count; can be revisited if the requirement is amended.

**Color: opt-in via `Options.Color`, forced off by `NO_COLOR`.** The
reporter checks `os.Getenv("NO_COLOR")` itself (any non-empty value
disables color, matching the [NO_COLOR](https://no-color.org) convention)
and ANDs that against `opts.Color`; there is no TTY auto-detection in this
change; no dependency beyond `os` and raw ANSI escape codes (`\x1b[31m` for
`error`, `\x1b[33m` for `warning`, `\x1b[0m` reset) — no new module
dependency, satisfying `TC-STACK-06`/the "no dependencies without approval"
boundary. Deciding *when* `Options.Color` is `true` (e.g. `isatty(stdout)`)
is deferred to whichever future change wires a `--color`/`--no-color` flag
or TTY check — likely `scan-command` or `configuration-support` — since
there's no CLI command in this change to attach that decision to yet.
Rationale for checking `NO_COLOR` inside the reporter rather than pushing
it entirely to the caller: BC-OUTPUT-01 explicitly calls out `NO_COLOR`
respect as a requirement of the output itself, so enforcing it at the
one place that emits ANSI codes is the same "checked in exactly one place"
principle applied in `schema-validation-rules`'s parser design.

**Success line short-circuits the whole format.** When `len(findings) ==
0`, `WriteText` writes only the FR-OUT-04 success line (`✓ AGENTS.md is
valid (N rules passed)`, using `ruleCount`) and returns — it does not also
print an empty per-line section or a `0 error(s), 0 warning(s)` summary.
Rationale: FR-OUT-04 describes "a single success line" as the entire output
for that case, and printing a redundant zero-summary under it would leave
two success signals in different words.

## Risks / Trade-offs

- **[Risk]** The literal `(s)` suffix in the summary line reads oddly for
  singular counts (`"1 error(s)"`). → **Mitigation**: this is exactly what
  FR-OUT-02 specifies; if it's later found to be a documentation shorthand
  rather than literal intent, fixing it is a one-line change isolated to
  the summary-formatting code path, and existing golden files make the
  before/after diff obvious.
- **[Risk]** `ruleCount` is an unchecked caller input — a future caller
  that passes the wrong count would silently produce a misleading success
  message (e.g. "(12 rules passed)" when only 5 ran). → **Mitigation**: the
  only near-term caller (`scan-command`) computes it from
  `len(rules.DefaultRules()) + 1`, a single source of truth documented
  above; `WriteText` doesn't attempt to validate it since it has no way to
  know the "true" count without importing `internal/lint/rules` (which
  would create a dependency this package doesn't otherwise need).
- **[Trade-off]** No TTY detection means `Options.Color` is inert until a
  future CLI flag sets it — acceptable since no command in this change
  actually calls `WriteText` yet; revisit when `scan-command` adds output
  flags.

## Open Questions

(none — deferred decisions, like when `Options.Color` defaults to `true`,
are explicitly punted to the future change that first has a CLI surface to
attach them to, per the "Color" decision above)
