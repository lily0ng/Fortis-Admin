#!/usr/bin/env bash
set -euo pipefail

# Safe-by-default parallel SSH executor. No destructive actions by itself; it runs the provided command.
# Reads targets from:
# - --hosts host (repeatable)
# - --file hosts.txt
# - --group <name> (requires inventory file)

inventory_file="${FORTIS_INVENTORY_FILE:-/etc/fortis/inventory.yaml}"
ssh_user="${FORTIS_SSH_USER:-root}"
ssh_port="${FORTIS_SSH_PORT:-22}"
ssh_key="${FORTIS_SSH_KEY:-}"
parallel=4
output="combined" # combined|separate|json
cmd=""

group=""
hosts=()
file=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --inventory-file) inventory_file="$2"; shift 2;;
    --ssh-user) ssh_user="$2"; shift 2;;
    --ssh-port) ssh_port="$2"; shift 2;;
    --ssh-key) ssh_key="$2"; shift 2;;
    --parallel) parallel="$2"; shift 2;;
    --output) output="$2"; shift 2;;
    --group) group="$2"; shift 2;;
    --hosts) hosts+=("$2"); shift 2;;
    --file) file="$2"; shift 2;;
    --command) cmd="$2"; shift 2;;
    *)
      # allow positional command
      if [[ -z "$cmd" ]]; then
        cmd="$1"
      else
        cmd="$cmd $1"
      fi
      shift;;
  esac
 done

if [[ -z "$cmd" ]]; then
  echo "error: --command is required" >&2
  exit 2
fi

# Build target list
if [[ -n "$file" ]]; then
  while IFS= read -r ln; do
    ln="${ln%%#*}"; ln="$(echo "$ln" | xargs)"
    [[ -z "$ln" ]] && continue
    hosts+=("$ln")
  done < "$file"
fi

if [[ -n "$group" && -f "$inventory_file" ]]; then
  # crude yaml parse: match lines under groups: [..]
  while IFS= read -r ln; do
    if [[ "$ln" =~ hostname:\ (.*)$ ]]; then
      current_host="${BASH_REMATCH[1]}"; current_host="$(echo "$current_host" | xargs)"
    fi
    if [[ "$ln" =~ ip:\ (.*)$ ]]; then
      current_ip="${BASH_REMATCH[1]}"; current_ip="$(echo "$current_ip" | xargs)"
    fi
    if [[ "$ln" =~ groups: ]]; then
      if echo "$ln" | grep -q "$group"; then
        if [[ -n "${current_host:-}" ]]; then
          hosts+=("$current_host")
        elif [[ -n "${current_ip:-}" ]]; then
          hosts+=("$current_ip")
        fi
      fi
    fi
  done < "$inventory_file"
fi

# de-dup
uniq_hosts=()
declare -A seen
for h in "${hosts[@]}"; do
  [[ -z "$h" ]] && continue
  if [[ -z "${seen[$h]:-}" ]]; then
    seen[$h]=1
    uniq_hosts+=("$h")
  fi
done

if [[ ${#uniq_hosts[@]} -eq 0 ]]; then
  echo "error: no target hosts" >&2
  exit 2
fi

ssh_args=(-p "$ssh_port" -o BatchMode=yes -o StrictHostKeyChecking=accept-new)
if [[ -n "$ssh_key" ]]; then
  ssh_args+=( -i "$ssh_key" )
fi

run_one() {
  local host="$1"
  local target="${ssh_user}@${host}"
  local out
  if out=$(ssh "${ssh_args[@]}" "$target" "$cmd" 2>&1); then
    echo "[$host] OK"
    echo "$out"
    return 0
  else
    echo "[$host] ERROR" >&2
    echo "$out" >&2
    return 1
  fi
}

# Basic parallelism with a job pool
fail=0
pids=()
for h in "${uniq_hosts[@]}"; do
  while [[ $(jobs -pr | wc -l | xargs) -ge $parallel ]]; do
    wait -n || fail=1
  done
  run_one "$h" &
  pids+=("$!")
done

for p in "${pids[@]}"; do
  wait "$p" || fail=1
 done

exit $fail
