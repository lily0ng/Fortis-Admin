#!/usr/bin/env bash
set -euo pipefail

# Transaction-safe database backup helper.
# Safe-by-default: prints a plan unless --apply --yes.

db=""          # mysql|postgres|mongo
host="localhost"
port=""
user=""
name=""
out=""
rotate=7
apply=0
yes=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --db) db="$2"; shift 2;;
    --host) host="$2"; shift 2;;
    --port) port="$2"; shift 2;;
    --user) user="$2"; shift 2;;
    --name) name="$2"; shift 2;;
    --output) out="$2"; shift 2;;
    --rotate) rotate="$2"; shift 2;;
    --apply) apply=1; shift;;
    --yes) yes=1; shift;;
    *) shift;;
  esac
 done

if [[ -z "$db" || -z "$name" ]]; then
  echo "error: --db and --name are required" >&2
  exit 2
fi

if [[ -z "$out" ]]; then
  out="./db-backups/${db}-${name}-$(date -u +%Y%m%d-%H%M%S).dump"
fi

mkdir -p "$(dirname "$out")"

echo "[fortis] database-backup"
echo "DB: $db"
echo "Host: $host"
echo "Name: $name"
echo "Output: $out"
echo "Rotate keep: $rotate"

echo "PLAN ONLY (safe-by-default). To apply: --apply --yes"
echo "Would run transaction-safe dump with credentials from env/secret store."
echo "Examples:"
echo "  postgres: pg_dump --format=custom --file=$out $name"
echo "  mysql: mysqldump --single-transaction $name > $out"
echo "  mongo: mongodump --archive=$out --db=$name"

echo "Rotation: not implemented (stub); would delete older than $rotate backups."

if [[ $apply -eq 1 && $yes -eq 1 ]]; then
  echo "APPLY requested but not implemented in this helper. Use native tools with proper secrets." >&2
fi
