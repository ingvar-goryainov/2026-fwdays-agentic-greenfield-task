## 1. Dockerfile (NFR-DIST-02)

- [x] 1.1 Create repo-root `Dockerfile` with a builder stage (`golang:1.26`, matching go.mod's `go 1.26.4` directive) that runs `CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X .../cmd.version=${VERSION}" -o /agents-lint ./cmd/agents-lint`, accepting a `VERSION` build-arg (default `dev`)
- [x] 1.2 Add a runtime stage from `gcr.io/distroless/static-debian12`, `COPY --from=builder /agents-lint /agents-lint`, `WORKDIR /workspace`, non-root `USER nonroot:nonroot`, `ENTRYPOINT ["/agents-lint"]`
- [x] 1.3 Add a repo-root `.dockerignore` excluding `dist/`, `.git/`, `openspec/`, and other non-build files to keep the build context small
- [x] 1.4 Build the image locally (`docker build -t agents-lint .`) and confirm it succeeds with no errors — image built, 5.4MB, confirmed no shell present

## 2. Makefile targets

- [x] 2.1 Add a `docker-build` target to `Makefile` that runs `docker build -t agents-lint .` (accept `VERSION` via `--build-arg` if set)
- [x] 2.2 Add a `docker-run` target (or documented example) demonstrating `docker run --rm -v $(pwd):/workspace agents-lint scan`

## 3. Behavioral parity verification

- [x] 3.1 Run the native binary (`go run ./cmd/agents-lint scan ./AGENTS.md`) and the container (`docker run --rm -v $(pwd):/workspace agents-lint scan`) against the same `AGENTS.md` and confirm identical stdout and exit code — verified against an `init`-generated valid AGENTS.md (both `✓ ... (7 rules passed)`, exit 0)
- [x] 3.2 Repeat against a deliberately invalid AGENTS.md fixture and confirm both exit 1 with matching finding lines and summary — verified against the repo's own `AGENTS.md` (both print identical 8 C001 findings, exit 1)
- [x] 3.3 Confirm `docker run --rm agents-lint --version` prints the version baked in via the `VERSION` build-arg (not `dev`, when built with one) — default build prints `dev`, `--build-arg VERSION=v9.9.9-test` prints `v9.9.9-test`
- [x] 3.4 Confirm `--format sarif` and `--config <path>` work identically through the container as they do natively — verified both produce identical SARIF output and identical config-driven rule-disable behavior

## 4. Documentation

- [x] 4.1 Add a Docker installation/usage section to `docs/README.md` with a concrete `docker build` (or `make docker-build`) command and a concrete `docker run` example with a volume mount
- [x] 4.2 Note the distroless runtime has no shell (debugging is via output/exit code, not `docker exec`)
- [x] 4.3 Confirm the repo-root `README.md` is unmodified (`git diff --stat README.md` shows no changes) — confirmed clean

## 5. Verify

- [x] 5.1 Run full suite: `go test ./...` — all packages pass
- [x] 5.2 Run `golangci-lint run` and fix any findings — 0 issues
- [x] 5.3 Run `make build`, `make test`, and `make release` to confirm nothing existing broke — all succeed, both dist binaries produced
- [x] 5.4 Re-read the new `docs/README.md` Docker section end-to-end and confirm every command runs as written — `make docker-build`, `make docker-build VERSION=...`, and all three `docker run` examples verified above
