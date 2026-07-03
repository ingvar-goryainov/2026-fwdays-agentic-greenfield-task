## Why

`agents-lint` currently has a working `go build`/`go test` scaffold
(`project-foundation`) but no actual validation logic — the tool cannot yet
open an AGENTS.md file and say whether it's valid. The schema rules (FR-S001
through FR-S005) are the product's layer-1 core: everything else in the
roadmap (`text-reporting`, `scan-command`, `init-command`,
`codebase-awareness-rules`) depends on a parser and rule engine that emits
`[]lint.Finding`. This change builds that core so later changes can focus on
wiring and output formatting rather than validation semantics.

## What Changes

- Add `goldmark` (TC-STACK-03) as the Markdown parser and `gopkg.in/yaml.v3`
  (TC-STACK-04) as the frontmatter parser — the two runtime dependencies this
  capability introduces.
- Implement a small rule engine in `internal/lint` that runs a set of
  registered rules against a parsed AGENTS.md file and returns `[]Finding`.
- Implement five schema rules, one file per rule under `internal/lint/rules/`:
  - `FR-S001` — AGENTS.md file must exist at the expected path.
  - `FR-S002` — required top-level H2 sections must be present (at minimum
    `## Agent` or `## Agents`).
  - `FR-S003` — each agent definition block must contain `name` (H3 heading),
    `role` (description text), and at least one of `instructions`, `tools`,
    or `context`.
  - `FR-S004` — no duplicate agent names within the same file
    (case-insensitive comparison).
  - `FR-S005` — YAML frontmatter (if present, delimited by `---`) must be
    valid YAML.
- Add fixtures under `testdata/<rule-id>/valid.md` and
  `testdata/<rule-id>/invalid.md` for each rule (NFR-TEST-01), following the
  convention documented in `testdata/README.md`.
- No CLI wiring yet — `scan`/`init` commands are out of scope; this change
  only produces the library functions later commands will call.

## Capabilities

### New Capabilities
- `schema-validation-rules`: goldmark-based AGENTS.md parsing plus the five
  layer-1 schema rules (FR-S001–FR-S005), exposed as a rule engine that
  returns `[]lint.Finding`.

### Modified Capabilities
(none — `internal/lint.Finding` and `Severity` already exist from
`project-foundation` and are reused as-is, not changed)

## Impact

- Affected code: `internal/lint/` (new parser and `Document` model only —
  stays free of rule logic to avoid an import cycle with `internal/lint/rules`,
  see design.md), `internal/lint/rules/` (five new rule files plus the rule
  engine/orchestrator, replacing the current doc-comment-only placeholder),
  `testdata/` (ten new fixture files).
- Dependencies introduced: `github.com/yuin/goldmark`, `gopkg.in/yaml.v3`
  (both already declared as constraints in `docs/requirements.md`
  TC-STACK-03/04, not yet in `go.mod`).
- No user-facing behavior yet: there is still no `scan` or `init` command to
  invoke this engine from the CLI — that's `scan-command` and `init-command`.
