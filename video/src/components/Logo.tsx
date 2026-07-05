import React from "react";
import { useCurrentFrame, useVideoConfig, interpolate } from "remotion";
import { colors, fonts } from "../theme";
import { pop } from "../util";

// The agents-lint wordmark: a check-mark badge + monospace name.
export const Logo: React.FC<{ start?: number; scale?: number }> = ({
  start = 0,
  scale = 1,
}) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  const badge = pop(frame, start, fps);
  const draw = interpolate(frame, [start + 6, start + 26], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });
  const nameOpacity = interpolate(frame, [start + 14, start + 30], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        gap: 30 * scale,
        transform: `scale(${scale})`,
      }}
    >
      <div
        style={{
          width: 116,
          height: 116,
          borderRadius: 26,
          background: `linear-gradient(150deg, ${colors.greenSoft}, #1a7f37)`,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          transform: `scale(${badge})`,
          boxShadow: `0 20px 60px ${colors.green}55`,
        }}
      >
        <svg width="72" height="72" viewBox="0 0 24 24" fill="none">
          <path
            d="M4 12.5 L10 18.5 L20 6"
            stroke="#0d1117"
            strokeWidth="3"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeDasharray={26}
            strokeDashoffset={26 * (1 - draw)}
          />
        </svg>
      </div>
      <div style={{ opacity: nameOpacity }}>
        <div
          style={{
            fontFamily: fonts.mono,
            fontSize: 82,
            fontWeight: 700,
            color: colors.text,
            letterSpacing: -1,
          }}
        >
          agents<span style={{ color: colors.green }}>-lint</span>
        </div>
      </div>
    </div>
  );
};
