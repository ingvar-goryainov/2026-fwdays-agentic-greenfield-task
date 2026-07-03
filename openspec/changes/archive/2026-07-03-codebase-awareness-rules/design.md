## Context

`internal/lint.Document` (built by `internal/lint.Parse`) currently only
records what `schema-validation-rules` needs: top-level `##` headings and
`###` agent blocks with per-field booleans (`internal/lint/parser.go`'s
`populateSections`, which walks only `root.FirstChild()...NextSibling()` —
i.e. top-level block nodes). It does not record inline-code (backtick) spans
anywhere in the document, which is what both new rules need to inspect.

The `lint.Rule` interface (`ID() string`, `Check(doc *Document) []Finding`)
is deliberately narrow — schema rules only ever need the parsed document.
C001 and C002 additionally need a **repo root** to resolve paths and look
for evidence files against, which no existing rule needs.

`rules.Run(path string) ([]lint.Finding, error)` is the single entry point
`cmd/agents-lint/cmd/scan.go` calls; it already runs S001 specially (outside
the `Rule` interface, since it must run *before* a `Document` exists) and
then the four `lint.Rule`s from `DefaultRules()`. `rules.KnownRuleIDs()` is
the source of truth `internal/config` validates `.agents-lint.yaml` rule IDs
against, and its own doc comment already anticipates this change: "Codebase-
awareness rules (FR-C001/FR-C002) are added by a later change."

The product brief documents `agents-lint scan .` as being run from the
repository root, and `defaultAgentsMDPath = "./AGENTS.md"` /
`config.DefaultFileName` are both already cwd-relative — there is an
established convention that the process's working directory *is* the repo
root, with no `--root`/`--repo-root` flag anywhere in `docs/requirements.md`.

## Goals / Non-Goals

**Goals:**
- FR-C001: every inline-code span in AGENTS.md whose text looks like a
  filesystem path is resolved relative to the repo root; a candidate that
  does not exist produces an error-severity `C001` finding at the span's
  line.
- FR-C002: every inline-code span whose text exactly matches a name in a
  fixed catalog of common developer tools (`terraform`, `npm`, `yarn`,
  `pnpm`, `docker`, `go`, `make`, `cargo`, `kubectl`, `git`) is checked
  against that tool's repo-root evidence patterns; no evidence produces a
  warning-severity `C002` finding at the span's line.
- Both rules run inside the same `rules.Run` pass `scan-command` already
  drives, so they appear in `agents-lint scan`'s output and success-line
  count with zero changes to `cmd/agents-lint/cmd/scan.go`.
- Both rule IDs are enable/disable/re-severity-able via `.agents-lint.yaml`
  with zero changes to `internal/config` (it only needs `KnownRuleIDs()` to
  list them).
- `internal/lint.Document` gains a `CodeSpans []CodeSpan` field populated by
  a full-tree walk, so both new rules (and any future rule that needs
  inline-code spans) get every span regardless of markdown nesting.

**Non-Goals:**
- No perfect path-vs-non-path classifier. FR-C001 itself describes this as
  "detected by backtick or inline-code patterns" — a heuristic, not exact
  parsing. Known false-positive/negative edges (below) are accepted, not
  solved.
- No recursive repo tree walk for C002 evidence. Evidence patterns are
  checked at the repo root only (`repoRoot` itself, non-recursive) — this
  keeps the check O(1) file stats (NFR-PERF-01) and avoids `agents-lint`
  quietly becoming a multi-file scanner (BC-SCOPE-01).
- No `--root`/`--repo-root` flag. `repoRoot` is always the process's working
  directory, matching every other cwd-relative convention already in the
  codebase (`defaultAgentsMDPath`, `config.DefaultFileName`). Not in
  `docs/requirements.md`; out of scope to introduce here.
- No expansion of the C002 tool catalog beyond the fixed list above. FR-C002
  names `terraform`/`kubectl`/`npm` as examples; the catalog is a plain Go
  map and trivially extensible later, but growing it isn't required by any
  `FR-*` ID today.
- No auto-fix of broken references (BC-SCOPE-02).

## Decisions

**`Document.CodeSpans` is populated by a document-wide `ast.Walk`, separate
from `populateSections`'s top-level-only traversal.** `populateSections`
intentionally only looks at direct children of the document root because
schema rules only care about top-level headings and the agent blocks they
contain. Code spans can appear anywhere (inside a list item, a blockquote, a
table cell), so collecting them needs `ast.Walk(root, ...)`, which visits
every node regardless of nesting. This is purely additive to `Parse` — a
second pass after `populateSections`, touching no existing field:
```go
type CodeSpan struct {
    Text string
    Line int
}

// on Document:
CodeSpans []CodeSpan
```
Line number comes from the `CodeSpan`'s first `*ast.Text` child segment
(inline nodes don't carry `.Lines()`, only block nodes do), reusing the
existing `lineNumber(src, pos)` helper — the same technique `nodeText`
already uses to walk into a node's text segments.

**A new `lint.CodebaseRule` interface, not a change to `lint.Rule`.**
```go
// CodebaseRule validates a parsed Document against the surrounding
// repository. repoRoot is the directory C001/C002 resolve paths and
// evidence-file lookups against.
type CodebaseRule interface {
    ID() string
    Check(doc *Document, repoRoot string) []Finding
}
```
Changing `lint.Rule.Check`'s signature to add `repoRoot` would force S002-
S005 to accept a parameter they never use, and would touch every existing
rule file and test for no benefit. A second, narrower interface keeps the
blast radius to exactly the two new rules. `registry.go` gains a parallel
`CodebaseAwareRules() []lint.CodebaseRule` alongside the existing
`DefaultRules() []lint.Rule`, following the same "typed slice of the
package's own rule structs" shape already established.

**`repoRoot` is `os.Getwd()`, computed once inside `rules.Run`, not threaded
in from `cmd`.** Every existing cwd-relative convention in this codebase
(`defaultAgentsMDPath`, `config.DefaultFileName`) already assumes the
process's working directory is the repo root; C001/C002 do the same instead
of introducing a new way to talk about "where is the repo." `Run`'s exported
signature (`Run(path string) ([]lint.Finding, error)`) does not need to
change — `os.Getwd()` failing is treated the same as any other
`rules.Run` infrastructure error (returned, not turned into a `Finding`),
mirroring how `lint.Parse`'s read errors are already handled.
*Rule-level* tests (`c001_test.go`/`c002_test.go`) call `Check(doc,
repoRoot)` directly with an explicit `repoRoot` (no `os.Getwd()`/`t.Chdir`
involved), the same pattern `s004_test.go` already uses for calling
`Check(doc)` directly against a fixture-parsed `Document`. Only the
`rules.Run`-level and `cmd`-level integration tests need `t.Chdir` into a
temp directory, mirroring `internal/config`'s `Resolve` tests.

**C001's "looks like a path" heuristic:** a code-span's text is a candidate
if all of:
- it contains no whitespace,
- it does not start with `-` (excludes CLI flags like `--config`),
- it does not contain `://` (excludes URLs),
- it does not contain `*` or `?` (excludes glob patterns like `**/*.go`),
- it does not contain `<` or `>` (excludes template placeholders like
  `<rule_id>.go` — a real pattern found dogfooding this rule against this
  project's own `AGENTS.md`, which uses angle-bracket placeholders
  throughout its examples),
- it does not start with `/` (excludes root-relative strings like REST
  paths — `/api/users` — which are common in AGENTS.md and are not
  filesystem paths relative to a repo root; every real path example in this
  project's own docs, e.g. `scripts/deploy.sh`, `internal/lint/rule.go`, is
  written without a leading slash),
- and either contains `/` (a multi-segment path) **or** matches a bare
  `name.ext`-shaped token whose final dot-segment *starts with a lowercase
  letter* (excludes version-looking tokens like `1.22`/`v1.2.3`, which are
  digit-first after the last dot, and Go-idiom qualified identifiers like
  `fmt.Print`/`os.ReadFile` — also found dogfooding against this project's
  own `AGENTS.md` — which are uppercase-first since exported Go identifiers
  are capitalized, while still catching `go.mod`, `README.md`,
  `.agents-lint.yaml`, whose extensions are conventionally lowercase).

A candidate is resolved with `filepath.Join(repoRoot, candidate)` and
checked with `os.Stat`; a miss is a `C001` finding. This is a plain-Go
heuristic (no new dependency) that directly matches FR-C001's own
"detected by backtick or inline-code patterns matching filesystem paths"
wording, erring toward the same real-world examples the project's own docs
use.

**C002's tool catalog is a fixed `map[string][]string]` of tool name → repo-
root evidence filenames/globs**, checked non-recursively against `repoRoot`:

| tool | evidence (any one present at repo root) |
|---|---|
| `terraform` | `*.tf`, `.terraform.lock.hcl` |
| `npm` | `package.json`, `package-lock.json` |
| `yarn` | `yarn.lock` |
| `pnpm` | `pnpm-lock.yaml` |
| `docker` | `Dockerfile` |
| `go` | `go.mod` |
| `make` | `Makefile` |
| `cargo` | `Cargo.toml` |
| `kubectl` | `k8s/`, `kubernetes/`, `kustomization.yaml` |
| `git` | `.git` |

Only code spans whose text exactly (case-sensitive) matches a catalog key
are checked at all — an arbitrary bare word in inline code (a rule ID like
`S004`, a Go identifier, a YAML key) is silently ignored rather than
misclassified as a tool reference. This bounds false positives to the fixed
list FR-C002 itself names examples from (`terraform`, `kubectl`, `npm`),
rather than trying to recognize "any command-looking word."

**Both rules run inside `rules.Run`, after the schema rules, only when
parsing succeeded** (same precondition S002-S005 already have — no `Document`
exists if S001 fails):
```go
func Run(path string) ([]lint.Finding, error) {
    if finding := CheckFileExists(path); finding != nil {
        return []lint.Finding{*finding}, nil
    }
    doc, err := lint.Parse(path)
    if err != nil {
        return nil, err
    }

    var findings []lint.Finding
    for _, rule := range DefaultRules() {
        findings = append(findings, rule.Check(doc)...)
    }

    repoRoot, err := os.Getwd()
    if err != nil {
        return nil, fmt.Errorf("rules: resolving repo root: %w", err)
    }
    for _, rule := range CodebaseAwareRules() {
        findings = append(findings, rule.Check(doc, repoRoot)...)
    }
    return findings, nil
}
```
`KnownRuleIDs()` grows to append `CodebaseAwareRules()`'s IDs after the
schema rules', so `.agents-lint.yaml` validation (FR-CFG-04) accepts
`C001`/`C002` with no other change to `internal/config`.

## Risks / Trade-offs

- **[Risk]** C001's heuristic can false-positive on tokens that happen to
  look path-like but aren't (e.g. a date written `2024/01/15`, or a
  multi-segment identifier that coincidentally contains `/`).
  → **Mitigation**: accepted as a known heuristic limitation, consistent
  with FR-C001's own "detected by ... patterns" phrasing; not solved with a
  larger exclusion list, since every additional carve-out increases the
  chance of also excluding a real broken reference (the actual bug this
  rule exists to catch).
- **[Risk]** C001's heuristic can false-negative on genuine bare-filename
  references with no extension (e.g. `Makefile`, `LICENSE` written without
  a path prefix) — they contain no `/` and no dot, so they're never
  flagged as candidates either way (never a false positive, but also never
  checked). → **Mitigation**: none added; extending the bare-filename list
  is a bounded, additive follow-up if it turns out to matter in practice,
  not a blocker for shipping FR-C001's core "paths matching filesystem
  paths" behavior.
- **[Risk]** C002's evidence table is necessarily incomplete/opinionated
  (e.g. `kubectl` evidence is looser than the others since Kubernetes
  manifests have no single canonical top-level file). → **Mitigation**:
  FR-C002's status is `proposed`, not a hard contract; the table is a plain
  Go map, trivially extended without touching the rule's control flow, and
  the rule's severity is `warning` (never fails a scan on its own) precisely
  because this signal is inherently softer than C001's.
- **[Trade-off]** `repoRoot` is always `os.Getwd()`, so `agents-lint scan`
  run from anywhere other than the actual repo root will resolve C001/C002
  against the wrong base directory. → **Mitigation**: this matches every
  other cwd-relative assumption already shipped (`defaultAgentsMDPath`,
  `config.DefaultFileName`); introducing a different convention just for
  these two rules would be inconsistent, and no `FR-*` asks for a
  `--repo-root` flag.

## Open Questions

(none)
