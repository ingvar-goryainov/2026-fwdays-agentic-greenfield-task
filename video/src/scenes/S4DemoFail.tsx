import React from "react";
import { AbsoluteFill, useCurrentFrame, useVideoConfig } from "remotion";
import { Background } from "../components/Background";
import { Terminal, Line } from "../components/Terminal";
import { Caption } from "../components/Caption";
import { colors, fonts } from "../theme";
import { pop } from "../util";

const lines: Line[] = [
  { kind: "comment", text: "rename a script the file points to — real drift", at: 20 },
  {
    kind: "prompt",
    text: "mv scripts/deploy.sh scripts/release.sh",
    at: 70,
  },
  { kind: "prompt", text: "agents-lint scan ./AGENTS.md", at: 190 },
  {
    kind: "out",
    text: 'error  C001  ./AGENTS.md:15 — referenced path',
    at: 320,
    color: colors.red,
  },
  {
    kind: "out",
    text: '       "scripts/deploy.sh" does not exist at the repo root',
    at: 352,
    color: colors.red,
  },
  {
    kind: "out",
    text: "1 error(s), 0 warning(s)",
    at: 405,
    color: colors.muted,
  },
];

export const S4DemoFail: React.FC = () => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  const stamp = pop(frame, 440, fps);

  return (
    <AbsoluteFill>
      <Background accent={colors.red} />
      <AbsoluteFill style={{ alignItems: "center", justifyContent: "center" }}>
        <div style={{ position: "relative" }}>
          <Terminal lines={lines} width={1240} />
          {frame >= 440 && (
            <div
              style={{
                position: "absolute",
                right: -30,
                top: -34,
                transform: `scale(${stamp}) rotate(-8deg)`,
                background: colors.red,
                color: "#0d1117",
                fontFamily: fonts.mono,
                fontWeight: 700,
                fontSize: 30,
                padding: "10px 22px",
                borderRadius: 10,
                boxShadow: `0 12px 40px ${colors.red}66`,
              }}
            >
              exit 1
            </div>
          )}
        </div>
      </AbsoluteFill>
      <Caption
        text="C001 catches the drift · fails CI before the agent sees it"
        start={460}
        accent={colors.red}
      />
    </AbsoluteFill>
  );
};
