# CLAUDE.md

## Context

Read `AGENTS.md` at the repo root for full project context, agent roles, and
conventions. That file is the source of truth for architecture, coding standards,
and agent responsibilities. This file adds Claude Code-specific workflow
instructions.

## Project

**agents-lint** — A Go CLI that validates AGENTS.md files against a formal schema
and the surrounding codebase. Pass/fail validator (like a compiler), not a scorer.

- **Stack:** Go 1.22+, Cobra, goldmark, gopkg.in/yaml.v3
- **Spec:** `docs/requirements.md` (traceable IDs: FR-*, NFR-*, TC-*, BC-*)
- **Brief:** `docs/product-brief.md`
- **Capabilities:** `docs/openspec-capability-plan.md`

## Commands

```bash
# Build
go build ./cmd/agents-lint

# Test (run after every change)
go test ./...

# Test with coverage
go test -cover ./internal/...

# Lint
golangci-lint run

# Run the tool
go run ./cmd/agents-lint scan .
go run ./cmd/agents-lint init
```

## Workflow (OpenSpec)

This project uses OpenSpec for structured change management. Each capability
from `docs/openspec-capability-plan.md` becomes one change:

1. `openspec new change <slug>` — create the change
2. Proposal → Design → Tasks (spec before code)
3. Apply (implement tasks in order)
4. Verify: `go test ./...` + `golangci-lint run`

Follow the dependency graph — don't start a capability before its dependencies
are `apply`-complete.

### Build sequence

```
Phase 0: project-foundation
Phase 1: schema-validation-rules → text-reporting
Phase 2: scan-command ← first demoable slice
Phase 3: init-command, configuration-support
Phase 4: codebase-awareness-rules
Phase 5: sarif-output
Phase 6: cli-version-command, rule-docs-command (stretch)
Phase 7: packaging-and-docs
```

## Loop protocol

When implementing a rule or feature:

```
1. Read the requirement (FR-XXXX) from docs/requirements.md
2. Generate fixtures in testdata/<rule_id>/ (valid + invalid)
3. Write the implementation
4. Run: go test ./internal/rules/ -run TestRule<ID>
5. If red → fix → goto 4
6. If green → run full suite: go test ./...
7. If green → run linter: golangci-lint run
8. If clean → commit. If not → fix → goto 7
```

Do NOT skip steps. Do NOT declare "done" without a green test suite.

## Commit conventions

- Prefix: `feat:`, `fix:`, `test:`, `docs:`, `refactor:`, `chore:`
- Reference requirement IDs: `feat: implement rule S002 (FR-S002)`
- Reference OpenSpec capability when relevant: `feat(schema-validation-rules): ...`
- One rule per commit when possible

## Code conventions

- Follow `gofmt` + `golangci-lint` defaults — no custom style
- Rule files: `internal/rules/<rule_id>.go` (lowercase)
- Tests: table-driven, fixtures in `testdata/<rule_id>/`
- Errors: wrap with `fmt.Errorf("rule %s: %w", id, err)` — never panic
- No globals, no CGO, no network calls in `internal/`
- Library-first: no `fmt.Print` in `internal/` — return typed results

## Boundaries — do NOT

- Add dependencies without asking (`go get` requires human approval)
- Modify `docs/requirements.md` or `docs/product-brief.md`
- Implement features not in `docs/requirements.md`
- Refactor across multiple packages in a single change
- Skip tests or mark them as passing without actually running them
- Use CGO or make network calls anywhere in `internal/`
- Add `--fix`, multi-file scanning, or scoring — these are explicitly out of scope

## Definition of Done

A rule/feature is "done" when:

1. ✅ Implementation satisfies the requirement (FR-XXXX)
2. ✅ All fixtures pass (`go test ./...`)
3. ✅ Coverage ≥ 90% for the rule's file
4. ✅ `golangci-lint run` is clean
5. ✅ Error messages are actionable (say what's wrong + how to fix)
6. ✅ Committed with proper prefix and requirement ID
