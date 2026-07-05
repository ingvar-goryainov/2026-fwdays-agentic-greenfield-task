import React from "react";
import { AbsoluteFill, useCurrentFrame, useVideoConfig } from "remotion";
import { Background } from "../components/Background";
import { colors, fonts } from "../theme";
import { appear, pop } from "../util";

const Panel: React.FC<{
  role: string;
  who: string;
  detail: string;
  color: string;
  start: number;
  icon: string;
}> = ({ role, who, detail, color, start, icon }) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();
  const p = pop(frame, start, fps);
  return (
    <div
      style={{
        transform: `scale(${p})`,
        width: 460,
        background: colors.bgSoft,
        border: `1px solid ${color}66`,
        borderRadius: 18,
        padding: "34px 36px",
      }}
    >
      <div style={{ fontSize: 52, marginBottom: 14 }}>{icon}</div>
      <div
        style={{
          fontFamily: fonts.mono,
          fontSize: 22,
          letterSpacing: 3,
          color,
        }}
      >
        {role}
      </div>
      <div
        style={{
          fontFamily: fonts.sans,
          fontSize: 40,
          fontWeight: 700,
          color: colors.text,
          margin: "6px 0 10px",
        }}
      >
        {who}
      </div>
      <div style={{ fontFamily: fonts.sans, fontSize: 25, color: colors.muted }}>
        {detail}
      </div>
    </div>
  );
};

export const S8MakerChecker: React.FC = () => {
  const frame = useCurrentFrame();
  const title = appear(frame, 4, 12);
  const neq = pop(frame, 70, 30);

  return (
    <AbsoluteFill>
      <Background accent={colors.yellow} />
      <AbsoluteFill
        style={{
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
          gap: 44,
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
          The maker was never the checker
        </div>
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: 40,
          }}
        >
          <Panel
            role="MAKER"
            who="the building agent"
            detail="spec → code → green tests"
            color={colors.green}
            start={20}
            icon="🛠️"
          />
          <div
            style={{
              transform: `scale(${neq})`,
              fontFamily: fonts.mono,
              fontSize: 80,
              fontWeight: 700,
              color: colors.yellow,
            }}
          >
            ≠
          </div>
          <Panel
            role="CHECKER"
            who="a separate review"
            detail="CodeRabbit on every PR · code-graph MCP"
            color={colors.blue}
            start={40}
            icon="🔍"
          />
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};
