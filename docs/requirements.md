# PRD — agents-lint

Last updated: 2026-07-04

This document is the **single source of truth** for what the product does and
what constraints govern it. Every requirement has a stable ID. Specs, tests,
PRs, and recordings reference these IDs to keep traceability intact.

Refer to [docs/product-brief.md](product-brief.md) for narrative context.

## ID conventions

| Prefix   | Meaning                    | Example                                      |
| -------- | -------------------------- | -------------------------------------------- |
| `FR-*`   | Functional Requirement     | `FR-SCAN-01` — scan finds AGENTS.md          |
| `NFR-*`  | Non-Functional Requirement | `NFR-PERF-01` — scan completes < 500 ms      |
| `TC-*`   | Technical Constraint       | `TC-STACK-01` — Go 1.22+                     |
| `BC-*`   | Business / UX Constraint   | `BC-OUTPUT-01` — human-readable terminal     |

Status values: `proposed` · `accepted` · `shipped` · `dropped`.

---

## Functional requirements

### CLI interface

| ID          | Description                                                                                           | Status   |
| ----------- | ----------------------------------------------------------------------------------------------------- | -------- |
| FR-CLI-01   | `agents-lint scan [path]` validates the target AGENTS.md (defaults to `./AGENTS.md`)                  | accepted |
| FR-CLI-02   | `agents-lint scan` exits 0 if no errors, exits 1 if any error-severity rule fails                    | accepted |
| FR-CLI-03   | `agents-lint init` generates a minimal AGENTS.md skeleton that passes all schema rules               | accepted |
| FR-CLI-04   | `agents-lint --config <path>` overrides the default config file location                             | proposed |
| FR-CLI-05   | `agents-lint --format [text|sarif]` selects the output format (default: text)                        | proposed |
| FR-CLI-06   | `agents-lint --version` prints the version and exits                                                 | proposed |
| FR-CLI-07   | `agents-lint docs <rule-id>` prints the full rule specification: description, severity, examples of valid and invalid input, and fix guidance | stretch  |

### Schema validation rules (layer 1)

| ID          | Description                                                                                           | Status   |
| ----------- | ----------------------------------------------------------------------------------------------------- | -------- |
| FR-S001     | AGENTS.md file must exist at the expected path                                                        | accepted |
| FR-S002     | Required top-level H2 sections must be present: at minimum `## Agent` (or `## Agents`)               | accepted |
| FR-S003     | Each agent definition block must contain: `name` (H3 heading), `role` (description text), and at least one of `instructions`, `tools`, or `context` | accepted |
| FR-S004     | No duplicate agent names within the same file (case-insensitive comparison)                           | accepted |
| FR-S005     | YAML frontmatter (if present, delimited by `---`) must be valid YAML                                 | accepted |

### Codebase-awareness rules (layer 2)

| ID          | Description                                                                                           | Status   |
| ----------- | ----------------------------------------------------------------------------------------------------- | -------- |
| FR-C001     | File paths referenced in AGENTS.md (detected by backtick or inline-code patterns matching filesystem paths) must resolve to existing files relative to the repo root | accepted |
| FR-C002     | Tool/command references (e.g., `terraform`, `kubectl`, `npm`) should have corresponding evidence in the repo: config files, lockfiles, or scripts that use them | proposed |

### Configuration

| ID          | Description                                                                                           | Status   |
| ----------- | ----------------------------------------------------------------------------------------------------- | -------- |
| FR-CFG-01   | `.agents-lint.yaml` at the repo root is the default config location                                  | accepted |
| FR-CFG-02   | Config allows: specifying the AGENTS.md path, enabling/disabling individual rules by ID, overriding severity (error → warning or warning → error) | accepted |
| FR-CFG-03   | If no config file exists, all rules run with default severities                                      | accepted |
| FR-CFG-04   | Config schema is documented and validated on load; invalid config produces a clear error              | proposed |

### Output and reporting

| ID          | Description                                                                                           | Status   |
| ----------- | ----------------------------------------------------------------------------------------------------- | -------- |
| FR-OUT-01   | Text output format: one line per finding — `severity  rule-id  file:line — message`                  | accepted |
| FR-OUT-02   | Text output ends with a summary: `N error(s), M warning(s)`                                          | accepted |
| FR-OUT-03   | SARIF output conforms to SARIF v2.1.0 schema for integration with GitHub Code Scanning, VS Code, etc.| proposed |
| FR-OUT-04   | When no issues are found, output a single success line: `✓ AGENTS.md is valid (N rules passed)`      | accepted |

---

## Non-functional requirements

| ID          | Description                                                                                           | Status   |
| ----------- | ----------------------------------------------------------------------------------------------------- | -------- |
| NFR-PERF-01 | Full scan of a single AGENTS.md (≤ 500 lines) completes in < 100 ms on commodity hardware            | proposed |
| NFR-PERF-02 | Binary size ≤ 15 MB (stripped, single platform)                                                      | proposed |
| NFR-TEST-01 | Every rule has at least 2 fixture files: one that passes, one that fails                             | accepted |
| NFR-TEST-02 | `go test ./...` achieves ≥ 90% line coverage on the `internal/` packages                             | proposed |
| NFR-DX-01   | `go build ./cmd/agents-lint` completes in < 10 s on a clean checkout                                | proposed |
| NFR-DX-02   | README includes: installation, usage, rule catalog, configuration reference                          | accepted |
| NFR-DIST-01 | Single statically-linked binary per platform (linux/amd64, darwin/arm64 at minimum)                  | proposed |
| NFR-DIST-02 | A Docker image is buildable from a repo-root `Dockerfile` that packages the `agents-lint` binary; running the container mounts a host directory and invokes `scan` against it, exiting with the same codes as the native binary | proposed |

---

## Technical constraints

| ID          | Description                                                                                           | Status   |
| ----------- | ----------------------------------------------------------------------------------------------------- | -------- |
| TC-STACK-01 | Go 1.22 or later; modules enabled                                                                    | accepted |
| TC-STACK-02 | CLI framework: `cobra` (github.com/spf13/cobra)                                                      | accepted |
| TC-STACK-03 | Markdown parsing: `goldmark` (github.com/yuin/goldmark) with AST access                              | accepted |
| TC-STACK-04 | YAML parsing: `gopkg.in/yaml.v3`                                                                     | accepted |
| TC-STACK-05 | No external network calls at runtime; all validation is local and offline                            | accepted |
| TC-STACK-06 | No CGO dependencies; pure Go for cross-compilation                                                   | accepted |
| TC-TEST-01  | Test framework: standard `testing` package + `testify/assert` for readability                        | accepted |
| TC-TEST-02  | Fixtures stored in `testdata/` per Go convention, organized by rule ID                               | accepted |

---

## Business / UX constraints

| ID          | Description                                                                                           | Status   |
| ----------- | ----------------------------------------------------------------------------------------------------- | -------- |
| BC-OUTPUT-01| Terminal output is human-readable without requiring a pager; color is optional and respects `NO_COLOR`| accepted |
| BC-OUTPUT-02| Error messages are actionable: they say what is wrong AND what the correct state would look like      | accepted |
| BC-SCOPE-01 | MVP validates a single AGENTS.md file; multi-file workspace scanning is deferred                     | accepted |
| BC-SCOPE-02 | No auto-fix in MVP; the tool diagnoses, the human (or their agent) fixes                             | accepted |
| BC-LICENSE-01| MIT license                                                                                          | accepted |

---

## Out of scope (MVP)

- Auto-fix / `--fix` flag
- Multi-file workspace scanning (CLAUDE.md, skills/, rules/ directories)
- Plugin system for custom rule authoring
- MCP server mode
- Watch mode / file system events
- Scoring or grading (this is pass/fail, not 0–100)
- Content quality heuristics (vagueness detection, token counting)
- Integration with specific AI agent platforms
- Web UI or dashboard

> **Architecture note:** Core validation logic is structured as a library
> (`internal/`) with typed results (`[]Finding`) to enable future consumers
> (HTTP API, VS Code extension, web dashboard) without refactoring.
