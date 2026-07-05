import React from "react";
import { AbsoluteFill, useCurrentFrame, interpolate } from "remotion";
import { Background } from "../components/Background";
import { colors, fonts } from "../theme";
import { appear } from "../util";

const steps = ["fixture", "implement", "test", "lint", "commit"];
const R = 190;

const Counter: React.FC<{
  label: string;
  target: number;
  start: number;
  color: string;
}> = ({ label, target, start, color }) => {
  const frame = useCurrentFrame();
  const v = Math.round(
    interpolate(frame, [start, start + 40], [0, target], {
      extrapolateLeft: "clamp",
      extrapolateRight: "clamp",
    })
  );
  return (
    <div style={{ textAlign: "center" }}>
      <div
        style={{
          fontFamily: fonts.mono,
          fontSize: 72,
          fontWeight: 700,
          color,
        }}
      >
        {v}
      </div>
      <div style={{ fontFamily: fonts.sans, fontSize: 26, color: colors.muted }}>
        {label}
      </div>
    </div>
  );
};

export const S7Loops: React.FC = () => {
  const frame = useCurrentFrame();
  const title = appear(frame, 4, 12);
  // one full highlight rotation across the 5 nodes
  const active = Math.floor(((frame - 40) / 18)) % steps.length;

  return (
    <AbsoluteFill>
      <Background accent={colors.green} />
      <AbsoluteFill
        style={{
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "row",
          gap: 120,
        }}
      >
        {/* loop diagram */}
        <div style={{ position: "relative", width: 520, height: 520 }}>
          <svg
            width={520}
            height={520}
            style={{ position: "absolute", inset: 0 }}
          >
            <circle
              cx={260}
              cy={260}
              r={R}
              fill="none"
              stroke={colors.border}
              strokeWidth={2}
            />
            <circle
              cx={260}
              cy={260}
              r={R}
              fill="none"
              stroke={colors.green}
              strokeWidth={4}
              strokeDasharray={`${2 * Math.PI * R}`}
              strokeDashoffset={interpolate(
                frame,
                [40, 150],
                [2 * Math.PI * R, 0],
                { extrapolateLeft: "clamp", extrapolateRight: "clamp" }
              )}
              strokeLinecap="round"
              transform="rotate(-90 260 260)"
            />
          </svg>
          {steps.map((s, i) => {
            const ang = (i / steps.length) * Math.PI * 2 - Math.PI / 2;
            const x = 260 + R * Math.cos(ang);
            const y = 260 + R * Math.sin(ang);
            const on = frame >= 40 && active === i;
            const shown = frame >= 40 + i * 8;
            return (
              <div
                key={s}
                style={{
                  position: "absolute",
                  left: x,
                  top: y,
                  transform: "translate(-50%,-50%)",
                  opacity: shown ? 1 : 0,
                  padding: "12px 20px",
                  borderRadius: 999,
                  fontFamily: fonts.mono,
                  fontSize: 24,
                  fontWeight: 600,
                  color: on ? colors.panel : colors.text,
                  background: on ? colors.green : colors.bgSoft,
                  border: `1px solid ${colors.green}66`,
                  boxShadow: on ? `0 0 30px ${colors.green}88` : "none",
                  transition: "none",
                }}
              >
                {s}
              </div>
            );
          })}
          <div
            style={{
              position: "absolute",
              inset: 0,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              flexDirection: "column",
            }}
          >
            <div
              style={{
                fontFamily: fonts.sans,
                fontSize: 30,
                fontWeight: 700,
                color: colors.text,
              }}
            >
              the loop
            </div>
            <div
              style={{ fontFamily: fonts.sans, fontSize: 20, color: colors.muted }}
            >
              run to green
            </div>
          </div>
        </div>

        {/* counters */}
        <div style={{ display: "flex", flexDirection: "column", gap: 40 }}>
          <div
            style={{
              opacity: title.opacity,
              fontFamily: fonts.sans,
              fontSize: 40,
              fontWeight: 700,
              color: colors.text,
              maxWidth: 520,
              lineHeight: 1.3,
            }}
          >
            A loop the agent ran — not step-by-step prompting
          </div>
          <div style={{ display: "flex", gap: 56 }}>
            <Counter label="commits" target={28} start={150} color={colors.green} />
            <Counter label="tests" target={76} start={175} color={colors.blue} />
            <Counter
              label="fixtures"
              target={17}
              start={200}
              color={colors.purple}
            />
          </div>
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};
