<div align="center">

# FORTIS-ADMIN

<img alt="FORTIS-ADMIN" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=120&section=header&text=FORTIS-ADMIN&fontColor=ffffff&fontSize=46&animation=fadeIn" />

<p><b>Server Hardening</b> · <b>Incident Response</b> · <b>Backups</b> · <b>Cluster Ops</b></p>

<p>
  <img alt="Go" src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white" />
  <img alt="CLI" src="https://img.shields.io/badge/CLI-Cobra-0A0A0A?logo=terminal&logoColor=white" />
  <img alt="Safe" src="https://img.shields.io/badge/Safe--By--Default-yes-2ea44f" />
  <img alt="Linux" src="https://img.shields.io/badge/Linux-Supported-111827?logo=linux&logoColor=white" />
  <img alt="macOS" src="https://img.shields.io/badge/macOS-Build/Test-111827?logo=apple&logoColor=white" />
  <img alt="Windows" src="https://img.shields.io/badge/Windows-Limited-111827?logo=windows&logoColor=white" />
  <img alt="License" src="https://img.shields.io/badge/License-Open%20Source-blue" />
</p>

<p>
  <img alt="Linux" src="https://img.icons8.com/color/48/linux.png" height="28" />
  <img alt="macOS" src="https://img.icons8.com/color/48/mac-os.png" height="28" />
  <img alt="Windows" src="https://img.icons8.com/color/48/windows-10.png" height="28" />
</p>

<p>
  <img alt="Creator" src="https://img.shields.io/badge/Creator-Lily%20Yang-FF2D95" />
</p>

<p>
  <a href="#quickstart">Quickstart</a> ·
  <a href="#features">Features</a> ·
  <a href="#system-design">System Design</a> ·
  <a href="#commands--examples">Commands & Examples</a> ·
  <a href="#community--support">Community</a>
</p>

</div>


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## 🧭 About

<img alt="ABOUT" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=ABOUT&fontColor=00E5FF&fontSize=26" width="100%" />

FORTIS-ADMIN is a modular SysAdmin automation CLI for:

- Server hardening automation
- Incident response toolkit
- Backup & recovery
- Multi-server (cluster) management

It combines a Go-based CLI with safe-by-default Bash helpers.

**Creator:** Lily Yang

<table>
  <tr>
    <td><b>🛡️ Hardening</b><br/>CIS-style audits, firewall/kernel tuning, SSH/user policy helpers</td>
    <td><b>🕵️ Incident</b><br/>Capture, triage, analyze, timeline, reports</td>
  </tr>
  <tr>
    <td><b>💾 Backup</b><br/>Create/verify/restore, catalogs, schedulers (safe)</td>
    <td><b>🖧 Cluster</b><br/>Inventory, exec, monitoring, patch planning</td>
  </tr>
</table>

> **Safety model:** Commands that can change system state require explicit confirmation (e.g., `--yes`, `--apply`).

**Design principles**

- **Safe-by-default**: actions that can change system state require explicit confirmation (e.g., `--yes`, `--apply`)
- **Portable**: minimal dependencies, works well over SSH
- **Extensible**: add scripts and/or Go modules without rewriting the CLI


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## 🗂️ Table of Contents

<img alt="TABLE OF CONTENTS" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=TABLE%20OF%20CONTENTS&fontColor=00E5FF&fontSize=26" width="100%" />

- [Quickstart](#quickstart)
- [Features](#features)
- [System Design](#system-design)
- [Commands & Examples](#commands--examples)
- [OS Support](#os-support)
- [Git Status (Real Time)](#git-status-real-time)
- [Community & Support](#community--support)
- [License](#license)
- [Plugins](#plugins)


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## ⚡ Quickstart

<img alt="QUICKSTART" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=QUICKSTART&fontColor=00E5FF&fontSize=26" width="100%" />

### Build

```bash
go mod tidy
go build -o fortis ./cmd/fortis
```

### Run

```bash
./fortis --help
./fortis --version
```

### Configuration

- **Default config path:** `/etc/fortis/config.yaml`
- **Template:** `configs/config.yaml`


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## ✨ Features

<img alt="FEATURES" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=FEATURES&fontColor=00E5FF&fontSize=26" width="100%" />

<details open>
<summary><b>🛡️ Server Hardening Automation</b></summary>

- `fortis harden audit` (Go): security audit, scoring, JSON/YAML/HTML output
- `fortis harden apply` (Go): profile application with dry-run + rollback
- `fortis harden firewall` (Go): firewall planning and apply gate with `--yes`
- `fortis harden kernel` (Go): sysctl planning and apply gate with `--yes`
- `fortis harden filesystem` (Go): filesystem permission audits
- `fortis harden package-audit` (Go): package inventory/audit
- `fortis harden ssh` (Bash): safe-by-default SSH hardening helper
- `fortis harden users` (Bash): safe-by-default password/sudo/session policy helper
- `fortis harden compliance` (Go): compliance reports (JSON/HTML; PDF intentionally stubbed)

</details>

<details>
<summary><b>🕵️ Incident Response Toolkit</b></summary>

- `fortis incident capture` (Go): evidence capture + integrity manifest + compression
- `fortis incident triage` (Go): quick/full triage with `--output`
- `fortis incident analyze` (Go): IOC matching and analysis report
- `fortis incident timeline` (Go): timeline generation (csv/json/html)
- `fortis incident report` (Go): markdown report generator
- `fortis incident ioc` (Go): local IOC store manager
- `fortis incident logs analyze` (Go): log analysis with IOC matching

</details>

<details>
<summary><b>💾 Backup & Recovery</b></summary>

- `fortis backup create` (Go): tar/tar.gz backups + `.meta.json` sidecar + SHA256
- `fortis backup list` (Go): lists backups from sidecar metadata
- `fortis backup verify` (Go): checksum validation + optional restore simulation (`--full`)
- `fortis backup restore` (Go): restore archives
- `fortis backup catalog` (Go): list/search archive contents

Advanced (hidden from `--help` to keep the CLI surface minimal):

- `fortis backup snapshot` (Go): snapshot manager planner (backend detection; apply stub)
- `fortis backup restore-wizard` (Go): interactive restore wizard (dry-run by default)

</details>

<details>
<summary><b>🖧 Multi-Server (Cluster) Management</b></summary>

- `fortis cluster init` (Go): creates an inventory template + SSH guidance (safe)
- `fortis cluster inventory` (Go): read inventory YAML and output text/json
- `fortis cluster exec` (Go): parallel SSH command execution (text/json output)
- `fortis cluster monitor` (Go): basic metrics + health score + JSON export
- `fortis cluster patch` (Go): patch plan (dry-run) and gated apply (`--apply --yes`)

</details>


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## 🧩 System Design

<img alt="SYSTEM DESIGN" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=SYSTEM%20DESIGN&fontColor=00E5FF&fontSize=26" width="100%" />

### High-level architecture

```mermaid
flowchart LR
  U[Operator / CI] --> CLI[fortis CLI (Go)]
  CLI -->|Safe core| GO[internal/* (Go engines)]
  CLI -->|OS helpers| SH[scripts/*.sh]
  GO --> FS[(Local filesystem)]
  SH --> FS
  CLI -->|SSH fan-out| R[(Remote nodes)]
```

### Safety model (why `--yes` exists)

```mermaid
flowchart TD
  A[Command invoked] --> B{Potentially destructive?}
  B -- no --> C[Run normally]
  B -- yes --> D{Explicit confirmation provided?}
  D -- no --> E[Plan / dry-run output]
  D -- yes --> F[Apply changes]
```

### Data flow: backups

```mermaid
sequenceDiagram
  participant U as Operator
  participant C as fortis backup (Go)
  participant F as Filesystem
  U->>C: backup create --source ... --target ...
  C->>F: write archive (.tar/.tar.gz)
  C->>F: compute sha256
  C->>F: write .meta.json
  U->>C: backup verify --full
  C->>F: restore simulation (temp dir)
```

### Data flow: incident response

```mermaid
flowchart TD
  A[incident capture] --> B[evidence folder]
  B --> C[triage]
  B --> D[analyze]
  B --> E[timeline]
  C --> F[report]
  D --> F
  E --> F
```

### Data flow: cluster execution

```mermaid
sequenceDiagram
  participant U as Operator
  participant C as fortis cluster (Go)
  participant I as Inventory (YAML)
  participant H as Hosts
  U->>C: cluster exec --group webservers --command "..."
  C->>I: load + filter hosts
  C->>H: SSH fan-out
  H-->>C: stdout/stderr + status
  C-->>U: combined/separate/json output
```


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## 📘 Command Usage Guide

<img alt="COMMAND USAGE" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=COMMAND%20USAGE&fontColor=00E5FF&fontSize=26" width="100%" />

### Safe-by-default flags

- `--yes`:
  - Used to confirm actions that may change system state.
  - Common in hardening and cluster operations.
- `--apply`:
  - Used when a command supports generating a plan first and applying later.
  - Common in patching/snapshots.

### Output formats

- Prefer `--output json` (or similar flags) for automation.
- Use `--verbose` for interactive debugging.
- Use `--quiet` for CI/log-friendly output.

### Inventory-driven operations

Cluster commands typically take inventory/SSH flags:

- `--inventory-file` (defaults to `/etc/fortis/inventory.yaml`)
- `--ssh-user`, `--ssh-port`, `--ssh-key`

### Scripts vs Go

Some commands call scripts in `./scripts` for OS-level configuration changes.

- Keep scripts audited and version-controlled.
- Run on a test host first.


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## ⌨️ Commands & Examples

<img alt="COMMANDS & EXAMPLES" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=COMMANDS%20%26%20EXAMPLES&fontColor=00E5FF&fontSize=26" width="100%" />

### Hardening: audit + apply (safe gate)

```bash
./fortis harden audit --profile cis --output ./audit.html
./fortis harden firewall --profile webserver
./fortis harden firewall --profile webserver --yes
```

Sample output (illustrative):

```text
Audit complete. Score: 78/100 (Medium). Report: ./audit.html
```

### Incident response: capture + triage

```bash
./fortis incident capture --case incident-001 --type all --integrity --compress
./fortis incident triage --full --output ./evidence/triage.txt
```

Sample output (illustrative):

```json
{
  "case_id": "incident-001",
  "output_dir": "./evidence/incident-001",
  "timestamp": "2024-01-15T14:30:00Z",
  "notes": ["OS=linux", "HOST=web01"]
}
```

### Backup: create + verify

```bash
./fortis backup create --source /etc --target ./backups
./fortis backup verify --backup ./backups/backup-YYYYMMDD-HHMMSS.tar.gz --full
```

Sample output (illustrative):

```json
{
  "backup_path": "./backups/backup-20240115-143000.tar.gz",
  "ok": true,
  "reason": "checksum validated",
  "sha256": "<sha256>"
}
```

### Cluster: inventory + exec

```bash
./fortis cluster init --inventory-file ./configs/inventory.yaml
./fortis cluster inventory --inventory-file ./configs/inventory.yaml --output json
./fortis cluster exec --inventory-file ./configs/inventory.yaml --group webservers --command "uname -a" --output json
```


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## 🧪 OS Support

<img alt="OS SUPPORT" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=OS%20SUPPORT&fontColor=00E5FF&fontSize=26" width="100%" />

- Linux (primary target)
- macOS: build/test support; many hardening/forensic actions are Linux-oriented
- Windows: not currently supported for most workflows


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## 🧷 Git Status (Real Time)

<img alt="GIT STATUS" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=GIT%20STATUS&fontColor=00E5FF&fontSize=26" width="100%" />

```bash
git status -sb
git diff
git diff --stat
```

Live refresh (macOS/Linux):

```bash
watch -n 1 git status -sb
```


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## 🤝 Community & Support

<img alt="COMMUNITY" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=COMMUNITY&fontColor=00E5FF&fontSize=26" width="100%" />

- **Documentation:** https://fortis-admin.readthedocs.io
- **Source:** https://github.com/lily0ng/Fortis-Admin
- **Issues:** https://github.com/lily0ng/Fortis-Admin/issues
- **Discussions:** GitHub Discussions (recommended)
- **Issue reports:** include OS/version + `fortis --version` + command output


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## 🧾 License

<img alt="LICENSE" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=LICENSE&fontColor=00E5FF&fontSize=26" width="100%" />

Open source.

- If your repository includes a `LICENSE` file, that is the authoritative license.
- If not, add one (MIT/Apache-2.0/GPL-3.0 are common choices).


<img alt="divider" src="https://capsule-render.vercel.app/api?type=rect&color=0:00E5FF,100:8A2BE2&height=3&section=header&text=" width="100%" />

## 🧩 Plugins

<img alt="PLUGINS" src="https://capsule-render.vercel.app/api?type=rect&color=0:0B1220,100:8A2BE2&height=56&section=header&text=PLUGINS&fontColor=00E5FF&fontSize=26" width="100%" />

Drop executables into `./plugins` and run:

```bash
./fortis plugin list
./fortis plugin run <plugin> [args...]
```

<br/>

<div align="center">
  <img alt="footer" src="https://capsule-render.vercel.app/api?type=rect&color=0:8A2BE2,100:00E5FF&height=90&section=footer&text=Created%20by%20Lily%20Yang&fontColor=ffffff&fontSize=22&animation=fadeIn" width="100%" />
</div>
