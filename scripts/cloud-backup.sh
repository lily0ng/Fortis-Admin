#!/usr/bin/env bash
set -euo pipefail

# Cloud backup helper (S3/Backblaze/Wasabi) - safe-by-default.

provider="s3"     # s3|b2|wasabi
source=""
dest=""
bandwidth=""
lifecycle=""
cost=0
apply=0
yes=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --provider) provider="$2"; shift 2;;
    --source) source="$2"; shift 2;;
    --dest) dest="$2"; shift 2;;
    --bandwidth) bandwidth="$2"; shift 2;;
    --lifecycle) lifecycle="$2"; shift 2;;
    --cost) cost=1; shift;;
    --apply) apply=1; shift;;
    --yes) yes=1; shift;;
    *) shift;;
  esac
 done

if [[ -z "$source" || -z "$dest" ]]; then
  echo "error: --source and --dest are required" >&2
  exit 2
fi

echo "[fortis] cloud-backup"
echo "Provider: $provider"
echo "Source: $source"
echo "Dest: $dest"
echo "Bandwidth: ${bandwidth:-<none>}"
echo "Lifecycle: ${lifecycle:-<none>}"

if [[ $cost -eq 1 ]]; then
  echo "Cost estimation: stub (depends on provider pricing + egress)."
fi

echo "PLAN ONLY (safe-by-default). To apply: --apply --yes"
echo "Recommended tooling: aws s3 sync / rclone / provider CLI with credentials in env."

if [[ $apply -eq 1 && $yes -eq 1 ]]; then
  echo "APPLY requested but not implemented in this helper. Use your approved cloud toolchain." >&2
fi
