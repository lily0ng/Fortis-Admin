package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/buildinfo"
	"fortis-admin/internal/ui"
)

func NewRootCmd(out, errOut io.Writer) *cobra.Command {
	a := app.New()

	root := &cobra.Command{
		Use:          "fortis",
		Short:        "Fortis-Admin - System Administration & Security Toolkit",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if a.ConfigPath == "" {
				a.ConfigPath = "/etc/fortis/config.yaml"
			}
			return a.Init()
		},
	}

	root.AddGroup(
		&cobra.Group{ID: "harden", Title: "SERVER HARDENING:"},
		&cobra.Group{ID: "incident", Title: "INCIDENT RESPONSE:"},
		&cobra.Group{ID: "backup", Title: "BACKUP & RECOVERY:"},
		&cobra.Group{ID: "cluster", Title: "CLUSTER MANAGEMENT:"},
		&cobra.Group{ID: "utilities", Title: "UTILITIES:"},
	)

	root.SetOut(out)
	root.SetErr(errOut)

	root.PersistentFlags().StringVarP(&a.ConfigPath, "config", "c", "/etc/fortis/config.yaml", "Configuration file")
	root.PersistentFlags().BoolVar(&a.Debug, "debug", false, "Enable debug mode")
	root.PersistentFlags().BoolVarP(&a.Quiet, "quiet", "q", false, "Quiet mode (minimal output)")
	root.PersistentFlags().BoolVarP(&a.Verbose, "verbose", "v", false, "Verbose output")
	root.PersistentFlags().BoolVar(&forceColor, "color", false, "Force color output")
	root.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")

	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		ui.Banner(cmd.OutOrStdout(), "FORTIS-ADMIN "+buildinfo.Version)
		_ = cmd.HelpFunc()(cmd, args)
	})

	root.AddCommand(newVersionCmd(a))
	root.AddCommand(newCompletionCmd())

	root.AddCommand(newConfigCmd(a))
	root.AddCommand(newLogsCmd(a))
	root.AddCommand(newUpdateCmd(a))

	root.AddCommand(newHardenCmd(a))
	root.AddCommand(newIncidentCmd(a))
	root.AddCommand(newBackupCmd(a))
	root.AddCommand(newClusterCmd(a))
	root.AddCommand(newPluginCmd(a))

	root.SetVersionTemplate(fmt.Sprintf("Fortis-Admin %s\n", buildinfo.Version))

	return root
}

var (
	forceColor bool
	noColor    bool
)

func colorEnabled() bool {
	if noColor {
		return false
	}
	if forceColor {
		return true
	}
	_, ok := os.LookupEnv("NO_COLOR")
	return !ok
}
