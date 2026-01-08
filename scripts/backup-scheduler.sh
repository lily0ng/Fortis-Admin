#!/usr/bin/env bash
set -euo pipefail

store="${FORTIS_SCHEDULE_STORE:-/etc/fortis/backup-schedules.txt}"

add=""
list=0
remove=""
enable=""
disable=""
run_now=""

retention="${FORTIS_RETENTION:-}"
pre_hook=""
post_hook=""
notify=""  # email/slack stub

apply=0
yes=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --add) add="$2"; shift 2;;
    --list) list=1; shift;;
    --remove) remove="$2"; shift 2;;
    --enable) enable="$2"; shift 2;;
    --disable) disable="$2"; shift 2;;
    --run-now) run_now="$2"; shift 2;;
    --retention) retention="$2"; shift 2;;
    --pre-hook) pre_hook="$2"; shift 2;;
    --post-hook) post_hook="$2"; shift 2;;
    --notify) notify="$2"; shift 2;;
    --apply) apply=1; shift;;
    --yes) yes=1; shift;;
    *) shift;;
  esac
done

mkdir -p "$(dirname "$store")"
touch "$store"

echo "[fortis] backup scheduler"
echo "Store: $store"

if [[ $list -eq 1 ]]; then
  echo "Schedules:"
  cat "$store" || true
  exit 0
fi

if [[ -n "$run_now" ]]; then
  echo "Run-now requested for schedule id: $run_now"
  echo "This is a stub. Use your cron/systemd timer or run 'fortis backup create ...' directly."
  exit 0
fi

if [[ $apply -eq 0 || $yes -eq 0 ]]; then
  echo "PLAN ONLY (safe-by-default). To apply: --apply --yes"
fi

if [[ -n "$add" ]]; then
  echo "Add schedule: $add"
  echo "Retention: ${retention:-<none>}"
  echo "Pre-hook: ${pre_hook:-<none>}"
  echo "Post-hook: ${post_hook:-<none>}"
  echo "Notify: ${notify:-<none>}"
  echo "Cron guidance (example):"
  echo "  0 2 * * * fortis backup create --source /etc --target /backups"
  if [[ $apply -eq 1 && $yes -eq 1 ]]; then
    echo "ENABLED\t$(date -u +%Y%m%d%H%M%S)\t$add" >> "$store"
    echo "Saved."
  fi
  exit 0
fi

if [[ -n "$remove" ]]; then
  echo "Remove schedule id contains: $remove"
  if [[ $apply -eq 1 && $yes -eq 1 ]]; then
    tmp="${store}.tmp"
    grep -v "$remove" "$store" > "$tmp" || true
    mv "$tmp" "$store"
    echo "Removed."
  fi
  exit 0
fi

if [[ -n "$enable" ]]; then
  echo "Enable schedule id contains: $enable"
  echo "Stub: stored schedules are informational; use cron/systemd timers to enable/disable."
  exit 0
fi

if [[ -n "$disable" ]]; then
  echo "Disable schedule id contains: $disable"
  echo "Stub: stored schedules are informational; use cron/systemd timers to enable/disable."
  exit 0
fi

echo "No action specified. Use --list or --add."
exit 0
