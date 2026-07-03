## 1. Version variable (FR-CLI-06)

- [x] 1.1 Add `cmd/agents-lint/cmd/version.go` with `var version = "dev"` (package-level, `-ldflags -X`-injectable)
- [x] 1.2 In `version.go`'s `init()`, set `rootCmd.Version = version`

## 2. `--version` flag wiring

- [x] 2.1 Verify Cobra registers `--version` (and `-v` shorthand) on `rootCmd` once `Version` is non-empty (no extra flag-registration code needed)
- [x] 2.2 Add a command-level test (e.g. `cmd/agents-lint/cmd/version_test.go`) asserting `agents-lint --version` prints the version string and exits 0 without invoking `scan`/`init` `RunE`
- [x] 2.3 Add a test asserting the fallback value is `"dev"` when no ldflags override is set
- [x] 2.4 Add a test asserting `agents-lint scan --version` fails with an unrecognized-flag error (confirms `--version` is root-only, not persistent)

## 3. Verify

- [x] 3.1 Run full suite: `go test ./...`
- [x] 3.2 Run `go test -cover ./cmd/...` and confirm the new file stays ≥ 90% line coverage
- [x] 3.3 Run `golangci-lint run` and fix any findings
- [x] 3.4 Manually run `go run ./cmd/agents-lint --version` and confirm it prints `agents-lint version dev` and exits 0
