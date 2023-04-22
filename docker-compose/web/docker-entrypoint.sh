#!/bin/bash

set -e

[ -f /var/lib/tailscale/dhparam.pem ] || \
  openssl dhparam -out /var/lib/tailscale/dhparam.pem 2048

(
    until tailscale up --hostname ${TS_HOSTNAME} --accept-dns=false ; do
        sleep 1
    done
    tailscale cert --cert-file /var/lib/tailscale/cert.pem --key-file /var/lib/tailscale/key.pem $(tailscale status --json | jq -r .Self.HostName).$(tailscale status --json | jq -r .MagicDNSSuffix)
    sed -e "s/__SERVER_NAME__/$(tailscale status --json | jq -r .Self.HostName).$(tailscale status --json | jq -r .MagicDNSSuffix)/" /etc/nginx/nginx.conf.tmpl > /etc/nginx/nginx.conf
    supervisorctl start nginx
) &

exec /usr/local/bin/supervisord --nodaemon -c /etc/supervisord.conf
