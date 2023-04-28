FROM golang:1.20-bullseye AS go-build
WORKDIR /app/ofelia
RUN curl -fsSL https://github.com/mcuadros/ofelia/archive/refs/tags/v0.3.7.tar.gz \
    | tar xzvpf - --strip-components=1
RUN go build -o /usr/local/bin/ofelia

FROM nginx:latest AS runner

ENV DEBIAN_FRONTEND=noninteractive

RUN set -ex && \
    apt-get update && \
    apt-get install --yes --no-install-recommends \
    ca-certificates \
    curl \
    jq \
    git \
    gnupg \
    openssl \
    python3 \
    python3-pip \
    && \
    pip3 install supervisor && \
    curl -fsSL https://tailscale.com/install.sh | sh && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /var/log/apt /var/log/dpkg.log

RUN mv /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.orig

COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

COPY --from=go-build /usr/local/bin/ofelia /usr/local/bin/
COPY --from=ghcr.io/kounoike/dtv-discord-go:0.13.0-pre3 /frontend/out /www
COPY supervisord.conf ofelia.ini /etc/
COPY nginx.conf.tmpl /etc/nginx/

ENTRYPOINT ["/docker-entrypoint.sh"]
