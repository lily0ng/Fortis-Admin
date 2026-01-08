#!/usr/bin/env bash
set -euo pipefail

scan=0
output="text"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --scan)
      scan=1; shift ;;
    --output)
      output="${2:-text}"; shift 2 ;;
    *)
      shift ;;
  esac
done

if [[ $scan -eq 1 && "$output" == "json" ]]; then
  cat <<'JSON'
{
  "timestamp": "2024-01-15T14:30:00Z",
  "servers": [
    {
      "hostname": "web01",
      "ip": "192.168.1.101",
      "os": "Ubuntu 22.04",
      "status": "online",
      "groups": ["webservers", "production"],
      "last_seen": "2024-01-15T14:29:45Z"
    },
    {
      "hostname": "db01",
      "ip": "192.168.1.102",
      "os": "CentOS 8",
      "status": "online",
      "groups": ["databases", "production"],
      "last_seen": "2024-01-15T14:29:50Z"
    }
  ],
  "total": 2,
  "online": 2,
  "offline": 0
}
JSON
  exit 0
fi

echo "[fortis] cluster inventory args: $*"
