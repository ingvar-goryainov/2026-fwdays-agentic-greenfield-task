## 1. Core doc type

- [x] 1.1 Add `RuleDoc` struct to `internal/lint` (ID, Description, Severity,
      Valid, Invalid, FixGuidance) per design.md decision 1

## 2. Per-rule documentation metadata

- [x] 2.1 Add `Doc() lint.RuleDoc` to S001 (file-exists) with a short
      valid/invalid example and fix guidance
- [x] 2.2 Add `Doc() lint.RuleDoc` to S002 (required sections)
- [x] 2.3 Add `Doc() lint.RuleDoc` to S003 (agent block fields)
- [x] 2.4 Add `Doc() lint.RuleDoc` to S004 (duplicate agent names)
- [x] 2.5 Add `Doc() lint.RuleDoc` to S005 (YAML frontmatter validity)
- [x] 2.6 Add `Doc() lint.RuleDoc` to C001 (file paths must resolve)
- [x] 2.7 Add `Doc() lint.RuleDoc` to C002 (tool/command evidence)

## 3. Registry lookup

- [x] 3.1 Add `RuleDocs() map[string]lint.RuleDoc` to
      `internal/lint/rules/registry.go`, mirroring `KnownRuleIDs()`
      (iterates `DefaultRules()` + `CodebaseAwareRules()`, plus S001)
- [x] 3.2 Unit test: `RuleDocs()` returns one entry per ID in
      `KnownRuleIDs()`, none with empty Description/FixGuidance/Valid/Invalid

## 4. `docs` CLI command

- [x] 4.1 Implement `runDocs(w io.Writer, ruleID string) error` in
      `cmd/agents-lint/cmd/docs.go`: look up `ruleID` in `rules.RuleDocs()`,
      print description/severity/valid/invalid/fix guidance on a hit, return
      a clear error on a miss (FR-CLI-07, BC-OUTPUT-02). No `exitCode`
      return needed (unlike `runScan`/`runInit`): `docs` has no
      success-but-nonzero-exit case, so a plain `error` plus main.go's
      existing `Execute() != nil → os.Exit(1)` is sufficient.
- [x] 4.2 Register `docsCmd` (`Use: "docs <rule-id>"`, `Args: cobra.ExactArgs(1)`)
      on `rootCmd`, following `scan.go`'s RunE → runDocs pattern
- [x] 4.3 Unit tests for `runDocs`: known rule ID (S002, C001, S001) prints
      all fields; unknown rule ID (`S999`) errors and writes no output;
      lowercase rule ID (`s002`) is treated as unknown
- [x] 4.4 Cobra-level tests: missing rule ID argument (`agents-lint docs`)
      produces a usage error; known/unknown rule ID via the full command
      tree

## 5. Verify

- [x] 5.1 Run `go test ./...` — full suite green
- [x] 5.2 Run `go test -cover ./internal/...` — confirm ≥ 90% coverage on
      touched files (internal/lint 98.5%, internal/lint/rules 98.2%, all
      new Doc()/RuleDocs()/s001Doc code at 100%)
- [x] 5.3 Run `golangci-lint run` — clean (0 issues)
- [x] 5.4 Manually ran `go run ./cmd/agents-lint docs S002` and
      `go run ./cmd/agents-lint docs S999` — output matches spec scenarios;
      also fixed multi-line example indentation (`indentBlock`) discovered
      during this manual check
- [ ] 5.5 Update README rule catalog / usage section to mention `docs
      <rule-id>` (NFR-DX-02) — **deferred**: the repo root `README.md` is
      the course-assignment brief, not agents-lint's own README. A proper
      README with a rule catalog is `packaging-and-docs` (phase 7), which
      hasn't started yet; writing one now would duplicate/conflict with
      that capability's scope.
