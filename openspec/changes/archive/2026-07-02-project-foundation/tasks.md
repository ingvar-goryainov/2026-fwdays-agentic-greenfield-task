## 1. Module setup

- [x] 1.1 Confirm the Go module path (GitHub org/repo) with the user (open
      question from design.md) before running `go mod init`
- [x] 1.2 Run `go mod init <module-path>` targeting Go 1.22+
- [x] 1.3 Add `github.com/spf13/cobra` as a dependency (`go get`)

## 2. Package layout

- [x] 2.1 Create `cmd/agents-lint/main.go` as the binary entrypoint
- [x] 2.2 Create `cmd/agents-lint/cmd/root.go` with the Cobra root command
      (`Execute()` function, no subcommands yet)
- [x] 2.3 Wire `main.go` to call `cmd.Execute()`
- [x] 2.4 Create `internal/lint/` package (empty except for the `Finding`
      type)
- [x] 2.5 Create `internal/lint/rules/` package placeholder (doc comment
      only, no rules yet)

## 3. Finding type

- [x] 3.1 Define `Finding` struct in `internal/lint` with fields: `RuleID
      string`, `Severity Severity`, `File string`, `Line int`, `Message
      string`
- [x] 3.2 Define a `Severity` type (e.g., `error`/`warning`) the `Finding`
      struct uses

## 4. Test and fixture scaffolding

- [x] 4.1 Create `testdata/` directory at repo root with a `.gitkeep` (or
      README note) documenting the "organized by rule ID" convention
      (TC-TEST-02)
- [x] 4.2 Add `github.com/stretchr/testify` as a test dependency (TC-TEST-01)
- [x] 4.3 Verify `go test ./...` runs and exits 0 with zero test files

## 5. Build tooling and license

- [x] 5.1 Add `Makefile` with `build`, `test`, and `lint` targets
- [x] 5.2 Add `LICENSE` file with MIT license text (BC-LICENSE-01)
- [x] 5.3 Verify `make build` (or `go build ./cmd/agents-lint`) completes in
      under 10 seconds on a clean checkout (NFR-DX-01)

## 6. Verification

- [x] 6.1 Run the built binary with no arguments — confirm exit 0 and help
      text printed
- [x] 6.2 Run the built binary with `--help` — confirm exit 0 and help text
      printed
- [x] 6.3 Build with `CGO_ENABLED=0` — confirm it still succeeds
- [x] 6.4 Confirm `internal/lint` compiles independently without importing
      `cmd/`
