# OpenSpec Capability Plan — agents-lint

> Derived from [docs/requirements.md](requirements.md) and
> [docs/product-brief.md](product-brief.md). Splits the PRD into
> implementation-sized capabilities, each intended to become one OpenSpec
> change (`openspec new change <slug>`). Requirement IDs are preserved for
> traceability — every ID below should map back to a task in its change's
> `tasks.md`.

## How to read this

- **Capability** — a cohesive unit of product behavior, sized to ship as one
  OpenSpec change (proposal → design → tasks → apply).
- **Requirements covered** — the `FR-*` / `NFR-*` / `TC-*` / `BC-*` IDs this
  capability is responsible for.
- **Depends on** — capabilities that must be `apply`-complete first, because
  this capability's design assumes their types, commands, or rule engine
  exist.
- **Cross-cutting rows** at the bottom apply *inside* every capability's
  tasks rather than standing alone (e.g., fixture/coverage requirements).

---

## Capability list

### 1. `project-foundation`
Go module setup, Cobra root command skeleton, build/lint/test scaffolding,
license file. Nothing else can start without this.

| Requirements | Status |
|---|---|
| TC-STACK-01 (Go 1.22+), TC-STACK-02 (cobra), TC-STACK-05 (offline), TC-STACK-06 (no CGO) | accepted |
| TC-TEST-01 (testing + testify), TC-TEST-02 (testdata/ convention) | accepted |
| NFR-DX-01 (build < 10s) | proposed |
| BC-LICENSE-01 (MIT) | accepted |

**Depends on:** none.

---

### 2. `schema-validation-rules`
The layer-1 validation core: goldmark AST parsing of AGENTS.md and the five
schema rules, emitting typed `[]Finding` results.

| Requirements | Status |
|---|---|
| TC-STACK-03 (goldmark), TC-STACK-04 (yaml.v3) | accepted |
| FR-S001 file exists, FR-S002 required H2 sections, FR-S003 agent block shape, FR-S004 no duplicate names, FR-S005 valid frontmatter YAML | accepted |

**Depends on:** `project-foundation`.

---

### 3. `text-reporting`
Renders `[]Finding` to the terminal in the human-readable format, plus the
success line. Deliberately decoupled from any single CLI command so both
`scan` and future consumers can reuse it.

| Requirements | Status |
|---|---|
| FR-OUT-01 (per-line format), FR-OUT-02 (summary line), FR-OUT-04 (success line) | accepted |
| BC-OUTPUT-01 (readable, `NO_COLOR`), BC-OUTPUT-02 (actionable messages) | accepted |

**Depends on:** `schema-validation-rules` (needs the `Finding` type to render).

---

### 4. `scan-command`
Wires the parser, rule engine, and reporter into `agents-lint scan [path]`
with correct exit-code semantics.

| Requirements | Status |
|---|---|
| FR-CLI-01 (`scan [path]`, defaults `./AGENTS.md`), FR-CLI-02 (exit 0/1 on error severity) | accepted |
| NFR-PERF-01 (< 100ms for ≤500 lines) — verify here | proposed |

**Depends on:** `schema-validation-rules`, `text-reporting`.

---

### 5. `init-command`
Generates a minimal AGENTS.md skeleton guaranteed to pass every schema rule.
Must be built against the *finished* rule set so the skeleton never drifts
from what `scan` actually checks.

| Requirements | Status |
|---|---|
| FR-CLI-03 (`init` generates passing skeleton) | accepted |

**Depends on:** `schema-validation-rules`.

---

### 6. `configuration-support`
`.agents-lint.yaml` loading: path override, per-rule enable/disable, severity
override, validation of the config file itself, and the `--config` flag.

| Requirements | Status |
|---|---|
| FR-CFG-01 (default location), FR-CFG-02 (path/enable/severity), FR-CFG-03 (no-config defaults) | accepted |
| FR-CFG-04 (schema-validated config, clear errors), FR-CLI-04 (`--config <path>`) | proposed |

**Depends on:** `scan-command` (needs a stable rule registry to enable/disable/re-severity against).

---

### 7. `codebase-awareness-rules`
The layer-2 rules that check AGENTS.md claims against the actual repo:
file-path resolution and tool/command evidence.

| Requirements | Status |
|---|---|
| FR-C001 (file paths resolve to real files) | accepted |
| FR-C002 (tool/command references have repo evidence) | proposed |

**Depends on:** `schema-validation-rules` (same rule engine/`Finding` contract), `scan-command` (needs to run inside the same pass).

---

### 8. `sarif-output`
Adds the `--format` flag and a SARIF v2.1.0 writer for CI-native tooling
(GitHub Code Scanning, VS Code).

| Requirements | Status |
|---|---|
| FR-OUT-03 (SARIF v2.1.0), FR-CLI-05 (`--format [text\|sarif]`) | proposed |

**Depends on:** `text-reporting` (parallel output backend for the same `[]Finding`), `codebase-awareness-rules` (SARIF should cover both rule layers before shipping).

---

### 9. `cli-version-command`
Trivial, independent, but low-value — sequenced late so it never blocks
higher-value work.

| Requirements | Status |
|---|---|
| FR-CLI-06 (`--version`) | proposed |

**Depends on:** `project-foundation` only (can technically ship anytime).

---

### 10. `rule-docs-command` (stretch)
`agents-lint docs <rule-id>` — needs every rule to already carry structured
metadata (description, severity, valid/invalid examples, fix guidance), so it
must come after the rule catalog is stable.

| Requirements | Status |
|---|---|
| FR-CLI-07 (`docs <rule-id>`) | stretch |

**Depends on:** `schema-validation-rules`, `codebase-awareness-rules`.

---

### 11. `packaging-and-docs`
Cross-compiled release binaries and the README rule catalog/config
reference. Wraps up the MVP.

| Requirements | Status |
|---|---|
| NFR-DIST-01 (static binaries, linux/amd64 + darwin/arm64 min) | proposed |
| NFR-PERF-02 (binary ≤ 15MB) | proposed |
| NFR-DX-02 (README: install, usage, rule catalog, config reference) | accepted |

**Depends on:** everything above that ships in the MVP (2–9, optionally 10).

---

## Cross-cutting (apply inside every capability, not standalone)

| Requirement | How to apply |
|---|---|
| NFR-TEST-01 — 2 fixtures per rule (pass/fail) | Add to every capability that introduces a rule (2, 7). |
| NFR-TEST-02 — ≥90% line coverage on `internal/` | Track continuously; don't defer to a final "testing" change. |
| BC-OUTPUT-01/02 | Already folded into `text-reporting`; re-check when `sarif-output` and `codebase-awareness-rules` add new message strings. |
| BC-SCOPE-01/02 | Guardrails, not tasks — reject any change proposal that adds `--fix` or multi-file scanning; these are explicitly out of scope for MVP. |

---

## Implementation order

```
Phase 0  project-foundation
             │
Phase 1  schema-validation-rules ──► text-reporting
             │                            │
Phase 2      └──────────► scan-command ◄──┘
             │                  │
Phase 3      ▼                  ▼
        init-command    configuration-support
                                │
Phase 4                 codebase-awareness-rules
                                │
Phase 5                    sarif-output
                                │
Phase 6      cli-version-command (anytime after Phase 0)
             rule-docs-command (stretch, after Phase 4)
                                │
Phase 7                 packaging-and-docs
```

**Recommended build sequence (linear):**

1. `project-foundation`
2. `schema-validation-rules`
3. `text-reporting`
4. `scan-command` — first end-to-end usable slice; MVP is demoable here
5. `init-command`
6. `configuration-support`
7. `codebase-awareness-rules`
8. `sarif-output`
9. `cli-version-command` (slot in whenever convenient — no dependents wait on it)
10. `rule-docs-command` (stretch — only if time remains)
11. `packaging-and-docs`

**Why this order:**
- Steps 1–4 produce the smallest usable product (`scan` against schema rules
  with readable output) as early as possible — matches the product brief's
  "pass or fail like a compiler" philosophy.
- `init` (5) is sequenced right after schema rules are locked, since its
  skeleton output is a direct function of those rules — building it earlier
  risks rework if rules change.
- `configuration-support` (6) comes after `scan-command` because enabling
  /disabling/re-severity-ing rules requires a stable rule registry to act on.
- `codebase-awareness-rules` (7) is layer 2 by design in the PRD and reuses
  the same `Finding`/engine contract validated in layer 1 — building it
  first would mean validating an unproven contract.
- `sarif-output` (8) is deferred until both rule layers exist so the SARIF
  writer is tested against the full finding surface, not a partial one.
- `cli-version-command` (9) and `rule-docs-command` (10) are intentionally
  low priority: one is trivial and dependency-free, the other is explicitly
  a stretch goal gated on rule metadata stabilizing.
- `packaging-and-docs` (11) closes the MVP once behavior is frozen, since
  README and cross-compiled binaries should describe the shipped tool, not a
  moving target.

## Next step

Each capability above is sized to become one OpenSpec change. To start the
first one:

```
/opsx:propose project-foundation
```

or describe it to Claude directly and it will run the propose workflow
(`proposal.md` → `design.md` → `tasks.md`) before implementation begins.
