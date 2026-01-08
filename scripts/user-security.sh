#!/usr/bin/env bash
set -euo pipefail

lock_inactive=0
password_policy=0
sudo_secure=0
audit=0
yes=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --lock-inactive) lock_inactive=1; shift;;
    --password-policy) password_policy=1; shift;;
    --sudo-secure) sudo_secure=1; shift;;
    --audit) audit=1; shift;;
    --yes) yes=1; shift;;
    *) shift;;
  esac
done

echo "[fortis] user security"

if [[ $audit -eq 1 ]]; then
  echo "Audit report (best-effort):"
  echo "- sudoers files: /etc/sudoers and /etc/sudoers.d/*"
  echo "- passwd policy: /etc/login.defs and PAM configuration"
  echo "- inactive lock: chage defaults"
  exit 0
fi

plan=()
if [[ $lock_inactive -eq 1 ]]; then
  plan+=("Set inactive account lock policy (chage -I 30) (stub)")
fi
if [[ $password_policy -eq 1 ]]; then
  plan+=("Set PASS_MIN_LEN, PASS_MAX_DAYS in /etc/login.defs (best-effort)")
  plan+=("Enable pam_pwquality if available (not enforced automatically)")
fi
if [[ $sudo_secure -eq 1 ]]; then
  plan+=("Harden sudoers: require tty, log output (stub)")
  plan+=("Set session timeout in /etc/profile.d/fortis-timeout.sh")
fi

if [[ ${#plan[@]} -eq 0 ]]; then
  echo "No action requested. Use flags like --password-policy or --sudo-secure." >&2
  exit 2
fi

echo "Plan:"
for p in "${plan[@]}"; do
  echo "- $p"
done

if [[ $yes -ne 1 ]]; then
  echo "[DRY-RUN] Re-run with 'fortis harden users ... --yes' to apply."
  exit 0
fi

if [[ $password_policy -eq 1 ]]; then
  login_defs="/etc/login.defs"
  if [[ -f "$login_defs" ]]; then
    cp "$login_defs" "$login_defs.fortis.bak.$(date -u +%Y%m%d%H%M%S)"
    perl -0777 -pe 's/^\s*PASS_MIN_LEN\s+\d+/PASS_MIN_LEN\t12/gm; s/^\s*PASS_MAX_DAYS\s+\d+/PASS_MAX_DAYS\t90/gm' -i "$login_defs" || true
  fi
fi

if [[ $sudo_secure -eq 1 ]]; then
  out="/etc/profile.d/fortis-timeout.sh"
  cp "$out" "$out.fortis.bak.$(date -u +%Y%m%d%H%M%S)" 2>/dev/null || true
  cat > "$out" <<'EOF'
export TMOUT=900
readonly TMOUT
EOF
  chmod 0644 "$out" || true
fi

echo "Applied user security changes (best-effort)."
exit 0
