#!/usr/bin/env bash
set -euo pipefail

# Safe-by-default network forensics helper.
# This script does NOT start packet captures automatically.
# It prints useful network context and instructions.

echo "[fortis] network forensics (safe mode)"

echo "\n=== Interfaces ==="
(ip link 2>/dev/null || ifconfig -a 2>/dev/null || true) | head -n 200

echo "\n=== Listening Ports ==="
(ss -tulpn 2>/dev/null || netstat -tulpn 2>/dev/null || true) | head -n 200

echo "\n=== Routes ==="
(ip route 2>/dev/null || netstat -rn 2>/dev/null || true) | head -n 200

echo "\nTo capture packets, consider running (requires root):"

echo "  sudo tcpdump -i <iface> -w capture.pcap"
