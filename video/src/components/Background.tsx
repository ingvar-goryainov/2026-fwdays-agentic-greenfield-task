import React from "react";
import { AbsoluteFill, useCurrentFrame } from "remotion";
import { colors } from "../theme";

export const Background: React.FC<{ accent?: string }> = ({
  accent = colors.green,
}) => {
  const frame = useCurrentFrame();
  const drift = Math.sin(frame / 90) * 40;
  return (
    <AbsoluteFill style={{ background: colors.bg }}>
      {/* soft radial glow that slowly drifts */}
      <AbsoluteFill
        style={{
          background: `radial-gradient(1200px 800px at ${
            50 + drift / 10
          }% 18%, ${accent}22, transparent 60%)`,
        }}
      />
      {/* faint dotted grid */}
      <AbsoluteFill
        style={{
          backgroundImage: `radial-gradient(${colors.border} 1px, transparent 1px)`,
          backgroundSize: "44px 44px",
          opacity: 0.28,
          maskImage:
            "radial-gradient(1400px 900px at 50% 40%, black, transparent 75%)",
          WebkitMaskImage:
            "radial-gradient(1400px 900px at 50% 40%, black, transparent 75%)",
        }}
      />
      {/* vignette */}
      <AbsoluteFill
        style={{
          boxShadow: "inset 0 0 400px rgba(0,0,0,0.75)",
        }}
      />
    </AbsoluteFill>
  );
};
