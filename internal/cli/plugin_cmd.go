package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
)

func newPluginCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Plugin system",
	}
	cmd.GroupID = "utilities"

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginsDir := filepath.Join(".", "plugins")
			entries, err := os.ReadDir(pluginsDir)
			if err != nil {
				return err
			}
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				fmt.Fprintln(cmd.OutOrStdout(), e.Name())
			}
			return nil
		},
	})

	runCmd := &cobra.Command{
		Use:   "run [name] [args...]",
		Short: "Run a plugin",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			pluginsDir := filepath.Join(".", "plugins")
			p := filepath.Join(pluginsDir, name)
			if _, err := os.Stat(p); err != nil {
				return fmt.Errorf("plugin not found: %s", name)
			}

			c := exec.Command(p, args[1:]...)
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.ErrOrStderr()
			c.Env = append(os.Environ(), "FORTIS_CONFIG="+a.ConfigPath)
			return c.Run()
		},
	}
	cmd.AddCommand(runCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install a plugin (not implemented)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("plugin install is not implemented yet")
		},
	})

	return cmd
}
