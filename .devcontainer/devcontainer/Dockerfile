FROM mcr.microsoft.com/devcontainers/go:0-1.20-bullseye
RUN go install github.com/rubenv/sql-migrate/...@latest
RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
RUN chmod -R og+w /go
