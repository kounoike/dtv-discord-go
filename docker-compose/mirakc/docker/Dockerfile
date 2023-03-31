FROM node:18-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    curl \
    dvb-tools \
    libpcsclite-dev \
    pcscd \
    pkg-config \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN npm install arib-b25-stream-test -g --unsafe 
RUN mkdir /etc/mirakc
RUN curl -sSL -o /etc/mirakc/dvbv5_channels_isdbs.conf https://raw.githubusercontent.com/Chinachu/dvbconf-for-isdb/master/conf/dvbv5_channels_isdbs.conf
RUN curl -sSL -o /etc/mirakc/dvbv5_channels_isdbt.conf https://raw.githubusercontent.com/Chinachu/dvbconf-for-isdb/master/conf/dvbv5_channels_isdbt.conf

COPY --from=docker.io/mirakc/mirakc:debian /usr/local/bin/* /usr/local/bin/
COPY --from=docker.io/mirakc/mirakc:debian /etc/mirakc/strings.yml /etc/mirakc/strings.yml
ENV MIRAKC_CONFIG=/etc/mirakc/config.yml

ENTRYPOINT ["/usr/local/bin/mirakc"]
