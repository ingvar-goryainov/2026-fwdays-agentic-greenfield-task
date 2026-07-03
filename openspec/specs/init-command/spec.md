## Purpose

The `agents-lint init` CLI command, generating a static AGENTS.md skeleton
(`internal/lint/skeleton`) guaranteed to pass rules S001-S005, so a new
user has a correct starting point instead of a blank page.

## Requirements

### Requirement: FR-CLI-03 — `init` generates a schema-passing AGENTS.md skeleton
`agents-lint init [path]` SHALL write a static AGENTS.md skeleton to `path`, defaulting `path` to `./AGENTS.md` when no positional argument is given. The generated skeleton MUST produce zero error-severity findings when validated against rules S001-S005.

#### Scenario: Default path, no existing file
- **WHEN** a user runs `agents-lint init` in a directory with no
  `./AGENTS.md`
- **THEN** the command creates `./AGENTS.md` containing a skeleton that
  passes S001-S005, exits with code 0, and prints a confirmation message
  naming the path written

#### Scenario: Explicit path argument
- **WHEN** a user runs `agents-lint init docs/AGENTS.md` and
  `docs/AGENTS.md` does not yet exist
- **THEN** the command writes the skeleton to `docs/AGENTS.md` instead of
  the default path

#### Scenario: Generated skeleton passes every schema rule
- **WHEN** the AGENTS.md skeleton produced by `init` is validated by the
  S001-S005 rule set
- **THEN** no error-severity findings are reported

### Requirement: `init` refuses to overwrite an existing AGENTS.md by default
`agents-lint init` SHALL NOT overwrite an existing file at the target path unless the `--force` flag is passed, to avoid silently destroying a hand-authored AGENTS.md.

#### Scenario: Existing file, no --force
- **WHEN** a user runs `agents-lint init` and a file already exists at the
  target path
- **THEN** the command makes no changes to the existing file, exits with a
  non-zero code, and prints an actionable error explaining that the file
  already exists and that `--force` will overwrite it

#### Scenario: Existing file, with --force
- **WHEN** a user runs `agents-lint init --force` and a file already exists
  at the target path
- **THEN** the command overwrites the existing file with the skeleton,
  exits with code 0, and prints a confirmation message naming the path
  written
