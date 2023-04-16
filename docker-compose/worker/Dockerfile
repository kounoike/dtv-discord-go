FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends \
    wget ca-certificates libva-dev libmfx-dev intel-media-va-driver-non-free vainfo xz-utils \
    python3 python-is-python3

RUN wget https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz && \
    tar Jxf ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz -C /usr/local --strip-component 1 && \
    rm ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz

COPY --from=ghcr.io/kounoike/dtv-discord-go:0.12.0-pre3 /dtv-discord-go /
COPY run_whisper.py /

CMD ["/dtv-discord-go", "worker"]
