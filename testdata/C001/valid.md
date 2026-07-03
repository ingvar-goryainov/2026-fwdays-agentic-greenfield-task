# AGENTS.md

## Agents

### maintainer

Keeps the repo tidy.

**Context:** See `go.mod` and `internal/lint/rule.go` for the module layout.

Non-path inline code that must NOT be treated as a file reference:
a flag like `--config`, a URL like `https://example.com`, a glob like
`**/*.go`, and a root-relative string like `/api/users`.
