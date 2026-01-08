# Fortis-Admin

Fortis-Admin is a modular SysAdmin automation CLI for server hardening, incident response, backup/recovery, and multi-server operations.

## Quick Start

Build:

```bash
go build -o fortis ./cmd/fortis
```

Run:

```bash
./fortis --help
```

## Configuration

Default config path:

- `/etc/fortis/config.yaml`

A starter template is provided in `configs/config.yaml`.

## Scripts

The Go CLI calls Bash scripts from `./scripts` by default.

If you want to override the scripts directory, set it in config:

- `scripts_dir: /path/to/scripts`

## Plugins

Drop executables into `./plugins` and run:

```bash
./fortis plugin list
./fortis plugin run <plugin> [args...]
```
