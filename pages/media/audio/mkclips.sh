#!/usr/bin/env bash
# Usage: mkclips.sh <out-dir>. A 12-second rising whole-tone scale from A3, one
# note per second with a short fade in and out, rendered once to WAV and then
# encoded to each format.
set -euo pipefail
out=$1
mkdir -p "$out"
expr='0.35*sin(2*PI*220*pow(2\,floor(t)*2/12)*t)*min(1\,mod(t\,1)*40)*min(1\,(1-mod(t\,1))*8)'
master=$(mktemp -d)/master.wav
ffmpeg -v error -y -f lavfi -i "aevalsrc=${expr}:s=44100" -t 12 -ac 1 -c:a pcm_s16le "$master"
enc() { local name=$1; shift; ffmpeg -v error -y -i "$master" "$@" "$out/$name"; }
enc scale.mp3  -c:a libmp3lame -b:a 64k
enc scale.m4a  -c:a aac -b:a 64k
enc scale.ogg  -ac 2 -c:a vorbis -strict -2 -q:a 2
enc scale.opus -ar 48000 -c:a libopus -b:a 32k
enc scale.flac -ar 22050 -sample_fmt s16 -c:a flac
enc scale.wav  -ar 22050 -c:a pcm_s16le
