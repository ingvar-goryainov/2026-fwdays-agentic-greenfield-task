import React from "react";
import { useCurrentFrame } from "remotion";
import { colors, fonts } from "../theme";
import { appear } from "../util";

// Lower-third caption bar that reinforces the narration.
export const Caption: React.FC<{
  text: string;
  start?: number;
  accent?: string;
}> = ({ text, start = 0, accent = colors.green }) => {
  const frame = useCurrentFrame();
  const { opacity, translateY } = appear(frame, start, 14);
  return (
    <div
      style={{
        position: "absolute",
        bottom: 78,
        left: 0,
        right: 0,
        display: "flex",
        justifyContent: "center",
        opacity,
        transform: `translateY(${translateY}px)`,
      }}
    >
      <div
        style={{
          display: "flex",
          alignItems: "center",
          gap: 18,
          padding: "16px 34px",
          borderRadius: 999,
          background: "rgba(1,4,9,0.72)",
          border: `1px solid ${colors.border}`,
          backdropFilter: "blur(6px)",
        }}
      >
        <span
          style={{
            width: 12,
            height: 12,
            borderRadius: "50%",
            background: accent,
            boxShadow: `0 0 16px ${accent}`,
          }}
        />
        <span
          style={{
            fontFamily: fonts.sans,
            fontSize: 32,
            fontWeight: 600,
            color: colors.text,
            letterSpacing: 0.2,
          }}
        >
          {text}
        </span>
      </div>
    </div>
  );
};
