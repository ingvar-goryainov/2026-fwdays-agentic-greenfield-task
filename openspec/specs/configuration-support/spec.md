## Purpose

`.agents-lint.yaml` configuration loading and validation
(`internal/config`), letting a repo override the scanned AGENTS.md path,
enable/disable individual rules, and override rule severities without
forking the binary. Covers the default config location, the `--config`
flag, schema validation with actionable errors, and the no-config-file
fallback that reproduces `agents-lint scan`'s unconfigured behavior.

## Requirements

### Requirement: FR-CFG-01 — Default config location
`agents-lint` SHALL treat `.agents-lint.yaml` in the current directory as the default configuration file location when no `--config` flag is given.

#### Scenario: Default file is loaded automatically
- **WHEN** `agents-lint scan` is run in a directory containing a `.agents-lint.yaml` file and no `--config` flag is given
- **THEN** that file's settings are applied to the scan

#### Scenario: No default file present
- **WHEN** `agents-lint scan` is run in a directory with no `.agents-lint.yaml` file and no `--config` flag
- **THEN** the scan proceeds with no config applied (see FR-CFG-03)

### Requirement: FR-CFG-02 — Configurable path, rule enablement, and severity
The config file SHALL support overriding the target AGENTS.md path, enabling or disabling individual rules by ID, and overriding an individual rule's severity between `error` and `warning`.

#### Scenario: Config overrides the scanned path
- **WHEN** a config file sets `path: docs/AGENTS.md` and `agents-lint scan` is run with no positional path argument
- **THEN** `docs/AGENTS.md` is validated instead of `./AGENTS.md`

#### Scenario: A positional path argument still wins over config
- **WHEN** a config file sets `path: docs/AGENTS.md` and `agents-lint scan other/AGENTS.md` is run
- **THEN** `other/AGENTS.md` is validated

#### Scenario: Disabling a rule removes its findings
- **WHEN** a config file sets `rules.S004.enabled: false` and the scanned file would otherwise produce an S004 finding
- **THEN** no S004 finding appears in the output, and it does not count toward the error/warning summary

#### Scenario: Disabling a rule reduces the reported rule count
- **WHEN** a config file sets `rules.S004.enabled: false` and the scanned file is otherwise clean
- **THEN** the success line reports one fewer rule passed than it would with no config applied

#### Scenario: Severity override changes a finding's severity
- **WHEN** a config file sets `rules.S004.severity: warning` and the scanned file would otherwise produce an S004 finding with severity `error`
- **THEN** the S004 finding is reported with severity `warning`, and it does not cause a non-zero exit code on its own

#### Scenario: Severity can be raised as well as lowered
- **WHEN** a config file sets a rule's `severity` to `error` for a rule whose built-in severity is `warning`
- **THEN** findings from that rule are reported with severity `error`, and at least one such finding causes a non-zero exit code

### Requirement: FR-CFG-03 — No config file means default behavior
When no config file is loaded (no `.agents-lint.yaml` at the default location and no `--config` flag), `agents-lint scan` SHALL run every rule at its built-in severity and SHALL apply no path override, exactly as if configuration support did not exist.

#### Scenario: Absent config reproduces unconfigured behavior
- **WHEN** `agents-lint scan` is run with no config file present and no `--config` flag
- **THEN** the findings, severities, exit code, and success-line rule count are identical to running the same scan before configuration support existed

### Requirement: FR-CFG-04 — Config schema validation with clear errors
The config file SHALL be validated against its schema when loaded. Unknown top-level keys, unknown keys under a rule entry, references to a rule ID that does not exist, and severity values other than `error` or `warning` SHALL each produce a clear error that identifies the offending file, key, or value, and SHALL prevent the scan from running.

#### Scenario: Unknown top-level key is rejected
- **WHEN** a config file contains a top-level key that is not `path` or `rules`
- **THEN** loading the config fails with an error naming the unknown key and the file path, and no scan is performed

#### Scenario: Unknown rule ID is rejected
- **WHEN** a config file's `rules` section references a rule ID that does not exist (e.g. a typo like `S0O3`)
- **THEN** loading the config fails with an error naming the unknown rule ID and listing the valid rule IDs, and no scan is performed

#### Scenario: Invalid severity value is rejected
- **WHEN** a config file sets a rule's `severity` to a value other than `error` or `warning`
- **THEN** loading the config fails with an error naming the rule and the invalid value, and no scan is performed

#### Scenario: Malformed YAML is rejected
- **WHEN** a config file is not syntactically valid YAML
- **THEN** loading the config fails with an error that includes the file path, and no scan is performed

#### Scenario: Empty config file is valid
- **WHEN** a config file exists at the default location (or the `--config` path) but is empty
- **THEN** the config loads successfully with no path override and no rule overrides, equivalent to no config file being present

### Requirement: FR-CLI-04 — `--config` flag overrides the default location
`agents-lint --config <path>` SHALL load the config file at the given path instead of looking for `.agents-lint.yaml` at the default location. The given path SHALL exist and be a valid config file, or the command SHALL fail with a clear error.

#### Scenario: Explicit config path is used
- **WHEN** `agents-lint --config custom.yaml scan` is run
- **THEN** `custom.yaml` is loaded and applied, regardless of whether a `.agents-lint.yaml` file exists at the default location

#### Scenario: Explicit config path takes priority over the default location
- **WHEN** both `custom.yaml` and a default-location `.agents-lint.yaml` exist, and `agents-lint --config custom.yaml scan` is run
- **THEN** only `custom.yaml`'s settings are applied

#### Scenario: Missing explicit config path fails clearly
- **WHEN** `agents-lint --config missing.yaml scan` is run and `missing.yaml` does not exist
- **THEN** the command fails with an error naming `missing.yaml`, writes that error to stderr, exits with a non-zero status, and does not run the scan
