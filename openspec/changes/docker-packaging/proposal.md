## Why

`agents-lint` ships as cross-compiled binaries (`packaging-and-docs`), but users
who want a zero-install, reproducible way to run it — e.g. in CI pipelines or
on machines without a Go toolchain — have no container image to pull. Docker
is the de facto standard for that use case, and packaging one is a small,
self-contained addition that doesn't touch validation logic. This closes the
gap now documented as `NFR-DIST-02` in `docs/requirements.md`.

## What Changes

- Add a repo-root `Dockerfile`: a multi-stage build where stage 1 compiles a
  static Go binary (`CGO_ENABLED=0`, same flags as the `release` Makefile
  target) and stage 2 copies only that binary into a minimal runtime base
  (`scratch` or `gcr.io/distroless/static`), running as a non-root user, with
  `ENTRYPOINT ["agents-lint"]`.
- Running the container mounts a host directory and invokes `scan` against
  it (e.g. `docker run --rm -v $(pwd):/workspace -w /workspace agents-lint
  scan`), producing identical output and exit codes (`FR-CLI-02`) to the
  native binary, including `--format` and `--config`.
- Add `docker-build` (and optionally `docker-run`) targets to `Makefile`
  alongside the existing `release` target.
- Document image build/run instructions in `docs/README.md`'s installation
  and usage sections, alongside the existing binary-download and `go
  install` instructions.
- No changes to CLI behavior, rule logic, exit codes, or flags — packaging
  only, same boundary `packaging-and-docs` followed.

## Capabilities

### New Capabilities
- `docker-packaging`: a `Dockerfile` producing a minimal container image
  that runs `agents-lint` against a mounted directory, plus the
  Makefile targets and README documentation to build and use it.

### Modified Capabilities
- (none — no existing spec's requirements change)

## Impact

- **Affected code:** repo-root `Dockerfile` (new), `Makefile` (new
  `docker-build`/`docker-run` targets), `docs/README.md` (new section); no
  changes to `internal/` or `cmd/`.
- **Dependencies:** none added to the Go module. Uses only the Go toolchain
  and a public minimal base image (`scratch` or `distroless/static`) — no
  new build tool.
- **CI/systems:** none — this only adds a local Docker build target and
  docs; no registry push or `.github/workflows` changes are in scope unless
  requested separately.
