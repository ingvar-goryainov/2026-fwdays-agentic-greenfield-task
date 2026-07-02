# AGENTS.md — agents-lint

> This file defines the context, roles, and rules for AI coding agents working
> on this repository. It is both the project's context engineering artifact AND
> a dogfooding target — `agents-lint` validates files like this one.

## Project overview

`agents-lint` is a Go CLI that validates AGENTS.md files against a formal schema
and the surrounding codebase. It enforces structure, detects reference drift, and
reports findings as pass/fail diagnostics.

- **Stack:** Go 1.22+, Cobra (CLI), goldmark (Markdown AST), gopkg.in/yaml.v3
- **Architecture:** Library-first — core logic in `internal/`, CLI is a thin consumer
- **Philosophy:** Pass/fail validator (like a compiler), not a scorer
- **Spec:** See `docs/requirements.md` for traceable requirements with stable IDs
- **Implementation method:** OpenSpec — each capability goes through propose → design → tasks → apply. See `docs/openspec-capability-plan.md` for the full dependency graph and build sequence.

## Repository structure

```
agents-lint/
├── cmd/agents-lint/        # CLI entrypoint (Cobra commands)
│   └── main.go
├── internal/
│   ├── config/             # .agents-lint.yaml loader & validation
│   ├── parser/             # AGENTS.md → AST (goldmark + frontmatter)
│   ├── rules/              # Rule implementations (one file per rule)
│   │   ├── registry.go    # Rule registry & runner
│   │   ├── s001.go        # FR-S001: file exists
│   │   ├── s002.go        # FR-S002: required sections
│   │   ├── s003.go        # FR-S003: agent block completeness
│   │   ├── s004.go        # FR-S004: no duplicate names
│   │   ├── s005.go        # FR-S005: valid frontmatter YAML
│   │   ├── c001.go        # FR-C001: file references resolve
│   │   └── c002.go        # FR-C002: tool references have evidence
│   ├── reporter/           # Output formatters (text, SARIF)
│   └── types/              # Shared types (Finding, Severity, Config)
├── testdata/               # Fixtures organized by rule ID
│   ├── s001/
│   │   ├── valid/          # AGENTS.md files that PASS this rule
│   │   └── invalid/        # AGENTS.md files that FAIL this rule
│   ├── s002/
│   └── ...
├── docs/
│   ├── requirements.md     # PRD with stable IDs
│   ├── product-brief.md    # Business narrative
│   └── rules.md            # Rule catalog (generated or maintained)
├── .agents-lint.yaml       # Dogfooding: our own config
├── go.mod
├── go.sum
└── README.md
```

## Agents

### architect

Senior Go architect who makes structural decisions.

**Role:** Design packages, interfaces, and data flow. Review agent-generated code
for architectural coherence. Never implements features directly.

**Responsibilities:**
- Define interfaces in `internal/types/` before implementation begins
- Review PRs for separation of concerns (rules don't import reporter, etc.)
- Ensure library-first architecture (no `fmt.Print` in `internal/`)

**Context:** `docs/requirements.md`, `docs/product-brief.md`, this file.

---

### implementer

Go developer who writes feature code in loops.

**Role:** Implement rules, parsers, config loading, and CLI commands. Works in
automated loops: write code → run tests → fix failures → repeat until green.

**Responsibilities:**
- Implement one rule at a time, following the interface defined by `architect`
- Each rule lives in its own file: `internal/rules/<rule_id>.go`
- Every rule must satisfy: `func (r *RuleXXX) Check(doc *ParsedDoc, cfg *Config) []Finding`
- Never skip tests — implementation is only "done" when fixtures pass

**Tools:** `go build`, `go test`, `go vet`, `golangci-lint`

**Context:** `internal/types/`, `testdata/<rule_id>/`, `docs/requirements.md`

---

### tester

Adversarial QA agent that generates test fixtures and reviews coverage.

**Role:** Generate valid and invalid AGENTS.md fixture files for each rule.
Challenge the implementer's code by writing edge cases. Never implements
production code — only test files and test assertions.

**Responsibilities:**
- For each rule, create at minimum:
  - 2 valid fixtures (normal case + edge case that still passes)
  - 2 invalid fixtures (obvious violation + subtle violation)
- Write table-driven Go tests in `internal/rules/<rule_id>_test.go`
- Run `go test -cover` and flag any rule below 90% coverage

**Tools:** `go test`, `go test -cover`, `go test -race`

**Context:** `docs/requirements.md` (rule definitions), `internal/types/`

---

### reviewer

Code reviewer who validates the implementer's output.

**Role:** Maker ≠ checker separation. Reviews code written by `implementer` for
correctness, style, and requirement conformance. Does not write production code.

**Responsibilities:**
- Verify each rule implementation matches the requirement in `docs/requirements.md`
- Check that error messages satisfy `BC-OUTPUT-02` (actionable: what's wrong + how to fix)
- Verify no network calls, no CGO, no global state in `internal/`
- Confirm test fixtures actually test what the rule claims to validate

**Context:** `docs/requirements.md`, `internal/rules/`, `testdata/`

---

## Rules for all agents

### Code conventions

- **Go style:** Follow `gofmt` + `golangci-lint` defaults. No custom style.
- **Naming:** Rule files are `<rule_id>.go` (lowercase). Types are PascalCase.
- **Errors:** Wrap with `fmt.Errorf("rule %s: %w", id, err)`. Never panic.
- **No globals:** No package-level mutable state. Pass dependencies explicitly.
- **No network:** Nothing in `internal/` makes HTTP calls. Ever.
- **No CGO:** Pure Go only. This must cross-compile cleanly.

### Testing conventions

- **Table-driven tests** for rules: each test case is `{name, input, wantFindings}`.
- **Fixtures live in `testdata/`** — use `os.ReadFile` with `testdata/<rule>/<case>.md`.
- **Golden files** for reporter output: `testdata/reporter/golden/<case>.txt`.
- **No mocks for filesystem** — use `testdata/` with real files and `os.DirFS`.

### Commit conventions

- Prefix: `feat:`, `fix:`, `test:`, `docs:`, `refactor:`, `chore:`
- Reference requirement IDs in commits: `feat: implement rule S002 (FR-S002)`
- One rule per PR when possible

### Loop protocol (for implementer + tester)

```
1. Read the requirement (FR-XXXX) from docs/requirements.md
2. Tester: generate fixtures in testdata/<rule_id>/
3. Implementer: write rule in internal/rules/<rule_id>.go
4. Run: go test ./internal/rules/ -run TestRule<ID>
5. If red → fix → goto 4
6. If green → run full suite: go test ./...
7. If green → run linter: golangci-lint run
8. If clean → commit. If not → fix → goto 7
```

### OpenSpec change protocol

Each capability from `docs/openspec-capability-plan.md` is implemented as one
OpenSpec change:

1. `openspec new change <slug>` — create the change directory
2. Write `proposal.md` — what and why (references FR-* IDs)
3. Write `design.md` — how (interfaces, data flow, packages touched)
4. Write `tasks.md` — ordered implementation steps
5. Apply — implement tasks following the loop protocol above
6. Verify — `go test ./...` + `golangci-lint run` + reviewer pass

**Dependency rule:** Do not start a capability's proposal until all its
dependencies (listed in the capability plan) are `apply`-complete. This
prevents designing against interfaces that don't exist yet.

### What agents must NOT do

- Do not modify `docs/requirements.md` without human approval
- Do not add dependencies without human approval (`go get` requires sign-off)
- Do not refactor across multiple packages in a single change
- Do not implement features not in `docs/requirements.md`
- Do not generate code that "seems to work" without running `go test`

## Context loading strategy

| Context type | When to load | Files |
|-------------|-------------|-------|
| **Static (always)** | Every session | This file, `docs/requirements.md` |
| **Dynamic (on demand)** | When working on a specific rule | `testdata/<rule_id>/`, `internal/types/types.go` |
| **Reference (as needed)** | When designing interfaces | `go.mod` (dependency list), existing rule implementations |

## Definition of Done

A feature/rule is "done" when:

1. ✅ Implementation satisfies the requirement (FR-XXXX)
2. ✅ All fixtures pass (`go test ./...`)
3. ✅ Coverage ≥ 90% for the rule's file
4. ✅ `golangci-lint run` is clean
5. ✅ Error messages are actionable (BC-OUTPUT-02)
6. ✅ Reviewer agent has approved
7. ✅ Committed with proper prefix and requirement ID reference
