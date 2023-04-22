FROM python:3.10-slim

RUN sed -ie 's/bullseye main/bullseye main non-free/' /etc/apt/sources.list

RUN apt-get update && apt-get install -y --no-install-recommends \
    wget ca-certificates libva-dev libmfx-dev intel-media-va-driver-non-free vainfo xz-utils

RUN wget https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz && \
    tar Jxf ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz -C /usr/local --strip-component 1 && \
    rm ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz

RUN pip install --force-reinstall "faster-whisper @ https://github.com/guillaumekln/faster-whisper/archive/refs/heads/master.tar.gz"

COPY --from=ghcr.io/kounoike/dtv-discord-go:latest /dtv-discord-go /
COPY run_whisper.py /

RUN useradd -u 1000 -m -s /bin/bash dtv
RUN usermod -aG video dtv
USER dtv
