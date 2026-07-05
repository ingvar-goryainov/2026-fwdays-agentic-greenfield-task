import React from "react";
import { AbsoluteFill, useCurrentFrame, interpolate } from "remotion";
import { Background } from "../components/Background";
import { Logo } from "../components/Logo";
import { colors, fonts } from "../theme";
import { appear } from "../util";

export const S1Hook: React.FC = () => {
  const frame = useCurrentFrame();
  const tag = appear(frame, 34, 16);
  const sub = appear(frame, 150, 16);
  // Gentle exit drift so it hands off to the problem scene.
  const exit = interpolate(frame, [270, 300], [0, -30], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

  return (
    <AbsoluteFill>
      <Background />
      <AbsoluteFill
        style={{
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
          gap: 34,
          transform: `translateY(${exit}px)`,
        }}
      >
        <Logo start={8} />
        <div
          style={{
            opacity: tag.opacity,
            transform: `translateY(${tag.translateY}px)`,
            fontFamily: fonts.sans,
            fontSize: 40,
            color: colors.muted,
            letterSpacing: 0.5,
          }}
        >
          a compiler for your{" "}
          <span style={{ color: colors.text, fontFamily: fonts.mono }}>
            AGENTS.md
          </span>
        </div>
        <div
          style={{
            opacity: sub.opacity,
            transform: `translateY(${sub.translateY}px)`,
            marginTop: 26,
            fontFamily: fonts.sans,
            fontSize: 28,
            color: colors.muted,
            maxWidth: 900,
            textAlign: "center",
            lineHeight: 1.5,
          }}
        >
          Agent config drifts. It points at code that no longer exists —
          and nobody notices until the agent goes wrong.
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};
