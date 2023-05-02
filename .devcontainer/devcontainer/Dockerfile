ARG VARIANT="11.7.1"
FROM nvidia/cuda:${VARIANT}-cudnn8-runtime-ubuntu22.04

ARG INSTALL_ZSH="false"
ARG UPGRADE_PACKAGES="false"

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID
COPY library-scripts/common-debian.sh /tmp/

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    wget \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts

RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash -

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && /bin/bash /tmp/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" \
    && apt-get install -y --no-install-recommends \
    ca-certificates \
    intel-media-va-driver-non-free \
    libva-dev \
    libmfx-dev \
    mariadb-client \
    nodejs \
    python-is-python3 \
    python3-pip \
    vainfo \
    wget \
    xz-utils \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/common-debian.sh

RUN wget https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz && \
    tar Jxf ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz -C /usr/local --strip-component 1 && \
    rm ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz


RUN npm install -g yarn

RUN wget https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz && \
    tar Jxf ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz -C /usr/local --strip-component 1 && \
    rm ffmpeg-n6.0-latest-linux64-gpl-6.0.tar.xz


RUN wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz && \
    tar zxf go1.20.3.linux-amd64.tar.gz -C /usr/local && \
    rm go1.20.3.linux-amd64.tar.gz
RUN echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> /home/vscode/.bashrc

ENV GOPATH=/usr/local/go
RUN /usr/local/go/bin/go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
RUN /usr/local/go/bin/go install github.com/rubenv/sql-migrate/...@latest
RUN /usr/local/go/bin/go install golang.org/x/tools/gopls@latest
RUN /usr/local/go/bin/go install github.com/go-delve/delve/cmd/dlv@latest

RUN pip3 install --force-reinstall "faster-whisper @ https://github.com/guillaumekln/faster-whisper/archive/refs/heads/master.tar.gz"

ENV GOPATH=
USER vscode
