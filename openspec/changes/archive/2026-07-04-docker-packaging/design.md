## Context

`agents-lint` is a pure-Go, CGO-free CLI (`TC-STACK-06`) already producing
static binaries via the `release` Makefile target
(`CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go build -trimpath -ldflags="-s -w"`).
That target already satisfies the hard part of containerizing it — a static
binary with no runtime dependencies is exactly what a minimal container base
needs. This design only decides the base image, the Dockerfile shape, and how
the container's mount/entrypoint contract maps onto the existing `scan`
command semantics (`FR-CLI-01`, `FR-CLI-02`).

## Goals / Non-Goals

**Goals:**
- Build a container image containing only the `agents-lint` binary and
  nothing else needed to run it.
- Preserve exact CLI behavior: same flags (`--config`, `--format`), same exit
  codes (0/1 per `FR-CLI-02`), same output.
- Keep the image buildable offline aside from the one-time base-image pull
  (consistent with `TC-STACK-05`'s "no network calls at runtime" — the
  constraint is on the tool's own execution, not the Docker build step).
- Fit naturally next to the existing `release` target rather than replacing
  it.

**Non-Goals:**
- No registry publishing / CI workflow for pushing images (out of scope per
  the proposal; can be a follow-up).
- No `docker-compose`, no daemon/server mode, no HTTP API — `agents-lint`
  remains a one-shot CLI invoked per `docker run`.
- No multi-arch manifest lists (`buildx`) — a single-arch build matching the
  host, consistent with how `release` targets specific `GOOS`/`GOARCH` pairs
  rather than a universal image.

## Decisions

**Multi-stage build, builder = `golang:1.22`, runtime = `gcr.io/distroless/static-debian12`.**
Distroless-static over `scratch`: it includes CA certificates and `/etc/passwd`
entries for a non-root user out of the box, which `scratch` lacks. Since
`agents-lint` makes no network calls (`TC-STACK-05`), CA certs aren't
functionally required, but distroless-static's non-root user support avoids
having to hand-roll a `USER` entry with numeric UID in `scratch` — a small
correctness win for negligible size cost (~2MB vs 0MB). Either meets
`NFR-PERF-02`'s spirit; distroless-static was chosen for the built-in non-root
user.

**Builder stage builds for the image's own `TARGETOS`/`TARGETARCH`** (Docker
buildkit's automatic build args), not a fixed platform. This means `docker
build` produces a binary matching whatever platform Docker is building for,
independent of the `release` target's fixed linux/amd64 + darwin/arm64 pair
(container images are Linux-only regardless of host darwin/arm64 support via
Docker Desktop's VM).

**`ENTRYPOINT ["/agents-lint"]`, no default `CMD`.** Users pass the subcommand
(`scan`, `init`, `docs`, `--version`) as `docker run` arguments, mirroring how
the native binary is invoked. No default `CMD scan` — an explicit subcommand
keeps behavior unsurprising and matches how the binary works standalone
(running `agents-lint` with no subcommand shows help, not an implicit scan).

**Working directory contract: `WORKDIR /workspace`.** Users mount their
project at `/workspace` (`-v $(pwd):/workspace`) and the container's default
CWD is `/workspace`, so `docker run --rm -v $(pwd):/workspace agents-lint
scan` resolves `./AGENTS.md` the same way the native binary would from the
host's CWD.

**Makefile targets `docker-build` and `docker-run`**, mirroring the existing
`build`/`test`/`release` naming pattern rather than introducing a separate
tool (e.g. `docker-compose`, `just`).

## Risks / Trade-offs

- **[Risk] Distroless-static has no shell**, so a user cannot `docker exec` in
  to debug. → **Mitigation**: acceptable trade-off for a one-shot CLI tool;
  documented in README as a known limitation; `docker run` output/exit code
  is the intended debugging surface, matching `BC-OUTPUT-02`'s "actionable
  messages" for troubleshooting.
- **[Risk] Version string defaults to `"dev"` in the image** unless the
  builder stage passes `-ldflags "-X .../cmd.version=..."` like `release`
  does. → **Mitigation**: Dockerfile accepts a `VERSION` build-arg (default
  `dev`) threaded into the same `-ldflags` the `release` target uses, so
  `--version` inside the container is meaningful when built with `--build-arg
  VERSION=vX.Y.Z`.
- **[Risk] Path confusion**: findings printed with `/workspace`-relative
  paths could look unfamiliar next to host paths in CI logs. → **Mitigation**:
  this matches how any containerized linter reports paths (relative to the
  mount point); documented in README with a worked example.

## Migration Plan

Additive only — no existing artifacts change behavior. `make docker-build`
and `docs/README.md`'s new section are net-new; the native binary and
`release` target are untouched.
