## 1. `internal/config` package — types and validation (FR-CFG-02, FR-CFG-04)

- [x] 1.1 Create `internal/config/config.go` with the public `Config` and
      `RuleConfig` types (`Path string`, `Rules map[string]RuleConfig`;
      `Enabled *bool`, `Severity *lint.Severity`), plus unexported
      `rawConfig`/`rawRuleConfig` decode targets (`Path string`,
      `Rules map[string]rawRuleConfig`; `Enabled *bool`, `Severity *string`)
- [x] 1.2 Implement `Load(path string, validRuleIDs []string) (*Config,
      error)`: open the file, decode with
      `yaml.NewDecoder(f); dec.KnownFields(true)`, treat a decode error that
      `errors.Is(err, io.EOF)` as a valid empty config (`&Config{}, nil`),
      wrap all other decode errors as `fmt.Errorf("config %s: %w", path,
      err)`
- [x] 1.3 Implement `rawConfig.toConfig(validRuleIDs []string) (*Config,
      error)`: reject rule IDs not present in `validRuleIDs` (error names
      the bad ID and lists valid IDs); reject `severity` values other than
      `"error"`/`"warning"` (error names the rule ID and the bad value);
      copy `Enabled` through unchanged; wire `toConfig`'s errors into
      `Load`'s `fmt.Errorf("config %s: %w", ...)` wrapping
- [x] 1.4 Implement `Resolve(explicitPath string, validRuleIDs []string)
      (*Config, error)`: `DefaultFileName = ".agents-lint.yaml"`; if
      `explicitPath != ""`, call `Load(explicitPath, validRuleIDs)`
      directly (must exist); else `os.Stat(DefaultFileName)` and return
      `&Config{}, nil` if absent, otherwise `Load(DefaultFileName,
      validRuleIDs)`
- [x] 1.5 Implement `(c *Config) Apply(findings []lint.Finding)
      []lint.Finding`: nil/empty-`Rules` short-circuit returns `findings`
      unchanged; otherwise drop findings whose rule has `Enabled != nil &&
      !*Enabled`, and overwrite `Severity` for findings whose rule has a
      non-nil `Severity` override

## 2. `internal/config` tests

- [x] 2.1 Fixtures under `internal/config/testdata/`: valid config with
      `path` only, valid config with rule `enabled`/`severity` overrides,
      empty file, malformed YAML, unknown top-level key, unknown rule ID,
      invalid severity value
- [x] 2.2 Table-driven tests for `Load` covering every fixture in 2.1:
      assert the parsed `Config` for valid cases, and assert the error
      message contains the offending file path/key/value for each invalid
      case
- [x] 2.3 Tests for `Resolve`: explicit path used when given (including
      the "explicit path missing → error" case); default file loaded when
      present at `DefaultFileName`; `&Config{}, nil` returned when absent
      and no explicit path given (use `t.Chdir`/a temp working directory
      so the test doesn't depend on the repo's own working directory)
- [x] 2.4 Table-driven tests for `Apply`: no config → unchanged findings;
      disabled rule → its findings removed, others untouched; severity
      override → matching findings' `Severity` changed, others untouched;
      a finding for a rule not mentioned in `Rules` passes through
      unchanged
- [x] 2.5 Run `go test -cover ./internal/config/...` — confirm ≥90% line
      coverage

## 3. Wire config into `rules` and `cmd` (FR-CFG-01, FR-CFG-03, FR-CLI-04)

- [x] 3.1 Add `KnownRuleIDs() []string` to
      `internal/lint/rules/registry.go`: returns `RuleS001` followed by
      `ID()` of each rule in `DefaultRules()`, in that order
- [x] 3.2 Test `KnownRuleIDs()` returns exactly
      `["S001", "S002", "S003", "S004", "S005"]`
- [x] 3.3 Add a persistent `--config` string flag to `rootCmd` in
      `cmd/agents-lint/cmd/root.go` (default `""`), bound to a package
      variable `configFlag` readable by `scan.go`
- [x] 3.4 In `cmd/agents-lint/cmd/scan.go`, add `resolvePath(argPath,
      cfgPath string) string`: returns `argPath` if non-empty, else
      `cfgPath` if non-empty, else `defaultAgentsMDPath`
- [x] 3.5 In `cmd/agents-lint/cmd/scan.go`, add `effectiveRuleCount(cfg
      *config.Config, knownIDs []string) int`: starts at
      `len(knownIDs)`, decrements once per ID in `knownIDs` whose config
      entry has `Enabled != nil && !*Enabled`
- [x] 3.6 Change `runScan`'s signature to `(w io.Writer, argPath,
      configPath string) (exitCode int, err error)`: call
      `config.Resolve(configPath, rules.KnownRuleIDs())`, propagate a
      non-nil error as `(1, err)`; resolve the scan path via
      `resolvePath`; after `rules.Run`, apply `cfg.Apply(findings)`
      before computing the exit code and before calling
      `reporter.WriteText`; pass `effectiveRuleCount(cfg,
      rules.KnownRuleIDs())` as `WriteText`'s `ruleCount` argument
- [x] 3.7 Update `scanCmd`'s `RunE` to pass `args[0]` (or `""` when no
      positional argument) and `configFlag` into `runScan`, instead of
      pre-resolving `path` to `defaultAgentsMDPath` before the call

## 4. `cmd` wiring tests

- [x] 4.1 Fixtures under `cmd/agents-lint/cmd/testdata/config/` (or reuse
      `internal/config/testdata/` fixtures via a relative path): a config
      that disables a rule, a config that overrides a rule's severity to
      `warning`, a config that overrides a rule's severity to `error`, a
      config that sets `path:` to an alternate AGENTS.md fixture
- [x] 4.2 Table-driven tests for `resolvePath`: arg wins over config path,
      config path wins over default, default used when both empty
- [x] 4.3 Table-driven tests for `effectiveRuleCount`: no config → full
      count; one disabled rule → count minus one; disabling a rule twice
      via re-assignment still only counts once (map semantics)
- [x] 4.4 `runScan` tests: config disabling a rule removes its finding
      from stdout and adjusts the success-line count; config overriding a
      rule's severity to `warning` changes the exit code from 1 to 0 for a
      file that only trips that rule; config overriding severity to
      `error` changes the exit code from 0 to 1; config `path:` is used
      when no positional argument is given; a positional argument still
      overrides config `path:`; an invalid `--config` file causes
      `runScan` to return `(1, non-nil error)` without calling
      `reporter.WriteText`
- [x] 4.5 Manual check: run `go run ./cmd/agents-lint scan` with a
      `.agents-lint.yaml` in the repo root that disables a currently
      non-triggering rule (e.g. `rules.S004.enabled: false`) — confirm the
      success line's rule count drops by one, then remove the file
- [x] 4.6 Manual check: `go run ./cmd/agents-lint --config
      does-not-exist.yaml scan` — confirm a clear error on stderr naming
      the missing file and a non-zero exit code

## 5. Full verification

- [x] 5.1 Run `go test ./...` — confirm green
- [x] 5.2 Run `go test -cover ./internal/config/... ./cmd/...` — confirm
      `internal/config` and the new `cmd` helper functions
      (`resolvePath`, `effectiveRuleCount`) are at or above 90% line
      coverage
- [x] 5.3 Run `golangci-lint run` — confirm clean
- [x] 5.4 Confirm no new entries were added to `go.mod` (yaml.v3 is
      already a dependency)
- [x] 5.5 Confirm `internal/lint/rules` and `internal/reporter` behavior
      is otherwise unchanged — only the additive `KnownRuleIDs()` export
      was added to `registry.go`
