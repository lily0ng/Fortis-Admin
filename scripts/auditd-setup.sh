#!/usr/bin/env bash
set -euo pipefail

# Minimal auditd rule deployment (linux).
# Safe-by-default: prints what it would do unless --apply --yes.

apply=0
yes=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --apply) apply=1; shift ;;
    --yes) yes=1; shift ;;
    *) shift ;;
  esac
done

if [[ "${OSTYPE:-}" != linux* ]]; then
  echo "auditd-setup: supported on linux only" >&2
  exit 0
fi

rules_path="/etc/audit/rules.d/fortis.rules"

cat_rules() {
  cat <<'RULES'
-w /etc/passwd -p wa -k identity
-w /etc/shadow -p wa -k identity
-w /etc/group  -p wa -k identity
-w /etc/sudoers -p wa -k scope
-w /var/log/ -p wa -k logs
RULES
}

if [[ $apply -ne 1 ]]; then
  echo "[DRY-RUN] Would write rules to: $rules_path"
  cat_rules
  exit 0
fi

if [[ $yes -ne 1 ]]; then
  echo "Refusing to apply without --yes" >&2
  exit 2
fi

mkdir -p "$(dirname "$rules_path")"
cat_rules > "$rules_path"
chmod 0640 "$rules_path"

if command -v augenrules >/dev/null 2>&1; then
  augenrules --load || true
elif command -v auditctl >/dev/null 2>&1; then
  auditctl -R "$rules_path" || true
fi

echo "auditd rules deployed: $rules_path"
