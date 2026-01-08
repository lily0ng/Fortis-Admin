#!/usr/bin/env bash
set -euo pipefail

config_path=""
target_path=""
validate=0
backup=0
rollback=""
diff=0
apply=0
yes=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --config) config_path="$2"; shift 2;;
    --target) target_path="$2"; shift 2;;
    --validate) validate=1; shift;;
    --backup) backup=1; shift;;
    --rollback) rollback="$2"; shift 2;;
    --diff) diff=1; shift;;
    --apply) apply=1; shift;;
    --yes) yes=1; shift;;
    *) shift;;
  esac
done

if [[ -z "$config_path" || -z "$target_path" ]]; then
  echo "error: --config and --target are required" >&2
  exit 2
fi

echo "[fortis] config deploy"
echo "Config: $config_path"
echo "Target: $target_path"

if [[ $validate -eq 1 ]]; then
  echo "Validate: enabled"
  if [[ -d "$config_path" ]]; then
    echo "Validation: directory exists"
  else
    if [[ -f "$config_path" ]]; then
      echo "Validation: file exists"
    else
      echo "Validation: config path not found" >&2
      exit 2
    fi
  fi
fi

if [[ $diff -eq 1 ]]; then
  echo "Diff: requested"
  echo "Drift detection: not implemented (hint: compare remote $target_path with version-controlled config)"
fi

if [[ $backup -eq 1 ]]; then
  echo "Backup: requested"
  echo "Backup plan: would copy existing remote $target_path to $target_path.bak.<timestamp>"
fi

if [[ -n "$rollback" ]]; then
  echo "Rollback: requested to $rollback"
  echo "Rollback is a stub. Supply a previous backup artifact and apply via cluster exec or scp tooling."
fi

if [[ $apply -eq 0 || $yes -eq 0 ]]; then
  echo "PLAN ONLY (safe-by-default). To apply: --apply --yes"
  echo "Recommended execution:"
  echo "  1) fortis cluster exec --command \"sudo mkdir -p $(dirname "$target_path")\""
  echo "  2) Use scp/rsync to stage config to each host"
  echo "  3) fortis cluster exec --command \"sudo install -m 0644 <staged> $target_path\""
  exit 0
fi

echo "APPLY requested (--apply --yes)"
echo "This script does not push to remote hosts itself. Use 'fortis cluster exec' + your transfer tool for controlled deployment."
exit 0
