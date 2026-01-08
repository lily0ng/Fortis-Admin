package cli

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/backup"
)

func newBackupRestoreWizardCmd(a *app.App) *cobra.Command {
	var backupPath string
	var target string
	var items []string
	var dryRun bool
	var interactive bool

	cmd := &cobra.Command{
		Use:    "restore-wizard",
		Short:  "Interactive restore wizard (safe-by-default)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = args
			res, err := backup.RunRestoreWizard(cmd.Context(), backup.RestoreWizardOptions{
				BackupPath:  backupPath,
				TargetDir:   target,
				Items:       items,
				DryRun:      dryRun,
				Interactive: interactive,
			})
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(res)
		},
	}
	cmd.Flags().StringVar(&backupPath, "backup", "", "Backup archive path")
	cmd.Flags().StringVar(&target, "target", "", "Restore target directory")
	cmd.Flags().StringSliceVar(&items, "items", nil, "Items/prefixes to restore")
	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "Simulation mode")
	cmd.Flags().BoolVar(&interactive, "interactive", true, "Prompt for confirmation")
	_ = a
	_ = cmd.MarkFlagRequired("backup")
	return cmd
}
