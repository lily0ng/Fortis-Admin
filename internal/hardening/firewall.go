package hardening

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type FirewallOptions struct {
	Profile   string
	Ports     string
	Direction string
	Save      bool
	Yes       bool
	DryRun    bool
}

type FirewallResult struct {
	Backend string
	Plan    []string
}

func ConfigureFirewall(ctx context.Context, opts FirewallOptions) (FirewallResult, error) {
	_ = ctx
	if runtime.GOOS != "linux" {
		return FirewallResult{Backend: "unsupported", Plan: []string{"not supported on this OS"}}, nil
	}
	backend := detectFirewallBackend()
	res := FirewallResult{Backend: backend}

	ports := parsePorts(opts.Ports)
	if len(ports) == 0 {
		ports = defaultPortsForProfile(opts.Profile)
	}

	switch backend {
	case "ufw":
		res.Plan = append(res.Plan, "ufw status verbose")
		for _, p := range ports {
			res.Plan = append(res.Plan, fmt.Sprintf("ufw allow %s", p))
		}
		if opts.Save {
			res.Plan = append(res.Plan, "ufw enable")
		}
		if opts.DryRun {
			return res, nil
		}
		if !opts.Yes {
			return res, errors.New("refusing to apply firewall changes without --yes")
		}
		_ = exec.CommandContext(ctx, "ufw", "status", "verbose").Run()
		for _, p := range ports {
			_ = exec.CommandContext(ctx, "ufw", "allow", p).Run()
		}
		if opts.Save {
			_ = exec.CommandContext(ctx, "ufw", "--force", "enable").Run()
		}
		return res, nil
	default:
		res.Plan = append(res.Plan, "no supported firewall backend detected (ufw) - install ufw")
		return res, nil
	}
}

func detectFirewallBackend() string {
	if _, err := exec.LookPath("ufw"); err == nil {
		return "ufw"
	}
	if _, err := exec.LookPath("nft"); err == nil {
		return "nftables"
	}
	if _, err := exec.LookPath("iptables"); err == nil {
		return "iptables"
	}
	return "none"
}

func parsePorts(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

func defaultPortsForProfile(profile string) []string {
	switch strings.ToLower(strings.TrimSpace(profile)) {
	case "webserver":
		return []string{"80/tcp", "443/tcp"}
	case "database":
		return []string{"5432/tcp"}
	default:
		return nil
	}
}
