import React from "react";
import { AbsoluteFill, useCurrentFrame, useVideoConfig } from "remotion";
import { Background } from "../components/Background";
import { colors, fonts } from "../theme";
import { appear, pop } from "../util";

const schema = [
  ["S001", "file exists"],
  ["S002", "Agents section"],
  ["S003", "valid agent block"],
  ["S004", "no duplicate names"],
  ["S005", "valid YAML frontmatter"],
];
const codebase = [
  ["C001", "paths resolve to real files"],
  ["C002", "tools have repo evidence"],
];

const Chip: React.FC<{
  id: string;
  desc: string;
  start: number;
  accent: string;
}> = ({ id, desc, start, accent }) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  const p = pop(frame, start, fps);
  return (
    <div
      style={{
        transform: `scale(${p})`,
        display: "flex",
        alignItems: "center",
        gap: 16,
        background: colors.bgSoft,
        border: `1px solid ${accent}55`,
        borderRadius: 12,
        padding: "16px 22px",
        width: 430,
      }}
    >
      <span
        style={{
          fontFamily: fonts.mono,
          fontSize: 26,
          fontWeight: 700,
          color: accent,
        }}
      >
        {id}
      </span>
      <span style={{ fontFamily: fonts.sans, fontSize: 24, color: colors.text }}>
        {desc}
      </span>
    </div>
  );
};

export const S5Rules: React.FC = () => {
  const frame = useCurrentFrame();
  const title = appear(frame, 4, 12);
  const sarif = appear(frame, 210, 16);

  return (
    <AbsoluteFill>
      <Background accent={colors.blue} />
      <AbsoluteFill
        style={{
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
          gap: 34,
        }}
      >
        <div
          style={{
            opacity: title.opacity,
            transform: `translateY(${title.translateY}px)`,
            fontFamily: fonts.sans,
            fontSize: 44,
            fontWeight: 700,
            color: colors.text,
          }}
        >
          7 rules · one Go binary · zero deps
        </div>

        <div style={{ display: "flex", gap: 34, alignItems: "flex-start" }}>
          <div style={{ display: "flex", flexDirection: "column", gap: 14 }}>
            <Label text="SCHEMA" color={colors.green} />
            {schema.map(([id, d], i) => (
              <Chip
                key={id}
                id={id}
                desc={d}
                start={20 + i * 16}
                accent={colors.green}
              />
            ))}
          </div>
          <div style={{ display: "flex", flexDirection: "column", gap: 14 }}>
            <Label text="CODEBASE-AWARE" color={colors.blue} />
            {codebase.map(([id, d], i) => (
              <Chip
                key={id}
                id={id}
                desc={d}
                start={110 + i * 16}
                accent={colors.blue}
              />
            ))}
            <div
              style={{
                opacity: sarif.opacity,
                transform: `translateY(${sarif.translateY}px)`,
                marginTop: 22,
                background: colors.panel,
                border: `1px solid ${colors.border}`,
                borderRadius: 12,
                padding: "18px 22px",
                width: 430,
                fontFamily: fonts.mono,
                fontSize: 20,
                color: colors.muted,
                lineHeight: 1.5,
              }}
            >
              <div style={{ color: colors.purple, marginBottom: 6 }}>
                --format sarif
              </div>
              {"{ "}
              <span style={{ color: colors.blue }}>"version"</span>: "2.1.0",
              <br />
              {"  "}
              <span style={{ color: colors.blue }}>"ruleId"</span>:{" "}
              <span style={{ color: colors.green }}>"C001"</span> {"}"}
              <div style={{ color: colors.text, marginTop: 8 }}>
                → GitHub code scanning
              </div>
            </div>
          </div>
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};

const Label: React.FC<{ text: string; color: string }> = ({ text, color }) => (
  <div
    style={{
      fontFamily: fonts.mono,
      fontSize: 20,
      letterSpacing: 3,
      color,
      marginBottom: 4,
    }}
  >
    {text}
  </div>
);
