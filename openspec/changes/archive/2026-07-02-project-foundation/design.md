## Context

`agents-lint` starts from an empty repository (`docs/` only). This is the
first OpenSpec change and has no prior code, module, or package layout to
build on. Every later capability in the roadmap (`schema-validation-rules`,
`scan-command`, `configuration-support`, `sarif-output`, ...) assumes a
working `go build`, a `Finding` type, and a place to attach new Cobra
subcommands — none of that exists yet.

Per the product brief's architecture note: core validation logic must live
in `internal/` as a library with typed results (`[]Finding`) so future
consumers (HTTP API, VS Code extension, dashboard) can reuse it without
refactoring. This change establishes that boundary from day one, even though
no rules exist yet.

## Goals / Non-Goals

**Goals:**
- A `go build ./cmd/agents-lint` that succeeds and produces a binary that
  runs and prints help (TC-STACK-01, TC-STACK-02, NFR-DX-01).
- A package layout (`cmd/`, `internal/`) that later changes can add to
  without restructuring.
- A `Finding` type and package location settled now, since it's the
  contract every rule and reporter will depend on.
- `go test ./...` runs (zero tests, exits 0) and `testdata/` conventions are
  documented so `schema-validation-rules` can drop fixtures straight in.
- MIT `LICENSE` file at repo root (BC-LICENSE-01).

**Non-Goals:**
- No validation rules, no `scan`/`init` subcommands, no config loading —
  those are separate changes (`schema-validation-rules`, `scan-command`,
  `configuration-support`).
- No CI workflow files — out of scope for this change; only local
  build/test tooling.
- No cross-compilation / release packaging — that's `packaging-and-docs`.

## Decisions

**Module path:** `github.com/<org>/agents-lint` (placeholder org — confirm
with user before running `go mod init` if a GitHub org isn't already
implied by the repo). Rationale: `go install` (product brief step 1) needs a
real importable path.

**Package layout:**
```
cmd/agents-lint/main.go       -- entrypoint, builds and runs the Cobra root command
internal/lint/                -- reserved for validation library (Finding type, rule engine) — populated by schema-validation-rules
internal/lint/rules/          -- reserved for individual rule implementations
testdata/                     -- reserved, rule-scoped fixture directories added as rules land
```
Rationale: `internal/` prevents external consumers from depending on
unstable APIs pre-1.0 while still allowing every in-repo package (CLI,
future HTTP API) to share the library, matching the product brief directly.

**`Finding` type home:** define `internal/lint.Finding` now (even though no
rule produces one yet) — fields: `RuleID string`, `Severity Severity`,
`File string`, `Line int`, `Message string`. Rationale: `text-reporting` and
`schema-validation-rules` both depend on this exact type; agreeing on it in
the foundation change avoids an early breaking rename.

**CLI framework wiring:** Cobra root command only, no subcommands yet.
`cmd/agents-lint/main.go` calls a `cmd.Execute()` in a small `cmd` package
under `cmd/agents-lint/` (Cobra convention) so `scan`/`init`/etc. can
register themselves as subcommands later without touching `main.go`.

**Build/test tooling:** a `Makefile` with `build`, `test`, `lint` targets
rather than a bespoke script — smallest thing that lets NFR-DX-01 (< 10s
build) be checked with one command (`make build`) both locally and later in
CI.

**Dependency scope:** only add `github.com/spf13/cobra` in this change.
`goldmark` and `yaml.v3` are declared by TC-STACK-03/04 but stay out of
`go.mod` until `schema-validation-rules` actually imports them — keeps this
change's `go.mod` honest about what's actually used.

## Risks / Trade-offs

- **[Risk]** Guessing the module path wrong (e.g., placeholder org) means a
  find/replace across every file in later changes. → **Mitigation**:
  confirm the module path with the user before `go mod init` during
  implementation; treat it as a blocking question, not an assumption.
- **[Risk]** Defining `Finding` before any rule exists risks guessing wrong
  fields. → **Mitigation**: keep the struct minimal (the 5 fields implied
  by FR-OUT-01's `severity rule-id file:line — message` format) and treat
  it as extensible, not final — `schema-validation-rules` can add fields
  without breaking this change's contract.
- **[Trade-off]** No CI workflow in this change means NFR-DX-01 and
  `go build`/`go test` are only verified locally until a later change adds
  CI. Acceptable: CI setup isn't in the requirements doc as its own
  capability, and manual verification is enough to unblock later work.

## Open Questions

- Confirm the Go module path (GitHub org/repo) before running `go mod init`.
