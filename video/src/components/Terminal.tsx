import React from "react";
import { useCurrentFrame } from "remotion";
import { colors, fonts } from "../theme";

export type Line =
  | { kind: "prompt"; text: string; at: number } // typed command after `$`
  | { kind: "out"; text: string; at: number; color?: string }
  | { kind: "comment"; text: string; at: number };

const CPS = 42;

const Cursor: React.FC = () => {
  const frame = useCurrentFrame();
  const on = Math.floor(frame / 15) % 2 === 0;
  return (
    <span
      style={{
        display: "inline-block",
        width: 14,
        height: 26,
        marginLeft: 2,
        background: on ? colors.green : "transparent",
        transform: "translateY(4px)",
      }}
    />
  );
};

export const Terminal: React.FC<{
  title?: string;
  lines: Line[];
  width?: number;
}> = ({ title = "zsh — agents-lint", lines, width = 1180 }) => {
  const frame = useCurrentFrame();

  // Find the currently-typing prompt line (only one cursor at a time).
  let activePromptIndex = -1;
  for (let i = 0; i < lines.length; i++) {
    const l = lines[i];
    if (l.kind === "prompt") {
      const dur = (l.text.length / CPS) * 30;
      if (frame >= l.at && frame < l.at + dur + 8) activePromptIndex = i;
    }
  }
  const lastVisible = lines
    .map((l, i) => ({ l, i }))
    .filter(({ l }) => frame >= l.at)
    .pop();

  return (
    <div
      style={{
        width,
        borderRadius: 14,
        overflow: "hidden",
        background: colors.panel,
        border: `1px solid ${colors.border}`,
        boxShadow: "0 40px 120px rgba(0,0,0,0.6)",
        fontFamily: fonts.mono,
      }}
    >
      <div
        style={{
          height: 46,
          background: "#161b22",
          borderBottom: `1px solid ${colors.border}`,
          display: "flex",
          alignItems: "center",
          padding: "0 18px",
          gap: 9,
        }}
      >
        <Dot c="#ff5f56" />
        <Dot c="#ffbd2e" />
        <Dot c="#27c93f" />
        <span
          style={{
            color: colors.muted,
            fontSize: 18,
            marginLeft: 16,
            letterSpacing: 0.3,
          }}
        >
          {title}
        </span>
      </div>
      <div style={{ padding: "26px 30px", minHeight: 360 }}>
        {lines.map((l, i) => {
          if (frame < l.at) return null;
          const isActive = i === activePromptIndex;
          let shown = l.text;
          if (l.kind === "prompt") {
            const n = Math.floor(((frame - l.at) / 30) * CPS);
            shown = l.text.slice(0, Math.max(0, n));
          }
          const showCursor =
            (isActive && l.kind === "prompt") ||
            (lastVisible?.i === i && l.kind !== "prompt" && !anyTyping(lines, frame));
          return (
            <div
              key={i}
              style={{
                fontSize: 26,
                lineHeight: 1.62,
                color:
                  l.kind === "comment"
                    ? colors.muted
                    : l.kind === "out"
                    ? l.color ?? colors.text
                    : colors.text,
                whiteSpace: "pre-wrap",
              }}
            >
              {l.kind === "prompt" && (
                <span style={{ color: colors.green }}>$ </span>
              )}
              {l.kind === "comment" && (
                <span style={{ color: colors.muted }}># </span>
              )}
              {shown}
              {showCursor && <Cursor />}
            </div>
          );
        })}
      </div>
    </div>
  );
};

const anyTyping = (lines: Line[], frame: number): boolean =>
  lines.some(
    (l) =>
      l.kind === "prompt" &&
      frame >= l.at &&
      frame < l.at + (l.text.length / CPS) * 30
  );

const Dot: React.FC<{ c: string }> = ({ c }) => (
  <span
    style={{
      width: 15,
      height: 15,
      borderRadius: "50%",
      background: c,
      display: "inline-block",
    }}
  />
);
