## 1. Parser: document-wide code spans

- [x] 1.1 In `internal/lint/parser.go`, add `CodeSpan{Text string, Line
      int}` and a `CodeSpans []CodeSpan` field on `Document`
- [x] 1.2 Implement `populateCodeSpans(doc *Document, root ast.Node, src
      []byte)`: `ast.Walk(root, ...)` collecting every `*ast.CodeSpan`
      node's text (via the existing `nodeText` helper) and line (via the
      first `*ast.Text` child segment's start offset, passed through the
      existing `lineNumber` helper); call it from `Parse` after
      `populateSections`
- [x] 1.3 Test: a fixture with code spans inside a top-level paragraph, a
      list item, and a blockquote all appear in `doc.CodeSpans` with
      correct text and line numbers (confirms the full-tree walk reaches
      nesting `populateSections` doesn't)
- [x] 1.4 Test: a code span inside YAML frontmatter is NOT collected
      (frontmatter bytes are blanked before goldmark parses, per existing
      `extractFrontmatter` behavior)
- [x] 1.5 Run `go test ./internal/lint/...` — confirm existing parser tests
      (S00x-related) are unaffected

## 2. `lint.CodebaseRule` interface

- [x] 2.1 In `internal/lint/rule.go`, add the `CodebaseRule` interface:
      `ID() string`, `Check(doc *Document, repoRoot string) []Finding`

## 3. Rule C001 — file path references (FR-C001)

- [x] 3.1 Create `internal/lint/rules/c001.go` with `RuleC001 = "C001"` and
      a `C001` struct implementing `lint.CodebaseRule`
- [x] 3.2 Implement `looksLikePath(text string) bool`: reject if it
      contains whitespace, starts with `-`, contains `://`, contains `*` or
      `?`, or starts with `/`; accept if it contains `/`, OR if it matches
      `name.ext` where the substring after the final `.` contains at least
      one non-digit character
- [x] 3.3 Implement `C001.Check`: for each `doc.CodeSpans` entry where
      `looksLikePath(span.Text)` is true, `os.Stat(filepath.Join(repoRoot,
      span.Text))`; on any stat error, emit a `C001` finding — severity
      `error`, file `doc.Path`, line `span.Line`, message naming the
      missing path and repo-root-relative resolution
- [x] 3.4 Fixtures: `testdata/C001/valid.md` referencing real files that
      exist relative to the actual repo root (e.g. `` `go.mod` ``,
      `` `internal/lint/rule.go` ``) plus non-path inline code (a CLI flag,
      a URL, a glob, a root-relative string) that must NOT be flagged;
      `testdata/C001/invalid.md` referencing at least two distinct paths
      that do not exist
- [x] 3.5 Table-driven test in `internal/lint/rules/c001_test.go`: parse
      each fixture with `lint.Parse`, call `rules.C001{}.Check(doc,
      repoRoot)` with `repoRoot` set to the real repo root (`"../../.."`
      relative to the test file's package directory); assert zero findings
      for `valid.md` and exactly the expected count/lines for `invalid.md`
- [x] 3.6 Unit tests for `looksLikePath` covering every accept/reject
      branch in 3.2 individually (flags, URLs, globs, root-relative
      strings, multi-segment paths, bare `name.ext`, version-looking
      tokens like `1.22`/`v1.2.3`)
- [x] 3.7 Run `go test -cover ./internal/lint/rules/...` for `c001.go` —
      confirm ≥90% line coverage

## 4. Rule C002 — tool/command evidence (FR-C002)

- [x] 4.1 Create `internal/lint/rules/c002.go` with `RuleC002 = "C002"` and
      a `C002` struct implementing `lint.CodebaseRule`
- [x] 4.2 Define the fixed evidence table (unexported
      `map[string][]string`) per design.md: `terraform` → `*.tf`,
      `.terraform.lock.hcl`; `npm` → `package.json`, `package-lock.json`;
      `yarn` → `yarn.lock`; `pnpm` → `pnpm-lock.yaml`; `docker` →
      `Dockerfile`; `go` → `go.mod`; `make` → `Makefile`; `cargo` →
      `Cargo.toml`; `kubectl` → `k8s/`, `kubernetes/`,
      `kustomization.yaml`; `git` → `.git`
- [x] 4.3 Implement `hasEvidence(repoRoot, patterns []string) bool`: for
      each pattern, `filepath.Glob(filepath.Join(repoRoot, pattern))` for
      glob-shaped patterns and `os.Stat(filepath.Join(repoRoot, pattern))`
      for exact names/dirs; true if any pattern matches
- [x] 4.4 Implement `C002.Check`: for each `doc.CodeSpans` entry whose text
      is an exact key in the evidence table, if `hasEvidence` is false,
      emit a `C002` finding — severity `warning`, file `doc.Path`, line
      `span.Line`, message naming the tool and the evidence patterns
      checked
- [x] 4.5 Fixtures: `testdata/C002/valid.md` referencing `` `go` `` (repo
      already has `go.mod` at its real root) plus a non-catalog bare word
      (e.g. a rule ID) that must NOT be flagged; `testdata/C002/invalid.md`
      referencing `` `terraform` `` (no `.tf`/`.terraform.lock.hcl` at repo
      root)
- [x] 4.6 Table-driven test in `internal/lint/rules/c002_test.go`: parse
      each fixture, call `rules.C002{}.Check(doc, repoRoot)` with
      `repoRoot` set to the real repo root; assert zero findings for
      `valid.md` and one `warning`-severity finding for `invalid.md`
- [x] 4.7 Unit tests for `hasEvidence` against a `t.TempDir()`: a glob
      pattern that matches, an exact filename that matches, a directory
      name that matches, and no patterns matching
- [x] 4.8 Run `go test -cover ./internal/lint/rules/...` for `c002.go` —
      confirm ≥90% line coverage

## 5. Wire codebase-awareness rules into the registry

- [x] 5.1 In `internal/lint/rules/registry.go`, add `CodebaseAwareRules()
      []lint.CodebaseRule` returning `[]lint.CodebaseRule{C001{}, C002{}}`
- [x] 5.2 Extend `Run(path string) ([]lint.Finding, error)`: after the
      existing `DefaultRules()` loop, call `os.Getwd()` (wrap any error as
      `fmt.Errorf("rules: resolving repo root: %w", err)` and return it),
      then run each `CodebaseAwareRules()` rule with `Check(doc, repoRoot)`
      and append its findings
- [x] 5.3 Extend `KnownRuleIDs()` to append each `CodebaseAwareRules()`
      rule's `ID()` after the schema rule IDs
- [x] 5.4 Update `internal/lint/rules/doc.go`'s package comment — remove
      the "added by a later change" note now that C001/C002 exist
- [x] 5.5 Test `KnownRuleIDs()` returns exactly `["S001", "S002", "S003",
      "S004", "S005", "C001", "C002"]`
- [x] 5.6 Test `Run`, using `t.Chdir` into a temp directory containing an
      AGENTS.md fixture with a broken path reference: confirm the returned
      findings include both schema findings (if any) and the `C001`
      finding, and that `repoRoot` resolution uses the chdir'd directory
      (place an evidence file in the temp dir and confirm a `C002`
      reference to it produces no finding)
- [x] 5.7 Test `Run` when the target file doesn't exist (S001 fails):
      confirm only the S001 finding is returned and neither C001 nor C002
      ran (no panic from a nil `Document`) — covered by the existing
      `TestRun_MissingFileShortCircuits`
- [x] 5.8 (discovered while wiring `Run`) Existing schema-rule tests and
      fixtures that assumed `repoRoot` was irrelevant needed updating now
      that it's derived from `os.Getwd()`: `internal/lint/rules`'
      `TestRun_CleanFileHasNoFindings` and `cmd/agents-lint/cmd`'s
      `TestRunScan`/`TestRunScan_Config` cases that use
      `testdata/S003/valid.md` (which references `` `docs/requirements.md` ``)
      now `t.Chdir` to the repo root before calling `Run`/`runScan`;
      `testdata/config/path_override.yaml`'s `path:` value and the rule-count
      assertions (5 known rules → 7) were updated to match

## 6. `cmd`-level integration

- [x] 6.1 Confirm (via a `cmd/agents-lint/cmd/scan_test.go` case) that
      `agents-lint scan` against a fixture with a broken file reference
      reports a `C001` line in the expected `FR-OUT-01` format and exits
      `1`
- [x] 6.2 Confirm a `.agents-lint.yaml` with `rules.C001.enabled: false`
      suppresses the `C001` finding end-to-end through `runScan` (reuses
      `internal/config`'s existing `Apply`, no new config code)
- [x] 6.3 Manual check: run `go run ./cmd/agents-lint scan` against this
      repo's own `AGENTS.md` from the repo root — confirm no unexpected
      `C001`/`C002` findings against real content (or, if any surface,
      confirm they're genuine and not heuristic false positives before
      moving on). Surfaced two heuristic false-positive classes (template
      placeholders like `` `<rule_id>.go` ``, and Go-idiom qualified
      identifiers like `` `fmt.Print` ``) — fixed by tightening
      `looksLikePath`/`isBareFileName` (see design.md). The remaining 8
      findings are genuine: this repo's own `AGENTS.md` references
      `internal/types/` and `internal/rules/`, neither of which exist
      (actual code lives under `internal/lint/`), plus three OpenSpec
      artifact names (`proposal.md`/`design.md`/`tasks.md`) that exist only
      inside `openspec/changes/<change>/`, not at repo root (expected,
      given C001's non-recursive repo-root-only scope). Not fixed here —
      editing `AGENTS.md` content is outside this change's scope.

## 7. Full verification

- [x] 7.1 Run `go test ./...` — confirm green
- [x] 7.2 Run `go test -cover ./internal/lint/...` — confirm `c001.go`,
      `c002.go`, and the extended `parser.go`/`registry.go` are at or above
      90% line coverage (98.3% total; `c001.go`/`c002.go` 100%, `Run` 92.9%,
      `codeSpanLine` 87.5% — the untested branches are an unreachable
      `os.Getwd()` failure and a code span with no text child, both
      consistent with this codebase's existing convention of not
      simulating unlikely OS-level failures beyond what's already covered)
- [x] 7.3 Run `golangci-lint run` — confirm clean (0 issues)
- [x] 7.4 Confirm no new entries were added to `go.mod` (goldmark's `ast`
      package is already a dependency) — `git diff go.mod go.sum` is empty
- [x] 7.5 Confirm `internal/config`, `internal/reporter`, and
      `cmd/agents-lint/cmd/scan.go`'s existing exported signatures are
      unchanged — only `rules.Run`'s internals and `rules.KnownRuleIDs()`'s
      return value grew (`git diff --stat` confirms zero changes to those
      files)
