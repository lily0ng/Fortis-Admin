package hardening

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Result string

const (
	ResultPass Result = "pass"
	ResultFail Result = "fail"
	ResultWarn Result = "warn"
	ResultSkip Result = "skip"
)

type Finding struct {
	ID             string `json:"id" yaml:"id"`
	Title          string `json:"title" yaml:"title"`
	Result         Result `json:"result" yaml:"result"`
	Details        string `json:"details,omitempty" yaml:"details,omitempty"`
	Recommendation string `json:"recommendation,omitempty" yaml:"recommendation,omitempty"`
	Weight         int    `json:"weight" yaml:"weight"`
}

type Report struct {
	Timestamp   time.Time `json:"timestamp" yaml:"timestamp"`
	Profile     string    `json:"profile" yaml:"profile"`
	Level       string    `json:"level" yaml:"level"`
	Hostname    string    `json:"hostname" yaml:"hostname"`
	OS          string    `json:"os" yaml:"os"`
	Platform    string    `json:"platform" yaml:"platform"`
	Findings    []Finding `json:"findings" yaml:"findings"`
	Passed      int       `json:"passed" yaml:"passed"`
	Failed      int       `json:"failed" yaml:"failed"`
	Warnings    int       `json:"warnings" yaml:"warnings"`
	Skipped     int       `json:"skipped" yaml:"skipped"`
	Score       int       `json:"score" yaml:"score"`
	ScoreLabel  string    `json:"score_label" yaml:"score_label"`
	ReportHash  string    `json:"report_hash" yaml:"report_hash"`
	ReportPath  string    `json:"report_path" yaml:"report_path"`
	GeneratedBy string    `json:"generated_by" yaml:"generated_by"`
}

type AuditOptions struct {
	Profile string
	Level   string
	Output  string
	Fix     bool
	Yes     bool
	Verbose bool
}

func RunAudit(ctx context.Context, opts AuditOptions) (Report, error) {
	if opts.Profile == "" {
		opts.Profile = "cis"
	}
	if opts.Level == "" {
		opts.Level = "basic"
	}

	host, _ := os.Hostname()
	rep := Report{
		Timestamp:   time.Now(),
		Profile:     opts.Profile,
		Level:       opts.Level,
		Hostname:    host,
		OS:          runtime.GOOS,
		Platform:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		GeneratedBy: "fortis",
	}

	checks := defaultChecks()
	for _, c := range checks {
		f, err := c(ctx, opts)
		if err != nil {
			f.Result = ResultWarn
			f.Details = err.Error()
		}
		rep.Findings = append(rep.Findings, f)
	}

	totalWeight := 0
	scoreWeight := 0
	for _, f := range rep.Findings {
		if f.Result == ResultSkip {
			rep.Skipped++
			continue
		}
		totalWeight += f.Weight
		switch f.Result {
		case ResultPass:
			rep.Passed++
			scoreWeight += f.Weight
		case ResultFail:
			rep.Failed++
		case ResultWarn:
			rep.Warnings++
			// partial credit for warnings
			scoreWeight += f.Weight / 2
		}
	}

	if totalWeight <= 0 {
		rep.Score = 0
	} else {
		rep.Score = int((float64(scoreWeight) / float64(totalWeight)) * 100)
	}

	switch {
	case rep.Score >= 90:
		rep.ScoreLabel = "High"
	case rep.Score >= 70:
		rep.ScoreLabel = "Medium"
	default:
		rep.ScoreLabel = "Low"
	}

	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d", rep.Hostname, rep.Timestamp.UTC().Format(time.RFC3339Nano), rep.Score)))
	rep.ReportHash = hex.EncodeToString(h[:])
	return rep, nil
}

type checkFunc func(ctx context.Context, opts AuditOptions) (Finding, error)

func defaultChecks() []checkFunc {
	return []checkFunc{
		checkSSHRootLogin,
		checkSSHPasswordAuth,
		checkIPForwarding,
		checkFirewallPresence,
	}
}

func checkSSHRootLogin(ctx context.Context, opts AuditOptions) (Finding, error) {
	f := Finding{ID: "ssh.root_login", Title: "Ensure root SSH login is disabled", Weight: 30}
	if runtime.GOOS != "linux" {
		f.Result = ResultSkip
		f.Details = "not supported on this OS"
		return f, nil
	}

	cfgPath := "/etc/ssh/sshd_config"
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		f.Result = ResultWarn
		return f, err
	}

	val, found := sshdConfigValue(string(b), "PermitRootLogin")
	if found && strings.EqualFold(strings.TrimSpace(val), "no") {
		f.Result = ResultPass
		return f, nil
	}

	f.Result = ResultFail
	f.Recommendation = "Set PermitRootLogin no in /etc/ssh/sshd_config and reload sshd"

	if opts.Fix {
		if !opts.Yes {
			return f, errors.New("--fix requested but requires --yes to apply changes")
		}
		if err := setSSHDConfigKey(cfgPath, "PermitRootLogin", "no"); err != nil {
			return f, err
		}
		_ = exec.CommandContext(ctx, "systemctl", "reload", "sshd").Run()
		f.Result = ResultWarn
		f.Details = "applied remediation; re-run audit to verify"
	}

	return f, nil
}

func checkSSHPasswordAuth(ctx context.Context, opts AuditOptions) (Finding, error) {
	f := Finding{ID: "ssh.password_auth", Title: "Ensure SSH password authentication is disabled", Weight: 25}
	if runtime.GOOS != "linux" {
		f.Result = ResultSkip
		f.Details = "not supported on this OS"
		return f, nil
	}

	cfgPath := "/etc/ssh/sshd_config"
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		f.Result = ResultWarn
		return f, err
	}

	val, found := sshdConfigValue(string(b), "PasswordAuthentication")
	if found && strings.EqualFold(strings.TrimSpace(val), "no") {
		f.Result = ResultPass
		return f, nil
	}

	f.Result = ResultFail
	f.Recommendation = "Set PasswordAuthentication no in /etc/ssh/sshd_config and reload sshd"

	if opts.Fix {
		if !opts.Yes {
			return f, errors.New("--fix requested but requires --yes to apply changes")
		}
		if err := setSSHDConfigKey(cfgPath, "PasswordAuthentication", "no"); err != nil {
			return f, err
		}
		_ = exec.CommandContext(ctx, "systemctl", "reload", "sshd").Run()
		f.Result = ResultWarn
		f.Details = "applied remediation; re-run audit to verify"
	}

	return f, nil
}

func checkIPForwarding(ctx context.Context, opts AuditOptions) (Finding, error) {
	f := Finding{ID: "sysctl.ip_forward", Title: "Ensure IPv4 forwarding is disabled", Weight: 20}
	if runtime.GOOS != "linux" {
		f.Result = ResultSkip
		f.Details = "not supported on this OS"
		return f, nil
	}

	val, err := os.ReadFile("/proc/sys/net/ipv4/ip_forward")
	if err != nil {
		f.Result = ResultWarn
		return f, err
	}
	if strings.TrimSpace(string(val)) == "0" {
		f.Result = ResultPass
		return f, nil
	}

	f.Result = ResultFail
	f.Recommendation = "Set net.ipv4.ip_forward=0 via sysctl and persist in /etc/sysctl.d/99-fortis.conf"

	if opts.Fix {
		if !opts.Yes {
			return f, errors.New("--fix requested but requires --yes to apply changes")
		}
		_ = exec.CommandContext(ctx, "sysctl", "-w", "net.ipv4.ip_forward=0").Run()
		if err := persistSysctl("net.ipv4.ip_forward", "0"); err != nil {
			return f, err
		}
		f.Result = ResultWarn
		f.Details = "applied remediation; re-run audit to verify"
	}

	return f, nil
}

func checkFirewallPresence(ctx context.Context, opts AuditOptions) (Finding, error) {
	f := Finding{ID: "firewall.present", Title: "Ensure a firewall is installed", Weight: 25}
	if runtime.GOOS != "linux" {
		f.Result = ResultSkip
		f.Details = "not supported on this OS"
		return f, nil
	}

	if _, err := exec.LookPath("ufw"); err == nil {
		f.Result = ResultPass
		f.Details = "ufw detected"
		return f, nil
	}
	if _, err := exec.LookPath("nft"); err == nil {
		f.Result = ResultPass
		f.Details = "nftables detected"
		return f, nil
	}
	if _, err := exec.LookPath("iptables"); err == nil {
		f.Result = ResultWarn
		f.Details = "iptables detected (legacy)"
		return f, nil
	}

	f.Result = ResultFail
	f.Recommendation = "Install and enable ufw or nftables"
	_ = opts
	return f, nil
}

func persistSysctl(key, value string) error {
	path := "/etc/sysctl.d/99-fortis.conf"
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	// naive append/update: rewrite file with updated key
	lines := []string{}
	if b, err := os.ReadFile(path); err == nil {
		s := bufio.NewScanner(strings.NewReader(string(b)))
		for s.Scan() {
			lines = append(lines, s.Text())
		}
	}

	found := false
	for i := range lines {
		trim := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trim, key+"=") || strings.HasPrefix(trim, key+" =") {
			lines[i] = fmt.Sprintf("%s = %s", key, value)
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, fmt.Sprintf("%s = %s", key, value))
	}
	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(path, []byte(content), 0o644)
}

func sshdConfigValue(cfg string, key string) (string, bool) {
	s := bufio.NewScanner(strings.NewReader(cfg))
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.EqualFold(fields[0], key) {
			return strings.Join(fields[1:], " "), true
		}
	}
	return "", false
}

func setSSHDConfigKey(path, key, value string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	backup := fmt.Sprintf("%s.bak.%d", path, time.Now().Unix())
	if err := os.WriteFile(backup, b, 0o600); err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")
	updated := false
	for i := range lines {
		trim := strings.TrimSpace(lines[i])
		if trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}
		fields := strings.Fields(trim)
		if len(fields) >= 1 && strings.EqualFold(fields[0], key) {
			lines[i] = fmt.Sprintf("%s %s", key, value)
			updated = true
			break
		}
	}
	if !updated {
		lines = append(lines, fmt.Sprintf("%s %s", key, value))
	}

	out := strings.Join(lines, "\n")
	return os.WriteFile(path, []byte(out), 0o600)
}
