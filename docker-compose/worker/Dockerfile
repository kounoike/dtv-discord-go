FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ffmpeg

COPY --from=ghcr.io/kounoike/dtv-discord-go:latest /dtv-discord-go /

CMD ["/dtv-discord-go", "worker"]
