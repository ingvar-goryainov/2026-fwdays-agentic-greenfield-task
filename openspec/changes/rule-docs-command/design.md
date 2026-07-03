## Context

`agents-lint` has seven rules (S001–S005, C001, C002). Today each rule is a
`lint.Rule` or `lint.CodebaseRule` implementation with only an `ID()` and a
`Check(...)` method — the human-readable explanation of *what a rule checks*
and *why* lives only as inline comments in Go source and as one-line
`Finding.Message` strings produced at failure time. There is no structured,
queryable metadata a CLI command could print. `internal/lint/rules/registry.go`
already centralizes rule ID knowledge (`KnownRuleIDs()`), so it's the natural
place to add a second lookup: rule ID → documentation.

FR-CLI-07 requires `agents-lint docs <rule-id>` to print: description,
severity, a valid example, an invalid example, and fix guidance. This is a
CLI-only addition — no scan/init/config/reporter behavior changes.

## Goals / Non-Goals

**Goals:**
- Attach structured doc metadata to each of the seven existing rules.
- `agents-lint docs <rule-id>` prints that metadata for a known rule ID.
- Unknown rule IDs fail clearly and non-zero, per `BC-OUTPUT-02`.
- Keep `internal/` free of `fmt.Print` — `docs.go` in `cmd/` renders,
  `internal/lint/rules` only returns typed data.

**Non-Goals:**
- No `--fix` / auto-remediation (BC-SCOPE-02) — fix guidance is text only.
- No change to `Finding`, scan output, or exit codes.
- No new command listing *all* rules at once (`docs <rule-id>` is
  single-rule per FR-CLI-07's literal signature); a `docs` with no args
  can be a simple usage error, not a catalog dump.
- No SARIF/JSON output for `docs` — text only, matching its role as a
  human-facing reference (BC-OUTPUT-01).

## Decisions

**1. New `lint.RuleDoc` struct + `Doc() RuleDoc` method, not a Rule interface change.**
`Rule` and `CodebaseRule` stay untouched (`ID()`, `Check(...)`) so existing
call sites (`Run`, `KnownRuleIDs`) are unaffected. Each rule type instead
gets an additional `Doc() lint.RuleDoc` method:
```go
type RuleDoc struct {
    ID          string
    Description string
    Severity    Severity
    Valid       string // short example snippet that passes the rule
    Invalid     string // short example snippet that fails the rule
    FixGuidance string
}
```
Alternative considered: extend `Rule`/`CodebaseRule` with `Doc()` directly.
Rejected — it would force every future rule author to fill in doc fields
before `Check` even compiles, and ripples through `DefaultRules()` /
`CodebaseAwareRules()` return types unnecessarily. A separate method on the
concrete struct is additive and doesn't touch the interfaces `Run` depends on.

**2. Examples are short inline string literals in each rule file, not `testdata/` fixtures.**
`testdata/<rule>/valid.md` and `invalid.md` are full, runnable AGENTS.md
files sized for table-driven tests (some are multi-section). `docs`
examples should be the smallest snippet that demonstrates the rule, e.g.
for S002: `"## Agent\n"` (valid) vs a file with no `## Agent`/`## Agents`
section (invalid). Reusing `testdata/` directly (via `go:embed`) was
considered and rejected: it would mix test-fixture concerns with
product-facing docs, couple doc output to fixture file layout, and (via
embed) put test-only content in the shipped binary. Inline literals keep
the two concerns independent and the example short enough to read at a
glance.

**3. Doc lookup lives in `registry.go` as `RuleDocs() map[string]lint.RuleDoc`.**
Mirrors the existing `KnownRuleIDs()` pattern (iterate `DefaultRules()` +
`CodebaseAwareRules()`, plus S001 handled separately since it runs outside
the `Rule` interface, same as today). `cmd/docs.go` looks up by ID in this
map; a miss is the "unknown rule ID" error path.

**4. `cmd/agents-lint/cmd/docs.go` follows the existing `scan.go` shape.**
A pure, testable `runDocs(w io.Writer, ruleID string) error` function called
from a thin Cobra `RunE`, same separation `scan.go` already uses for
`runScan`. `docs` takes exactly one positional arg (`cobra.ExactArgs(1)`);
FR-CLI-07 is single-rule lookup, not a catalog.

## Risks / Trade-offs

- **Doc/behavior drift**: nothing enforces that `Description`/fix guidance
  stays in sync with `Check` logic if a rule's behavior changes later. →
  Mitigation: doc fields live in the same file as `Check`, right next to
  it, so a rule-logic change is a one-file diff that includes its doc.
- **Duplication with `Finding.Message`**: `FixGuidance` will overlap in
  content with the `Message` string `Check` already returns. → Accepted:
  they serve different audiences (one-line inline vs. full reference); kept
  as plain string duplication rather than introducing a shared-template
  abstraction, per the project's "no premature abstraction" convention.

## Migration Plan

Purely additive; no existing behavior changes, no config/schema migration.
Roll out as one change touching all seven rule files, the registry, and one
new `cmd/docs.go`. Rollback is a revert of this single commit/branch.
