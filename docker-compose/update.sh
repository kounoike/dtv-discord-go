#!/bin/sh

set -e

docker pull ghcr.io/kounoike/dtv-discord-go:latest
docker pull mirakc/mirakc:debian
docker pull mariadb:latest
docker pull redis:latest

docker compose build --pull
