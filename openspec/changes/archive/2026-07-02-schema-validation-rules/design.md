## Context

`project-foundation` established `internal/lint.Finding`/`Severity` and an
empty `internal/lint/rules` package, but nothing parses AGENTS.md or produces
a `Finding` yet. This change adds that: a goldmark-based parser and five
rules (FR-S001–FR-S005) implementing the "layer 1" schema described in the
product brief.

The repo's own `AGENTS.md` doubles as the closest thing to a grammar
reference (it's the dogfooding target called out in the brief). Its `##
Agents` section and per-agent blocks (`### <name>`, a description paragraph,
then `**Role:**`, `**Responsibilities:**`, `**Tools:**`, `**Context:**`
bold-labeled paragraphs) are the concrete shape FR-S003 has to recognize, so
this design treats that file as the primary real-world fixture in addition
to the synthetic `testdata/` fixtures.

## Goals / Non-Goals

**Goals:**
- A `Parse(path string) (*Document, error)` that turns an AGENTS.md file into
  a small, rule-agnostic `Document` model (frontmatter text, H2 sections, H3
  agent blocks with detected fields) — goldmark's AST stays inside the parser
  package; no rule imports `goldmark` directly.
- Five `Rule` implementations (S001–S005) that consume `*Document` (or, for
  S001, just a path) and return `[]lint.Finding`, each in its own file under
  `internal/lint/rules/` per the naming convention (`s001.go` … `s005.go`).
- A `Run(path string) ([]lint.Finding, error)` orchestrator in
  `internal/lint/rules` that wires parsing + rules together in a fixed,
  deterministic order, ready for `scan-command` to call directly.
- Fixtures per rule under `testdata/<rule-id>/{valid,invalid}.md` (NFR-TEST-01).

**Non-Goals:**
- No CLI command, no text/SARIF rendering — `Run` returns `[]Finding`, full
  stop. Wiring to `scan`/`init` is `scan-command`/`init-command`.
- No config-driven severity or rule enable/disable — `configuration-support`
  adds that later; this change hardcodes each rule's default severity.
- No codebase-awareness (file/tool reference resolution) — that's
  `codebase-awareness-rules` (FR-C001/C002), layer 2, out of scope here.

## Decisions

**Package split: `internal/lint` stays dependency-free.** `internal/lint`
holds only types with no behavior beyond parsing: `Finding`, `Severity`
(already existed), plus the new `Document`/`Heading`/`AgentBlock`/
`Frontmatter` model and `Parse`, and the `Rule` interface. It imports
nothing from `internal/lint/rules`. All five rule implementations, the
`CheckFileExists` (S001) function, `DefaultRules() []lint.Rule`, and the
`Run(path string) ([]lint.Finding, error)` orchestrator live in
`internal/lint/rules` instead, which imports `internal/lint` for the shared
types. Rationale: `internal/lint/rules` necessarily imports `internal/lint`
for `Document`/`Finding`/`Rule`; if the orchestrator instead lived in
`internal/lint` and called into `internal/lint/rules` to assemble the
default rule set, that would be an import cycle (`lint` → `rules` → `lint`).
Keeping `internal/lint` a one-way dependency (parser + types only, no rule
logic) and putting all behavior — including orchestration — in `rules`
avoids the cycle entirely. `scan-command` will call `rules.Run(path)`
directly. Alternative considered: inject the rule set into `lint.Run` as a
`[]Rule` parameter so `Run` itself never imports `rules`, with `rules`
supplying the slice at the call site — workable, but it splits "the
orchestrator" and "the default rule set" across two packages for no benefit
in this change (nothing outside `rules` needs a custom rule set yet), so the
simpler single-package-for-behavior split was chosen instead.

**Rule interface.** `type Rule interface { ID() string; Check(doc
*lint.Document) []lint.Finding }`, defined in `internal/lint`, one small
concrete implementation per file in `internal/lint/rules`, registered in ID
order by `rules.DefaultRules() []lint.Rule`. S001 is deliberately *not* part
of this interface — it answers "does the file exist at all," which has to
run before there's anything to parse. `rules.Run` calls
`rules.CheckFileExists` standalone first and short-circuits (returns just
that one `Finding`, no error) if it fails, since every other rule is
meaningless without a readable file. Alternative considered: make `Check`
take `(path string, doc *lint.Document)` for every rule so S001 fits the
same interface — rejected because it forces S002–S005 to carry a `path`
parameter they never use, and forces `doc` to be nilable everywhere instead
of only at the one call site that needs it.

**Document model, not raw AST.** The parser walks the goldmark AST once and
produces a plain struct:
```go
type Document struct {
    Frontmatter *Frontmatter // nil if absent
    H2Sections  []Heading    // all top-level "## ..." headings, for S002
    AgentBlocks []AgentBlock // H3 blocks found under an Agent(s) section
}
type Heading struct { Text string; Line int }
type AgentBlock struct {
    Name                          string
    Line                          int
    HasRole                       bool
    HasInstructions, HasTools, HasContext bool
}
```
Rules operate on this struct only. Rationale: keeps goldmark's node-walking
and byte-offset-to-line-number conversion in exactly one place, and makes
rule files (and their tests) trivial to write without importing
`goldmark/ast`. Alternative considered: hand rules the raw `ast.Node` tree
and let each rule walk it — rejected, it would duplicate traversal logic
across five files and couple every rule to goldmark's API.

**Agent block field detection (FR-S003 grammar).** Within an H3 block (from
one `### Name` heading up to the next `###`/`##` or EOF):
- `HasRole` is true if the block contains *either* a plain paragraph (not a
  recognized bold-label paragraph, not a list, not a heading) *or* a
  paragraph starting with the bold label `**Role:**` (case-insensitive).
- `HasInstructions` is true if a paragraph or list starts with a bold label
  matching `Instructions` or `Responsibilities` (case-insensitive) — the
  repo's own `AGENTS.md` uses "Responsibilities" as the instructions field,
  so both spellings are treated as the same concept.
- `HasTools` / `HasContext` follow the same bold-label pattern for `Tools:` /
  `Context:`.
Rationale: this is the loosest rule that still satisfies FR-S003's literal
wording ("role: description text", "at least one of instructions, tools, or
context") and passes against the real `AGENTS.md` without inventing stricter
syntax the requirement doesn't specify. Alternative considered: require an
exact `**Role:**` label for every block — rejected, it would fail an
`init`-generated minimal skeleton that just uses a plain description
paragraph, which FR-S003's wording ("role (description text)") allows.

**Frontmatter handling (FR-S005).** The parser detects frontmatter by raw
byte scanning: if the file starts with a line that is exactly `---`, capture
everything up to the next line that is exactly `---` as the raw frontmatter
text (not fed to goldmark), and parse the *remainder* of the file as
Markdown. Rationale: goldmark has no built-in frontmatter support, and
pulling in `goldmark-meta` would add a dependency beyond TC-STACK-03/04's
`goldmark` + `yaml.v3`. Doing it manually also avoids `---` being
misinterpreted as a thematic break by goldmark's core parser. The parser
only extracts the raw text and line range (`Frontmatter.Raw`,
`Frontmatter.Line`); the S005 rule owns calling `yaml.Unmarshal` and turning
a YAML error into an actionable `Finding` message, keeping the
dependency-specific error formatting inside the rule file that owns it.

**Section matching (FR-S002).** An H2 heading counts as satisfying FR-S002
if its trimmed text case-insensitively equals `Agent` or `Agents`. Agent
blocks (H3) are only collected from *inside* such a section (rationale:
FR-S003/S004 describe "agent definition blocks," implying they live under
the required section, not anywhere in the file — an unrelated `### `
subheading under, say, `## Repository structure` shouldn't be mistaken for
an agent).

**Default severities.** All five schema rules default to `SeverityError`.
Rationale: FR-CLI-02 ties exit code 1 to "any error-severity rule fails," and
the product brief frames schema conformance as pass/fail, not advisory — a
malformed AGENTS.md should fail CI by default. `configuration-support` adds
the ability to downgrade individual rules to `warning` later; nothing here
forecloses that.

**Line numbers.** Goldmark AST nodes expose byte offsets (`Lines()` /
`Segment`), not line numbers. The parser converts the offset of a node's
first byte to a 1-based line number by counting `\n` bytes in the source up
to that offset, once per node that needs a `Finding` line — cheap enough at
AGENTS.md's expected size (NFR-PERF-01: < 100ms for ≤ 500 lines) without
needing a precomputed offset index.

## Risks / Trade-offs

- **[Risk]** The bold-label field detection (`**Role:**`, `**Tools:**`, …) is
  inferred from one example file, not a formal grammar in
  `docs/requirements.md`. A real-world AGENTS.md using a different
  convention (e.g., `## Role` as an H4, or plain `Role:` without bold) could
  be misclassified as missing a field it actually has. → **Mitigation**:
  keep the detection intentionally permissive (paragraph-or-list, several
  label spellings) rather than strict, and cover both the labeled and
  unlabeled "role" forms in fixtures so the boundary is test-documented, not
  just described in prose.
- **[Risk]** Treating `Instructions` and `Responsibilities` as synonyms is a
  judgment call not stated in FR-S003. → **Mitigation**: this is the change
  most likely to need revisiting once real-world AGENTS.md files are tested
  against the tool; the mapping lives in one place (`s003.go`) so it's a
  small, isolated fix if wrong.
- **[Trade-off]** Byte-offset-to-line-number conversion is O(offset) per
  call rather than O(1) with a precomputed line-start index. Acceptable at
  the target file size (NFR-PERF-01); revisit only if a future change lifts
  the size assumption.

## Open Questions

(none — resolved above against the repo's own `AGENTS.md` as the reference
grammar; flag during `apply` if a fixture surfaces a shape this design
doesn't account for)

## Amendments

- **2026-07-02, during apply:** the original text of the "Rule interface"
  decision put `Run`/`DefaultRules` in `internal/lint` while every concrete
  rule lived in `internal/lint/rules` (which must import `internal/lint` for
  `Document`/`Finding`/`Rule`). Having `internal/lint` call back into
  `internal/lint/rules` to assemble the default rule set would have been an
  import cycle. Fixed by moving all behavior (rules, `CheckFileExists`,
  `DefaultRules`, `Run`) into `internal/lint/rules`, leaving `internal/lint`
  as a pure types-and-parser package with no dependency on `rules`. See the
  updated "Package split" and "Rule interface" decisions above.
