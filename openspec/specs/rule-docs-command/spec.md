## Purpose

`docs <rule-id>` subcommand for `agents-lint`, letting users print the full
specification for a known rule (description, default severity, valid/invalid
examples, and fix guidance) without consulting external documentation.

## Requirements

### Requirement: FR-CLI-07 — `docs <rule-id>` prints the full rule specification
`agents-lint docs <rule-id>` SHALL print, for a known rule ID, its
description, default severity, one example of input that passes the rule,
one example of input that fails the rule, and fix guidance explaining how
to resolve a violation.

#### Scenario: Known rule ID prints the full specification
- **WHEN** a user runs `agents-lint docs S002`
- **THEN** the command prints S002's description, its default severity
  (`error`), a valid example, an invalid example, and fix guidance to
  stdout, and exits with code 0

#### Scenario: Every shipped rule has documentation
- **WHEN** a user runs `agents-lint docs <rule-id>` for any rule ID
  returned by the rule registry (`S001`–`S005`, `C001`, `C002`)
- **THEN** the command succeeds and prints non-empty description, severity,
  valid example, invalid example, and fix guidance for that rule

#### Scenario: Codebase-awareness rule documentation
- **WHEN** a user runs `agents-lint docs C001` or `agents-lint docs C002`
- **THEN** the command prints that rule's specification the same way it
  does for a schema rule, since `docs` does not distinguish rule layers

### Requirement: FR-CLI-07 — Unknown rule ID fails clearly
`agents-lint docs <rule-id>` SHALL exit with a non-zero code and print an
actionable error to stderr when `<rule-id>` is not one of the rule IDs the
tool knows about, per `BC-OUTPUT-02`.

#### Scenario: Unrecognized rule ID
- **WHEN** a user runs `agents-lint docs S999`
- **THEN** the command prints an error to stderr stating that `S999` is not
  a known rule ID and exits with a non-zero code, without printing partial
  documentation

#### Scenario: Missing rule ID argument
- **WHEN** a user runs `agents-lint docs` with no arguments
- **THEN** the command reports a usage error indicating a rule ID is
  required and exits with a non-zero code, without attempting a lookup

#### Scenario: Rule ID is case-sensitive
- **WHEN** a user runs `agents-lint docs s002` (lowercase)
- **THEN** the command treats it as an unrecognized rule ID, since all
  known rule IDs are uppercase (`S002`, not `s002`), and exits non-zero
