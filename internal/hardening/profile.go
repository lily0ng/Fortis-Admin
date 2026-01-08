package hardening

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ApplyOptions struct {
	Profile    string
	DryRun     bool
	Rollback   bool
	SkipChecks bool
	Yes        bool
}

type PlanItem struct {
	Description string
	ChangedFile string
}

type ApplyResult struct {
	Profile     string
	Plan        []PlanItem
	RollbackDir string
}

func ApplyProfile(ctx context.Context, opts ApplyOptions) (ApplyResult, error) {
	_ = ctx
	if opts.Profile == "" {
		p, err := promptProfile()
		if err != nil {
			return ApplyResult{}, err
		}
		opts.Profile = p
	}

	res := ApplyResult{Profile: opts.Profile}

	// Very small starter set of actions; safe and expandable.
	// Only applies when --yes is provided.
	filesToBackup := map[string]struct{}{}

	applySSHD := func(key, value string) error {
		path := "/etc/ssh/sshd_config"
		res.Plan = append(res.Plan, PlanItem{Description: fmt.Sprintf("Set %s %s", key, value), ChangedFile: path})
		filesToBackup[path] = struct{}{}
		if opts.DryRun {
			return nil
		}
		if !opts.Yes {
			return errors.New("refusing to apply without --yes")
		}
		if err := setSSHDConfigKey(path, key, value); err != nil {
			return err
		}
		return nil
	}

	applySysctl := func(key, value string) error {
		res.Plan = append(res.Plan, PlanItem{Description: fmt.Sprintf("Persist sysctl %s=%s", key, value), ChangedFile: "/etc/sysctl.d/99-fortis.conf"})
		if opts.DryRun {
			return nil
		}
		if !opts.Yes {
			return errors.New("refusing to apply without --yes")
		}
		return persistSysctl(key, value)
	}

	if !opts.SkipChecks {
		res.Plan = append(res.Plan, PlanItem{Description: "Pre-check: ensure running on supported OS", ChangedFile: ""})
	}

	switch strings.ToLower(opts.Profile) {
	case "cis", "cis-level1", "baseline":
		if err := applySSHD("PermitRootLogin", "no"); err != nil {
			return res, err
		}
		if err := applySSHD("PasswordAuthentication", "no"); err != nil {
			return res, err
		}
		if err := applySysctl("net.ipv4.ip_forward", "0"); err != nil {
			return res, err
		}
	case "webserver":
		if err := applySSHD("PermitRootLogin", "no"); err != nil {
			return res, err
		}
		if err := applySSHD("PasswordAuthentication", "no"); err != nil {
			return res, err
		}
	case "database":
		if err := applySSHD("PermitRootLogin", "no"); err != nil {
			return res, err
		}
	case "desktop":
		if err := applySSHD("PermitRootLogin", "no"); err != nil {
			return res, err
		}
	default:
		return res, fmt.Errorf("unknown profile: %s", opts.Profile)
	}

	if opts.Rollback && !opts.DryRun {
		if !opts.Yes {
			return res, errors.New("--rollback requires --yes to create and store backup")
		}
		rb, err := backupFiles(filesToBackup)
		if err != nil {
			return res, err
		}
		res.RollbackDir = rb
	}

	return res, nil
}

func promptProfile() (string, error) {
	// Minimal interactive selector.
	fmt.Fprintln(os.Stdout, "Select hardening profile:")
	fmt.Fprintln(os.Stdout, "  1) cis")
	fmt.Fprintln(os.Stdout, "  2) webserver")
	fmt.Fprintln(os.Stdout, "  3) database")
	fmt.Fprintln(os.Stdout, "  4) desktop")
	fmt.Fprint(os.Stdout, "Enter choice [1-4]: ")
	br := bufio.NewReader(os.Stdin)
	line, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	switch line {
	case "1", "cis":
		return "cis", nil
	case "2", "webserver":
		return "webserver", nil
	case "3", "database":
		return "database", nil
	case "4", "desktop":
		return "desktop", nil
	default:
		return "", errors.New("invalid selection")
	}
}

func backupFiles(files map[string]struct{}) (string, error) {
	ts := time.Now().Format("20060102-150405")
	rbDir := filepath.Join("/var/lib/fortis/rollback", ts)
	if err := os.MkdirAll(rbDir, 0o755); err != nil {
		// fallback
		rbDir = filepath.Join(".", ".fortis-rollback", ts)
		if err2 := os.MkdirAll(rbDir, 0o755); err2 != nil {
			return "", err
		}
	}

	for p := range files {
		b, err := os.ReadFile(p)
		if err != nil {
			// if file missing, skip backup
			continue
		}
		name := strings.ReplaceAll(strings.TrimPrefix(p, "/"), "/", "_")
		if err := os.WriteFile(filepath.Join(rbDir, name), b, 0o600); err != nil {
			return "", err
		}
	}
	return rbDir, nil
}
