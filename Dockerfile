# syntax=docker/dockerfile:1

FROM golang:1.20-bullseye AS build

WORKDIR /app
COPY ./ ./
RUN go mod download

RUN go build -o /dtv-discord-go

FROM debian:bullseye
COPY --from=build /dtv-discord-go /
COPY --from=build /app/db/migrations/* /db/migrations/

CMD ["/dtv-discord-go"]
