#!/usr/bin/env bash
set -euo pipefail

# Minimal centralized logging setup stub.
# Safe-by-default: shows config content unless --apply --yes.

apply=0
yes=0
remote=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --apply) apply=1; shift ;;
    --yes) yes=1; shift ;;
    --remote) remote="${2:-}"; shift 2 ;;
    *) shift ;;
  esac
done

if [[ "${OSTYPE:-}" != linux* ]]; then
  echo "configure-logging: supported on linux only" >&2
  exit 0
fi

rsyslog_conf="/etc/rsyslog.d/50-fortis.conf"
logrotate_conf="/etc/logrotate.d/fortis"

rsyslog_body() {
  if [[ -n "$remote" ]]; then
    echo "*.* @@$remote"
  else
    echo "# No remote configured. Set with --remote host:port"
  fi
}

logrotate_body() {
  cat <<'LR'
/var/log/fortis/*.log {
  daily
  rotate 14
  compress
  missingok
  notifempty
  create 0640 root adm
}
LR
}

if [[ $apply -ne 1 ]]; then
  echo "[DRY-RUN] Would write: $rsyslog_conf"
  rsyslog_body
  echo
  echo "[DRY-RUN] Would write: $logrotate_conf"
  logrotate_body
  exit 0
fi

if [[ $yes -ne 1 ]]; then
  echo "Refusing to apply without --yes" >&2
  exit 2
fi

mkdir -p /var/log/fortis
mkdir -p "$(dirname "$rsyslog_conf")" "$(dirname "$logrotate_conf")"
rsyslog_body > "$rsyslog_conf"
logrotate_body > "$logrotate_conf"

systemctl restart rsyslog 2>/dev/null || true

echo "Logging configured."
