package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
)

func newClusterCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Multi-server management platform",
	}
	cmd.GroupID = "cluster"

	cmd.AddCommand(newClusterExecCmd(a))
	cmd.AddCommand(newClusterDeployCmd(a))
	cmd.AddCommand(newClusterMonitorCmd(a))
	cmd.AddCommand(newClusterInventoryCmd(a))
	cmd.AddCommand(newClusterPatchCmd(a))
	cmd.AddCommand(newClusterSyncCmd(a))
	cmd.AddCommand(newClusterReportCmd(a))
	cmd.AddCommand(newClusterAlertCmd(a))

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
		Use:   "exec",
		Short: "Execute command on multiple servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if command == "" {
				return errors.New("--command is required")
			}
			argv := []string{"--group", group, "--file", file, "--command", command, "--parallel", itoa(parallel), "--timeout", timeout, "--output", output}
			for _, h := range hosts {
				argv = append(argv, "--hosts", h)
			}
			return a.RunScript(cmd.Context(), "cluster-exec.sh", argv...)
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
	)
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor cluster health",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--refresh", refresh, "--export", export}
			for _, m := range metrics {
				argv = append(argv, "--metrics", m)
			}
			if dashboard {
				argv = append(argv, "--dashboard")
			}
			if alerts {
				argv = append(argv, "--alerts")
			}
			return a.RunScript(cmd.Context(), "cluster-monitor.sh", argv...)
		},
	}
	cmd.Flags().BoolVar(&dashboard, "dashboard", false, "Launch interactive dashboard")
	cmd.Flags().StringSliceVar(&metrics, "metrics", nil, "Metrics to monitor (cpu, memory, disk, network)")
	cmd.Flags().StringVar(&refresh, "refresh", "", "Refresh interval")
	cmd.Flags().BoolVar(&alerts, "alerts", false, "Show active alerts")
	cmd.Flags().StringVar(&export, "export", "", "Export metrics data")
	return cmd
}

func newClusterInventoryCmd(a *app.App) *cobra.Command {
	var (
		scan   bool
		add    string
		remove string
		groups bool
		tags   []string
		export string
	)
	cmd := &cobra.Command{
		Use:   "inventory",
		Short: "Manage server inventory",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--add", add, "--remove", remove, "--export", export}
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
		},
	}
	cmd.Flags().BoolVar(&scan, "scan", false, "Scan network for servers")
	cmd.Flags().StringVar(&add, "add", "", "Add server to inventory")
	cmd.Flags().StringVar(&remove, "remove", "", "Remove server from inventory")
	cmd.Flags().BoolVar(&groups, "groups", false, "Manage server groups")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Tag servers")
	cmd.Flags().StringVar(&export, "export", "", "Export inventory")
	return cmd
}

func newClusterPatchCmd(a *app.App) *cobra.Command {
	var (
		packages          []string
		strategy          string
		batchSize         int
		preCheck          bool
		postCheck         bool
		rollbackOnFailure bool
	)
	cmd := &cobra.Command{
		Use:   "patch",
		Short: "Orchestrate patching across cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--strategy", strategy, "--batch-size", itoa(batchSize)}
			for _, p := range packages {
				argv = append(argv, "--packages", p)
			}
			if preCheck {
				argv = append(argv, "--pre-check")
			}
			if postCheck {
				argv = append(argv, "--post-check")
			}
			if rollbackOnFailure {
				argv = append(argv, "--rollback-on-failure")
			}
			return a.RunScript(cmd.Context(), "patch-orchestrator.sh", argv...)
		},
	}
	cmd.Flags().StringSliceVar(&packages, "packages", nil, "Packages to update")
	cmd.Flags().StringVar(&strategy, "strategy", "rolling", "Patching strategy (rolling, parallel, canary)")
	cmd.Flags().IntVar(&batchSize, "batch-size", 0, "Batch size for rolling updates")
	cmd.Flags().BoolVar(&preCheck, "pre-check", false, "Run pre-patch checks")
	cmd.Flags().BoolVar(&postCheck, "post-check", false, "Run post-patch validation")
	cmd.Flags().BoolVar(&rollbackOnFailure, "rollback-on-failure", false, "Auto-rollback on failure")
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
