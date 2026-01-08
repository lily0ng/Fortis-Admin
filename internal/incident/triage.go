package incident

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func Triage(ctx context.Context, opts TriageOptions) (string, error) {
	_ = ctx
	out := opts.Output
	if out == "" {
		out = filepath.Join(".", fmt.Sprintf("triage-%s.txt", time.Now().Format("20060102-150405")))
	}
	_ = os.MkdirAll(filepath.Dir(out), 0o755)

	b := &strings.Builder{}
	b.WriteString("FORTIS TRIAGE REPORT\n")
	b.WriteString("Generated: " + time.Now().Format(time.RFC3339) + "\n")
	host, _ := os.Hostname()
	b.WriteString("Host: " + host + "\n")
	b.WriteString("Platform: " + runtime.GOOS + "\n\n")

	if runtime.GOOS == "linux" {
		b.WriteString("=== SYSTEM ===\n")
		b.WriteString(run("bash", "-lc", "uptime 2>/dev/null || true"))
		b.WriteString(run("bash", "-lc", "df -h 2>/dev/null | head -n 50 || true"))
		b.WriteString("\n")
	}

	if opts.Quick || (!opts.Full && !opts.Quick) {
		b.WriteString("=== QUICK ===\n")
		if runtime.GOOS == "linux" {
			b.WriteString(run("bash", "-lc", "who 2>/dev/null || true"))
			b.WriteString(run("bash", "-lc", "last -n 20 2>/dev/null || true"))
		}
		b.WriteString("\n")
	}

	if opts.Full {
		b.WriteString("=== FULL ===\n")
		if runtime.GOOS == "linux" {
			b.WriteString(run("bash", "-lc", "systemctl --failed 2>/dev/null || true"))
			b.WriteString(run("bash", "-lc", "journalctl -n 200 --no-pager 2>/dev/null || true"))
		}
		b.WriteString("\n")
	}

	if opts.Processes {
		b.WriteString("=== PROCESSES ===\n")
		if runtime.GOOS == "linux" {
			b.WriteString(run("bash", "-lc", "ps auxww 2>/dev/null | head -n 300 || true"))
		}
		b.WriteString("\n")
	}

	if opts.Network {
		b.WriteString("=== NETWORK ===\n")
		if runtime.GOOS == "linux" {
			b.WriteString(run("bash", "-lc", "ss -tulpn 2>/dev/null | head -n 200 || netstat -tulpn 2>/dev/null | head -n 200 || true"))
		}
		b.WriteString("\n")
	}

	if opts.Persistence {
		b.WriteString("=== PERSISTENCE ===\n")
		if runtime.GOOS == "linux" {
			b.WriteString(run("bash", "-lc", "crontab -l 2>/dev/null || true"))
			b.WriteString(run("bash", "-lc", "ls -la /etc/cron.* 2>/dev/null | head -n 200 || true"))
		}
		b.WriteString("\n")
	}

	if err := os.WriteFile(out, []byte(b.String()), 0o600); err != nil {
		return "", err
	}
	return out, nil
}

func run(name string, args ...string) string {
	p := filepath.Join(os.TempDir(), fmt.Sprintf("fortis-%d.txt", time.Now().UnixNano()))
	_ = writeCmdOutput(p, name, args...)
	b, _ := os.ReadFile(p)
	_ = os.Remove(p)
	return string(b)
}
