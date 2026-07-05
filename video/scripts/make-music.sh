#!/usr/bin/env bash
# Generate a subtle ambient pad bed (keynote-style) with ffmpeg — no external API.
# Output: public/music.mp3  (~32s, Remotion loops it across the 120s timeline)
set -euo pipefail
cd "$(dirname "$0")/.."
mkdir -p public tmp_music
cd tmp_music

# I–V–vi–IV pad progression, warm mid octaves. Each chord = 8s.
# chord <out> <f1> <f2> <f3> <f4>
chord () {
  local out=$1 f1=$2 f2=$3 f4=$5 f3=$4
  ffmpeg -y -loglevel error \
    -f lavfi -i "sine=frequency=${f1}:duration=8" \
    -f lavfi -i "sine=frequency=${f2}:duration=8" \
    -f lavfi -i "sine=frequency=${f3}:duration=8" \
    -f lavfi -i "sine=frequency=${f4}:duration=8" \
    -filter_complex "[0][1][2][3]amix=inputs=4:normalize=1,\
afade=t=in:st=0:d=2.2,afade=t=out:st=5.8:d=2.2" \
    "$out"
}

chord c1.wav 130.81 164.81 196.00 261.63   # C  major
chord c2.wav 146.83 196.00 246.94 293.66   # G  major
chord c3.wav 110.00 164.81 220.00 261.63   # A  minor
chord c4.wav 174.61 220.00 261.63 349.23   # F  major

# concat, then master: soften (lowpass), gentle movement (tremolo),
# space (echo), and keep it quiet.
ffmpeg -y -loglevel error \
  -i c1.wav -i c2.wav -i c3.wav -i c4.wav \
  -filter_complex "[0][1][2][3]concat=n=4:v=0:a=1[a];\
[a]lowpass=f=1600,tremolo=f=0.15:d=0.4,\
aecho=0.8:0.7:400|800:0.3|0.2,\
volume=0.9,afade=t=in:st=0:d=1.5,afade=t=out:st=30:d=2[m]" \
  -map "[m]" -ar 44100 -b:a 192k ../public/music.mp3

cd ..
rm -rf tmp_music
echo "wrote public/music.mp3"
ls -la public/music.mp3
