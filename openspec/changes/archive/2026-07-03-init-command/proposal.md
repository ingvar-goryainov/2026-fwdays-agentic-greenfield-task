## Why

Right now `agents-lint scan` is a pure validator: it tells you an AGENTS.md is
missing or broken but gives you nothing to start from. New users hit a wall of
S001-S005 findings with no reference example. `agents-lint init` (FR-CLI-03)
closes that loop by generating a minimal AGENTS.md skeleton that is
guaranteed to pass every schema rule, giving users a correct starting point
instead of a blank page.

## What Changes

- Add a new `agents-lint init [path]` subcommand, wired into the existing
  Cobra root command alongside `scan`.
- Generate a minimal, static AGENTS.md skeleton containing:
  - valid YAML frontmatter (satisfies S005)
  - a top-level `## Agent` section (satisfies S002)
  - one `###` agent block with a name, a role description, and an
    Instructions field (satisfies S003, and trivially S004 since there is
    only one agent name)
- Default output path is `./AGENTS.md` (same default as `scan`); an optional
  positional `path` argument overrides it.
- Refuse to overwrite an existing file at the target path by default, exiting
  non-zero with a clear message; add a `--force` flag to allow overwriting.
- After writing the file, print a short confirmation message with the path
  written.
- Add a regression test that runs the real S001-S005 rule set against the
  generated skeleton to guarantee it always passes (this is the core
  guarantee of FR-CLI-03 and must never silently drift).

## Capabilities

### New Capabilities
- `init-command`: `agents-lint init` command that generates a schema-passing
  AGENTS.md skeleton, including the skeleton content, file-write/overwrite
  behavior, and CLI wiring.

### Modified Capabilities
- (none — this does not change scan-command, schema-validation-rules, or
  text-reporting behavior)

## Impact

- New file: `cmd/agents-lint/cmd/init.go` (+ `init_test.go`), following the
  same `runInit`/Cobra-`RunE` split as `scan.go`.
- New file: `internal/lint/skeleton/skeleton.go` (or similar) holding the
  embedded/generated AGENTS.md template, so the content is unit-testable
  independent of file I/O and CLI plumbing.
- Depends on the finished `schema-validation-rules` capability (already
  `apply`-complete) — the skeleton is validated against the real S001-S005
  rules, not a hand-maintained assumption of what they check.
- No changes to `go.mod` dependencies, `scan` behavior, or output formats.
