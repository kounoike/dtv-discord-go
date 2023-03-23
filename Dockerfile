# syntax=docker/dockerfile:1
FROM gcr.io/distroless/static-debian11

COPY dtv-discord-go /
COPY db/migrations/* /db/migrations/

CMD ["/dtv-discord-go"]
