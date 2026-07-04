## ADDED Requirements

### Requirement: NFR-DIST-02 — Docker image build
The project SHALL provide a repo-root `Dockerfile` that builds a container image containing the `agents-lint` binary via a multi-stage build (a Go build stage producing a static, `CGO_ENABLED=0` binary, and a minimal runtime stage containing only that binary and its non-root runtime user).

#### Scenario: Image builds from a clean checkout
- **WHEN** a developer runs `make docker-build` (or the equivalent `docker build` command) from a clean checkout with Docker available
- **THEN** an image is produced containing the `agents-lint` binary and no Go toolchain or source code

#### Scenario: Runtime stage has no shell or package manager
- **WHEN** the image is inspected (e.g. `docker run --rm --entrypoint sh <image>` fails to find a shell)
- **THEN** the runtime stage contains only the binary, CA certs/user files provided by the minimal base image, and no shell, package manager, or build tooling

### Requirement: Container run behavior matches the native binary
Running the image with a mounted host directory SHALL invoke `agents-lint` against that directory and produce identical output and exit codes to running the native binary on the same input (`FR-CLI-01`, `FR-CLI-02`).

#### Scenario: Scan against a mounted valid AGENTS.md exits 0
- **WHEN** a user runs `docker run --rm -v <host-dir>:/workspace <image> scan` against a directory whose `AGENTS.md` passes all rules
- **THEN** the container exits 0 and prints the same success line the native binary would print (`FR-OUT-04`)

#### Scenario: Scan against a mounted invalid AGENTS.md exits 1
- **WHEN** a user runs `docker run --rm -v <host-dir>:/workspace <image> scan` against a directory whose `AGENTS.md` fails an error-severity rule
- **THEN** the container exits 1 and prints the same per-finding lines and summary the native binary would print (`FR-OUT-01`, `FR-OUT-02`)

#### Scenario: CLI flags pass through unchanged
- **WHEN** a user runs the container with `--format sarif` or `--config <path>` appended after the subcommand
- **THEN** the container honors the flag identically to the native binary (`FR-CLI-04`, `FR-CLI-05`)

### Requirement: Makefile Docker targets
The project SHALL provide a `docker-build` Makefile target that builds the image, alongside the existing `build`/`test`/`release` targets.

#### Scenario: docker-build target exists and succeeds
- **WHEN** a developer runs `make docker-build`
- **THEN** the target invokes `docker build` against the repo-root `Dockerfile` and completes successfully when Docker is available

### Requirement: README documents Docker usage
`docs/README.md` SHALL document how to build and run the Docker image, consistent with how it documents native binary installation and usage (`NFR-DX-02`).

#### Scenario: README shows a build and run example
- **WHEN** a new user reads `docs/README.md`'s installation/usage sections
- **THEN** it includes a concrete `docker build` (or `make docker-build`) command and a concrete `docker run` example with a volume mount that scans a project directory
