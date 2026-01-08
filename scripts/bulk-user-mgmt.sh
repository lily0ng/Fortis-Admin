#!/usr/bin/env bash
set -euo pipefail

# Cross-server user management helper.
# Safe-by-default: prints plan unless --apply AND --yes are set.

apply=0
yes=0
user=""
groups=()
password_policy=""
report=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --apply) apply=1; shift;;
    --yes) yes=1; shift;;
    --user) user="$2"; shift 2;;
    --group) groups+=("$2"); shift 2;;
    --password-policy) password_policy="$2"; shift 2;;
    --report) report="$2"; shift 2;;
    *) shift;;
  esac
 done

if [[ -z "$user" ]]; then
  echo "error: --user is required" >&2
  exit 2
fi

echo "[fortis] bulk-user-mgmt"
echo "User: $user"
echo "Groups: ${groups[*]:-}" 
echo "Password policy: ${password_policy:-}" 

if [[ $apply -eq 0 || $yes -eq 0 ]]; then
  echo "PLAN ONLY (safe-by-default). To apply: --apply --yes"
  echo "Would ensure user exists and groups match across target hosts (execution handled externally)."
else
  echo "APPLY requested. This script is a helper; integrate with your cluster exec tooling."
  echo "Recommended: fortis cluster exec --command 'id -u $user || useradd $user'"
fi

if [[ -n "$report" ]]; then
  mkdir -p "$(dirname "$report")"
  {
    echo "bulk-user-mgmt report"
    echo "user=$user"
    echo "groups=${groups[*]:-}"
    echo "applied=$apply"
  } > "$report"
  echo "Report saved: $report"
fi
