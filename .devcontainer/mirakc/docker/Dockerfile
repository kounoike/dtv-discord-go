FROM node:18-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    # RUN apt-get update && apt-get install -y \
    build-essential \
    ca-certificates \
    curl \
    dvb-tools \
    libpcsclite-dev \
    pcscd \
    pkg-config

RUN npm install arib-b25-stream-test -g --unsafe 

RUN mkdir -p /etc/mirakc
RUN curl -sSL -o /etc/mirakc/dvbv5_channels_isdbs.conf https://raw.githubusercontent.com/Chinachu/dvbconf-for-isdb/master/conf/dvbv5_channels_isdbs.conf
RUN curl -sSL -o /etc/mirakc/dvbv5_channels_isdbt.conf https://raw.githubusercontent.com/Chinachu/dvbconf-for-isdb/master/conf/dvbv5_channels_isdbt.conf

COPY --from=docker.io/mirakc/mirakc:debian /usr/local/bin/* /usr/local/bin/
COPY --from=docker.io/mirakc/mirakc:debian /etc/mirakc/strings.yml /etc/mirakc/strings.yml
# COPY --from=docker.io/mirakc/mirakc:debian /etc/mirakurun.openapi.json /etc/mirakurun.openapi.json
ENV MIRAKC_CONFIG=/etc/mirakc/config.yml

USER node

ENTRYPOINT ["/usr/local/bin/mirakc"]
