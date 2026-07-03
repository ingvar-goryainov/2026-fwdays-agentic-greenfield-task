## Why

`agents-lint` has no way to report its own version: users filing bugs,
scripting CI, or checking what's installed have no `--version` flag to
query. Per `docs/openspec-capability-plan.md`, `cli-version-command` is
trivial and independent of every other capability (only `project-foundation`
is a prerequisite), so it can ship now without blocking or depending on
higher-value work.

## What Changes

- Implement FR-CLI-06: add a `--version` flag on the root command. Running
  `agents-lint --version` prints the tool's version and exits 0 without
  running any subcommand logic.
- Introduce a single version string, defaulting to `"dev"` for local/`go run`
  builds, that a later packaging step (`packaging-and-docs`) can override at
  build time via `-ldflags` without any change to this capability's contract.
- Use Cobra's built-in `Version` field on `rootCmd` so `--version` is
  handled by the existing command framework rather than a hand-rolled
  flag/branch. Cobra registers this flag only on `rootCmd` itself (it is not
  a persistent flag), matching FR-CLI-06's scope: `agents-lint --version`,
  not a flag on every subcommand.

## Capabilities

### New Capabilities
- `cli-version-command`: the `--version` flag (FR-CLI-06) on the root
  command and the version string it reports.

### Modified Capabilities
(none — this adds a new root-level flag; no existing capability's
requirements change)

## Impact

- Affected code: `cmd/agents-lint/cmd/root.go` (set `rootCmd.Version`);
  possibly a new small `version.go` holding the version constant/var.
- Dependencies introduced: none — uses Cobra's existing built-in version
  support (already a dependency).
- User-facing behavior: `agents-lint --version` prints a version string and
  exits 0. No existing command's behavior changes.
