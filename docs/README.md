# agents-lint

`agents-lint` validates `AGENTS.md` files against a formal schema and checks
their claims against the surrounding repository. It is a pass/fail
validator — like a compiler — not a scoring tool.

## Installation

**Download a release binary.** Cross-compiled, statically-linked binaries
(no runtime dependencies) are built with:

```bash
make release
```

This produces `dist/agents-lint-linux-amd64` and
`dist/agents-lint-darwin-arm64`. Copy the one matching your platform
somewhere on your `PATH`, e.g.:

```bash
cp dist/agents-lint-darwin-arm64 /usr/local/bin/agents-lint
```

**Or install with Go** (requires Go 1.22+):

```bash
go install github.com/ingvar-goryainov/agents-lint/cmd/agents-lint@latest
```

## Usage

### `scan` — validate an AGENTS.md file

```bash
agents-lint scan                # scans ./AGENTS.md (the default path)
agents-lint scan path/to/AGENTS.md
```

Exits `0` unless at least one finding has severity `error` (in which case
it exits `1`); warning-severity findings alone do not fail the scan.

Select the output format with `--format`:

```bash
agents-lint scan --format sarif   # SARIF v2.1.0, for CI/code-scanning tools
agents-lint scan --format text    # human-readable (default)
```

### `init` — bootstrap a new AGENTS.md

```bash
agents-lint init                # writes ./AGENTS.md
agents-lint init --force        # overwrite an existing file
```

Generates a minimal skeleton guaranteed to pass every schema rule.

### `docs` — look up a rule's full specification

```bash
agents-lint docs S002
```

Prints the rule's description, default severity, a passing example, a
failing example, and fix guidance.

### `--version`

```bash
agents-lint --version
```

## Rule catalog

| ID | Layer | Default severity | Checks |
|----|-------|-------------------|--------|
| S001 | schema | error | The target `AGENTS.md` file exists at the resolved path. |
| S002 | schema | error | The document has a top-level `## Agent`/`## Agents` section (case-insensitive). |
| S003 | schema | error | Every agent block (a `###` heading under Agent(s)) has a non-empty name, a role/description, and at least one of instructions, tools, or context. |
| S004 | schema | error | No two agent blocks share the same name (case-insensitive, whitespace-trimmed). |
| S005 | schema | error | If YAML frontmatter is present, it parses as valid YAML. Files without frontmatter are not flagged. |
| C001 | codebase-awareness | error | Every inline-code span that looks like a filesystem path resolves to a real file or directory relative to the repo root. |
| C002 | codebase-awareness | warning | Every inline-code span matching a recognized tool/command (terraform, npm, yarn, pnpm, docker, go, make, cargo, kubectl, git) has corresponding evidence at the repo root (a config file, lockfile, or marker file). |

Run `agents-lint docs <rule-id>` for the full description, examples, and
fix guidance for any rule above.

## Configuration reference

`agents-lint` looks for `.agents-lint.yaml` in the current directory by
default. Override the location with `--config <path>`:

```bash
agents-lint scan --config custom.yaml
```

With no config file present, every rule runs at its built-in severity and
the default path (`./AGENTS.md`) is used.

### Supported keys

| Key | Type | Description |
|-----|------|-------------|
| `path` | string | Overrides the default `./AGENTS.md` path. A positional argument to `scan` still wins over this. |
| `rules.<ID>.enabled` | bool | Disables a rule (`false`) so its findings are omitted and it no longer counts toward the reported rule total. |
| `rules.<ID>.severity` | `error` \| `warning` | Overrides a rule's default severity. |

Unknown top-level keys, unknown rule IDs, and invalid severity values are
rejected with an error naming the offending key/value — the scan does not
run until the config is valid.

### Example

```yaml
path: docs/AGENTS.md

rules:
  S004:
    enabled: false
  C002:
    severity: error
```
