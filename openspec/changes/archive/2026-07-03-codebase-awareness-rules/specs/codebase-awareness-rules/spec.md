## ADDED Requirements

### Requirement: FR-C001 — File path references resolve to real files
File paths referenced in AGENTS.md, detected as backtick/inline-code spans whose text looks like a filesystem path, SHALL resolve to an existing file or directory relative to the repo root. A reference that does not resolve SHALL produce an error-severity `C001` finding at the line the reference appears on.

#### Scenario: Existing file reference passes
- **WHEN** AGENTS.md contains an inline-code reference to a path that exists relative to the repo root (e.g. `` `internal/lint/rule.go` ``)
- **THEN** no `C001` finding is produced for that reference

#### Scenario: Missing file reference fails
- **WHEN** AGENTS.md contains an inline-code reference to a path that does not exist relative to the repo root (e.g. a script that was renamed or deleted)
- **THEN** a `C001` finding is produced with severity `error`, the file set to the AGENTS.md path, and the line set to where the reference appears

#### Scenario: Multiple broken references each get their own finding
- **WHEN** AGENTS.md contains two different inline-code path references that both fail to resolve
- **THEN** two separate `C001` findings are produced, one per reference, each at its own line

#### Scenario: Non-path inline code is ignored
- **WHEN** AGENTS.md contains inline-code spans that do not look like filesystem paths (e.g. a CLI flag like `` `--config` ``, a URL like `` `https://example.com` ``, a glob like `` `**/*.go` ``, or a root-relative string like `` `/api/users` ``)
- **THEN** no `C001` finding is produced for those spans, regardless of whether a file happens to exist at that location

#### Scenario: C001 only runs when AGENTS.md parses successfully
- **WHEN** the target AGENTS.md file does not exist (FR-S001 fails)
- **THEN** `C001` does not run and produces no findings

### Requirement: FR-C002 — Tool/command references have repo evidence
Tool or command references in AGENTS.md, detected as backtick/inline-code spans whose text exactly matches a name in the tool's recognized catalog (including `terraform`, `kubectl`, and `npm`), SHALL have corresponding evidence in the repo: a config file, lockfile, or other marker file associated with that tool, present at the repo root. A referenced tool with no such evidence SHALL produce a warning-severity `C002` finding at the line the reference appears on.

#### Scenario: Tool reference with matching evidence passes
- **WHEN** AGENTS.md references a recognized tool by name (e.g. `` `go` ``) and the repo root contains that tool's evidence file (e.g. `go.mod`)
- **THEN** no `C002` finding is produced for that reference

#### Scenario: Tool reference with no evidence produces a warning
- **WHEN** AGENTS.md references a recognized tool by name (e.g. `` `terraform` ``) and the repo root contains none of that tool's evidence files
- **THEN** a `C002` finding is produced with severity `warning`, the file set to the AGENTS.md path, and the line set to where the reference appears

#### Scenario: A C002 finding does not fail the scan by itself
- **WHEN** the only findings produced by a scan are `C002` findings
- **THEN** the scan's exit code is `0` (FR-CLI-02 only fails on error-severity findings)

#### Scenario: Unrecognized bare words are ignored
- **WHEN** AGENTS.md contains an inline-code span whose text is not in the recognized tool catalog (e.g. a rule ID like `` `S004` `` or an arbitrary identifier)
- **THEN** no `C002` finding is produced for that span

#### Scenario: C002 only runs when AGENTS.md parses successfully
- **WHEN** the target AGENTS.md file does not exist (FR-S001 fails)
- **THEN** `C002` does not run and produces no findings

### Requirement: Codebase-awareness rules are configurable
`C001` and `C002` SHALL be included in the rule registry's known rule IDs, so `.agents-lint.yaml` can enable, disable, or override their severity using the same mechanism as the schema rules (FR-CFG-02, FR-CFG-04), with no changes to the config package itself.

#### Scenario: C001 or C002 can be disabled via config
- **WHEN** a config file sets `rules.C001.enabled: false` and the scanned AGENTS.md would otherwise produce a `C001` finding
- **THEN** no `C001` finding appears in the output, and it does not count toward the error/warning summary

#### Scenario: C002's severity can be raised
- **WHEN** a config file sets `rules.C002.severity: error` and the scanned AGENTS.md would otherwise produce a `C002` finding
- **THEN** the finding is reported with severity `error` and causes a non-zero exit code

#### Scenario: Referencing C001 or C002 in config does not fail validation
- **WHEN** a config file's `rules` section references `C001` or `C002`
- **THEN** the config loads successfully, since both IDs are part of `rules.KnownRuleIDs()`
