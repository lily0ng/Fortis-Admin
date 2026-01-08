package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/cluster"
)

func newClusterCmd(a *app.App) *cobra.Command {
	var (
		inventoryFile string
		sshKey        string
		sshUser       string
		sshPort       int
		sshTimeout    int
	)

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Multi-server management platform",
	}
	cmd.GroupID = "cluster"
	cmd.PersistentFlags().StringVar(&inventoryFile, "inventory-file", "/etc/fortis/inventory.yaml", "Inventory file")
	cmd.PersistentFlags().StringVar(&sshKey, "ssh-key", "", "SSH key to use")
	cmd.PersistentFlags().StringVar(&sshUser, "ssh-user", "root", "SSH username")
	cmd.PersistentFlags().IntVar(&sshPort, "ssh-port", 22, "SSH port")
	cmd.PersistentFlags().IntVar(&sshTimeout, "ssh-timeout", 30, "SSH timeout in seconds")

	cmd.AddCommand(newClusterInitCmd(a))
	cmd.AddCommand(newClusterExecCmd(a))
	cmd.AddCommand(newClusterDeployCmd(a))
	cmd.AddCommand(newClusterMonitorCmd(a))
	cmd.AddCommand(newClusterInventoryCmd(a))
	cmd.AddCommand(newClusterPatchCmd(a))
	cmd.AddCommand(newClusterSyncCmd(a))
	cmd.AddCommand(newClusterReportCmd(a))
	cmd.AddCommand(newClusterAlertCmd(a))
	setGroupHelp(cmd, "CLUSTER MANAGEMENT COMMANDS", "fortis cluster [command] [flags]", func(w io.Writer) {
		_ = inventoryFile
		_ = sshKey
		_ = sshUser
		_ = sshPort
		_ = sshTimeout

		io.WriteString(w, "COMMANDS:\n")
		io.WriteString(w, "  exec [flags]                   Execute command on multiple servers\n")
		io.WriteString(w, "    --group string               Server group to target\n")
		io.WriteString(w, "    --hosts strings              Specific hosts to target\n")
		io.WriteString(w, "    --file string                Read hosts from file\n")
		io.WriteString(w, "    --command string             Command to execute (required)\n")
		io.WriteString(w, "    --parallel int               Parallel execution limit\n")
		io.WriteString(w, "    --timeout duration           Command timeout\n")
		io.WriteString(w, "    --output string              Output format (combined, separate, json)\n\n")

		io.WriteString(w, "  deploy [flags]                 Deploy configuration to cluster\n")
		io.WriteString(w, "    --config string              Configuration file/directory\n")
		io.WriteString(w, "    --target string              Target path on remote servers\n")
		io.WriteString(w, "    --validate                   Validate before deploying\n")
		io.WriteString(w, "    --backup                     Backup existing configuration\n")
		io.WriteString(w, "    --rollback string            Rollback to previous version\n")
		io.WriteString(w, "    --diff                       Show differences\n\n")

		io.WriteString(w, "  monitor [flags]                Monitor cluster health\n")
		io.WriteString(w, "    --dashboard                  Launch interactive dashboard\n")
		io.WriteString(w, "    --metrics strings            Metrics to monitor (cpu, memory, disk, network)\n")
		io.WriteString(w, "    --refresh duration           Refresh interval\n")
		io.WriteString(w, "    --alerts                     Show active alerts\n")
		io.WriteString(w, "    --export string              Export metrics data\n\n")

		io.WriteString(w, "  inventory [flags]              Manage server inventory\n")
		io.WriteString(w, "    --scan                       Scan network for servers\n")
		io.WriteString(w, "    --add string                 Add server to inventory\n")
		io.WriteString(w, "    --remove string              Remove server from inventory\n")
		io.WriteString(w, "    --groups                     Manage server groups\n")
		io.WriteString(w, "    --tags strings               Tag servers\n")
		io.WriteString(w, "    --export string              Export inventory\n\n")

		io.WriteString(w, "  patch [flags]                  Orchestrate patching across cluster\n")
		io.WriteString(w, "    --packages strings           Packages to update\n")
		io.WriteString(w, "    --strategy string            Patching strategy (rolling, parallel, canary)\n")
		io.WriteString(w, "    --batch-size int             Batch size for rolling updates\n")
		io.WriteString(w, "    --pre-check                  Run pre-patch checks\n")
		io.WriteString(w, "    --post-check                 Run post-patch validation\n")
		io.WriteString(w, "    --rollback-on-failure        Auto-rollback on failure\n\n")

		io.WriteString(w, "  sync [flags]                  Synchronize files across nodes\n")
		io.WriteString(w, "    --source string              Source file/directory\n")
		io.WriteString(w, "    --destination string         Destination path\n")
		io.WriteString(w, "    --delete                    Delete extraneous files\n")
		io.WriteString(w, "    --checksum                  Use checksum comparison\n")
		io.WriteString(w, "    --dry-run                   Show what would be synced\n\n")

		io.WriteString(w, "  report [flags]                Generate cluster status report\n")
		io.WriteString(w, "    --format string             Output format (html, pdf, markdown, json)\n")
		io.WriteString(w, "    --sections strings          Report sections (summary, health, inventory, alerts)\n")
		io.WriteString(w, "    --schedule string           Schedule regular reports\n")
		io.WriteString(w, "    --email string              Email report to address\n\n")

		io.WriteString(w, "  alert [flags]                 Configure cluster alerts\n")
		io.WriteString(w, "    --add string                Add new alert rule\n")
		io.WriteString(w, "    --list                     List alert rules\n")
		io.WriteString(w, "    --remove string            Remove alert rule\n")
		io.WriteString(w, "    --test string              Test alert rule\n")
		io.WriteString(w, "    --silence duration         Silence alerts for duration\n\n")

		io.WriteString(w, "FLAGS:\n")
		io.WriteString(w, "  --inventory-file string      Inventory file (default \"/etc/fortis/inventory.yaml\")\n")
		io.WriteString(w, "  --ssh-key string             SSH key to use\n")
		io.WriteString(w, "  --ssh-user string            SSH username (default \"root\")\n")
		io.WriteString(w, "  --ssh-port int               SSH port (default 22)\n")
		io.WriteString(w, "  --ssh-timeout int            SSH timeout in seconds (default 30)\n\n")

		io.WriteString(w, "EXAMPLES:\n")
		io.WriteString(w, "  fortis cluster exec --group webservers --command \"systemctl restart nginx\"\n")
		io.WriteString(w, "  fortis cluster deploy --config nginx.conf --target /etc/nginx/nginx.conf\n")
		io.WriteString(w, "  fortis cluster monitor --dashboard --refresh 10s\n")
		io.WriteString(w, "  fortis cluster patch --packages nginx openssl --strategy rolling --batch-size 2\n")
		io.WriteString(w, "  fortis cluster sync --source /etc/ssl/certs --destination /etc/ssl/certs --delete\n")
	})

	return cmd
}

func newClusterExecCmd(a *app.App) *cobra.Command {
	var (
		group    string
		hosts    []string
		file     string
		command  string
		parallel int
		timeout  string
		output   string
	)
	cmd := &cobra.Command{
		Use:   "exec [command]",
		Short: "Execute command on multiple servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if command == "" && len(args) > 0 {
				command = strings.Join(args, " ")
			}
			if command == "" {
				return errors.New("--command is required")
			}
			invPath := getStringFlag(cmd, "inventory-file")
			sshUser := getStringFlag(cmd, "ssh-user")
			sshKey := getStringFlag(cmd, "ssh-key")
			sshPort := getIntFlag(cmd, "ssh-port")
			sshTimeout := time.Duration(getIntFlag(cmd, "ssh-timeout")) * time.Second

			par := parallel
			if par <= 0 {
				par = 4
			}
			to := sshTimeout
			if strings.TrimSpace(timeout) != "" {
				if d, err := time.ParseDuration(timeout); err == nil {
					to = d
				}
			}

			res, err := cluster.Exec(cmd.Context(), cluster.ExecOptions{
				Command:       command,
				Group:         group,
				Hosts:         hosts,
				HostsFile:     file,
				InventoryPath: invPath,
				SSHUser:       sshUser,
				SSHPort:       sshPort,
				SSHKey:        sshKey,
				SSHTimeout:    to,
				Parallel:      par,
				Output:        output,
			})
			if err != nil {
				return err
			}

			if output == "json" {
				b, err := cluster.EncodeExecResultsJSON(res)
				if err != nil {
					return err
				}
				_, err = cmd.OutOrStdout().Write(b)
				if err == nil {
					_, _ = io.WriteString(cmd.OutOrStdout(), "\n")
				}
				return err
			}
			for _, r := range res {
				if r.OK {
					fmt.Fprintf(cmd.OutOrStdout(), "[%s] OK\n%s\n", r.Host, strings.TrimRight(r.Output, "\n"))
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "[%s] ERROR: %s\n%s\n", r.Host, r.Error, strings.TrimRight(r.Output, "\n"))
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&group, "group", "", "Server group to target")
	cmd.Flags().StringSliceVar(&hosts, "hosts", nil, "Specific hosts to target")
	cmd.Flags().StringVar(&file, "file", "", "Read hosts from file")
	cmd.Flags().StringVar(&command, "command", "", "Command to execute (required)")
	cmd.Flags().IntVar(&parallel, "parallel", 0, "Parallel execution limit")
	cmd.Flags().StringVar(&timeout, "timeout", "", "Command timeout")
	cmd.Flags().StringVar(&output, "output", "combined", "Output format (combined, separate, json)")
	return cmd
}

func newClusterDeployCmd(a *app.App) *cobra.Command {
	var (
		configPath string
		target     string
		validate   bool
		backup     bool
		rollback   string
		diff       bool
	)
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy configuration to cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--config", configPath, "--target", target, "--rollback", rollback}
			if validate {
				argv = append(argv, "--validate")
			}
			if backup {
				argv = append(argv, "--backup")
			}
			if diff {
				argv = append(argv, "--diff")
			}
			return a.RunScript(cmd.Context(), "config-deploy.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&configPath, "config", "", "Configuration file/directory")
	cmd.Flags().StringVar(&target, "target", "", "Target path on remote servers")
	cmd.Flags().BoolVar(&validate, "validate", false, "Validate before deploying")
	cmd.Flags().BoolVar(&backup, "backup", false, "Backup existing configuration")
	cmd.Flags().StringVar(&rollback, "rollback", "", "Rollback to previous version")
	cmd.Flags().BoolVar(&diff, "diff", false, "Show differences")
	return cmd
}

func newClusterMonitorCmd(a *app.App) *cobra.Command {
	var (
		dashboard bool
		metrics   []string
		refresh   string
		alerts    bool
		export    string
		group     string
		hosts     []string
		file      string
		output    string
	)
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor cluster health",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = refresh
			_ = dashboard
			_ = alerts

			invPath := getStringFlag(cmd, "inventory-file")
			sshUser := getStringFlag(cmd, "ssh-user")
			sshKey := getStringFlag(cmd, "ssh-key")
			sshPort := getIntFlag(cmd, "ssh-port")
			sshTimeout := time.Duration(getIntFlag(cmd, "ssh-timeout")) * time.Second

			rep, err := cluster.Monitor(cmd.Context(), cluster.MonitorOptions{
				InventoryPath: invPath,
				Group:         group,
				Hosts:         hosts,
				HostsFile:     file,
				SSHUser:       sshUser,
				SSHPort:       sshPort,
				SSHKey:        sshKey,
				SSHTimeout:    sshTimeout,
				Parallel:      4,
				Metrics:       metrics,
				Output:        output,
			})
			if err != nil {
				return err
			}

			if output == "json" || export != "" {
				b, err := cluster.EncodeMonitorReportJSON(rep)
				if err != nil {
					return err
				}
				if export != "" {
					if err := writeFile0600(export, b); err != nil {
						return err
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Exported: %s\n", export)
					return nil
				}
				_, err = cmd.OutOrStdout().Write(b)
				if err == nil {
					_, _ = io.WriteString(cmd.OutOrStdout(), "\n")
				}
				return err
			}

			_, _ = io.WriteString(cmd.OutOrStdout(), cluster.MonitorReportToText(rep))
			return nil
		},
	}
	cmd.Flags().StringVar(&group, "group", "", "Server group to target")
	cmd.Flags().StringSliceVar(&hosts, "hosts", nil, "Specific hosts to target")
	cmd.Flags().StringVar(&file, "file", "", "Read hosts from file")
	cmd.Flags().BoolVar(&dashboard, "dashboard", false, "Launch interactive dashboard")
	cmd.Flags().StringSliceVar(&metrics, "metrics", nil, "Metrics to monitor (cpu, memory, disk, network)")
	cmd.Flags().StringVar(&refresh, "refresh", "", "Refresh interval")
	cmd.Flags().BoolVar(&alerts, "alerts", false, "Show active alerts")
	cmd.Flags().StringVar(&export, "export", "", "Export metrics data")
	cmd.Flags().StringVar(&output, "output", "text", "Output format (text, json)")
	_ = cmd.Flags().MarkHidden("group")
	_ = cmd.Flags().MarkHidden("hosts")
	_ = cmd.Flags().MarkHidden("file")
	_ = cmd.Flags().MarkHidden("output")
	return cmd
}

func newClusterInventoryCmd(a *app.App) *cobra.Command {
	var (
		scan   bool
		add    string
		remove string
		groups bool
		tags   []string
		output string
		export string
	)
	cmd := &cobra.Command{
		Use:   "inventory",
		Short: "Manage server inventory",
		RunE: func(cmd *cobra.Command, args []string) error {
			if export != "" {
				output = export
			}
			// Keep intrusive/mutating operations script-based for now.
			if scan || add != "" || remove != "" || groups || len(tags) > 0 {
				argv := []string{"--add", add, "--remove", remove, "--output", output}
				for _, t := range tags {
					argv = append(argv, "--tags", t)
				}
				if scan {
					argv = append(argv, "--scan")
				}
				if groups {
					argv = append(argv, "--groups")
				}
				return a.RunScript(cmd.Context(), "inventory-manager.sh", argv...)
			}

			invPath := getStringFlag(cmd, "inventory-file")
			inv, err := cluster.LoadInventory(invPath)
			if err != nil {
				return err
			}
			if output == "json" {
				b, err := cluster.InventoryToJSON(inv)
				if err != nil {
					return err
				}
				_, err = cmd.OutOrStdout().Write(b)
				if err == nil {
					_, _ = io.WriteString(cmd.OutOrStdout(), "\n")
				}
				return err
			}
			_, _ = io.WriteString(cmd.OutOrStdout(), cluster.InventoryToText(inv))
			return nil
		},
	}
	cmd.Flags().BoolVar(&scan, "scan", false, "Scan network for servers")
	cmd.Flags().StringVar(&add, "add", "", "Add server to inventory")
	cmd.Flags().StringVar(&remove, "remove", "", "Remove server from inventory")
	cmd.Flags().BoolVar(&groups, "groups", false, "Manage server groups")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Tag servers")
	cmd.Flags().StringVar(&output, "output", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&export, "export", "", "Export inventory")
	_ = cmd.Flags().MarkHidden("output")
	return cmd
}

func getStringFlag(cmd *cobra.Command, name string) string {
	if cmd == nil {
		return ""
	}
	if f := cmd.Flags().Lookup(name); f != nil {
		v, _ := cmd.Flags().GetString(name)
		if v != "" {
			return v
		}
	}
	if f := cmd.InheritedFlags().Lookup(name); f != nil {
		v, _ := cmd.InheritedFlags().GetString(name)
		if v != "" {
			return v
		}
	}
	return getStringFlag(cmd.Parent(), name)
}

func getIntFlag(cmd *cobra.Command, name string) int {
	if cmd == nil {
		return 0
	}
	if f := cmd.Flags().Lookup(name); f != nil {
		v, _ := cmd.Flags().GetInt(name)
		if v != 0 {
			return v
		}
	}
	if f := cmd.InheritedFlags().Lookup(name); f != nil {
		v, _ := cmd.InheritedFlags().GetInt(name)
		if v != 0 {
			return v
		}
	}
	return getIntFlag(cmd.Parent(), name)
}

func newClusterPatchCmd(a *app.App) *cobra.Command {
	var (
		packages          []string
		strategy          string
		batchSize         int
		preCheck          bool
		postCheck         bool
		rollbackOnFailure bool
		group             string
		hosts             []string
		file              string
		apply             bool
		yes               bool
	)
	cmd := &cobra.Command{
		Use:   "patch",
		Short: "Orchestrate patching across cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			invPath := getStringFlag(cmd, "inventory-file")
			sshUser := getStringFlag(cmd, "ssh-user")
			sshKey := getStringFlag(cmd, "ssh-key")
			sshPort := getIntFlag(cmd, "ssh-port")
			sshTimeout := time.Duration(getIntFlag(cmd, "ssh-timeout")) * time.Second

			rep, err := cluster.OrchestratePatches(cmd.Context(), cluster.PatchOptions{
				InventoryPath:     invPath,
				Group:             group,
				Hosts:             hosts,
				HostsFile:         file,
				SSHUser:           sshUser,
				SSHPort:           sshPort,
				SSHKey:            sshKey,
				SSHTimeout:        sshTimeout,
				Parallel:          4,
				Packages:          packages,
				Strategy:          strategy,
				BatchSize:         batchSize,
				PreCheck:          preCheck,
				PostCheck:         postCheck,
				RollbackOnFailure: rollbackOnFailure,
				Apply:             apply,
				Yes:               yes,
			})
			if err != nil {
				return err
			}
			b, err := cluster.EncodePatchReportJSON(rep)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(b)
			if err == nil {
				_, _ = io.WriteString(cmd.OutOrStdout(), "\n")
			}
			return err
		},
	}
	cmd.Flags().StringVar(&group, "group", "", "Server group to target")
	cmd.Flags().StringSliceVar(&hosts, "hosts", nil, "Specific hosts to target")
	cmd.Flags().StringVar(&file, "file", "", "Read hosts from file")
	cmd.Flags().StringSliceVar(&packages, "packages", nil, "Packages to update")
	cmd.Flags().StringVar(&strategy, "strategy", "rolling", "Patching strategy (rolling, parallel, canary)")
	cmd.Flags().IntVar(&batchSize, "batch-size", 0, "Batch size for rolling updates")
	cmd.Flags().BoolVar(&preCheck, "pre-check", false, "Run pre-patch checks")
	cmd.Flags().BoolVar(&postCheck, "post-check", false, "Run post-patch validation")
	cmd.Flags().BoolVar(&rollbackOnFailure, "rollback-on-failure", false, "Auto-rollback on failure")
	cmd.Flags().BoolVar(&apply, "apply", false, "Apply patch plan (requires --yes)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Auto-confirm patch application")
	_ = cmd.Flags().MarkHidden("group")
	_ = cmd.Flags().MarkHidden("hosts")
	_ = cmd.Flags().MarkHidden("file")
	_ = cmd.Flags().MarkHidden("apply")
	_ = cmd.Flags().MarkHidden("yes")
	return cmd
}

func newClusterSyncCmd(a *app.App) *cobra.Command {
	var (
		source      string
		destination string
		deleteExtra bool
		checksum    bool
		dryRun      bool
	)
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize files across nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--source", source, "--destination", destination}
			if deleteExtra {
				argv = append(argv, "--delete")
			}
			if checksum {
				argv = append(argv, "--checksum")
			}
			if dryRun {
				argv = append(argv, "--dry-run")
			}
			return a.RunScript(cmd.Context(), "cluster-sync.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&source, "source", "", "Source file/directory")
	cmd.Flags().StringVar(&destination, "destination", "", "Destination path")
	cmd.Flags().BoolVar(&deleteExtra, "delete", false, "Delete extraneous files")
	cmd.Flags().BoolVar(&checksum, "checksum", false, "Use checksum comparison")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be synced")
	return cmd
}

func newClusterReportCmd(a *app.App) *cobra.Command {
	var (
		format   string
		sections []string
		schedule string
		email    string
	)
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate cluster status report",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--format", format, "--schedule", schedule, "--email", email}
			for _, s := range sections {
				argv = append(argv, "--sections", s)
			}
			return a.RunScript(cmd.Context(), "cluster-report.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&format, "format", "markdown", "Output format (html, pdf, markdown, json)")
	cmd.Flags().StringSliceVar(&sections, "sections", nil, "Report sections (summary, health, inventory, alerts)")
	cmd.Flags().StringVar(&schedule, "schedule", "", "Schedule regular reports")
	cmd.Flags().StringVar(&email, "email", "", "Email report to address")
	return cmd
}

func newClusterAlertCmd(a *app.App) *cobra.Command {
	var (
		add     string
		list    bool
		remove  string
		test    string
		silence string
	)
	cmd := &cobra.Command{
		Use:   "alert",
		Short: "Configure cluster alerts",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--add", add, "--remove", remove, "--test", test, "--silence", silence}
			if list {
				argv = append(argv, "--list")
			}
			return a.RunScript(cmd.Context(), "cluster-alert.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&add, "add", "", "Add new alert rule")
	cmd.Flags().BoolVar(&list, "list", false, "List alert rules")
	cmd.Flags().StringVar(&remove, "remove", "", "Remove alert rule")
	cmd.Flags().StringVar(&test, "test", "", "Test alert rule")
	cmd.Flags().StringVar(&silence, "silence", "", "Silence alerts for duration")
	return cmd
}
