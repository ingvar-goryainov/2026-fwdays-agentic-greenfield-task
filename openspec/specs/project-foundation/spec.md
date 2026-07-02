## Purpose

Foundational Go project scaffolding for agents-lint: a buildable module,
Cobra CLI entrypoint, the `internal/` package layout, the shared `Finding`
result type, test/fixture scaffolding, and licensing. Every other capability
builds on this.

## Requirements

### Requirement: Buildable Go module
The system SHALL provide a Go module (Go 1.22+, modules enabled) that builds
successfully with `go build ./cmd/agents-lint` in under 10 seconds on a
clean checkout, with no CGO dependencies.

#### Scenario: Clean checkout build
- **WHEN** a developer runs `go build ./cmd/agents-lint` on a freshly cloned
  repository
- **THEN** the build completes successfully in under 10 seconds and produces
  a single executable binary

#### Scenario: No CGO required
- **WHEN** the binary is built with `CGO_ENABLED=0`
- **THEN** the build still succeeds, confirming no CGO dependency was
  introduced

### Requirement: Cobra CLI entrypoint
The system SHALL expose a runnable CLI entrypoint built on `cobra`
(`github.com/spf13/cobra`) that executes without error and prints help text
when invoked with no arguments or `--help`.

#### Scenario: Running with no arguments
- **WHEN** a user runs the built binary with no arguments
- **THEN** the process exits 0 and prints usage/help text without crashing

#### Scenario: Running --help
- **WHEN** a user runs the built binary with `--help`
- **THEN** the process exits 0 and prints the command's help output

### Requirement: Internal package layout for validation library
The system SHALL structure core logic under `internal/` as an importable
library, separate from the `cmd/` entrypoint, so future consumers (HTTP
API, editor extensions) can depend on it without depending on the CLI.

#### Scenario: Library package is independently importable
- **WHEN** another package in the module imports the `internal/lint`
  package
- **THEN** it compiles without importing anything from `cmd/`

### Requirement: Shared Finding type
The system SHALL define a `Finding` type in the `internal/lint` package
representing a single validation result, with at minimum: rule ID,
severity, file path, line number, and message fields.

#### Scenario: Finding type is constructible
- **WHEN** code constructs an `internal/lint.Finding` with rule ID,
  severity, file, line, and message
- **THEN** all five fields are set and readable, with no compile error

### Requirement: Test scaffolding runs clean
The system SHALL support `go test ./...` running successfully (exit 0) on a
module with zero rule implementations, and SHALL document the `testdata/`
fixture convention (organized by rule ID) that later rule changes will
populate.

#### Scenario: Empty test suite passes
- **WHEN** a developer runs `go test ./...` immediately after this change is
  applied
- **THEN** the command exits 0, whether or not any test files exist yet

### Requirement: MIT license present
The system SHALL include a `LICENSE` file at the repository root containing
the MIT license text.

#### Scenario: License file exists
- **WHEN** a developer inspects the repository root
- **THEN** a `LICENSE` file is present and its content is the standard MIT
  license text
