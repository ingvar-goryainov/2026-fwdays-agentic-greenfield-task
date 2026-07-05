# agents-lint — demo video (Remotion)

The 1–2 minute demo video for `agents-lint` is itself built agentically: the
whole thing is code. This directory is a [Remotion](https://remotion.dev)
project that renders a 118-second, 1080p walkthrough — product demo plus the
"how it was built" story.

Every fact on screen is grounded in the real repo: the terminal output was
captured from the actual binary (`agents-lint init` / `scan`), the rule roster
comes from `agents-lint docs <id>`, and the counts (28 commits, 76 tests,
17 fixtures, 15 traceable requirements) are pulled from git and `testdata/`.

## Structure

| Path | What |
|------|------|
| `SCRIPT.md` | Scene-by-scene script: on-screen text, visuals, narration |
| `src/theme.ts` | Colors, fonts, and the scene timeline (frames sized to narration) |
| `src/Video.tsx` | Composition — sequences the 9 scenes + soundtrack |
| `src/scenes/S1..S9` | One file per scene |
| `src/components/` | Reusable: animated `Terminal`, `Logo`, `Background`, `Caption` |
| `src/Soundtrack.tsx` | Narration-per-scene + looped music bed |
| `scripts/narrate.mjs` | Generates per-scene narration via Eleven Labs |
| `scripts/make-music.sh` | Generates the ambient music bed with ffmpeg (no API) |

## Build it

```bash
npm install

# 1. Narration (needs an Eleven Labs key; writes public/narration/*.mp3)
ELEVENLABS_API_KEY=sk_... VOICE_ID=EXAVITQu4vr4xnSDxMaL node scripts/narrate.mjs

# 2. Music bed (writes public/music.mp3)
bash scripts/make-music.sh

# 3. Preview in the studio, or render
npm start                 # interactive studio at localhost:3000
npm run render            # -> out/agents-lint-demo.mp4
```

Audio files under `public/` are generated artifacts and are git-ignored — run
the two scripts above to reproduce them. The narration text lives in
`scripts/narrate.mjs` and `SCRIPT.md`; no API keys are stored in the repo.

## Requirements

- Node 18+ and npm
- `ffmpeg` on PATH (for the music bed)
- An Eleven Labs API key for narration (premade voices work on the free tier;
  library voices require a paid plan)
