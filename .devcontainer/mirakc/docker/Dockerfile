FROM node:18-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
# RUN apt-get update && apt-get install -y \
    build-essential \
    libpcsclite-dev \
    pcscd \
    pkg-config \
    python3 \
    python3-pip

RUN npm install arib-b25-stream-test -g --unsafe 
RUN pip3 install supervisor

COPY --from=docker.io/mirakc/mirakc:debian /usr/local/bin/* /usr/local/bin/
COPY --from=docker.io/mirakc/mirakc:debian /etc/mirakc/strings.yml /etc/mirakc/strings.yml
# COPY --from=docker.io/mirakc/mirakc:debian /etc/mirakurun.openapi.json /etc/mirakurun.openapi.json
ENV MIRAKC_CONFIG=/etc/mirakc/config.yml

COPY supervisord.conf /etc/supervisord.conf

ENTRYPOINT ["/usr/local/bin/mirakc"]
