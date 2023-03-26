FROM node:18-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
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
ENV MIRAKC_CONFIG=/etc/mirakc/config.yml

COPY supervisord.conf /etc/supervisord.conf

CMD ["/usr/local/bin/supervisord", "--nodaemon", "-c", "/etc/supervisord.conf"]
