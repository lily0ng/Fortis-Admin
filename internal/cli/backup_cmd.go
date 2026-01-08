package cli

import (
	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
)

func newBackupCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup & recovery system",
	}
	cmd.GroupID = "backup"

	cmd.AddCommand(newBackupCreateCmd(a))
	cmd.AddCommand(newBackupListCmd(a))
	cmd.AddCommand(newBackupVerifyCmd(a))
	cmd.AddCommand(newBackupRestoreCmd(a))
	cmd.AddCommand(newBackupScheduleCmd(a))
	cmd.AddCommand(newBackupCatalogCmd(a))
	cmd.AddCommand(newBackupMonitorCmd(a))
	cmd.AddCommand(newBackupTestDRCmd(a))

	return cmd
}

func newBackupCreateCmd(a *app.App) *cobra.Command {
	var (
		target   string
		sources  []string
		btype    string
		exclude  []string
		encrypt  bool
		compress string
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--target", target, "--type", btype, "--compress", compress}
			for _, s := range sources {
				argv = append(argv, "--source", s)
			}
			for _, e := range exclude {
				argv = append(argv, "--exclude", e)
			}
			if encrypt {
				argv = append(argv, "--encrypt")
			}
			return a.RunScript(cmd.Context(), "backup-create.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&target, "target", "", "Backup target directory")
	cmd.Flags().StringSliceVar(&sources, "source", nil, "Source directories/files")
	cmd.Flags().StringVar(&btype, "type", "full", "Backup type (full, incremental, differential)")
	cmd.Flags().StringSliceVar(&exclude, "exclude", nil, "Patterns to exclude")
	cmd.Flags().BoolVar(&encrypt, "encrypt", false, "Encrypt backup")
	cmd.Flags().StringVar(&compress, "compress", "gzip", "Compression algorithm (gzip, zstd, lz4, none)")
	return cmd
}

func newBackupListCmd(a *app.App) *cobra.Command {
	var (
		detailed bool
		sortBy   string
		filter   string
		jsonOut  bool
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available backups",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--sort", sortBy, "--filter", filter}
			if detailed {
				argv = append(argv, "--detailed")
			}
			if jsonOut {
				argv = append(argv, "--json")
			}
			return a.RunScript(cmd.Context(), "backup-list.sh", argv...)
		},
	}
	cmd.Flags().BoolVar(&detailed, "detailed", false, "Show detailed information")
	cmd.Flags().StringVar(&sortBy, "sort", "date", "Sort by (date, size, name)")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter backups by pattern")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output in JSON format")
	return cmd
}

func newBackupVerifyCmd(a *app.App) *cobra.Command {
	var (
		backup string
		quick  bool
		full   bool
		repair bool
	)
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify backup integrity",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--backup", backup}
			if quick {
				argv = append(argv, "--quick")
			}
			if full {
				argv = append(argv, "--full")
			}
			if repair {
				argv = append(argv, "--repair")
			}
			return a.RunScript(cmd.Context(), "backup-verify.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&backup, "backup", "", "Backup to verify")
	cmd.Flags().BoolVar(&quick, "quick", false, "Quick verification (checksums only)")
	cmd.Flags().BoolVar(&full, "full", false, "Full verification (restore test)")
	cmd.Flags().BoolVar(&repair, "repair", false, "Attempt to repair corrupt backups")
	return cmd
}

func newBackupRestoreCmd(a *app.App) *cobra.Command {
	var (
		backup string
		target string
		items  []string
		timePt string
		dryRun bool
	)
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore from backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--backup", backup, "--target", target, "--time", timePt}
			for _, it := range items {
				argv = append(argv, "--items", it)
			}
			if dryRun {
				argv = append(argv, "--dry-run")
			}
			return a.RunScript(cmd.Context(), "backup-restore.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&backup, "backup", "", "Backup to restore from")
	cmd.Flags().StringVar(&target, "target", "", "Restore target location")
	cmd.Flags().StringSliceVar(&items, "items", nil, "Specific items to restore")
	cmd.Flags().StringVar(&timePt, "time", "", "Point-in-time recovery")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulation mode")
	return cmd
}

func newBackupScheduleCmd(a *app.App) *cobra.Command {
	var (
		add     string
		list    bool
		remove  string
		enable  string
		disable string
		runNow  string
	)
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Manage backup schedules",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--add", add, "--remove", remove, "--enable", enable, "--disable", disable, "--run-now", runNow}
			if list {
				argv = append(argv, "--list")
			}
			return a.RunScript(cmd.Context(), "backup-scheduler.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&add, "add", "", "Add new schedule")
	cmd.Flags().BoolVar(&list, "list", false, "List schedules")
	cmd.Flags().StringVar(&remove, "remove", "", "Remove schedule")
	cmd.Flags().StringVar(&enable, "enable", "", "Enable schedule")
	cmd.Flags().StringVar(&disable, "disable", "", "Disable schedule")
	cmd.Flags().StringVar(&runNow, "run-now", "", "Run schedule immediately")
	return cmd
}

func newBackupCatalogCmd(a *app.App) *cobra.Command {
	var (
		backup  string
		search  string
		tree    bool
		stats   bool
		extract string
	)
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "Browse backup contents",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--backup", backup, "--search", search, "--extract", extract}
			if tree {
				argv = append(argv, "--tree")
			}
			if stats {
				argv = append(argv, "--stats")
			}
			return a.RunScript(cmd.Context(), "backup-catalog.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&backup, "backup", "", "Backup to examine")
	cmd.Flags().StringVar(&search, "search", "", "Search for files")
	cmd.Flags().BoolVar(&tree, "tree", false, "Show directory tree")
	cmd.Flags().BoolVar(&stats, "stats", false, "Show backup statistics")
	cmd.Flags().StringVar(&extract, "extract", "", "Extract specific file")
	return cmd
}

func newBackupMonitorCmd(a *app.App) *cobra.Command {
	var (
		watch   bool
		alerts  bool
		stats   bool
		history bool
	)
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor backup status",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{}
			if watch {
				argv = append(argv, "--watch")
			}
			if alerts {
				argv = append(argv, "--alerts")
			}
			if stats {
				argv = append(argv, "--stats")
			}
			if history {
				argv = append(argv, "--history")
			}
			return a.RunScript(cmd.Context(), "backup-monitor.sh", argv...)
		},
	}
	cmd.Flags().BoolVar(&watch, "watch", false, "Watch mode")
	cmd.Flags().BoolVar(&alerts, "alerts", false, "Configure alert thresholds")
	cmd.Flags().BoolVar(&stats, "stats", false, "Show performance statistics")
	cmd.Flags().BoolVar(&history, "history", false, "Show backup history")
	return cmd
}

func newBackupTestDRCmd(a *app.App) *cobra.Command {
	var (
		scenario    string
		environment string
		automated   bool
		report      bool
	)
	cmd := &cobra.Command{
		Use:   "test-dr",
		Short: "Test disaster recovery",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--scenario", scenario, "--environment", environment}
			if automated {
				argv = append(argv, "--automated")
			}
			if report {
				argv = append(argv, "--report")
			}
			return a.RunScript(cmd.Context(), "disaster-recovery.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&scenario, "scenario", "", "DR test scenario")
	cmd.Flags().StringVar(&environment, "environment", "", "Test environment")
	cmd.Flags().BoolVar(&automated, "automated", false, "Automated test")
	cmd.Flags().BoolVar(&report, "report", false, "Generate test report")
	return cmd
}
