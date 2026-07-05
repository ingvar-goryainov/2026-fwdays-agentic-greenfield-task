import React from "react";
import { AbsoluteFill, Sequence } from "remotion";
import "./fonts";
import { scenes } from "./theme";
import { S1Hook } from "./scenes/S1Hook";
import { S2Problem } from "./scenes/S2Problem";
import { S3DemoPass } from "./scenes/S3DemoPass";
import { S4DemoFail } from "./scenes/S4DemoFail";
import { S5Rules } from "./scenes/S5Rules";
import { S6SpecFirst } from "./scenes/S6SpecFirst";
import { S7Loops } from "./scenes/S7Loops";
import { S8MakerChecker } from "./scenes/S8MakerChecker";
import { S9Outro } from "./scenes/S9Outro";
import { Soundtrack } from "./Soundtrack";

const Scene: React.FC<{
  cfg: { from: number; durationInFrames: number };
  children: React.ReactNode;
}> = ({ cfg, children }) => (
  <Sequence from={cfg.from} durationInFrames={cfg.durationInFrames}>
    {children}
  </Sequence>
);

export const Video: React.FC = () => {
  return (
    <AbsoluteFill style={{ background: "#0d1117" }}>
      <Scene cfg={scenes.hook}>
        <S1Hook />
      </Scene>
      <Scene cfg={scenes.problem}>
        <S2Problem />
      </Scene>
      <Scene cfg={scenes.demoPass}>
        <S3DemoPass />
      </Scene>
      <Scene cfg={scenes.demoFail}>
        <S4DemoFail />
      </Scene>
      <Scene cfg={scenes.rules}>
        <S5Rules />
      </Scene>
      <Scene cfg={scenes.specFirst}>
        <S6SpecFirst />
      </Scene>
      <Scene cfg={scenes.loops}>
        <S7Loops />
      </Scene>
      <Scene cfg={scenes.makerChecker}>
        <S8MakerChecker />
      </Scene>
      <Scene cfg={scenes.outro}>
        <S9Outro />
      </Scene>
      <Soundtrack />
    </AbsoluteFill>
  );
};
