package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"fortis-admin/internal/app"
	"fortis-admin/internal/hardening"
)

func newHardenCmd(a *app.App) *cobra.Command {
	var (
		backup    bool
		yes       bool
		logFile   string
		configDir string
	)

	cmd := &cobra.Command{
		Use:   "harden",
		Short: "Server hardening automation",
	}
	cmd.GroupID = "harden"
	cmd.PersistentFlags().BoolVar(&backup, "backup", false, "Create backup before making changes")
	cmd.PersistentFlags().BoolVar(&yes, "yes", false, "Auto-confirm all prompts")
	cmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Log file location")
	cmd.PersistentFlags().StringVar(&configDir, "config-dir", "", "Configuration directory")

	cmd.AddCommand(newHardenAuditCmd(a))
	cmd.AddCommand(newHardenApplyCmd(a))
	cmd.AddCommand(newHardenSSHCmdd(a))
	cmd.AddCommand(newHardenFirewallCmd(a))
	cmd.AddCommand(newHardenKernelCmd(a))
	cmd.AddCommand(newHardenUsersCmd(a))
	cmd.AddCommand(newHardenComplianceCmd(a))
	cmd.AddCommand(newHardenAutoFixCmd(a))
	cmd.AddCommand(newHardenFilesystemCmd(a))
	cmd.AddCommand(newHardenPackageAuditCmd(a))
	cmd.AddCommand(newHardenAuditdCmd(a))
	cmd.AddCommand(newHardenLoggingCmd(a))
	cmd.AddCommand(newHardenServicesCmd(a))
	setGroupHelp(cmd, "SERVER HARDENING COMMANDS", "fortis harden [command] [flags]", func(w io.Writer) {
		_ = backup
		_ = yes
		_ = logFile
		_ = configDir

		io.WriteString(w, "COMMANDS:\n")
		io.WriteString(w, "  audit [flags]                    Run security audit and generate report\n")
		io.WriteString(w, "    --profile string               Audit profile (cis, pci, hipaa, custom)\n")
		io.WriteString(w, "    --output string                Output format/file (json, yaml, html, pdf)\n")
		io.WriteString(w, "    --level string                 Audit level (basic, medium, strict)\n")
		io.WriteString(w, "    --fix                          Auto-fix low-risk issues\n\n")

		io.WriteString(w, "  apply [flags]                    Apply hardening configuration\n")
		io.WriteString(w, "    --profile string               Hardening profile to apply\n")
		io.WriteString(w, "    --dry-run                      Show changes without applying\n")
		io.WriteString(w, "    --rollback                     Create rollback point\n")
		io.WriteString(w, "    --skip-checks                  Skip pre-application checks\n\n")

		io.WriteString(w, "  ssh [flags]                      Configure SSH security\n")
		io.WriteString(w, "    --disable-root                 Disable root SSH login\n")
		io.WriteString(w, "    --port number                  Change SSH port\n")
		io.WriteString(w, "    --key-only                     Enforce key-based authentication\n")
		io.WriteString(w, "    --banner string                Set SSH warning banner\n\n")

		io.WriteString(w, "  firewall [flags]                 Configure system firewall\n")
		io.WriteString(w, "    --profile string               Firewall profile (webserver, database, desktop)\n")
		io.WriteString(w, "    --ports string                 Comma-separated list of ports to allow\n")
		io.WriteString(w, "    --direction string             Rule direction (incoming, outgoing, both)\n")
		io.WriteString(w, "    --save                         Save rules to persist after reboot\n\n")

		io.WriteString(w, "  kernel [flags]                   Harden kernel parameters\n")
		io.WriteString(w, "    --apply string                 Apply preset (network, memory, security)\n")
		io.WriteString(w, "    --param string                 Set specific kernel parameter\n")
		io.WriteString(w, "    --value string                 Parameter value\n")
		io.WriteString(w, "    --persist                      Make changes persistent\n\n")

		io.WriteString(w, "  users [flags]                    Manage user security\n")
		io.WriteString(w, "    --lock-inactive                Lock inactive accounts\n")
		io.WriteString(w, "    --password-policy              Set password policies\n")
		io.WriteString(w, "    --sudo-secure                  Secure sudo configuration\n")
		io.WriteString(w, "    --audit                        Audit user permissions\n\n")

		io.WriteString(w, "  compliance [flags]               Generate compliance reports\n")
		io.WriteString(w, "    --standard string              Compliance standard (pci-dss, hipaa, gdpr, iso27001)\n")
		io.WriteString(w, "    --evidence                     Collect evidence for compliance\n")
		io.WriteString(w, "    --gap-analysis                 Show compliance gaps\n")
		io.WriteString(w, "    --export string                Export format (csv, pdf, json)\n\n")

		io.WriteString(w, "  auto-fix [flags]                 Automatically fix security issues\n")
		io.WriteString(w, "    --level string                 Fix level (low, medium, high, critical)\n")
		io.WriteString(w, "    --exclude strings              Issues to exclude from auto-fix\n")
		io.WriteString(w, "    --confirm                      Skip confirmation prompts\n")
		io.WriteString(w, "    --log-only                     Log issues without fixing\n\n")

		io.WriteString(w, "FLAGS:\n")
		io.WriteString(w, "  --backup                Create backup before making changes\n")
		io.WriteString(w, "  --yes                   Auto-confirm all prompts\n")
		io.WriteString(w, "  --log-file string       Log file location\n")
		io.WriteString(w, "  --config-dir string     Configuration directory\n\n")

		io.WriteString(w, "EXAMPLES:\n")
		io.WriteString(w, "  fortis harden audit --profile cis --output html\n")
		io.WriteString(w, "  fortis harden apply --profile webserver --dry-run\n")
		io.WriteString(w, "  fortis harden ssh --disable-root --key-only\n")
		io.WriteString(w, "  fortis harden auto-fix --level medium --confirm\n")
	})
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
			yes := getBoolFlag(cmd, "yes")

			rep, err := hardening.RunAudit(cmd.Context(), hardening.AuditOptions{
				Profile: profile,
				Level:   level,
				Output:  output,
				Fix:     fix,
				Yes:     yes,
				Verbose: a.Verbose,
			})
			if err != nil {
				return err
			}

			format := hardening.DetectFormat(output)
			path, outErr := resolveAuditOutputPath(output, format, rep.Timestamp)
			if outErr != nil {
				return outErr
			}

			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			defer f.Close()

			if err := hardening.Render(f, rep, format); err != nil {
				return err
			}

			if a.Verbose {
				fmt.Fprintf(cmd.OutOrStdout(), "📊  [STATS] Passed: %d | Failed: %d | Warnings: %d | Skipped: %d\n", rep.Passed, rep.Failed, rep.Warnings, rep.Skipped)
				fmt.Fprintf(cmd.OutOrStdout(), "🎯  [SCORE] Security Score: %d/100 (%s)\n", rep.Score, rep.ScoreLabel)
				fmt.Fprintf(cmd.OutOrStdout(), "📁  [SAVE]  Report saved to: %s\n", path)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Audit complete. Score: %d/100 (%s). Report: %s\n", rep.Score, rep.ScoreLabel, path)
			return nil
		},
	}
	cmd.Flags().StringVar(&profile, "profile", "cis", "Audit profile (cis, pci, hipaa, custom)")
	cmd.Flags().StringVar(&output, "output", "", "Output format/file (json, yaml, html, pdf)")
	cmd.Flags().StringVar(&level, "level", "basic", "Audit level (basic, medium, strict)")
	cmd.Flags().BoolVar(&fix, "fix", false, "Auto-fix low-risk issues")
	return cmd
}

func resolveAuditOutputPath(output string, format hardening.OutputFormat, ts time.Time) (string, error) {
	if output != "" {
		// If looks like a file path (has an extension or contains a slash), respect it.
		if strings.Contains(output, "/") || strings.Contains(output, "\\") || strings.Contains(output, ".") {
			return output, nil
		}
	}

	ext := string(format)
	name := fmt.Sprintf("audit-%s.%s", ts.Format("20060102-150405"), ext)

	preferred := filepath.Join("/var/log/fortis", name)
	if canWriteDir("/var/log/fortis") {
		return preferred, nil
	}
	return filepath.Join(".", name), nil
}

func getBoolFlag(cmd *cobra.Command, name string) bool {
	if f := cmd.Flags().Lookup(name); f != nil {
		v, err := cmd.Flags().GetBool(name)
		if err == nil {
			return v
		}
	}
	if f := cmd.InheritedFlags().Lookup(name); f != nil {
		v, err := cmd.InheritedFlags().GetBool(name)
		if err == nil {
			return v
		}
	}
	return false
}

func canWriteDir(dir string) bool {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false
	}
	f, err := os.CreateTemp(dir, ".perm")
	if err != nil {
		return false
	}
	name := f.Name()
	_ = f.Close()
	_ = os.Remove(name)
	return true
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
			yes := getBoolFlag(cmd, "yes")
			res, err := hardening.ApplyProfile(cmd.Context(), hardening.ApplyOptions{
				Profile:    profile,
				DryRun:     dryRun,
				Rollback:   rollback,
				SkipChecks: skipChecks,
				Yes:        yes,
			})
			if err != nil {
				return err
			}
			for _, p := range res.Plan {
				if p.ChangedFile != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "- %s (%s)\n", p.Description, p.ChangedFile)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", p.Description)
				}
			}
			if res.RollbackDir != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Rollback saved to: %s\n", res.RollbackDir)
			}
			if dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), "[DRY-RUN] No changes applied.")
			}
			return nil
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
			if getBoolFlag(cmd, "yes") {
				argv = append(argv, "--yes")
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
			yes := getBoolFlag(cmd, "yes")
			res, err := hardening.ConfigureFirewall(cmd.Context(), hardening.FirewallOptions{
				Profile:   profile,
				Ports:     ports,
				Direction: direction,
				Save:      save,
				Yes:       yes,
				DryRun:    !yes,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Firewall backend: %s\n", res.Backend)
			for _, line := range res.Plan {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", line)
			}
			if !yes {
				fmt.Fprintln(cmd.OutOrStdout(), "[DRY-RUN] Re-run with --yes to apply.")
			}
			return nil
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
			yes := getBoolFlag(cmd, "yes")
			res, err := hardening.ApplyKernel(cmd.Context(), hardening.KernelOptions{
				ApplyPreset: apply,
				Param:       param,
				Value:       value,
				Persist:     persist,
				Yes:         yes,
				DryRun:      !yes,
			})
			if err != nil {
				return err
			}
			for _, line := range res.Plan {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", line)
			}
			if !yes {
				fmt.Fprintln(cmd.OutOrStdout(), "[DRY-RUN] Re-run with --yes to apply.")
			}
			return nil
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
			if getBoolFlag(cmd, "yes") {
				argv = append(argv, "--yes")
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
			_ = args
			outPath, fmtDetected, err := hardening.ResolveComplianceOutputPath(exportFmt, time.Now())
			if err != nil {
				return err
			}
			// Explicitly stub pdf.
			if strings.HasSuffix(strings.ToLower(outPath), ".pdf") || strings.ToLower(exportFmt) == "pdf" {
				return errors.New("pdf export not implemented")
			}
			rep, err := hardening.GenerateComplianceReport(cmd.Context(), hardening.ComplianceOptions{
				Standard:        standard,
				CollectEvidence: evidence,
				GapAnalysis:     gap,
				Format:          fmtDetected,
			})
			if err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				return err
			}
			f, err := os.Create(outPath)
			if err != nil {
				return err
			}
			defer f.Close()
			if err := hardening.RenderCompliance(f, rep, fmtDetected); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Compliance report saved to: %s\n", outPath)
			return nil
		},
	}
	cmd.Flags().StringVar(&standard, "standard", "", "Compliance standard (pci-dss, hipaa, gdpr, iso27001)")
	cmd.Flags().BoolVar(&evidence, "evidence", false, "Collect evidence for compliance")
	cmd.Flags().BoolVar(&gap, "gap-analysis", false, "Show compliance gaps")
	cmd.Flags().StringVar(&exportFmt, "export", "", "Export format (csv, pdf, json)")
	_ = a
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

func newHardenFilesystemCmd(a *app.App) *cobra.Command {
	var (
		rootPath string
		output   string
	)
	cmd := &cobra.Command{
		Use:    "filesystem",
		Short:  "Filesystem permission auditing",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			rep, err := hardening.ScanFilesystem(cmd.Context(), hardening.FilesystemOptions{Root: rootPath})
			if err != nil {
				return err
			}
			w := io.Writer(cmd.OutOrStdout())
			var file *os.File
			if output != "" {
				if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
					return err
				}
				f, err := os.Create(output)
				if err != nil {
					return err
				}
				defer f.Close()
				file = f
				w = f
			}

			if strings.HasSuffix(output, ".yaml") || strings.HasSuffix(output, ".yml") {
				b, err := yaml.Marshal(rep)
				if err != nil {
					return err
				}
				_, err = w.Write(b)
				return err
			}

			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			if err := enc.Encode(rep); err != nil {
				return err
			}
			if file != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Report saved to: %s\n", output)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&rootPath, "root", "/", "Root path to scan")
	cmd.Flags().StringVar(&output, "output", "", "Output file (json|yaml inferred by extension)")
	_ = a
	return cmd
}

func newHardenPackageAuditCmd(a *app.App) *cobra.Command {
	var (
		listPkgs bool
		output   string
	)
	cmd := &cobra.Command{
		Use:    "package-audit",
		Short:  "Installed package vulnerability check",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			rep, err := hardening.AuditPackages(cmd.Context(), hardening.PackageAuditOptions{List: listPkgs})
			if err != nil {
				return err
			}
			w := io.Writer(cmd.OutOrStdout())
			var file *os.File
			if output != "" {
				if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
					return err
				}
				f, err := os.Create(output)
				if err != nil {
					return err
				}
				defer f.Close()
				file = f
				w = f
			}

			if strings.HasSuffix(output, ".yaml") || strings.HasSuffix(output, ".yml") {
				b, err := yaml.Marshal(rep)
				if err != nil {
					return err
				}
				_, err = w.Write(b)
				return err
			}

			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			if err := enc.Encode(rep); err != nil {
				return err
			}
			if file != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Report saved to: %s\n", output)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&listPkgs, "list", false, "Include package list")
	cmd.Flags().StringVar(&output, "output", "", "Output file (json|yaml inferred by extension)")
	_ = a
	return cmd
}

func newHardenAuditdCmd(a *app.App) *cobra.Command {
	var apply bool
	cmd := &cobra.Command{
		Use:    "auditd",
		Short:  "Auditd rule deployment",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{}
			if apply {
				argv = append(argv, "--apply")
			}
			if getBoolFlag(cmd, "yes") {
				argv = append(argv, "--yes")
			}
			return a.RunScript(cmd.Context(), "auditd-setup.sh", argv...)
		},
	}
	cmd.Flags().BoolVar(&apply, "apply", false, "Apply changes")
	return cmd
}

func newHardenLoggingCmd(a *app.App) *cobra.Command {
	var (
		apply  bool
		remote string
	)
	cmd := &cobra.Command{
		Use:    "logging",
		Short:  "Centralized logging setup",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{}
			if remote != "" {
				argv = append(argv, "--remote", remote)
			}
			if apply {
				argv = append(argv, "--apply")
			}
			if getBoolFlag(cmd, "yes") {
				argv = append(argv, "--yes")
			}
			return a.RunScript(cmd.Context(), "configure-logging.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&remote, "remote", "", "Remote syslog host:port")
	cmd.Flags().BoolVar(&apply, "apply", false, "Apply changes")
	return cmd
}

func newHardenServicesCmd(a *app.App) *cobra.Command {
	var (
		list    bool
		disable string
	)
	cmd := &cobra.Command{
		Use:    "services",
		Short:  "Identify/disable unnecessary services",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{}
			if list {
				argv = append(argv, "--list")
			}
			if disable != "" {
				argv = append(argv, "--disable", disable)
			}
			if getBoolFlag(cmd, "yes") {
				argv = append(argv, "--yes")
			}
			return a.RunScript(cmd.Context(), "disable-services.sh", argv...)
		},
	}
	cmd.Flags().BoolVar(&list, "list", true, "List enabled services")
	cmd.Flags().StringVar(&disable, "disable", "", "Disable a specific service")
	return cmd
}
