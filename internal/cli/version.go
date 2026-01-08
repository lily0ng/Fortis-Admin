package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/buildinfo"
)

func newVersionCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "Fortis-Admin %s\n", buildinfo.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "Build Date: %s\n", buildinfo.BuildDate)
			fmt.Fprintf(cmd.OutOrStdout(), "Commit: %s\n", buildinfo.Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "Go Version: %s\n", buildinfo.GoVersion())
			fmt.Fprintf(cmd.OutOrStdout(), "Platform: %s\n", buildinfo.Platform())
		},
	}
	cmd.GroupID = "utilities"
	return cmd
}
