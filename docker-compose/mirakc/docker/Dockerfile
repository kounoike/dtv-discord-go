FROM node:18-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libpcsclite-dev \
    pcscd \
    pkg-config

RUN npm install arib-b25-stream-test -g --unsafe 

COPY --from=docker.io/mirakc/mirakc:debian /usr/local/bin/* /usr/local/bin/
COPY --from=docker.io/mirakc/mirakc:debian /etc/mirakc/strings.yml /etc/mirakc/strings.yml
ENV MIRAKC_CONFIG=/etc/mirakc/config.yml

ENTRYPOINT ["/usr/local/bin/mirakc"]

