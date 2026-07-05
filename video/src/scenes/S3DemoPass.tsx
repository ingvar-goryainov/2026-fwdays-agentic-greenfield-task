import React from "react";
import { AbsoluteFill, useCurrentFrame, useVideoConfig } from "remotion";
import { Background } from "../components/Background";
import { Terminal, Line } from "../components/Terminal";
import { Caption } from "../components/Caption";
import { colors, fonts } from "../theme";
import { pop } from "../util";

const lines: Line[] = [
  { kind: "prompt", text: "agents-lint init", at: 20 },
  { kind: "out", text: "created ./AGENTS.md", at: 78, color: colors.green },
  { kind: "comment", text: "seven rules — schema + codebase", at: 120 },
  { kind: "prompt", text: "agents-lint scan ./AGENTS.md", at: 165 },
  {
    kind: "out",
    text: "✓ AGENTS.md is valid (7 rules passed)",
    at: 300,
    color: colors.green,
  },
];

export const S3DemoPass: React.FC = () => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  const stamp = pop(frame, 330, fps);

  return (
    <AbsoluteFill>
      <Background />
      <AbsoluteFill style={{ alignItems: "center", justifyContent: "center" }}>
        <div style={{ position: "relative" }}>
          <Terminal lines={lines} />
          {frame >= 330 && (
            <div
              style={{
                position: "absolute",
                right: -34,
                top: -34,
                transform: `scale(${stamp}) rotate(-8deg)`,
                background: colors.green,
                color: colors.panel,
                fontFamily: fonts.mono,
                fontWeight: 700,
                fontSize: 30,
                padding: "10px 22px",
                borderRadius: 10,
                boxShadow: `0 12px 40px ${colors.green}66`,
              }}
            >
              exit 0
            </div>
          )}
        </div>
      </AbsoluteFill>
      <Caption text="init scaffolds it · scan runs 7 rules · all green" start={345} />
    </AbsoluteFill>
  );
};
