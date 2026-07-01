# Product Brief — agents-lint

> Companion to `docs/requirements.md`. The requirements document is the numbered,
> traceable source of truth; this brief is the business narrative behind it.

## What this is

`agents-lint` is a Go CLI tool that validates AGENTS.md files and agent context
workspaces against a formal schema and the surrounding codebase. It treats agent
configuration as code — enforcing structure, detecting drift from reality, and
catching misconfigurations before they produce vague or broken agent behavior.

It is not a scoring tool or a style checker. It is a structural validator: either
your AGENTS.md conforms to the schema and references real things, or it does not.

## Who it is for

The single actor is a **developer or DevOps engineer** who uses AI coding agents
(Claude Code, Cursor, Kiro, Copilot) and maintains AGENTS.md / CLAUDE.md files
in their repositories. They run `agents-lint` locally or in CI to catch problems
before they affect agent behavior.

## The pain it addresses

Agent context files (AGENTS.md, CLAUDE.md, rules, skills) are written once and
forgotten. They reference files that get renamed, tools that get removed, and
agents that no longer exist. They lack required sections, duplicate agent names,
and contain broken YAML frontmatter. Existing tools (AgentLinter, Agent Lint)
focus on content quality scoring and heuristics — neither enforces a formal
structural schema or validates references against the actual repository.

The result: agents silently degrade because their context is stale or malformed,
and nobody knows until the output is wrong.

## End-to-end usage

1. **Install.** Download a single Go binary or `go install`. No runtime
   dependencies, no Node.js, no API keys.
2. **Run.** `agents-lint scan .` from the repository root. The tool finds
   AGENTS.md (or the configured path), parses it, and runs all enabled rules.
3. **Read results.** Terminal output shows pass/fail per rule with file, line,
   and a clear message. Exit code is non-zero if any error-severity rule fails.
4. **Fix.** Developer corrects the AGENTS.md based on the diagnostic. Re-runs
   until clean.
5. **CI gate.** Add `agents-lint scan .` to CI. PRs that break the schema or
   reference deleted files fail the check.

## Key workflows in prose

- **Validate on save.** Developer edits AGENTS.md, runs `agents-lint scan .`,
  gets instant feedback on structural errors — missing sections, duplicate names,
  invalid frontmatter.
- **Catch drift in CI.** A PR renames `scripts/deploy.sh` to `scripts/release.sh`.
  AGENTS.md still references `scripts/deploy.sh`. CI fails with rule `C001`:
  "File reference `scripts/deploy.sh` does not exist."
- **Bootstrap a new project.** `agents-lint init` generates a minimal valid
  AGENTS.md skeleton that passes all schema rules out of the box.
- **Configure per-project.** A `.agents-lint.yaml` at the repo root lets teams
  disable rules, change severity, or specify which file to validate.

## Differentiators vs existing tools

| Dimension | AgentLinter (seojoonkim) | Agent Lint (samilozturk) | agents-lint (this) |
|-----------|-------------------------|--------------------------|-------------------|
| Approach | Heuristic scoring (0–100) | MCP advisor + scanning | Schema enforcement + codebase validation |
| Stack | Node.js / npx | Node.js monorepo + MCP | Go binary, zero deps |
| Structural validation | Pattern-based | None | Formal schema (required sections, agent grammar) |
| Codebase-awareness | Stale file/date detection | Code drift detection | File/tool reference resolution against repo |
| Output | Terminal + web report | Terminal + MCP responses | Terminal + SARIF (CI-native) |
| Philosophy | "Score and improve" | "Advise and maintain" | "Pass or fail" — like a compiler |

## MVP vs Future boundary

**In the MVP:** schema validation (5 rules), codebase-awareness (2 rules),
`.agents-lint.yaml` configuration, terminal output with clear diagnostics, SARIF
output for CI, `scan` and `init` commands, full test coverage via fixtures.

**Future (deferred):**

- Auto-fix capabilities (`--fix` flag)
- Inter-agent contract validation (input/output matching between agents)
- Plugin system for custom rule authoring
- MCP server mode (agents can query the linter programmatically)
- Watch mode for continuous validation
- Multi-file workspace scanning (CLAUDE.md, skills/, rules/)
- Web UI / dashboard — core validation logic is structured as a library
  (`internal/`) with typed results to enable future consumers (HTTP API,
  VS Code extension, web dashboard) without refactoring
- `agents-lint docs <rule-id>` — inline rule documentation (stretch goal for MVP)

## Operating principles

- **Pass or fail.** This is a validator, not a scorer. Rules are binary: the file
  conforms or it does not. Severity (error/warning) controls exit codes, not
  "how bad" the problem is.
- **Zero configuration to start.** Running `agents-lint scan .` with no config
  file uses sensible defaults. Configuration is opt-in refinement.
- **Fast and local.** No network calls, no LLM, no API keys. The binary runs in
  milliseconds on any repository.
- **Traceable.** Every rule has a stable ID (S001, C001). CI logs, PR comments,
  and documentation reference these IDs.
