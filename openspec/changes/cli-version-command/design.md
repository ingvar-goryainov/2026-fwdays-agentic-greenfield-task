## Context

`cmd/agents-lint/cmd/root.go` defines `rootCmd` (a `*cobra.Command`) and two
persistent flags (`--config`, `--format`), each validated by the subcommand
that uses them (`scan`). There is no existing version string anywhere in
the codebase, and no build-time version injection (`Makefile` has only
`build`/`test`/`lint` targets; `go build ./cmd/agents-lint` with no
`-ldflags`). `packaging-and-docs` (Phase 7, after this capability) is where
cross-compiled release binaries and `-ldflags` version injection are
expected to land, per `docs/openspec-capability-plan.md`; this change must
not block on that work, only leave a seam for it.

Cobra (already a dependency) has built-in support for a version flag:
setting `Command.Version` to a non-empty string makes Cobra register a
`--version` flag on that command, which — when passed — prints
`<name> version <Version>` and returns before `RunE` executes.

## Goals / Non-Goals

**Goals:**
- Implement FR-CLI-06: `agents-lint --version` prints the version and exits
  (code 0), doing no scan/config/file work.
- Leave a single, obvious seam (`version` var) for `packaging-and-docs` to
  override via `-ldflags "-X ...=vX.Y.Z"` later, without this change
  guessing at release tooling it doesn't own yet.

**Non-Goals:**
- Actual release versioning/tagging scheme, cross-compilation, or
  `-ldflags` wiring in `Makefile`/CI — that's `packaging-and-docs` (Phase 7).
  This change only defines the variable those future ldflags will target.
- `agents-lint version` as a subcommand — FR-CLI-06 specifies a `--version`
  flag, not a subcommand; adding both would be two ways to do the same
  thing with no requirement asking for it.
- A `--verbose`/`-v` flag — out of scope; noted only because Cobra's
  default version-flag registration also claims `-v` as shorthand when free
  (see Decision 3).

## Decisions

**1. `rootCmd.Version` (Cobra's built-in mechanism), not a hand-rolled flag
and branch in `RunE`.**
`--format`/`--config` need custom validation logic (accepted-value checks,
file-path resolution) so they're hand-rolled persistent flags. `--version`
needs none of that — it's a static string with print-and-exit semantics,
which is exactly what Cobra's `Version` field already does. Reimplementing
it as a custom flag would duplicate framework behavior the project already
depends on.
  - *Alternative considered:* a `version` subcommand alongside `scan`/`init`
    — rejected, FR-CLI-06 explicitly specifies a flag (`--version`), not a
    subcommand.

**2. Version stored as `var version = "dev"` in a new
`cmd/agents-lint/cmd/version.go`, assigned to `rootCmd.Version` in `init()`.**
A package-level `var` (not `const`) is required for `-ldflags -X` injection
to work at all (Go's linker can only overwrite package-level string vars).
`"dev"` is the fallback for `go run`/`go build` with no ldflags, mirroring
common Go CLI convention (e.g. `goreleaser`-built tools default to `"dev"`
or `"(devel)"` absent injected values). A separate file (rather than adding
to `root.go`) keeps the one line `packaging-and-docs` will need to touch
(`-X github.com/.../cmd.version=...`) isolated from the flag-heavy
`root.go`.
  - *Alternative considered:* reading version from `debug.BuildInfo` (Go
    1.18+ embeds module version for `go install`-installed binaries) —
    rejected for this change: it only populates useful data for `go
    install`-fetched binaries with a tagged module version, not for the
    local `go build ./cmd/agents-lint` this repo's `Makefile` uses today,
    and adds a second version source to reconcile later. Can be revisited
    in `packaging-and-docs` if useful.

**3. Accept Cobra's default `-v` shorthand for `--version` rather than
suppressing it.**
Cobra registers `-v` as shorthand for `--version` automatically unless
another flag already claims it. No flag in this project currently uses
`-v`, so `agents-lint -v` will also print the version. This is Cobra's
default behavior for any `Command.Version`-bearing command and matches
common CLI convention; suppressing it would need extra code for no
requirement-driven reason.

**4. `--version` is registered only on `rootCmd`, not as a persistent
flag.**
Cobra's `InitDefaultVersionFlag` adds `--version` to the command's local
(non-persistent) flag set. `agents-lint --version` works; `agents-lint scan
--version` does not (unrecognized flag), which matches FR-CLI-06's wording
("`agents-lint --version`") — it describes top-level invocation, not a
flag every subcommand must also accept.

## Risks / Trade-offs

- **[Risk]** `"dev"` as a permanent default could look like a bug if a user
  runs a real release build and still sees `agents-lint version dev`. →
  **Mitigation:** out of scope for this change; `packaging-and-docs` owns
  wiring `-ldflags` in the release build path. This change's job is only to
  make that wiring possible (a package-level `var`), not to perform a
  release build itself.
- **[Risk]** Cobra's default version output template (`"<name> version
  <Version>\n"`) is fixed unless explicitly overridden. → **Mitigation**:
  FR-CLI-06 only requires "prints the version"; the default template
  satisfies that without needing a custom `VersionTemplate`.

## Open Questions

None — Cobra's built-in `Version` field is a direct, dependency-free fit
for FR-CLI-06, and the `var version = "dev"` seam defers all release-
specific decisions to `packaging-and-docs` without blocking this change.
