ARG BASE_IMAGE=debian
ARG TAG_NAME=bullseye-slim
FROM ${BASE_IMAGE}:${TAG_NAME}

RUN sed -ie 's/bullseye main/bullseye main non-free/g' /etc/apt/sources.list

RUN apt-get update && apt-get install -y --no-install-recommends \
    wget ca-certificates libva-dev libmfx-dev intel-media-va-driver-non-free vainfo xz-utils \
    python3 python-is-python3 python3-pip

RUN wget https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz && \
    tar Jxf ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz -C /usr/local --strip-component 1 && \
    rm ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz

RUN pip3 install --force-reinstall "faster-whisper @ https://github.com/guillaumekln/faster-whisper/archive/refs/heads/master.tar.gz"

COPY --from=ghcr.io/kounoike/dtv-discord-go:0.13.0-pre3 /dtv-discord-go /
COPY run_whisper_cuda.py run_whisper_cpu.py /
ARG BASE_IMAGE=debian
RUN [ "${BASE_IMAGE}" = "nvidia/cuda" ] && mv run_whisper_cuda.py run_whisper.py || mv run_whisper_cpu.py run_whisper.py

RUN useradd -u 1000 -m -s /bin/bash dtv
RUN usermod -aG video dtv
USER dtv
