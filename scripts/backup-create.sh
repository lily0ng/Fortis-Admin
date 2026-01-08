#!/usr/bin/env bash
set -euo pipefail

if [[ "${FORTIS_QUIET:-}" == "1" ]]; then
  echo "Backup completed successfully: backup-20240115-143000.tar.gz"
  echo "Size: 4.2GB, Duration: 2m15s, Integrity: Verified"
  exit 0
fi

echo "[fortis] backup create args: $*"
