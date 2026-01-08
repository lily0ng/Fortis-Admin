#!/usr/bin/env bash
set -euo pipefail

# Safe-by-default: only lists candidate services unless --disable is used.
# Usage:
#   disable-services.sh --list
#   disable-services.sh --disable sshd --yes

mode="list"
whitelist=""
blacklist=""
disable_name=""
yes=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --list) mode="list"; shift ;;
    --whitelist) whitelist="${2:-}"; shift 2 ;;
    --blacklist) blacklist="${2:-}"; shift 2 ;;
    --disable) mode="disable"; disable_name="${2:-}"; shift 2 ;;
    --yes) yes=1; shift ;;
    *) shift ;;
  esac
done

if [[ "${OSTYPE:-}" != linux* ]]; then
  echo "disable-services: supported on linux only" >&2
  exit 0
fi

if [[ "$mode" == "list" ]]; then
  echo "Candidate services (manual review recommended):"
  systemctl list-unit-files --type=service --state=enabled 2>/dev/null | awk 'NR>1 {print $1}' | sed '/^$/d' | head -n 200
  exit 0
fi

if [[ "$mode" == "disable" ]]; then
  if [[ -z "$disable_name" ]]; then
    echo "--disable requires a service name" >&2
    exit 2
  fi
  if [[ $yes -ne 1 ]]; then
    echo "Refusing to disable service without --yes" >&2
    exit 2
  fi
  systemctl disable --now "$disable_name"
  echo "Disabled: $disable_name"
  exit 0
fi
