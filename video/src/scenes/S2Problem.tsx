import React from "react";
import { AbsoluteFill, useCurrentFrame, interpolate } from "remotion";
import { Background } from "../components/Background";
import { colors, fonts } from "../theme";
import { appear } from "../util";

type Row = { text: string; badAt?: number; note?: string };

const rows: Row[] = [
  { text: "## Agents" },
  { text: "" },
  { text: "### Deploy Agent" },
  { text: "**Instructions:** run `scripts/deploy.sh`", badAt: 70, note: "file renamed" },
  { text: "" },
  { text: "### Deploy Agent", badAt: 96, note: "duplicate name" },
  { text: "**Tools:** `terraform`, `kubectl`", badAt: 122, note: "no evidence in repo" },
  { text: "" },
  { text: "---", badAt: 0 },
  { text: "name: [unclosed", badAt: 148, note: "broken YAML" },
];

export const S2Problem: React.FC = () => {
  const frame = useCurrentFrame();
  const title = appear(frame, 6, 14);

  return (
    <AbsoluteFill>
      <Background accent={colors.red} />
      <AbsoluteFill
        style={{
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
          gap: 40,
        }}
      >
        <div
          style={{
            opacity: title.opacity,
            transform: `translateY(${title.translateY}px)`,
            fontFamily: fonts.sans,
            fontSize: 46,
            fontWeight: 700,
            color: colors.text,
          }}
        >
          Written once, forgotten —{" "}
          <span style={{ color: colors.red }}>then it rots</span>
        </div>

        <div
          style={{
            width: 1120,
            background: colors.panel,
            border: `1px solid ${colors.border}`,
            borderRadius: 14,
            padding: "34px 44px",
            fontFamily: fonts.mono,
            boxShadow: "0 40px 120px rgba(0,0,0,0.6)",
          }}
        >
          <div style={{ color: colors.muted, fontSize: 22, marginBottom: 18 }}>
            AGENTS.md
          </div>
          {rows.map((r, i) => {
            const bad = r.badAt !== undefined && r.badAt > 0 && frame >= r.badAt;
            const strike = bad
              ? interpolate(frame, [r.badAt!, r.badAt! + 18], [0, 1], {
                  extrapolateLeft: "clamp",
                  extrapolateRight: "clamp",
                })
              : 0;
            return (
              <div
                key={i}
                style={{
                  position: "relative",
                  fontSize: 27,
                  lineHeight: 1.68,
                  color: bad ? colors.red : colors.text,
                  minHeight: r.text === "" ? 20 : undefined,
                  display: "flex",
                  alignItems: "center",
                  gap: 20,
                }}
              >
                <span style={{ position: "relative" }}>
                  {r.text}
                  {bad && (
                    <span
                      style={{
                        position: "absolute",
                        left: 0,
                        top: "52%",
                        height: 3,
                        width: `${strike * 100}%`,
                        background: colors.red,
                      }}
                    />
                  )}
                </span>
                {bad && strike > 0.6 && r.note && (
                  <span
                    style={{
                      fontFamily: fonts.sans,
                      fontSize: 20,
                      color: colors.red,
                      background: `${colors.red}1e`,
                      border: `1px solid ${colors.red}55`,
                      borderRadius: 8,
                      padding: "3px 12px",
                    }}
                  >
                    ✗ {r.note}
                  </span>
                )}
              </div>
            );
          })}
        </div>
        <div
          style={{
            opacity: appear(frame, 210, 14).opacity,
            fontFamily: fonts.sans,
            fontSize: 30,
            color: colors.muted,
          }}
        >
          Other linters <span style={{ color: colors.muted }}>score</span> this.
          agents-lint{" "}
          <span style={{ color: colors.green, fontWeight: 700 }}>
            passes or fails
          </span>{" "}
          it.
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};
