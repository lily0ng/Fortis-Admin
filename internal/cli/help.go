package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"fortis-admin/internal/buildinfo"
)

func setRootHelp(cmd *cobra.Command) {
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printRootHelp(cmd.OutOrStdout())
	})
}

func shortVersion(v string) string {
	parts := strings.Split(v, ".")
	if len(parts) >= 2 {
		if len(parts) >= 3 && parts[2] == "0" {
			return parts[0] + "." + parts[1]
		}
		return parts[0] + "." + parts[1] + "." + parts[2]
	}
	return v
}

func setGroupHelp(cmd *cobra.Command, title string, usage string, body func(w io.Writer)) {
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		w := cmd.OutOrStdout()
		printGroupBanner(w, title)
		fmt.Fprintln(w, "USAGE:")
		fmt.Fprintf(w, "  %s\n\n", usage)
		body(w)
	})
}

func printRootHelp(w io.Writer) {
	printRootBanner(w)
	fmt.Fprintln(w, "USAGE:")
	fmt.Fprintln(w, "  fortis [command] [flags] [arguments]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "AVAILABLE COMMANDS:")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "SERVER HARDENING:")
	fmt.Fprintln(w, "  harden audit          Run comprehensive security audit")
	fmt.Fprintln(w, "  harden apply          Apply hardening configuration")
	fmt.Fprintln(w, "  harden ssh            Secure SSH configuration")
	fmt.Fprintln(w, "  harden firewall       Configure firewall rules")
	fmt.Fprintln(w, "  harden kernel         Optimize kernel security parameters")
	fmt.Fprintln(w, "  harden users          Manage user security policies")
	fmt.Fprintln(w, "  harden compliance     Generate compliance reports")
	fmt.Fprintln(w, "  harden auto-fix       Automatically remediate security issues")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "INCIDENT RESPONSE:")
	fmt.Fprintln(w, "  incident capture      Collect forensic evidence")
	fmt.Fprintln(w, "  incident triage       Perform rapid system triage")
	fmt.Fprintln(w, "  incident analyze      Analyze captured data")
	fmt.Fprintln(w, "  incident hunt         Hunt for threats/IOCs")
	fmt.Fprintln(w, "  incident report       Generate incident reports")
	fmt.Fprintln(w, "  incident timeline     Create forensic timeline")
	fmt.Fprintln(w, "  incident contain      Execute containment procedures")
	fmt.Fprintln(w, "  incident eradicate    Remove threats from system")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "BACKUP & RECOVERY:")
	fmt.Fprintln(w, "  backup create         Create new backup")
	fmt.Fprintln(w, "  backup list           List available backups")
	fmt.Fprintln(w, "  backup verify         Verify backup integrity")
	fmt.Fprintln(w, "  backup restore        Restore from backup")
	fmt.Fprintln(w, "  backup schedule       Manage backup schedules")
	fmt.Fprintln(w, "  backup catalog        Browse backup contents")
	fmt.Fprintln(w, "  backup monitor        Monitor backup status")
	fmt.Fprintln(w, "  backup test-dr        Test disaster recovery")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "CLUSTER MANAGEMENT:")
	fmt.Fprintln(w, "  cluster exec          Execute command on multiple servers")
	fmt.Fprintln(w, "  cluster deploy        Deploy configuration to cluster")
	fmt.Fprintln(w, "  cluster monitor       Monitor cluster health")
	fmt.Fprintln(w, "  cluster inventory     Manage server inventory")
	fmt.Fprintln(w, "  cluster patch         Orchestrate patching across cluster")
	fmt.Fprintln(w, "  cluster sync          Synchronize files across nodes")
	fmt.Fprintln(w, "  cluster report        Generate cluster status report")
	fmt.Fprintln(w, "  cluster alert         Configure cluster alerts")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "UTILITIES:")
	fmt.Fprintln(w, "  config view           View configuration")
	fmt.Fprintln(w, "  config set            Set configuration value")
	fmt.Fprintln(w, "  logs show             Show application logs")
	fmt.Fprintln(w, "  logs tail             Tail application logs")
	fmt.Fprintln(w, "  update                Update fortis to latest version")
	fmt.Fprintln(w, "  version               Display version information")
	fmt.Fprintln(w, "  completion            Generate shell completions")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "FLAGS:")
	fmt.Fprintln(w, "  -c, --config string      Configuration file (default \"/etc/fortis/config.yaml\")")
	fmt.Fprintln(w, "  -d, --debug              Enable debug mode")
	fmt.Fprintln(w, "  -q, --quiet              Quiet mode (minimal output)")
	fmt.Fprintln(w, "  -v, --verbose            Verbose output")
	fmt.Fprintln(w, "      --color              Force color output")
	fmt.Fprintln(w, "      --no-color           Disable color output")
	fmt.Fprintln(w, "      --version            Display version information")
	fmt.Fprintln(w, "  -h, --help               Help for command")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "EXAMPLES:")
	fmt.Fprintln(w, "  fortis harden audit --output report.html")
	fmt.Fprintln(w, "  fortis backup create --target /data --encrypt")
	fmt.Fprintln(w, "  fortis cluster exec --group webservers \"systemctl restart nginx\"")
	fmt.Fprintln(w, "  fortis incident capture --case incident-001")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "For more information about a specific command, run:")
	fmt.Fprintln(w, "  fortis [command] --help")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Documentation: https://fortis-admin.readthedocs.io")
	fmt.Fprintln(w, "Report issues: https://github.com/yourname/fortis-admin/issues")
}

func printRootBanner(w io.Writer) {
	title := fmt.Sprintf("FORTIS-ADMIN v%s", shortVersion(buildinfo.Version))
	fmt.Fprintln(w, "╔══════════════════════════════════════════════════════════╗")
	fmt.Fprintf(w, "║ %s ║\n", center(56, title))
	fmt.Fprintln(w, "║         System Administration & Security Toolkit         ║")
	fmt.Fprintln(w, "╚══════════════════════════════════════════════════════════╝")
	fmt.Fprintln(w)
}

func printGroupBanner(w io.Writer, title string) {
	fmt.Fprintln(w, "╔══════════════════════════════════════════════════════════╗")
	fmt.Fprintf(w, "║ %s ║\n", center(56, title))
	fmt.Fprintln(w, "╚══════════════════════════════════════════════════════════╝")
	fmt.Fprintln(w)
}

func center(width int, s string) string {
	if width <= 0 {
		return s
	}
	if len(s) >= width {
		return s[:width]
	}
	pad := width - len(s)
	left := pad / 2
	right := pad - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
