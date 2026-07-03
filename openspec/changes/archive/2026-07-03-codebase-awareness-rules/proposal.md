## Why

Today `agents-lint scan` only validates AGENTS.md's own structure (FR-S001-S005)
— it never checks whether the file's *claims* about the repository are still
true. An AGENTS.md that references a script that was renamed, or names a tool
the repo no longer uses, passes every existing rule while being silently
wrong. Per `docs/openspec-capability-plan.md`, `codebase-awareness-rules` is
the layer-2 capability that closes this gap: it reuses the same rule
engine/`Finding` contract validated by `schema-validation-rules` and
`scan-command`, but checks AGENTS.md against the surrounding repo instead of
just against itself.

## What Changes

- Add a new rule `C001` (FR-C001, error severity): every inline-code
  (backtick) span in AGENTS.md that looks like a filesystem path is resolved
  relative to the repo root; a path that does not exist produces a finding.
- Add a new rule `C002` (FR-C002, warning severity): every inline-code span
  that matches a known developer tool/command name (e.g. `terraform`,
  `npm`, `docker`) is checked against a fixed table of repo-root evidence
  files (lockfiles, manifests, config files); a referenced tool with no
  matching evidence produces a finding.
- Extend `internal/lint.Document` (produced by `internal/lint.Parse`) with a
  `CodeSpans []CodeSpan` field, populated by a full-document AST walk so
  codebase-awareness rules can see every inline-code span regardless of
  where it appears (heading, paragraph, list item, etc.) — not just inside
  agent blocks, which is all the current parser tracks.
- Introduce a `lint.CodebaseRule` interface (`ID() string`,
  `Check(doc *Document, repoRoot string) []Finding`), distinct from the
  existing `lint.Rule` interface, since C001/C002 need a repo root that
  schema rules don't.
- Extend `rules.Run` to also execute the codebase-awareness rules (using the
  process's working directory as `repoRoot`, matching the existing
  cwd-relative convention behind `defaultAgentsMDPath` and
  `config.DefaultFileName`) in the same pass as the schema rules, and extend
  `rules.KnownRuleIDs()` to include `C001`/`C002` so `configuration-support`
  can enable/disable/re-severity them with no changes to `internal/config`.

## Capabilities

### New Capabilities
- `codebase-awareness-rules`: the layer-2 rules that validate AGENTS.md's
  file-path and tool/command references against the actual repository
  (FR-C001, FR-C002).

### Modified Capabilities
(none — `schema-validation-rules`'s FR-S001-S005 contracts are unchanged;
the `Document.CodeSpans` addition and the new `CodebaseRule` interface are
purely additive plumbing that this capability's own spec covers. `scan-command`'s
FR-CLI-01/02 exit-code and default-path contracts are unchanged — this
capability only adds new finding sources that flow through the same
`rules.Run` → reporter pipeline.)

## Impact

- Affected code: `internal/lint/parser.go` (new `CodeSpans` field +
  document-wide walk), `internal/lint/rule.go` (new `CodebaseRule`
  interface), new `internal/lint/rules/c001.go` and
  `internal/lint/rules/c002.go`, `internal/lint/rules/registry.go`
  (`Run`/`KnownRuleIDs()` extended to include the codebase-awareness rules),
  new fixtures under `testdata/C001/` and `testdata/C002/`.
- Dependencies introduced: none — reuses `github.com/yuin/goldmark/ast`'s
  existing `Walk` function, already available via the module's existing
  goldmark dependency (TC-STACK-03).
- User-facing behavior: `agents-lint scan` now also flags AGENTS.md file
  references that don't resolve to real files (error, fails the scan) and
  tool/command references with no corresponding repo evidence (warning,
  does not fail the scan by itself). `.agents-lint.yaml` can enable,
  disable, or re-severity `C001`/`C002` exactly like the existing schema
  rules, with no config-schema changes required.
