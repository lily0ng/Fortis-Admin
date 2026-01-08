package cli

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/ui"
)

func NewRootCmd(out, errOut io.Writer) *cobra.Command {
	a := app.New()
	var showVersion bool

	root := &cobra.Command{
		Use:          "fortis",
		Short:        "Fortis-Admin - System Administration & Security Toolkit",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				printVersionLong(cmd.OutOrStdout())
				return
			}
			printRootHelp(cmd.OutOrStdout())
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if a.ConfigPath == "" {
				a.ConfigPath = "/etc/fortis/config.yaml"
			}
			if forceColor {
				a.ColorMode = ui.ColorAlways
			} else if noColor {
				a.ColorMode = ui.ColorNever
			} else {
				a.ColorMode = ui.ColorAuto
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
	root.PersistentFlags().BoolVarP(&a.Debug, "debug", "d", false, "Enable debug mode")
	root.PersistentFlags().BoolVarP(&a.Quiet, "quiet", "q", false, "Quiet mode (minimal output)")
	root.PersistentFlags().BoolVarP(&a.Verbose, "verbose", "v", false, "Verbose output")
	root.PersistentFlags().BoolVar(&forceColor, "color", false, "Force color output")
	root.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	root.PersistentFlags().BoolVar(&showVersion, "version", false, "Display version information")
	setRootHelp(root)

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
