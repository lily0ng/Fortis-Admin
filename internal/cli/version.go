package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/buildinfo"
)

func newVersionCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run: func(cmd *cobra.Command, args []string) {
			printVersionLong(cmd.OutOrStdout())
		},
	}
	cmd.GroupID = "utilities"
	return cmd
}

func printVersionLong(w io.Writer) {
	fmt.Fprintf(w, "Fortis-Admin v%s\n", buildinfo.Version)
	fmt.Fprintf(w, "Build Date: %s\n", buildinfo.BuildDate)
	fmt.Fprintf(w, "Commit: %s\n", buildinfo.Commit)
	fmt.Fprintf(w, "Go Version: %s\n", buildinfo.GoVersion())
	fmt.Fprintf(w, "Platform: %s\n\n", buildinfo.Platform())

	fmt.Fprintln(w, "Features Enabled:")
	fmt.Fprintln(w, "✓ Server Hardening Automation")
	fmt.Fprintln(w, "✓ Incident Response Toolkit")
	fmt.Fprintln(w, "✓ Backup & Recovery System")
	fmt.Fprintln(w, "✓ Multi-Server Management")
	fmt.Fprintln(w, "✓ Plugin System")
	fmt.Fprintln(w, "✓ API Server")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "License: Apache 2.0")
	fmt.Fprintln(w, "Documentation: https://fortis-admin.readthedocs.io")
	fmt.Fprintln(w, "Report issues: https://github.com/lily0ng/Fortis-Admin/issues")
}
