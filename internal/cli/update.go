package cli

import (
	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
)

func newUpdateCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update fortis to latest version",
		Run: func(cmd *cobra.Command, args []string) {
			a.Log.Infof("update is not implemented yet")
		},
	}
	cmd.GroupID = "utilities"
	return cmd
}
