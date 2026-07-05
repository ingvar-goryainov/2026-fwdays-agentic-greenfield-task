import React from "react";
import { AbsoluteFill, useCurrentFrame, useVideoConfig, interpolate } from "remotion";
import { Background } from "../components/Background";
import { colors, fonts } from "../theme";
import { appear, pop } from "../util";

const stages = ["Proposal", "Design", "Tasks", "Apply", "Archive"];
const ids = [
  "FR-CLI-01",
  "FR-S001",
  "FR-C001",
  "FR-OUT-03",
  "NFR-DIST-02",
  "FR-CFG-01",
];

export const S6SpecFirst: React.FC = () => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  const title = appear(frame, 4, 12);

  return (
    <AbsoluteFill>
      <Background accent={colors.purple} />
      <AbsoluteFill
        style={{
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
          gap: 30,
        }}
      >
        <div
          style={{
            opacity: title.opacity,
            transform: `translateY(${title.translateY}px)`,
            fontFamily: fonts.sans,
            fontSize: 30,
            color: colors.purple,
            letterSpacing: 4,
          }}
        >
          HOW IT WAS BUILT — AGENTICALLY
        </div>
        <div
          style={{
            fontFamily: fonts.sans,
            fontSize: 46,
            fontWeight: 700,
            color: colors.text,
            opacity: appear(frame, 20, 12).opacity,
          }}
        >
          Spec first — one OpenSpec change per capability
        </div>

        {/* pipeline */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: 8,
            marginTop: 20,
          }}
        >
          {stages.map((s, i) => {
            const start = 60 + i * 34;
            const p = pop(frame, start, fps);
            const active = frame >= start;
            return (
              <React.Fragment key={s}>
                <div
                  style={{
                    transform: `scale(${p})`,
                    padding: "18px 30px",
                    borderRadius: 12,
                    fontFamily: fonts.mono,
                    fontSize: 28,
                    fontWeight: 600,
                    color: active ? colors.panel : colors.muted,
                    background: active ? colors.purple : colors.bgSoft,
                    border: `1px solid ${colors.purple}66`,
                  }}
                >
                  {s}
                </div>
                {i < stages.length - 1 && (
                  <div
                    style={{
                      width: 40,
                      height: 3,
                      background:
                        frame >= start + 20 ? colors.purple : colors.border,
                    }}
                  />
                )}
              </React.Fragment>
            );
          })}
        </div>

        {/* streaming requirement IDs */}
        <div
          style={{
            display: "flex",
            gap: 14,
            flexWrap: "wrap",
            justifyContent: "center",
            maxWidth: 1100,
            marginTop: 26,
          }}
        >
          {ids.map((id, i) => {
            const start = 230 + i * 12;
            const o = interpolate(frame, [start, start + 14], [0, 1], {
              extrapolateLeft: "clamp",
              extrapolateRight: "clamp",
            });
            return (
              <span
                key={id}
                style={{
                  opacity: o,
                  fontFamily: fonts.mono,
                  fontSize: 24,
                  color: colors.green,
                  background: `${colors.green}18`,
                  border: `1px solid ${colors.green}44`,
                  borderRadius: 8,
                  padding: "6px 14px",
                }}
              >
                {id}
              </span>
            );
          })}
        </div>
        <div
          style={{
            opacity: appear(frame, 300, 14).opacity,
            fontFamily: fonts.sans,
            fontSize: 28,
            color: colors.muted,
          }}
        >
          15 traceable requirements · every commit references an ID
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};
