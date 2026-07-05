import { interpolate, spring } from "remotion";

// Fade + rise-in helper for a block appearing at `start` over `dur` frames.
export const appear = (
  frame: number,
  start: number,
  dur = 15
): { opacity: number; translateY: number } => {
  const opacity = interpolate(frame, [start, start + dur], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });
  const translateY = interpolate(frame, [start, start + dur], [18, 0], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });
  return { opacity, translateY };
};

// Spring pop for emphatic elements.
export const pop = (frame: number, start: number, fps: number): number =>
  spring({
    frame: frame - start,
    fps,
    config: { damping: 12, stiffness: 140, mass: 0.7 },
  });

// Number of characters to reveal for a typewriter effect.
export const typed = (
  frame: number,
  start: number,
  text: string,
  cps = 38
): string => {
  const n = Math.max(0, Math.floor(((frame - start) / 30) * cps));
  return text.slice(0, n);
};

export const isTyping = (
  frame: number,
  start: number,
  text: string,
  cps = 38
): boolean => {
  const full = (text.length / cps) * 30;
  return frame >= start && frame < start + full;
};
