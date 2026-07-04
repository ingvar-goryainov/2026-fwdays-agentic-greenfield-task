## ADDED Requirements

### Requirement: NFR-DIST-01 — Cross-compiled release binaries
The project SHALL provide a build target that produces statically-linked, `CGO_ENABLED=0` binaries for at least `linux/amd64` and `darwin/arm64` from a single command, with no external packaging tool beyond the Go toolchain.

#### Scenario: Release target builds both minimum platforms
- **WHEN** a developer runs the release build target from a clean checkout
- **THEN** `dist/` contains a binary for `linux/amd64` and a binary for `darwin/arm64`, each built with `CGO_ENABLED=0`

#### Scenario: Release build requires no network access
- **WHEN** the release build target is run with module dependencies already present in the local module cache
- **THEN** the build completes without any network call, per `TC-STACK-05`

### Requirement: NFR-PERF-02 — Binary size ceiling
Each cross-compiled release binary SHALL be no larger than 15 MB (stripped, single platform).

#### Scenario: Each platform binary is under the size ceiling
- **WHEN** the release build target completes
- **THEN** every binary in `dist/` is ≤ 15 MB as reported by the filesystem

### Requirement: NFR-DX-02 — Project README with install, usage, rule catalog, and config reference
The project SHALL provide a README (`docs/README.md`) documenting: how to install the binary (download or `go install`), how to use each CLI command (`scan`, `init`, `docs`, `--version`, `--format`), a catalog of every rule the tool ships (ID, description, default severity), and the `.agents-lint.yaml` configuration reference (supported keys, defaults, and an example file).

#### Scenario: README documents installation
- **WHEN** a new user reads `docs/README.md`
- **THEN** it describes at least one installation method (binary download or `go install`) with a concrete command

#### Scenario: README documents every CLI command
- **WHEN** a new user reads `docs/README.md`
- **THEN** it includes a usage example for `scan`, `init`, `docs`, `--version`, and the `--format` flag

#### Scenario: README's rule catalog covers every shipped rule
- **WHEN** a new user reads the rule catalog section of `docs/README.md`
- **THEN** it lists all rule IDs returned by `rules.KnownRuleIDs()` (`S001`-`S005`, `C001`, `C002`), each with a description and default severity matching `rules.RuleDocs()`

#### Scenario: README documents configuration
- **WHEN** a new user reads `docs/README.md`
- **THEN** it explains the `.agents-lint.yaml` default location, the `--config` flag, the `path` and `rules` keys, and shows a complete example config file

#### Scenario: Root README is untouched
- **WHEN** this capability is applied
- **THEN** the repo-root `README.md` (fwdays course assignment template) is unchanged
