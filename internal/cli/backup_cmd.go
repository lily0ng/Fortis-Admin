package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/backup"
)

func newBackupCmd(a *app.App) *cobra.Command {
	var (
		retention string
		threads   int
		bandwidth string
		resume    bool
	)

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup & recovery system",
	}
	cmd.GroupID = "backup"
	cmd.Flags().StringVar(&retention, "retention", "", "Retention policy (e.g., \"30d\", \"12M\", \"2y\")")
	cmd.Flags().IntVar(&threads, "threads", 0, "Number of parallel threads")
	cmd.Flags().StringVar(&bandwidth, "bandwidth", "", "Bandwidth limit (e.g., \"10M\", \"1G\")")
	cmd.Flags().BoolVar(&resume, "resume", false, "Resume interrupted backup")

	cmd.AddCommand(newBackupCreateCmd(a))
	cmd.AddCommand(newBackupListCmd(a))
	cmd.AddCommand(newBackupVerifyCmd(a))
	cmd.AddCommand(newBackupRestoreCmd(a))
	cmd.AddCommand(newBackupRestoreWizardCmd(a))
	cmd.AddCommand(newBackupScheduleCmd(a))
	cmd.AddCommand(newBackupCatalogCmd(a))
	cmd.AddCommand(newBackupSnapshotCmd(a))
	cmd.AddCommand(newBackupMonitorCmd(a))
	cmd.AddCommand(newBackupTestDRCmd(a))
	setGroupHelp(cmd, "BACKUP & RECOVERY COMMANDS", "fortis backup [command] [flags]", func(w io.Writer) {
		_ = retention
		_ = threads
		_ = bandwidth
		_ = resume

		io.WriteString(w, "COMMANDS:\n")
		io.WriteString(w, "  create [flags]                  Create new backup\n")
		io.WriteString(w, "    --target string               Backup target directory\n")
		io.WriteString(w, "    --source strings              Source directories/files\n")
		io.WriteString(w, "    --type string                 Backup type (full, incremental, differential)\n")
		io.WriteString(w, "    --exclude strings             Patterns to exclude\n")
		io.WriteString(w, "    --encrypt                     Encrypt backup\n")
		io.WriteString(w, "    --compress string             Compression algorithm (gzip, zstd, lz4, none)\n\n")

		io.WriteString(w, "  list [flags]                    List available backups\n")
		io.WriteString(w, "    --detailed                    Show detailed information\n")
		io.WriteString(w, "    --sort string                 Sort by (date, size, name)\n")
		io.WriteString(w, "    --filter string               Filter backups by pattern\n")
		io.WriteString(w, "    --json                        Output in JSON format\n\n")

		io.WriteString(w, "  verify [flags]                  Verify backup integrity\n")
		io.WriteString(w, "    --backup string               Backup to verify\n")
		io.WriteString(w, "    --quick                       Quick verification (checksums only)\n")
		io.WriteString(w, "    --full                        Full verification (restore test)\n")
		io.WriteString(w, "    --repair                      Attempt to repair corrupt backups\n\n")

		io.WriteString(w, "  restore [flags]                 Restore from backup\n")
		io.WriteString(w, "    --backup string               Backup to restore from\n")
		io.WriteString(w, "    --target string               Restore target location\n")
		io.WriteString(w, "    --items strings               Specific items to restore\n")
		io.WriteString(w, "    --time string                 Point-in-time recovery\n")
		io.WriteString(w, "    --dry-run                     Simulation mode\n\n")

		io.WriteString(w, "  schedule [flags]                Manage backup schedules\n")
		io.WriteString(w, "    --add string                  Add new schedule\n")
		io.WriteString(w, "    --list                        List schedules\n")
		io.WriteString(w, "    --remove string               Remove schedule\n")
		io.WriteString(w, "    --enable string               Enable schedule\n")
		io.WriteString(w, "    --disable string              Disable schedule\n")
		io.WriteString(w, "    --run-now string              Run schedule immediately\n\n")

		io.WriteString(w, "  catalog [flags]                 Browse backup contents\n")
		io.WriteString(w, "    --backup string               Backup to examine\n")
		io.WriteString(w, "    --search string               Search for files\n")
		io.WriteString(w, "    --tree                        Show directory tree\n")
		io.WriteString(w, "    --stats                       Show backup statistics\n")
		io.WriteString(w, "    --extract string              Extract specific file\n\n")

		io.WriteString(w, "  monitor [flags]                 Monitor backup status\n")
		io.WriteString(w, "    --watch                       Watch mode (continuous monitoring)\n")
		io.WriteString(w, "    --alerts                      Configure alert thresholds\n")
		io.WriteString(w, "    --stats                       Show performance statistics\n")
		io.WriteString(w, "    --history                     Show backup history\n\n")

		io.WriteString(w, "  test-dr [flags]                 Test disaster recovery\n")
		io.WriteString(w, "    --scenario string             DR test scenario\n")
		io.WriteString(w, "    --environment string          Test environment\n")
		io.WriteString(w, "    --automated                   Automated test\n")
		io.WriteString(w, "    --report                      Generate test report\n\n")

		io.WriteString(w, "FLAGS:\n")
		io.WriteString(w, "  --retention string            Retention policy (e.g., \"30d\", \"12M\", \"2y\")\n")
		io.WriteString(w, "  --threads int                 Number of parallel threads\n")
		io.WriteString(w, "  --bandwidth string            Bandwidth limit (e.g., \"10M\", \"1G\")\n")
		io.WriteString(w, "  --resume                      Resume interrupted backup\n\n")

		io.WriteString(w, "EXAMPLES:\n")
		io.WriteString(w, "  fortis backup create --source /home /etc --target /backups --encrypt\n")
		io.WriteString(w, "  fortis backup list --detailed --sort date\n")
		io.WriteString(w, "  fortis backup restore --backup backup-2024-01-01 --target /recovery\n")
		io.WriteString(w, "  fortis backup schedule --add \"daily at 2am\"\n")
		io.WriteString(w, "  fortis backup test-dr --scenario full-restore --automated\n")
	})

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
	_ = a
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(sources) == 0 {
				return errors.New("--source is required")
			}
			if strings.TrimSpace(target) == "" {
				target = "./backups"
			}
			meta, err := backup.Create(backup.CreateOptions{
				TargetDir: target,
				Sources:   sources,
				Type:      backup.BackupType(btype),
				Exclude:   exclude,
				Encrypt:   encrypt,
				Compress:  backup.Compression(compress),
			})
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(meta)
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
		target   string
		detailed bool
		sortBy   string
		filter   string
		jsonOut  bool
	)
	_ = a
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available backups",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(target) == "" {
				target = "./backups"
			}
			items, err := backup.List(backup.ListOptions{TargetDir: target, Detailed: detailed, SortBy: sortBy, Filter: filter, JSON: jsonOut})
			if err != nil {
				return err
			}
			if jsonOut {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(items)
			}
			for _, it := range items {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%d\t%s\n", it.ID, it.CreatedAt.Format(time.RFC3339), it.SizeBytes, it.Archive)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&target, "target", "", "Backup target directory")
	cmd.Flags().BoolVar(&detailed, "detailed", false, "Show detailed information")
	cmd.Flags().StringVar(&sortBy, "sort", "date", "Sort by (date, size, name)")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter backups by pattern")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output in JSON format")
	return cmd
}

func newBackupVerifyCmd(a *app.App) *cobra.Command {
	var (
		backupPath string
		quick      bool
		full       bool
		repair     bool
	)
	_ = a
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify backup integrity",
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := backup.Verify(backup.VerifyOptions{BackupPath: backupPath, Quick: quick, Full: full, Repair: repair})
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(res)
		},
	}
	cmd.Flags().StringVar(&backupPath, "backup", "", "Backup to verify")
	cmd.Flags().BoolVar(&quick, "quick", false, "Quick verification (checksums only)")
	cmd.Flags().BoolVar(&full, "full", false, "Full verification (restore test)")
	cmd.Flags().BoolVar(&repair, "repair", false, "Attempt to repair corrupt backups")
	return cmd
}

func newBackupRestoreCmd(a *app.App) *cobra.Command {
	var (
		backupPath string
		target     string
		items      []string
		timePt     string
		dryRun     bool
	)
	_ = a
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore from backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = timePt
			if err := backup.Restore(backup.RestoreOptions{BackupPath: backupPath, TargetDir: target, Items: items, DryRun: dryRun}); err != nil {
				return err
			}
			if dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), "Dry-run restore simulation completed")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Restore completed to: %s\n", target)
			return nil
		},
	}
	cmd.Flags().StringVar(&backupPath, "backup", "", "Backup to restore from")
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
	_ = a
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
		backupPath string
		search     string
		tree       bool
		stats      bool
		extract    string
	)
	_ = a
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "Browse backup contents",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = tree
			_ = stats
			_ = extract
			entries, err := backup.Catalog(backup.CatalogOptions{BackupPath: backupPath, Search: search})
			if err != nil {
				return err
			}
			for _, e := range entries {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%d\n", e.Path, e.Size)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&backupPath, "backup", "", "Backup to examine")
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
		backupPath  string
		target      string
		dryRun      bool
	)
	_ = a
	cmd := &cobra.Command{
		Use:   "test-dr",
		Short: "Test disaster recovery",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = scenario
			_ = environment
			_ = automated
			_ = report
			res, err := backup.TestDR(backup.TestDROptions{BackupPath: backupPath, TargetDir: target, DryRun: dryRun})
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(res)
		},
	}
	cmd.Flags().StringVar(&scenario, "scenario", "", "DR test scenario")
	cmd.Flags().StringVar(&environment, "environment", "", "Test environment")
	cmd.Flags().BoolVar(&automated, "automated", false, "Automated test")
	cmd.Flags().BoolVar(&report, "report", false, "Generate test report")
	cmd.Flags().StringVar(&backupPath, "backup", "", "Backup to restore from")
	cmd.Flags().StringVar(&target, "target", "", "Restore target location")
	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "Simulation mode")
	return cmd
}
