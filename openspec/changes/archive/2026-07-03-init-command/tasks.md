## 1. Skeleton content

- [x] 1.1 Create `internal/lint/skeleton/skeleton.go` with an exported
      `Content() []byte` (or `string`) function returning the static
      AGENTS.md skeleton described in design.md (frontmatter + `## Agent` +
      one `###` agent block with name, role, instructions)
- [x] 1.2 Add `internal/lint/skeleton/skeleton_test.go` that runs the real
      S001-S005 rules (via `rules.Run` against a temp file, or
      `rules.DefaultRules()` against a parsed `lint.Document`, plus
      `rules.CheckFileExists`) over `Content()` and asserts zero
      error-severity findings — this is the regression test guaranteeing
      FR-CLI-03's "always passes" property

## 2. `init` command

- [x] 2.1 Create `cmd/agents-lint/cmd/init.go` with a `runInit(w io.Writer,
      path string, force bool) (exitCode int, err error)` function
      following the `runScan` pattern in `scan.go`
- [x] 2.2 Implement the existing-file guard in `runInit`: if `path` exists
      and `force` is false, return exit code 1 and a wrapped, actionable
      error (mention `--force`); do not write the file
- [x] 2.3 Implement the write path in `runInit`: write `skeleton.Content()`
      to `path` (create parent behavior matches `os.WriteFile` defaults —
      no directory creation beyond what already exists), then write a
      confirmation line to `w` naming the path written
- [x] 2.4 Wire up the Cobra `initCmd`: `Use: "init [path]"`,
      `Args: cobra.MaximumNArgs(1)`, default path `./AGENTS.md` (reuse or
      mirror `defaultAgentsMDPath` from `scan.go`), a `--force` bool flag,
      `RunE` calling `runInit` and calling `os.Exit` on non-zero exit code
      (matching `scanCmd`'s pattern)
- [x] 2.5 Register `initCmd` on `rootCmd` via `init()`

## 3. Tests

- [x] 3.1 Add `cmd/agents-lint/cmd/init_test.go` covering: default path
      with no existing file (exit 0, file written, confirmation message),
      explicit path argument, existing file without `--force` (exit
      non-zero, file untouched, actionable error), existing file with
      `--force` (exit 0, file overwritten)
- [x] 3.2 Add a fixture/integration-style test asserting that running
      `agents-lint init` followed by `agents-lint scan` on the same path
      (or the equivalent in-process `rules.Run`) reports zero errors and
      the FR-OUT-04 success line

## 4. Verification

- [x] 4.1 Run `go test ./...` and confirm all green
- [x] 4.2 Run `go test -cover ./internal/lint/skeleton/...` and confirm
      coverage ≥ 90%
- [x] 4.3 Run `golangci-lint run` and confirm clean
- [x] 4.4 Manually run `go run ./cmd/agents-lint init` in a scratch
      directory, then `go run ./cmd/agents-lint scan` on the result, and
      confirm the success line prints
- [x] 4.5 Commit with `feat(init-command): implement agents-lint init
      (FR-CLI-03)`
