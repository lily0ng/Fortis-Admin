package cli

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
)

func newHardenCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "harden",
		Short: "Server hardening automation",
	}
	cmd.GroupID = "harden"

	cmd.AddCommand(newHardenAuditCmd(a))
	cmd.AddCommand(newHardenApplyCmd(a))
	cmd.AddCommand(newHardenSSHCmdd(a))
	cmd.AddCommand(newHardenFirewallCmd(a))
	cmd.AddCommand(newHardenKernelCmd(a))
	cmd.AddCommand(newHardenUsersCmd(a))
	cmd.AddCommand(newHardenComplianceCmd(a))
	cmd.AddCommand(newHardenAutoFixCmd(a))

	_ = context.Background()
	return cmd
}

func newHardenAuditCmd(a *app.App) *cobra.Command {
	var (
		profile string
		output  string
		level   string
		fix     bool
	)
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Run comprehensive security audit",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--profile", profile, "--output", output, "--level", level}
			if fix {
				argv = append(argv, "--fix")
			}
			return a.RunScript(cmd.Context(), "harden-audit.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&profile, "profile", "cis", "Audit profile (cis, pci, hipaa, custom)")
	cmd.Flags().StringVar(&output, "output", "", "Output format/file (json, yaml, html, pdf)")
	cmd.Flags().StringVar(&level, "level", "basic", "Audit level (basic, medium, strict)")
	cmd.Flags().BoolVar(&fix, "fix", false, "Auto-fix low-risk issues")
	return cmd
}

func newHardenApplyCmd(a *app.App) *cobra.Command {
	var (
		profile    string
		dryRun     bool
		rollback   bool
		skipChecks bool
	)
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply hardening configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--profile", profile}
			if dryRun {
				argv = append(argv, "--dry-run")
			}
			if rollback {
				argv = append(argv, "--rollback")
			}
			if skipChecks {
				argv = append(argv, "--skip-checks")
			}
			return a.RunScript(cmd.Context(), "harden-apply.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&profile, "profile", "", "Hardening profile to apply")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show changes without applying")
	cmd.Flags().BoolVar(&rollback, "rollback", false, "Create rollback point")
	cmd.Flags().BoolVar(&skipChecks, "skip-checks", false, "Skip pre-application checks")
	return cmd
}

func newHardenSSHCmdd(a *app.App) *cobra.Command {
	var (
		disableRoot bool
		port        int
		keyOnly     bool
		banner      string
	)
	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "Secure SSH configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{}
			if disableRoot {
				argv = append(argv, "--disable-root")
			}
			if port != 0 {
				argv = append(argv, "--port", itoa(port))
			}
			if keyOnly {
				argv = append(argv, "--key-only")
			}
			if banner != "" {
				argv = append(argv, "--banner", banner)
			}
			return a.RunScript(cmd.Context(), "harden-ssh.sh", argv...)
		},
	}
	cmd.Flags().BoolVar(&disableRoot, "disable-root", false, "Disable root SSH login")
	cmd.Flags().IntVar(&port, "port", 0, "Change SSH port")
	cmd.Flags().BoolVar(&keyOnly, "key-only", false, "Enforce key-based authentication")
	cmd.Flags().StringVar(&banner, "banner", "", "Set SSH warning banner")
	return cmd
}

func newHardenFirewallCmd(a *app.App) *cobra.Command {
	var (
		profile   string
		ports     string
		direction string
		save      bool
	)
	cmd := &cobra.Command{
		Use:   "firewall",
		Short: "Configure firewall rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--profile", profile, "--ports", ports, "--direction", direction}
			if save {
				argv = append(argv, "--save")
			}
			return a.RunScript(cmd.Context(), "harden-firewall.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&profile, "profile", "", "Firewall profile (webserver, database, desktop)")
	cmd.Flags().StringVar(&ports, "ports", "", "Comma-separated list of ports to allow")
	cmd.Flags().StringVar(&direction, "direction", "incoming", "Rule direction (incoming, outgoing, both)")
	cmd.Flags().BoolVar(&save, "save", false, "Save rules to persist after reboot")
	return cmd
}

func newHardenKernelCmd(a *app.App) *cobra.Command {
	var (
		apply   string
		param   string
		value   string
		persist bool
	)
	cmd := &cobra.Command{
		Use:   "kernel",
		Short: "Optimize kernel security parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--apply", apply, "--param", param, "--value", value}
			if persist {
				argv = append(argv, "--persist")
			}
			return a.RunScript(cmd.Context(), "harden-kernel.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&apply, "apply", "", "Apply preset (network, memory, security)")
	cmd.Flags().StringVar(&param, "param", "", "Set specific kernel parameter")
	cmd.Flags().StringVar(&value, "value", "", "Parameter value")
	cmd.Flags().BoolVar(&persist, "persist", false, "Make changes persistent")
	return cmd
}

func newHardenUsersCmd(a *app.App) *cobra.Command {
	var (
		lockInactive   bool
		passwordPolicy bool
		sudoSecure     bool
		audit          bool
	)
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage user security policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{}
			if lockInactive {
				argv = append(argv, "--lock-inactive")
			}
			if passwordPolicy {
				argv = append(argv, "--password-policy")
			}
			if sudoSecure {
				argv = append(argv, "--sudo-secure")
			}
			if audit {
				argv = append(argv, "--audit")
			}
			return a.RunScript(cmd.Context(), "user-security.sh", argv...)
		},
	}
	cmd.Flags().BoolVar(&lockInactive, "lock-inactive", false, "Lock inactive accounts")
	cmd.Flags().BoolVar(&passwordPolicy, "password-policy", false, "Set password policies")
	cmd.Flags().BoolVar(&sudoSecure, "sudo-secure", false, "Secure sudo configuration")
	cmd.Flags().BoolVar(&audit, "audit", false, "Audit user permissions")
	return cmd
}

func newHardenComplianceCmd(a *app.App) *cobra.Command {
	var (
		standard  string
		evidence  bool
		gap       bool
		exportFmt string
	)
	cmd := &cobra.Command{
		Use:   "compliance",
		Short: "Generate compliance reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--standard", standard, "--export", exportFmt}
			if evidence {
				argv = append(argv, "--evidence")
			}
			if gap {
				argv = append(argv, "--gap-analysis")
			}
			return a.RunScript(cmd.Context(), "compliance-report.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&standard, "standard", "", "Compliance standard (pci-dss, hipaa, gdpr, iso27001)")
	cmd.Flags().BoolVar(&evidence, "evidence", false, "Collect evidence for compliance")
	cmd.Flags().BoolVar(&gap, "gap-analysis", false, "Show compliance gaps")
	cmd.Flags().StringVar(&exportFmt, "export", "", "Export format (csv, pdf, json)")
	return cmd
}

func newHardenAutoFixCmd(a *app.App) *cobra.Command {
	var (
		level   string
		exclude []string
		confirm bool
		logOnly bool
	)
	cmd := &cobra.Command{
		Use:   "auto-fix",
		Short: "Automatically remediate security issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			if level == "" {
				return errors.New("--level is required")
			}
			argv := []string{"--level", level}
			for _, x := range exclude {
				argv = append(argv, "--exclude", x)
			}
			if confirm {
				argv = append(argv, "--confirm")
			}
			if logOnly {
				argv = append(argv, "--log-only")
			}
			return a.RunScript(cmd.Context(), "harden-auto-fix.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&level, "level", "", "Fix level (low, medium, high, critical)")
	cmd.Flags().StringSliceVar(&exclude, "exclude", nil, "Issues to exclude from auto-fix")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Skip confirmation prompts")
	cmd.Flags().BoolVar(&logOnly, "log-only", false, "Log issues without fixing")
	return cmd
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	b := make([]byte, 0, 16)
	for v > 0 {
		d := v % 10
		b = append(b, byte('0'+d))
		v /= 10
	}
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}
