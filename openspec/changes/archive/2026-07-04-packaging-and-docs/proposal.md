## Why

`agents-lint` is functionally complete: all seven rules (S001-S005, C001-C002),
both CLI commands (`scan`, `init`), configuration support, SARIF output, and
the `docs`/`--version` commands are implemented and archived. The MVP
described in `docs/product-brief.md` is behavior-complete, but it isn't yet
*distributable* or *discoverable* — there are no cross-compiled release
binaries and no documentation a new user could follow to install, run, and
understand it. This capability closes the MVP by shipping the artifacts a
real user needs to adopt the tool.

## What Changes

- Add a `make release` (or equivalent) build step that cross-compiles static
  binaries for `linux/amd64` and `darwin/arm64` (minimum), with `CGO_ENABLED=0`
  and build flags to strip debug symbols so each binary stays ≤ 15 MB
  (NFR-DIST-01, NFR-PERF-02).
- Add a project-specific `docs/README.md` covering: installation (binary
  download and `go install`), usage (`scan`, `init`, `docs`, `--version`,
  `--format`), the full rule catalog (all 7 rules, sourced from the existing
  `rules.RuleDocs()` metadata to avoid drift), and the `.agents-lint.yaml`
  configuration reference (NFR-DX-02).
  - The repo-root `README.md` (the fwdays course assignment template) is left
    untouched — it serves a different audience (course grading) and is out of
    scope for this capability.
- No changes to CLI behavior, rule logic, or exit codes — this capability is
  packaging and documentation only.

## Capabilities

### New Capabilities
- `packaging-and-docs`: cross-compiled release binaries and the standalone
  project README (rule catalog + config reference), closing out the MVP.

### Modified Capabilities
- (none — no existing spec's requirements change; this only adds build
  tooling and documentation)

## Impact

- **Affected code:** `Makefile` (new release/cross-compile target); no
  changes to `internal/` or `cmd/`.
- **New files:** `docs/README.md` (or generated from `RuleDocs()` at doc-write
  time — decided in design.md).
- **Dependencies:** none added. Cross-compilation uses the Go toolchain's
  built-in `GOOS`/`GOARCH` support — no GoReleaser or other new tool, per the
  "no dependencies without asking" boundary.
- **CI/systems:** none — this capability only adds a local build target and
  docs; no `.github/workflows` changes are in scope unless the user asks
  separately.
