## Context

`scan-command` fixed three things that this capability now needs to make
configurable: the scanned path (`cmd/agents-lint/cmd/scan.go`'s
`defaultAgentsMDPath` / positional arg), the rule set
(`rules.DefaultRules()` — S002-S005 — plus S001, checked separately by
`rules.Run`), and each finding's severity (hard-coded per rule, e.g.
`Severity: lint.SeverityError` inside `S002.Check`). `internal/lint.Finding`
already carries `RuleID` and `Severity` as plain fields, and
`gopkg.in/yaml.v3` is already a module dependency (TC-STACK-04), so no new
external dependency is needed.

The capability plan notes this change "depends on scan-command (needs a
stable rule registry to enable/disable/re-severity against)" — that
registry is `rules.DefaultRules()` plus the `RuleS001` constant.

## Goals / Non-Goals

**Goals:**
- FR-CFG-01: `.agents-lint.yaml` in the process's current directory is the
  default config location (mirrors `defaultAgentsMDPath`'s
  cwd-relative convention).
- FR-CFG-02: config supports `path:` (string), `rules.<ID>.enabled`
  (bool), `rules.<ID>.severity` (`"error"` or `"warning"`).
- FR-CFG-03: no config file at the default location, and no `--config`
  flag, reproduces `scan-command`'s existing behavior exactly — no
  filtering, no severity changes, no path override.
- FR-CFG-04: the config file is schema-validated on load: unknown
  top-level keys, unknown per-rule keys, unknown rule IDs, and invalid
  severity strings are all rejected before any scanning happens, with an
  error that names the offending key/value and (for unknown rule IDs) the
  valid ID list.
- FR-CLI-04: a persistent `--config <path>` flag on the root command that,
  when set, is loaded as-is (must exist) and takes priority over the
  default-location lookup.
- `agents-lint scan`'s path precedence becomes: positional argument →
  config `path:` → `./AGENTS.md`.
- Disabled rules' findings are removed from the report and the success
  line's `(N rules passed)` count reflects the reduced rule set.
- Severity overrides are applied to findings before they reach the text
  reporter, so `FR-CLI-02`'s exit-code decision (which only looks at
  `SeverityError`) sees the *overridden* severity.

**Non-Goals:**
- No `--fix`/auto-remediation of config problems — an invalid config file
  is a hard error (BC-SCOPE-02, and matches `scan-command`'s existing
  "unexpected error → stderr + non-zero exit" path).
- No config support for rules that don't exist yet (`codebase-awareness-rules`'
  C001/C002) — `KnownRuleIDs()` only covers the schema rules
  (S001-S005) shipped so far; that function is additive, so a later
  capability can extend it without redesigning config validation.
- No `--format`/SARIF interaction — out of scope for `sarif-output`.
- No merging of multiple config files or config inheritance — exactly one
  config source (explicit `--config` path, or the single default-location
  file) is loaded per run.

## Decisions

**New `internal/config` package, decoupled from both `cmd` and
`internal/lint/rules`.** `config.Config`/`config.RuleConfig` depend only on
`internal/lint` (for the `Severity` type and `Finding` in `Apply`).
Validation needs the set of valid rule IDs, but rather than importing
`internal/lint/rules` from `config` (which would make config aware of the
concrete rule catalog), `Resolve`/`Load` take `validRuleIDs []string` as a
parameter — the caller (`cmd/agents-lint/cmd/scan.go`) supplies
`rules.KnownRuleIDs()`. This keeps `internal/config` reusable by anything
that has *a* rule catalog, not just today's schema rules, and avoids a
config → rules → lint dependency chain that would make `internal/lint/rules`
harder to reason about in isolation.

```go
package config

type RuleConfig struct {
    Enabled  *bool
    Severity *lint.Severity
}

type Config struct {
    Path  string
    Rules map[string]RuleConfig
}

func Resolve(explicitPath string, validRuleIDs []string) (*Config, error)
func Load(path string, validRuleIDs []string) (*Config, error)
func (c *Config) Apply(findings []lint.Finding) []lint.Finding
```

`Enabled`/`Severity` are pointers so "not mentioned in the config" (leave
the rule's default behavior alone) is distinguishable from "explicitly set
to false/a given severity" — a plain `bool`/`Severity` value type can't
represent "unset" without a sentinel.

**`Resolve` owns the FR-CFG-01/03/CLI-04 precedence; `Load` owns FR-CFG-04
validation and is independently testable against arbitrary file paths.**
```go
const DefaultFileName = ".agents-lint.yaml"

func Resolve(explicitPath string, validRuleIDs []string) (*Config, error) {
    path := explicitPath
    if path == "" {
        path = DefaultFileName
        if _, err := os.Stat(path); err != nil {
            return &Config{}, nil // FR-CFG-03: no file, no error, all defaults
        }
    }
    return Load(path, validRuleIDs) // explicit path: must exist (FR-CLI-04)
}
```
An explicit `--config` path that doesn't exist is a hard error (surfaces
through `Load`'s `os.Open` failure) rather than silently falling back to
defaults — a user who typed `--config` almost certainly intended a
specific file, and silently ignoring a typo'd path would hide the mistake
FR-CFG-04 is meant to catch. Alternative considered: treat a missing
`--config` path the same as a missing default-location file (silent
fallback) — rejected because it makes `--config typo.yaml` behave
identically to omitting the flag, with no error to signal the typo.

**Validation: `yaml.Decoder.KnownFields(true)` for structural errors, a
manual pass for semantic errors (unknown rule ID, invalid severity
string).** `Load` decodes into an unexported `rawConfig`/`rawRuleConfig`
pair (plain `string`/`bool`/`*string` fields, no `lint.Severity` coupling)
with `KnownFields(true)` set on the decoder, so a typo like `enalbed:` or
an unrecognized top-level key fails immediately with yaml.v3's own
"field ... not found in type" message. A second pass
(`rawConfig.toConfig(validRuleIDs)`) converts to the public `Config` type
and checks the two things YAML's structural validation can't: that every
key under `rules:` is one of `validRuleIDs`, and that every `severity:`
value is exactly `"error"` or `"warning"`. Both failure messages name the
bad value and the valid alternatives (e.g. `rule "S0O3": unknown rule ID
— valid rule IDs: S001, S002, S003, S004, S005`), satisfying FR-CFG-04's
"clear error" requirement and the project's actionable-error-message
convention. All errors from `Load` are wrapped as `fmt.Errorf("config %s:
%w", path, err)` so the offending file path is always in the message.

**An empty (0-byte) config file is treated as a valid, empty config, not a
parse error.** `yaml.Decoder.Decode` returns `io.EOF` for a file with no
YAML documents, which would otherwise surface as a confusing "config
.agents-lint.yaml: EOF" error for a user who created an empty file to mean
"use defaults." `Load` special-cases `errors.Is(err, io.EOF)` on the
initial decode and returns `&Config{}, nil` in that case.

**`Config.Apply` filters and re-severitizes; it does not need to know
about S001's early-return short-circuit in `rules.Run`.** `rules.Run`
already returns *only* the S001 finding when the file doesn't exist (it
never runs S002-S005 without a parsed doc). Disabling S001 via config is
handled correctly by `Apply` running as a pure post-processing step over
whatever `[]lint.Finding` `Run` returned: if S001 fired and is disabled,
`Apply` drops it (findings end up empty, matching "the user opted out of
the file-existence check"); if S001 didn't fire (file exists), there's
nothing to drop, and `Run` proceeded to the other rules exactly as
before. No change to `rules.Run` or `rules.DefaultRules()` is needed —
config integration is entirely additive at the call site.
```go
func (c *Config) Apply(findings []lint.Finding) []lint.Finding {
    if c == nil || len(c.Rules) == 0 {
        return findings
    }
    out := make([]lint.Finding, 0, len(findings))
    for _, f := range findings {
        rc, ok := c.Rules[f.RuleID]
        if ok && rc.Enabled != nil && !*rc.Enabled {
            continue
        }
        if ok && rc.Severity != nil {
            f.Severity = *rc.Severity
        }
        out = append(out, f)
    }
    return out
}
```

**Path precedence and rule-count adjustment live in `scan.go` as small,
independently testable pure functions**, following the
`exitCodeForFindings` precedent set by `scan-command`'s design (pure
decision logic separated from the `os.Exit`/Cobra glue):
```go
func resolvePath(argPath, cfgPath string) string {
    switch {
    case argPath != "":
        return argPath
    case cfgPath != "":
        return cfgPath
    default:
        return defaultAgentsMDPath
    }
}

func effectiveRuleCount(cfg *config.Config, knownIDs []string) int {
    count := len(knownIDs)
    for _, id := range knownIDs {
        if rc, ok := cfg.Rules[id]; ok && rc.Enabled != nil && !*rc.Enabled {
            count--
        }
    }
    return count
}
```
`runScan`'s signature grows from `(w io.Writer, path string)` to `(w
io.Writer, argPath, configPath string)`, since it can no longer assume
"the caller already resolved the default path" — it needs the *raw*
positional argument (empty string if omitted) to correctly rank it against
config's `path:`. The `RunE` closure passes `args[0]` (or `""`) instead of
pre-defaulting it.

**`--config` is a persistent flag on `rootCmd`, not scoped to `scanCmd`.**
FR-CLI-04's literal syntax (`agents-lint --config <path>`) reads as a
root-level flag, and `PersistentFlags()` makes it automatically available
to `scan` (and any future subcommand) without re-declaring it. `init`
doesn't currently consume config, so the flag is simply unused there today
— declaring it at the root is forward-compatible with a subcommand that
does need it later, at zero cost now.

**`rules.KnownRuleIDs()` is the single source of truth for valid rule
IDs**, added to `internal/lint/rules/registry.go`:
```go
func KnownRuleIDs() []string {
    ids := []string{RuleS001}
    for _, r := range DefaultRules() {
        ids = append(ids, r.ID())
    }
    return ids
}
```
This is additive (no existing signature changes) and keeps the "what rule
IDs exist" answer in the one package that already owns `DefaultRules()`
and `RuleS001`, rather than duplicating the list in `internal/config` or
`cmd`.

## Risks / Trade-offs

- **[Risk]** A config file that disables every rule produces a scan that
  always exits 0 and reports "(0 rules passed)", which could mask a
  genuinely broken AGENTS.md. → **Mitigation**: this is the explicit,
  documented effect of FR-CFG-02's per-rule `enabled: false` — the same
  trust model as any linter's inline-disable comment. No additional
  guardrail is added since FR-CFG-02 doesn't specify one and adding an
  unrequested "can't disable everything" check would be scope creep.
- **[Risk]** `yaml.Decoder.KnownFields(true)` error messages come from the
  yaml.v3 library and aren't phrased consistently with this project's
  other error messages (e.g. they don't always name the exact field in a
  way that matches `rules.<ID>.enabled`'s YAML path). → **Mitigation**:
  wrapping with `fmt.Errorf("config %s: %w", path, err)` at least anchors
  every error to the file path; the semantic checks added in `toConfig`
  (unknown rule ID, bad severity) use this project's own message
  phrasing, since those are the two cases explicitly called out by
  FR-CFG-02/04.
- **[Trade-off]** `resolvePath`/`effectiveRuleCount` duplicate a small
  amount of precedence logic in `cmd/agents-lint/cmd/scan.go` rather than
  pushing it into `internal/config`. Kept in `cmd` because both functions
  need `defaultAgentsMDPath` (a `cmd`-level constant) and `rules.KnownRuleIDs()`
  (which `config` deliberately doesn't import) — moving them into
  `internal/config` would require passing both back in anyway, with no
  reduction in coupling.

## Open Questions

(none)
