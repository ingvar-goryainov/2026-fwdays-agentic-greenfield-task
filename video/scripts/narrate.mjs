#!/usr/bin/env node
// Generate per-scene narration with Eleven Labs.
// Usage: ELEVENLABS_API_KEY=... [VOICE_ID=...] node scripts/narrate.mjs
import { writeFile, mkdir } from "node:fs/promises";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";

const __dirname = dirname(fileURLToPath(import.meta.url));
const OUT = join(__dirname, "..", "public", "narration");

const API_KEY = process.env.ELEVENLABS_API_KEY;
// Default: "Rachel" — clear, neutral, narration-friendly. Override with VOICE_ID.
const VOICE_ID = process.env.VOICE_ID || "21m00Tcm4TlvDq8ikWAM";
const MODEL = process.env.MODEL_ID || "eleven_multilingual_v2";

if (!API_KEY) {
  console.error("ERROR: set ELEVENLABS_API_KEY in the environment.");
  process.exit(1);
}

const blocks = {
  n1: "Your AI agents follow a config file — AGENTS dot M-D. But that file drifts. It points at code that no longer exists, and nobody notices until the agent goes wrong.",
  n2: "AGENTS files are written once and forgotten. They reference files that got renamed, tools that got removed, agents that no longer exist. Existing linters just score them. This one doesn't score. It passes, or it fails.",
  n3: "So here it is. agents-lint init scaffolds a valid file. agents-lint scan runs seven rules against it — schema, and codebase. Green. Seven rules passed.",
  n4: "Now watch. Rename a script the file references — real drift. Re-scan. Rule C-oh-oh-one fails: the path doesn't exist. Non-zero exit — so this fails your C-I, before the agent ever sees a broken config.",
  n5: "Five schema rules, two codebase-aware rules, and SARIF output that drops straight into GitHub code scanning. One Go binary, zero dependencies, no network.",
  n6: "But the point of this project is how it was built — agentically. Spec first: fifteen traceable requirements. Then one OpenSpec change per capability — proposal, design, tasks, apply.",
  n7: "Each rule went through the same loop — write fixtures, implement, run the tests, lint, commit. Not step-by-step prompting: a loop the agent ran to green. Seventy-six tests, and every commit tied to a requirement I-D.",
  n8: "And the maker was never the checker. A separate review pass — CodeRabbit on every pull request, plus a code-graph M-C-P for structural review.",
  n9: "agents-lint. A small tool, taken through a full engineering loop — built agentically.",
};

const settings = {
  stability: 0.45,
  similarity_boost: 0.8,
  style: 0.35,
  use_speaker_boost: true,
};

async function tts(id, text) {
  const res = await fetch(
    `https://api.elevenlabs.io/v1/text-to-speech/${VOICE_ID}`,
    {
      method: "POST",
      headers: {
        "xi-api-key": API_KEY,
        "Content-Type": "application/json",
        Accept: "audio/mpeg",
      },
      body: JSON.stringify({
        text,
        model_id: MODEL,
        voice_settings: settings,
      }),
    }
  );
  if (!res.ok) {
    throw new Error(`${id}: HTTP ${res.status} — ${await res.text()}`);
  }
  const buf = Buffer.from(await res.arrayBuffer());
  await writeFile(join(OUT, `${id}.mp3`), buf);
  console.log(`  ✓ ${id}.mp3  (${(buf.length / 1024).toFixed(0)} KB)`);
}

await mkdir(OUT, { recursive: true });
console.log(`Voice ${VOICE_ID} · model ${MODEL}`);
for (const [id, text] of Object.entries(blocks)) {
  await tts(id, text);
}
console.log("Done. Narration written to public/narration/");
