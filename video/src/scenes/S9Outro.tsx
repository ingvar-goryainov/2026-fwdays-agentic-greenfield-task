import React from "react";
import { AbsoluteFill, useCurrentFrame } from "remotion";
import { Background } from "../components/Background";
import { Logo } from "../components/Logo";
import { colors, fonts } from "../theme";
import { appear } from "../util";

export const S9Outro: React.FC = () => {
  const frame = useCurrentFrame();
  const tag = appear(frame, 40, 16);

  return (
    <AbsoluteFill>
      <Background />
      <AbsoluteFill
        style={{
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
          gap: 34,
        }}
      >
        <Logo start={6} />
        <div
          style={{
            opacity: tag.opacity,
            transform: `translateY(${tag.translateY}px)`,
            fontFamily: fonts.sans,
            fontSize: 34,
            color: colors.muted,
            textAlign: "center",
          }}
        >
          A small tool, taken through a full engineering loop —{" "}
          <span style={{ color: colors.green, fontWeight: 700 }}>
            built agentically
          </span>
          <div
            style={{
              fontSize: 24,
              marginTop: 16,
              color: colors.muted,
              fontFamily: fonts.mono,
            }}
          >
            Go · Cobra · OpenSpec · Claude Code
          </div>
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};
