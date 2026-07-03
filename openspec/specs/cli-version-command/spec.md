## Purpose

`--version` flag for `agents-lint`, letting users print the tool's version
string and exit immediately from the root command, without running any
subcommand logic.

## Requirements

### Requirement: FR-CLI-06 — `--version` prints the version and exits
`agents-lint --version` SHALL print the tool's version string and exit with
code 0 without executing any subcommand's logic (no scan, no config
resolution, no file I/O beyond writing the version to stdout).

#### Scenario: `--version` at the root command
- **WHEN** a user runs `agents-lint --version`
- **THEN** the command prints a version string to stdout and exits with
  code 0, without running `scan` or `init` logic

#### Scenario: No build-time version injected
- **WHEN** the binary was built with `go build ./cmd/agents-lint` (no
  `-ldflags` version override)
- **THEN** `agents-lint --version` prints the fallback version `"dev"`

#### Scenario: `--version` takes precedence over other flags
- **WHEN** a user runs `agents-lint --version` together with other flags
  (e.g. `agents-lint --version --format sarif`)
- **THEN** the command still prints the version and exits 0, ignoring the
  other flags rather than attempting a scan

#### Scenario: `--version` is not accepted on subcommands
- **WHEN** a user runs `agents-lint scan --version`
- **THEN** the command reports an unrecognized flag error, since
  `--version` is registered on the root command only
