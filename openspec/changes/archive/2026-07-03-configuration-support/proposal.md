## Why

`scan-command` hard-codes the rule set (`rules.DefaultRules()` plus S001) and
the target path (`./AGENTS.md` or a positional argument) — there is no way
for a repo to opt a rule out, downgrade/upgrade its severity, or point at a
non-default AGENTS.md location without editing the tool's source. Per
`docs/openspec-capability-plan.md`, `configuration-support` is the next
capability after `scan-command`: it adds `.agents-lint.yaml` loading so
teams can tune the rule set to their repo without forking the binary.

## What Changes

- Add a new `internal/config` package that loads and validates
  `.agents-lint.yaml`: AGENTS.md path override, per-rule
  enable/disable, and per-rule severity override (`error` ↔ `warning`).
- Implement FR-CFG-01: `.agents-lint.yaml` in the current directory is the
  default config location.
- Implement FR-CFG-02: config supports `path:` (AGENTS.md location),
  `rules.<ID>.enabled` (bool), and `rules.<ID>.severity` (`error`/`warning`).
- Implement FR-CFG-03: when no config file is present at the default
  location (and `--config` was not passed), `scan` behaves exactly as it
  does today — every rule runs at its built-in severity.
- Implement FR-CFG-04: the config file is schema-validated on load —
  unknown top-level/rule keys, unknown rule IDs, and invalid severity
  values are all rejected with a clear, actionable error before any
  scanning happens.
- Implement FR-CLI-04: add a persistent `--config <path>` flag on the root
  command; when given, that exact path is loaded (and must exist), taking
  priority over the default `.agents-lint.yaml` lookup.
- Wire the resolved config into `agents-lint scan`: path precedence becomes
  positional argument → config `path:` → `./AGENTS.md`; findings are
  filtered/re-severitized per `rules.<ID>.enabled`/`severity` before being
  handed to the text reporter; the success line's rule count reflects any
  disabled rules.
- Add `rules.KnownRuleIDs()` to `internal/lint/rules` so config validation
  has a single source of truth for valid rule IDs (includes S001, which
  runs outside the `Rule` interface).

## Capabilities

### New Capabilities
- `configuration-support`: `.agents-lint.yaml` loading and validation
  (FR-CFG-01..04) plus the `--config` flag (FR-CLI-04), and its effect on
  `agents-lint scan`'s path resolution, rule set, and finding severities.

### Modified Capabilities
(none — `scan-command`'s FR-CLI-01/02 exit-code and default-path contracts
are unchanged; this change only adds a new input, config, that can shift
*which* path is scanned and *which* findings are reported, not the
contracts themselves)

## Impact

- Affected code: new `internal/config` package; `cmd/agents-lint/cmd/root.go`
  (new persistent `--config` flag); `cmd/agents-lint/cmd/scan.go` (config
  resolution, path precedence, finding filtering, rule-count adjustment);
  `internal/lint/rules/registry.go` (new `KnownRuleIDs()` export).
- Dependencies introduced: none — reuses `gopkg.in/yaml.v3`, already a
  module dependency (TC-STACK-04).
- User-facing behavior: `agents-lint scan` gains an optional
  `.agents-lint.yaml` (or `--config <path>`) that can redirect the scanned
  file, silence specific rules, and promote/demote severities. Behavior
  with no config file present is unchanged from today.
