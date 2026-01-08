#!/usr/bin/env bash
set -euo pipefail

disable_root=0
key_only=0
port=""
banner=""
yes=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --disable-root) disable_root=1; shift;;
    --key-only) key_only=1; shift;;
    --port) port="$2"; shift 2;;
    --banner) banner="$2"; shift 2;;
    --yes) yes=1; shift;;
    *) shift;;
  esac
done

verbose="${FORTIS_VERBOSE:-0}"

echo "[fortis] harden ssh"

plan=()
if [[ $disable_root -eq 1 ]]; then
  plan+=("Set PermitRootLogin no")
fi
if [[ $key_only -eq 1 ]]; then
  plan+=("Set PasswordAuthentication no")
  plan+=("Set PubkeyAuthentication yes")
fi
if [[ -n "$port" ]]; then
  plan+=("Set Port $port")
fi
if [[ -n "$banner" ]]; then
  plan+=("Set Banner $banner")
fi

if [[ ${#plan[@]} -eq 0 ]]; then
  echo "No changes requested. Use flags like --disable-root or --key-only." >&2
  exit 2
fi

echo "Plan:"
for p in "${plan[@]}"; do
  echo "- $p"
done

if [[ $yes -ne 1 ]]; then
  echo "[DRY-RUN] Re-run with 'fortis harden ssh ... --yes' to apply."
  exit 0
fi

sshd_config="/etc/ssh/sshd_config"
backup="$sshd_config.fortis.bak.$(date -u +%Y%m%d%H%M%S)"

if [[ ! -f "$sshd_config" ]]; then
  echo "error: $sshd_config not found" >&2
  exit 1
fi

cp "$sshd_config" "$backup"
tmp="$sshd_config.tmp"
cp "$sshd_config" "$tmp"

set_kv() {
  local key="$1" val="$2"
  if grep -qiE "^\s*${key}\b" "$tmp"; then
    perl -0777 -pe "s/^\s*${key}\\b.*$/${key} ${val}/gim" -i "$tmp"
  else
    echo "${key} ${val}" >> "$tmp"
  fi
}

if [[ $disable_root -eq 1 ]]; then
  set_kv "PermitRootLogin" "no"
fi
if [[ $key_only -eq 1 ]]; then
  set_kv "PasswordAuthentication" "no"
  set_kv "PubkeyAuthentication" "yes"
fi
if [[ -n "$port" ]]; then
  set_kv "Port" "$port"
fi
if [[ -n "$banner" ]]; then
  set_kv "Banner" "$banner"
fi

mv "$tmp" "$sshd_config"

if [[ "$verbose" == "1" ]]; then
  echo "Backup saved: $backup"
  echo "Updated: $sshd_config"
  echo "Note: restart sshd manually if appropriate (not performed automatically)."
fi

echo "Applied SSH hardening changes (restart not automatic)."
exit 0
