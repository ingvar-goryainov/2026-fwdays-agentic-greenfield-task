## Why

`agents-lint` currently tells users a rule failed but not what the rule
actually checks, why it matters, or how to fix it beyond the one-line
message in a Finding. `docs/requirements.md` reserves `agents-lint docs
<rule-id>` (FR-CLI-07) as a stretch goal for exactly this: a self-contained
rule reference the user (or their agent) can consult without leaving the
terminal or reading Go source. The rule catalog (S001–S005, C001–C002) is
now stable, so this is unblocked.

## What Changes

- Add structured, per-rule documentation metadata (description, default
  severity, one valid example, one invalid example, fix guidance) attached
  to each of the seven existing rules (S001–S005, C001, C002).
- Add `agents-lint docs <rule-id>` — prints the full specification for a
  single rule ID to stdout.
- Unknown rule IDs produce a clear, actionable error (exit non-zero) rather
  than empty output.
- No changes to `scan`, `init`, findings, or exit-code behavior of any
  existing command.

## Capabilities

### New Capabilities
- `rule-docs-command`: `agents-lint docs <rule-id>` prints a rule's
  description, severity, valid/invalid examples, and fix guidance
  (FR-CLI-07).

### Modified Capabilities
(none — existing rule behavior and Finding output are unchanged)

## Impact

- `internal/lint/rules/*.go`: each rule (S001–S005, C001, C002) gains
  documentation metadata alongside its existing `Check` logic.
- `internal/lint/rules/registry.go`: gains a lookup from rule ID to its
  documentation metadata, mirroring the existing `KnownRuleIDs()` pattern.
- `cmd/agents-lint/cmd/`: new `docs.go` registering the `docs <rule-id>`
  subcommand.
- No dependency changes, no CGO, no network calls.
