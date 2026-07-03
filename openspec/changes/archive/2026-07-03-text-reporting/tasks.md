## 1. Package scaffolding

- [x] 1.1 Create `internal/reporter/` package with a doc comment explaining
      its role (renders `[]lint.Finding` to text, no CLI/fmt.Print
      dependency, per `AGENTS.md`'s library-first architecture)
- [x] 1.2 Define `Options struct { Color bool }` in `internal/reporter/text.go`

## 2. Per-finding line rendering (FR-OUT-01)

- [x] 2.1 Implement the line formatter: `severity  rule-id  file:line —
      message` (two literal spaces between fields, ` — ` before message),
      reading `Severity`, `RuleID`, `File`, `Line`, `Message` verbatim from
      each `lint.Finding`
- [x] 2.2 Write findings to the output `io.Writer` in input-slice order (no
      sorting/grouping)
- [x] 2.3 Golden fixture: `testdata/reporter/golden/single_error.txt`
- [x] 2.4 Golden fixture: `testdata/reporter/golden/multi_mixed_order.txt`
      (findings from more than one rule, confirming order is preserved)

## 3. Summary line (FR-OUT-02)

- [x] 3.1 Implement the summary line: count `SeverityError` and
      `SeverityWarning` findings, write `N error(s), M warning(s)` after all
      per-finding lines, using the literal `(s)` suffix regardless of count
      (see design.md decision)
- [x] 3.2 Golden fixture: `testdata/reporter/golden/mixed_severities.txt`
      (2 errors, 1 warning)
- [x] 3.3 Golden fixture: `testdata/reporter/golden/only_warnings.txt` (0
      errors, N warnings)

## 4. Success line (FR-OUT-04)

- [x] 4.1 Implement the empty-findings short-circuit: when `len(findings)
      == 0`, write only `✓ AGENTS.md is valid (N rules passed)` using the
      `ruleCount` argument, and return — no per-finding lines, no summary
      line
- [x] 4.2 Golden fixture: `testdata/reporter/golden/no_findings.txt`
- [x] 4.3 Table-driven test asserting the success line is the *entire*
      output (not just a prefix) when findings is empty

## 5. Color and NO_COLOR (BC-OUTPUT-01)

- [x] 5.1 Implement ANSI severity coloring (red for `error`, yellow for
      `warning`, reset after) applied only when `Options.Color` is true
- [x] 5.2 Implement `NO_COLOR` override: check `os.Getenv("NO_COLOR")`
      inside `WriteText` and suppress color codes when it's set to any
      non-empty value, regardless of `Options.Color`
- [x] 5.3 Table-driven tests: `Color: false` emits no ANSI codes,
      `Color: true` + `NO_COLOR` unset emits ANSI codes around the severity
      token, `Color: true` + `NO_COLOR` set emits no ANSI codes
- [x] 5.4 Confirm default (zero-value) `Options{}` produces plain,
      uncolored output — no pager required (BC-OUTPUT-01)

## 6. Message fidelity (BC-OUTPUT-02)

- [x] 6.1 Table-driven test: a `Finding.Message` longer than 200 characters
      appears in the output unmodified and untruncated
- [x] 6.2 Table-driven test: a `Finding.Message` containing the literal
      substrings used by the formatter itself (e.g. `" — "` or double
      spaces) is not mangled or double-escaped by the line formatter

## 7. Full verification

- [x] 7.1 Run `go test ./internal/reporter/... -run TestWriteText` until
      green
- [x] 7.2 Run `go test ./...` — confirm green
- [x] 7.3 Run `go test -cover ./internal/...` — confirm `internal/reporter`
      is at or above 90% line coverage
- [x] 7.4 Run `golangci-lint run` — confirm clean
- [x] 7.5 Confirm no `fmt.Print*` calls were introduced in
      `internal/reporter` (library-first constraint — output only goes
      through the caller-supplied `io.Writer`)
- [x] 7.6 Confirm no new entries were added to `go.mod` (color/NO_COLOR use
      only `os` and raw ANSI escape sequences, per design.md)
