## 1. Dependencies

- [x] 1.1 Add `github.com/yuin/goldmark` (TC-STACK-03) via `go get`
- [x] 1.2 Add `gopkg.in/yaml.v3` (TC-STACK-04) via `go get`

## 2. Parser and Document model

- [x] 2.1 Define `Document`, `Heading`, `AgentBlock`, and `Frontmatter`
      types in `internal/lint/parser.go`
- [x] 2.2 Implement raw-byte frontmatter detection/extraction (leading
      `---`-delimited block) ahead of goldmark parsing
- [x] 2.3 Implement `Parse(path string) (*Document, error)`: reads the file,
      strips frontmatter, parses the remainder with goldmark, and walks the
      AST to populate `H2Sections`
- [x] 2.4 Implement H3 agent-block extraction scoped to sections matching
      `Agent`/`Agents` (case-insensitive), populating `AgentBlocks` with
      `Name` and `Line`
- [x] 2.5 Implement bold-label field detection within each agent block
      (`Role`, `Instructions`/`Responsibilities`, `Tools`, `Context`) to set
      `HasRole`/`HasInstructions`/`HasTools`/`HasContext`
- [x] 2.6 Implement byte-offset-to-line-number conversion helper used by all
      of the above
- [x] 2.7 Write parser unit tests covering: no frontmatter, valid
      frontmatter, malformed frontmatter delimiters, no Agent(s) section,
      multiple agent blocks, bold-label variants (`Instructions` vs
      `Responsibilities`)

## 3. Rule S001 — file exists (FR-S001)

- [x] 3.1 Generate fixtures: `testdata/S001/valid.md` (file present),
      `testdata/S001/invalid.md` documented as "path does not exist" case
      (test references a non-existent path directly; no fixture content
      needed for the invalid case itself)
- [x] 3.2 Implement `internal/lint/rules/s001.go`: `CheckFileExists(path
      string) *lint.Finding` returning nil or an actionable `SeverityError`
      finding
- [x] 3.3 Write table-driven tests in `internal/lint/rules/s001_test.go`
- [x] 3.4 Run `go test ./internal/lint/... -run TestS001` until green

## 4. Rule S002 — required sections (FR-S002)

- [x] 4.1 Generate fixtures: `testdata/S002/valid.md` (has `## Agents`),
      `testdata/S002/invalid.md` (no Agent/Agents H2 section)
- [x] 4.2 Implement `internal/lint/rules/s002.go` against `Document.H2Sections`
- [x] 4.3 Write table-driven tests in `internal/lint/rules/s002_test.go`,
      including the case-insensitive `## Agent` vs `## Agents` spelling
- [x] 4.4 Run `go test ./internal/lint/... -run TestS002` until green

## 5. Rule S003 — agent block completeness (FR-S003)

- [x] 5.1 Generate fixtures: `testdata/S003/valid.md` (complete block, using
      the repo's own `AGENTS.md` agent-block shape as the model),
      `testdata/S003/invalid.md` (block missing role and all optional
      fields)
- [x] 5.2 Implement `internal/lint/rules/s003.go` against
      `Document.AgentBlocks`, reporting one finding per missing piece
- [x] 5.3 Write table-driven tests in `internal/lint/rules/s003_test.go`
      covering: missing role only, missing all of
      instructions/tools/context, complete block, `Responsibilities` as an
      instructions synonym
- [x] 5.4 Run `go test ./internal/lint/... -run TestS003` until green
- [x] 5.5 Sanity-check `s003.go` against the repo's real `AGENTS.md` (should
      report zero findings) as a manual cross-check, not an automated test

## 6. Rule S004 — no duplicate agent names (FR-S004)

- [x] 6.1 Generate fixtures: `testdata/S004/valid.md` (distinct names),
      `testdata/S004/invalid.md` (two blocks with the same name in different
      casing)
- [x] 6.2 Implement `internal/lint/rules/s004.go` against
      `Document.AgentBlocks`, comparing trimmed, case-insensitive names
- [x] 6.3 Write table-driven tests in `internal/lint/rules/s004_test.go`
- [x] 6.4 Run `go test ./internal/lint/... -run TestS004` until green

## 7. Rule S005 — valid frontmatter YAML (FR-S005)

- [x] 7.1 Generate fixtures: `testdata/S005/valid.md` (well-formed
      frontmatter), `testdata/S005/invalid.md` (malformed YAML frontmatter);
      confirm a file with no frontmatter is also covered as a passing case
      in the test table (not necessarily a separate fixture file)
- [x] 7.2 Implement `internal/lint/rules/s005.go`: calls `yaml.Unmarshal` on
      `Document.Frontmatter.Raw` (skips entirely if `Frontmatter` is nil)
- [x] 7.3 Write table-driven tests in `internal/lint/rules/s005_test.go`
- [x] 7.4 Run `go test ./internal/lint/... -run TestS005` until green

## 8. Rule engine orchestration

- [x] 8.1 Implement `DefaultRules() []lint.Rule` in
      `internal/lint/rules/registry.go` returning S002–S005 in fixed ID
      order (lives in `rules`, not `internal/lint`, to avoid an import cycle
      — see design.md's "Package split" decision)
- [x] 8.2 Implement `Run(path string) ([]lint.Finding, error)` in
      `internal/lint/rules/registry.go`: S001 short-circuit via
      `CheckFileExists`, then `lint.Parse`, then iterate `DefaultRules()`
- [x] 8.3 Write `internal/lint/rules/registry_test.go` covering: missing
      file short-circuits with only an S001 finding, existing file runs all
      rules and aggregates findings in order

## 9. Full verification

- [x] 9.1 Run `go test ./...` — confirm green and coverage ≥ 90% for every
      new file in `internal/lint/` and `internal/lint/rules/` (`go test
      -cover ./internal/...`)
- [x] 9.2 Run `golangci-lint run` — confirm clean
- [x] 9.3 Confirm no `fmt.Print*` calls were introduced in `internal/`
      (library-first constraint)
- [x] 9.4 Confirm no network calls or CGO were introduced
- [x] 9.5 Manually run `Run("AGENTS.md")` against this repo's own AGENTS.md
      (e.g. via a throwaway `go run` snippet or existing test) and confirm
      it returns zero findings, validating the dogfooding assumption in
      design.md
