package cli

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/backup"
)

func newBackupSnapshotCmd(a *app.App) *cobra.Command {
	var volume string
	var name string
	var keep int
	var remote string
	var apply bool
	var yes bool

	cmd := &cobra.Command{
		Use:    "snapshot",
		Short:  "Snapshot manager (plan-only by default)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = args
			res, err := backup.ManageSnapshots(cmd.Context(), backup.SnapshotOptions{
				Volume: volume,
				Name:   name,
				Keep:   keep,
				Remote: remote,
				Apply:  apply,
				Yes:    yes,
			})
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(res)
		},
	}
	cmd.Flags().StringVar(&volume, "volume", "", "Volume/dataset to snapshot")
	cmd.Flags().StringVar(&name, "name", "", "Snapshot name")
	cmd.Flags().IntVar(&keep, "keep", 7, "Snapshots to keep")
	cmd.Flags().StringVar(&remote, "remote", "", "Remote sync target (stub)")
	cmd.Flags().BoolVar(&apply, "apply", false, "Apply snapshot operations (requires --yes)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Auto-confirm snapshot operations")
	_ = a
	_ = cmd.MarkFlagRequired("volume")
	return cmd
}
