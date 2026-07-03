## Purpose

Schema validation rules (S001–S005) and the rule engine orchestration that
runs them against an AGENTS.md file, producing deterministic `lint.Finding`
results.

## Requirements

### Requirement: Rule engine orchestration
`internal/lint/rules.Run(path string) ([]lint.Finding, error)` SHALL run the file-existence check first and SHALL only parse the file and run the remaining schema rules (S002–S005) if that check produces no finding, always evaluating the five rules in the fixed order S001, S002, S003, S004, S005 so output is deterministic; a parse or filesystem error unrelated to the file simply not existing (e.g., a read permission error) SHALL be returned as a Go `error`, distinct from a `Finding`.

#### Scenario: File does not exist short-circuits the run
- **WHEN** `Run` is called with a path that does not exist
- **THEN** it returns exactly one `Finding` (rule `S001`, `SeverityError`)
  and no other rules are evaluated

#### Scenario: File exists and all rules run
- **WHEN** `Run` is called with a path to a readable AGENTS.md file
- **THEN** it parses the file and evaluates rules S002 through S005 against
  it, returning the union of every rule's findings in `S002, S003, S004,
  S005` order

### Requirement: FR-S001 — AGENTS.md file must exist
The tool SHALL report an error-severity finding when the target AGENTS.md file does not exist at the given path.

#### Scenario: Missing file
- **WHEN** the configured AGENTS.md path does not exist on disk
- **THEN** rule `S001` reports a `SeverityError` finding whose message
  states the expected path and that the file is missing

#### Scenario: Existing file passes
- **WHEN** the configured AGENTS.md path exists and is a readable regular
  file
- **THEN** rule `S001` reports no finding

### Requirement: FR-S002 — Required top-level sections
The tool SHALL report an error-severity finding when the document has no top-level (`##`) heading whose trimmed text case-insensitively equals `Agent` or `Agents`.

#### Scenario: Missing Agent(s) section
- **WHEN** an AGENTS.md file has no `## Agent` or `## Agents` H2 heading
  anywhere in the document
- **THEN** rule `S002` reports a `SeverityError` finding naming the missing
  section and the accepted spellings

#### Scenario: Section present with either accepted spelling
- **WHEN** an AGENTS.md file contains a `## Agents` (or `## Agent`) H2
  heading, in any letter casing
- **THEN** rule `S002` reports no finding

### Requirement: FR-S003 — Agent block completeness
For every `###` heading found inside the required Agent(s) section, the tool SHALL report an error-severity finding if the block is missing a name (a non-empty heading), a role (a plain description paragraph or a `**Role:**` labeled paragraph), or all three of instructions/responsibilities, tools, and context (each detected via a bold-labeled paragraph or list, e.g. `**Tools:**`).

#### Scenario: Agent block missing all optional fields
- **WHEN** an agent block has an H3 name and a role paragraph but none of
  `**Instructions:**`/`**Responsibilities:**`, `**Tools:**`, or
  `**Context:**`
- **THEN** rule `S003` reports a `SeverityError` finding on that block's
  heading line naming which required piece is missing

#### Scenario: Agent block missing a role
- **WHEN** an agent block has an H3 name followed immediately by another
  heading or a labeled field, with no description paragraph and no
  `**Role:**` paragraph
- **THEN** rule `S003` reports a `SeverityError` finding stating the block
  has no role/description

#### Scenario: Complete agent block passes
- **WHEN** an agent block has an H3 name, a role (plain paragraph or
  `**Role:**` label), and at least one of
  `**Instructions:**`/`**Responsibilities:**`, `**Tools:**`, or
  `**Context:**`
- **THEN** rule `S003` reports no finding for that block

### Requirement: FR-S004 — No duplicate agent names
The tool SHALL report an error-severity finding for every agent block whose name (case-insensitive, trimmed) duplicates an earlier agent block's name in the same file.

#### Scenario: Duplicate agent names
- **WHEN** two `###` agent blocks in the Agent(s) section have the same name
  ignoring case and surrounding whitespace (e.g., `### architect` and `###
  Architect`)
- **THEN** rule `S004` reports a `SeverityError` finding on the second
  occurrence naming the duplicate

#### Scenario: All agent names unique
- **WHEN** every agent block in the file has a distinct name
  (case-insensitive)
- **THEN** rule `S004` reports no finding

### Requirement: FR-S005 — Valid frontmatter YAML
When the AGENTS.md file begins with a `---`-delimited frontmatter block, the tool SHALL report an error-severity finding if that block's content is not valid YAML; files with no frontmatter block SHALL NOT be flagged by this rule.

#### Scenario: Malformed frontmatter YAML
- **WHEN** the file starts with a `---`-delimited block whose content fails
  to parse as YAML
- **THEN** rule `S005` reports a `SeverityError` finding including the
  underlying YAML parse error

#### Scenario: Valid frontmatter YAML
- **WHEN** the file starts with a `---`-delimited block whose content parses
  as valid YAML
- **THEN** rule `S005` reports no finding

#### Scenario: No frontmatter present
- **WHEN** the file does not begin with a `---` line
- **THEN** rule `S005` reports no finding
