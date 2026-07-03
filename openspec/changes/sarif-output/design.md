## Context

`internal/lint.Finding` (`RuleID`, `Severity`, `File`, `Line`, `Message`) is
the single contract every rule already writes to and the only contract
`internal/reporter.WriteText` reads from (see `openspec/specs/text-reporting`).
`sarif-output` adds a second renderer for that same slice — no new data
model, no changes to `Finding` or the rule interfaces. `Line` is a 0-based
"unset" sentinel today: `FR-S001`'s `CheckFileExists` (`internal/lint/rules/s001.go`)
produces a `Finding` with `Line: 0` when the file itself is missing, since
there is no line to point at. Every other existing rule sets `Line` to a
1-based line number from the parsed document.

The CLI already has one persistent, validated-at-parse-time flag
(`--config`, `cmd/agents-lint/cmd/root.go`) and one command (`scan`,
`cmd/agents-lint/cmd/scan.go`) that resolves inputs before calling a
reporter. `--format` follows the same shape.

## Goals / Non-Goals

**Goals:**
- Render `[]lint.Finding` as a syntactically valid SARIF v2.1.0 log that
  GitHub Code Scanning and VS Code's SARIF viewer accept.
- Let `agents-lint scan --format sarif` opt into that output with no change
  to exit-code behavior (FR-CLI-02 stays format-independent).
- Keep `Finding` and both `Rule`/`CodebaseRule` interfaces untouched — the
  SARIF writer is a pure function of `[]Finding`, exactly like
  `WriteText`.

**Non-Goals:**
- Rule metadata (descriptions, help URIs, examples) in
  `tool.driver.rules[].fullDescription` / `helpUri` — that depends on
  `FR-CLI-07`'s `docs` command (stretch, not yet built) supplying a
  structured catalog. `sarif-output` populates `tool.driver.rules[].id`
  only; SARIF v2.1.0 does not require more.
- Multiple `run`s, `originalUriBaseIds`, or artifact-relative URI
  resolution beyond the `File` string a `Finding` already carries verbatim
  (matches `WriteText`, which does the same).
- A `--format` value beyond `text`/`sarif` (e.g. `json`) — out of scope per
  `FR-CLI-05`.

## Decisions

**1. New file `internal/reporter/sarif.go`, function
`WriteSARIF(w io.Writer, findings []lint.Finding) error`.**
Mirrors `WriteText`'s signature shape (`io.Writer` in, `[]lint.Finding` in,
`error` out) so `scan.go` can select between them with a single branch.
Unlike `WriteText`, it takes no `ruleCount`/`Options` — SARIF has no
"N rules passed" success line (an empty `findings` slice simply produces a
log with an empty `results` array) and no color option.

**2. Hand-rolled SARIF structs + `encoding/json`, no new dependency.**
SARIF v2.1.0 is a plain JSON schema; a handful of structs
(`sarifLog`, `sarifRun`, `sarifTool`, `sarifDriver`, `sarifReportingDescriptor`,
`sarifResult`, `sarifMessage`, `sarifLocation`, `sarifPhysicalLocation`,
`sarifArtifactLocation`, `sarifRegion`) cover exactly the fields this
project emits. Pulling in a SARIF library would be the first non-stdlib,
non-Cobra/goldmark/yaml dependency in `internal/` and CLAUDE.md requires
asking before adding one; a ~40-line struct set avoids that entirely and
keeps `internal/` dependency-free per the project's boundaries.
  - *Alternative considered:* `github.com/owenrumney/go-sarif` — rejected,
    unnecessary dependency for a fixed, small output shape.

**3. Severity mapping is a 2-case switch, matching `reporter.severityColor`'s
shape.**
`lint.SeverityError` → `"error"`, `lint.SeverityWarning` → `"warning"`.
SARIF also defines `"note"`/`"none"`, but `lint.Severity` only ever has two
values (`internal/lint/finding.go`), so there is no third case to map.

**4. `Line == 0` omits `region` entirely rather than emitting `region.startLine: 0`.**
SARIF's `region.startLine` is 1-based; `0` is not a valid line number and
some consumers (GitHub Code Scanning) reject or misrender it. When
`Finding.Line == 0` (today, only `FR-S001`'s file-not-found case), the
`physicalLocation` is emitted with just `artifactLocation.uri` and no
`region` field, rather than a fabricated line number.

**5. `tool.driver.rules[]` is built from the distinct `RuleID`s actually
present in `findings`, not from `rules.KnownRuleIDs()`.**
Keeps the writer a pure function of its input (no import from
`internal/lint/rules`, avoiding a new inbound dependency on the rules
package from `internal/reporter`) and avoids listing rules with zero
findings, which SARIF does not require. Order: first-seen order in
`findings`, deduplicated — deterministic because `findings` order is
already deterministic (rules run in a fixed order today).

**6. `--format` is a persistent flag on `rootCmd`, validated in `scan`'s
`RunE` before `runScan` does any work.**
Matches the existing `--config` flag's placement (`root.go`) and
precedent from `configuration-support` (FR-CLI-04 lived entirely in that
capability despite being a root-level flag). An unrecognized `--format`
value (e.g. `--format xml`) fails fast with
`unknown format "xml": must be "text" or "sarif"` on stderr, exit non-zero,
before any file I/O — same fail-fast shape as a missing `--config` path.

**7. `runScan` gains a `format string` parameter; reporter selection is a
2-way branch, not a registry/interface.**
Only two reporters exist and FR-CLI-05 only ever needs two; a
`map[string]func(...)` or a `Reporter` interface would be speculative
generality for a fixed, closed set (CLAUDE.md: no premature abstraction).

## Risks / Trade-offs

- **[Risk]** GitHub Code Scanning has stricter-than-spec expectations for
  some SARIF fields in practice (e.g. wanting non-empty
  `shortDescription.text`). → **Mitigation:** `tool.driver.rules[].id` is
  always populated (it's the one field this project actually has data
  for); richer descriptions can be layered on later via `FR-CLI-07`
  without changing the shape decided here.
- **[Risk]** `Finding.File` is not normalized to a repo-relative path
  today (it's whatever path string the CLI was invoked with), and SARIF's
  `artifactLocation.uri` is nominally relative to a base URI. →
  **Mitigation:** `WriteText` has the identical property and it hasn't
  been a problem in practice; `sarif-output` writes `File` verbatim rather
  than inventing new path-resolution logic outside this change's scope.

## Open Questions

None — SARIF v2.1.0's minimal-viable-log shape and this project's
existing `Finding` fields are a direct fit; no unresolved decisions block
`tasks.md`.
