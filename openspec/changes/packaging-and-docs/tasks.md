## 1. Release build target (NFR-DIST-01, NFR-PERF-02)

- [x] 1.1 Add a `release` target to `Makefile` that builds `dist/agents-lint-linux-amd64` and `dist/agents-lint-darwin-arm64` via `CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go build -trimpath -ldflags="-s -w" -o dist/agents-lint-<os>-<arch> ./cmd/agents-lint`
- [x] 1.2 Add a `clean` target (or extend an existing one) that removes `dist/`
- [x] 1.3 Add `dist/` to `.gitignore`
- [x] 1.4 Run `make release` from a clean checkout and confirm both binaries are produced with no network access
- [x] 1.5 Run `ls -lh dist/` and confirm each binary is ≤ 15 MB; if not, investigate `-ldflags`/build tags before proceeding
- [x] 1.6 Smoke-test the native-platform binary: run `dist/agents-lint-<native> --version` and `dist/agents-lint-<native> scan ./AGENTS.md`, confirm output matches `go run ./cmd/agents-lint`

## 2. Project README (NFR-DX-02)

- [x] 2.1 Create `docs/README.md` with an installation section (binary download from `dist/` or release page, and `go install github.com/ingvar-goryainov/agents-lint/cmd/agents-lint@latest`)
- [x] 2.2 Add a usage section with real command examples for `scan`, `init`, `docs <rule-id>`, `--version`, and `--format sarif`
- [x] 2.3 Add a rule catalog table (ID, description, default severity) for all rules; cross-check every row against the output of `rules.KnownRuleIDs()` and `rules.RuleDocs()` (write a throwaway `go run` snippet or existing `docs` command output to verify, don't hand-guess severities)
- [x] 2.4 Add a configuration reference section: default `.agents-lint.yaml` location, `--config` flag, `path` key, `rules.<ID>.enabled`/`rules.<ID>.severity` keys, and one complete example config file
- [x] 2.5 Confirm the repo-root `README.md` is unmodified (`git diff --stat README.md` shows no changes)

## 3. Verify

- [x] 3.1 Run full suite: `go test ./...`
- [x] 3.2 Run `golangci-lint run` and fix any findings
- [x] 3.3 Run `make build` and `make test` to confirm the new Makefile targets didn't break existing ones
- [x] 3.4 Re-read `docs/README.md` end-to-end as a first-time user would and fix any command that doesn't match actual CLI output
