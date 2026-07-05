export const FPS = 30;
export const WIDTH = 1920;
export const HEIGHT = 1080;

export const colors = {
  bg: "#0d1117",
  bgSoft: "#161b22",
  panel: "#010409",
  border: "#30363d",
  text: "#e6edf3",
  muted: "#8b949e",
  green: "#3fb950",
  greenSoft: "#2ea043",
  red: "#f85149",
  yellow: "#d29922",
  blue: "#58a6ff",
  purple: "#bc8cff",
  cyan: "#39c5cf",
};

export const fonts = {
  mono: "'JetBrains Mono', 'SF Mono', Menlo, monospace",
  sans: "'Inter', -apple-system, system-ui, sans-serif",
};

// Scene layout in frames (30fps). Durations are sized to each narration clip
// (narration is the master track). Keep in sync with SCRIPT.md + Soundtrack.
export const scenes = {
  hook: { from: 0, durationInFrames: 375 }, // n1 11.7s
  problem: { from: 375, durationInFrames: 495 }, // n2 15.7s
  demoPass: { from: 870, durationInFrames: 420 }, // n3 11.2s
  demoFail: { from: 1290, durationInFrames: 500 }, // n4 15.8s
  rules: { from: 1790, durationInFrames: 350 }, // n5 10.9s
  specFirst: { from: 2140, durationInFrames: 420 }, // n6 13.1s
  loops: { from: 2560, durationInFrames: 495 }, // n7 15.8s
  makerChecker: { from: 3055, durationInFrames: 300 }, // n8 9.3s
  outro: { from: 3355, durationInFrames: 195 }, // n9 5.8s
};

export const TOTAL_FRAMES = 3550;
