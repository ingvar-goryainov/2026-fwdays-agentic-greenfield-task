## Context

Every behavioral capability in `docs/openspec-capability-plan.md` (1-10) is
implemented and archived. What remains before the MVP is "done" per
`docs/product-brief.md` is purely deliverable-facing: a way to hand someone a
binary (NFR-DIST-01, NFR-PERF-02) and a way for them to learn the tool without
reading source (NFR-DX-02). Both are additive — no `internal/` or `cmd/`
code changes are needed.

The repo's root `README.md` is the fwdays course-assignment template
(submission instructions, in Ukrainian) and is explicitly out of scope per
user decision: it stays untouched. The project's own documentation goes in a
new `docs/README.md`.

## Goals / Non-Goals

**Goals:**
- A single `make` target that cross-compiles static, stripped binaries for
  `linux/amd64` and `darwin/arm64`, each ≤ 15 MB.
- A `docs/README.md` covering installation, usage, the full rule catalog, and
  the `.agents-lint.yaml` reference, accurate as of the current CLI surface.
- Rule catalog content sourced from the existing `rules.RuleDocs()` /
  `RuleDoc` metadata (already used by `agents-lint docs <rule-id>`) so the
  catalog can't silently drift from what the binary actually does.

**Non-Goals:**
- No new CI workflow (`.github/workflows`) — this capability ships the local
  build target; wiring it into CI/release automation is a separate concern
  the user hasn't asked for.
- No additional target platforms beyond the two minimums in NFR-DIST-01
  (windows/amd64, linux/arm64, etc. are not required).
- No GoReleaser or other packaging dependency — `go build` with `GOOS`/
  `GOARCH` env vars is sufficient and keeps `TC-STACK-05` (offline) and the
  "no new deps without asking" boundary intact.
- No changes to rule behavior, CLI flags, exit codes, or output formats.

## Decisions

**Cross-compilation via Makefile, not GoReleaser.**
A `make release` target loops `GOOS`/`GOARCH` pairs and runs
`CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o dist/agents-lint-<os>-<arch> ./cmd/agents-lint`
for each. `-ldflags="-s -w"` strips debug/symbol info to help satisfy the
15 MB ceiling; `-trimpath` keeps builds reproducible. Alternative considered:
GoReleaser — rejected, it's a new external dependency requiring approval and
network access during release, which conflicts with `TC-STACK-05`.

**Rule catalog is hand-written prose in `docs/README.md`, cross-checked
against `RuleDocs()`, not code-generated.**
Generating the catalog from `RuleDocs()` at build time would need a new
`go generate` step and a template — more machinery than a static list of 7
rules justifies. Instead, the catalog table is written by hand and a task
step in `tasks.md` explicitly cross-checks it against
`rules.RuleDocs()`/`rules.KnownRuleIDs()` output so it can't miss a rule or
misstate a severity. Alternative considered: a `docs` subcommand flag that
dumps Markdown — deferred as unrequested scope creep (not in
`docs/requirements.md`).

**Binary size verified by measurement, not estimation.**
After `make release`, `ls -lh dist/` is checked against the 15 MB ceiling as
part of the task loop, rather than assumed from `-ldflags="-s -w"` alone,
since actual size depends on the Go version and linked stdlib.

## Risks / Trade-offs

- **Risk:** `-ldflags="-s -w"` strips symbols needed for debugging crash
  reports from a released binary. → **Mitigation:** acceptable for a CLI
  linter with no telemetry; users can rebuild from source with symbols if
  needed, and this is standard practice for size-constrained Go release
  binaries.
- **Risk:** Two READMEs (root + `docs/README.md`) could confuse a reader who
  lands on the repo root expecting project docs. → **Mitigation:** explicitly
  the user's chosen trade-off for this repo's dual purpose (course submission
  + real project); not this capability's problem to solve further.
- **Risk:** Rule catalog in `docs/README.md` drifts from `RuleDocs()` after a
  future rule change, since it's hand-written. → **Mitigation:** low risk —
  the rule set is frozen (MVP scope is closed per `product-brief.md`), and the
  cross-check task step catches drift at write time.
