## Why

`agents-lint` has no buildable Go project yet — no module, no CLI entrypoint,
no test scaffolding. Every other capability in the roadmap (schema
validation, scan/init commands, config, SARIF output) needs a working
`go build` and a Cobra command tree to attach to. This change establishes
that foundation so subsequent changes can focus purely on behavior.

## What Changes

- Initialize the Go module (`go.mod`) targeting Go 1.22+.
- Add `cmd/agents-lint/main.go` wired to a Cobra root command that currently
  does nothing but print help and exit cleanly (no subcommands yet — those
  arrive in later changes).
- Establish the `internal/` package layout that later changes will fill in
  (validation library, typed `Finding` results referenced by the product
  brief's architecture note).
- Add `testdata/` convention (empty, rule-scoped subfolders reserved for
  fixtures) and confirm `go test ./...` runs clean with zero tests.
- Add `LICENSE` (MIT).
- Add a `Makefile` or equivalent scripts for `build`/`test`/`lint` so build
  time (< 10s) is easy to verify going forward.

## Capabilities

### New Capabilities
- `project-foundation`: Go module, Cobra root command skeleton, package
  layout, and build/test tooling that every later capability depends on.

### Modified Capabilities
(none — first change in the project)

## Impact

- Affected code: entire repo (new `go.mod`, `cmd/`, `internal/`, `LICENSE`,
  build scripts).
- Dependencies introduced: `github.com/spf13/cobra` (TC-STACK-02). No other
  runtime dependencies yet — `goldmark` and `yaml.v3` are added in
  `schema-validation-rules`, not here.
- No user-facing behavior yet: this change produces a binary that builds and
  runs but has no rules, no `scan`/`init` commands, and no config support.
