# agents-lint — Demo Video Script (~120s)

Style: Mixed (motion-graphic framing + real animated terminal centerpiece).
Narration: Eleven Labs. Pace ~150 wpm → ~300 words total.
Format: 1920×1080, 30fps → 3600 frames total.

| # | Scene | Time | Frames | On-screen text | Visuals | Narration |
|---|-------|------|--------|----------------|---------|-----------|
| 1 | Hook | 0–10s | 0–300 | `agents-lint` — *a compiler for your AGENTS.md* | Logo assembles; terminal cursor blink; dark IDE bg | "Your AI agents follow a config file — AGENTS dot md. But that file drifts. It points at code that no longer exists, and nobody notices until the agent goes wrong." |
| 2 | Problem | 10–22s | 300–660 | Stale references · Missing sections · Broken YAML | A mock AGENTS.md with lines decaying/striking through: `scripts/deploy.sh` ✗, duplicate agent, bad frontmatter | "AGENTS files are written once and forgotten. They reference files that got renamed, tools that got removed, agents that no longer exist. Existing linters just score them. This one doesn't score — it passes or fails." |
| 3 | Demo: init + pass | 22–48s | 660–1440 | (terminal) | Animated terminal, typed commands, real output: `init` → `created ./AGENTS.md`; `scan ./AGENTS.md` → `✓ AGENTS.md is valid (7 rules passed)` exit 0 | "So here it is. `agents-lint init` scaffolds a valid file. `agents-lint scan` runs seven rules against it — schema and codebase. Green. Seven rules passed." |
| 4 | Demo: drift + fail | 48–70s | 1440–2100 | (terminal) | `mv scripts/deploy.sh scripts/release.sh`; re-run scan → red `error C001 … referenced path "scripts/deploy.sh" does not exist …` exit 1 | "Now watch. Rename a script the file references — real drift. Re-scan. Rule C001 fails: the path doesn't exist. Non-zero exit — so this fails your CI, before the agent ever sees a broken config." |
| 5 | Rules + SARIF | 70–82s | 2100–2460 | S001–S005 schema · C001–C002 codebase · SARIF for CI | Grid of the 7 rule chips lighting up; SARIF JSON snippet slides in | "Five schema rules, two codebase-aware rules, and SARIF output that drops straight into GitHub code scanning. One Go binary, zero dependencies, no network." |
| 6 | How: spec-first | 82–95s | 2460–2850 | Spec → Design → Tasks → Apply → Archive · 15 traceable IDs | OpenSpec pipeline animates; requirement IDs (FR-*/NFR-*) stream | "But the point of this project is *how* it was built — agentically. Spec first: fifteen traceable requirements, then one OpenSpec change per capability — proposal, design, tasks, apply." |
| 7 | How: loops + verify | 95–108s | 2850–3240 | fixture → implement → test → lint → commit · 76 tests · 17 fixtures | A loop diagram spins; counters roll up: 28 commits, 76 tests | "Each rule went through the same loop — write fixtures, implement, run the tests, lint, commit. Not step-by-step prompting: a loop the agent ran to green. Seventy-six tests, every commit tied to a requirement ID." |
| 8 | How: maker≠checker | 108–116s | 3240–3480 | maker ≠ checker · CodeRabbit · code-review-graph MCP | Two-agent split: builder vs reviewer; review comments appear | "And the maker was never the checker. A separate review pass — CodeRabbit on every PR, plus a code-graph MCP for structural review." |
| 9 | Outro | 116–120s | 3480–3600 | `agents-lint` · Built agentically with Claude Code | Logo re-forms, tagline, green cursor | "agents-lint. A small tool, taken through a full engineering loop — built agentically." |

## Narration blocks (for Eleven Labs, one file per scene)

- **n1**: Your AI agents follow a config file — AGENTS dot md. But that file drifts. It points at code that no longer exists, and nobody notices until the agent goes wrong.
- **n2**: AGENTS files are written once and forgotten. They reference files that got renamed, tools that got removed, agents that no longer exist. Existing linters just score them. This one doesn't score — it passes or fails.
- **n3**: So here it is. agents-lint init scaffolds a valid file. agents-lint scan runs seven rules against it — schema and codebase. Green. Seven rules passed.
- **n4**: Now watch. Rename a script the file references — real drift. Re-scan. Rule C-oh-oh-one fails: the path doesn't exist. Non-zero exit — so this fails your C-I, before the agent ever sees a broken config.
- **n5**: Five schema rules, two codebase-aware rules, and SARIF output that drops straight into GitHub code scanning. One Go binary, zero dependencies, no network.
- **n6**: But the point of this project is how it was built — agentically. Spec first: fifteen traceable requirements, then one OpenSpec change per capability — proposal, design, tasks, apply.
- **n7**: Each rule went through the same loop — write fixtures, implement, run the tests, lint, commit. Not step-by-step prompting: a loop the agent ran to green. Seventy-six tests, every commit tied to a requirement ID.
- **n8**: And the maker was never the checker. A separate review pass — CodeRabbit on every pull request, plus a code-graph MCP for structural review.
- **n9**: agents-lint. A small tool, taken through a full engineering loop — built agentically.
