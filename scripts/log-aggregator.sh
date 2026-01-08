#!/usr/bin/env bash
set -euo pipefail

# Centralized log collection helper.
# Safe-by-default: supports listing/streaming commands; collection is opt-in.

mode="list"   # list|stream|collect
output=""
pattern=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --list) mode="list"; shift;;
    --stream) mode="stream"; shift;;
    --collect) mode="collect"; shift;;
    --output) output="$2"; shift 2;;
    --pattern) pattern="$2"; shift 2;;
    *) shift;;
  esac
 done

echo "[fortis] log-aggregator ($mode)"

echo "This script is a helper. Use 'fortis cluster exec' to run remote commands and pipe results into a central store."

echo "Examples:"
echo "  fortis cluster exec --group webservers --command 'tail -n 200 /var/log/syslog'"
echo "  fortis cluster exec --hosts web01 --command 'journalctl -n 200 --no-pager'"

echo "Pattern: ${pattern:-<none>}"

if [[ "$mode" == "collect" ]]; then
  if [[ -z "$output" ]]; then
    echo "error: --output is required for --collect" >&2
    exit 2
  fi
  mkdir -p "$(dirname "$output")"
  {
    echo "collection not implemented in-script (use cluster exec + redirect)"
    echo "timestamp=$(date -u +%FT%TZ)"
    echo "pattern=${pattern:-}"
  } > "$output"
  echo "Wrote collection stub: $output"
fi
