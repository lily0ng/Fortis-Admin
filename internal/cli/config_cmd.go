package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"fortis-admin/internal/app"
	"fortis-admin/internal/config"
)

func newConfigCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration utilities",
	}
	cmd.GroupID = "utilities"

	cmd.AddCommand(&cobra.Command{
		Use:   "view",
		Short: "View configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := a.Config
			b, err := yaml.Marshal(&cfg)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(b))
			return nil
		},
	})

	var (
		key   string
		value string
	)
	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Set configuration value",
		RunE: func(cmd *cobra.Command, args []string) error {
			if key == "" {
				return errors.New("--key is required")
			}
			if a.ConfigPath == "" {
				return errors.New("--config must be set")
			}

			cfg, err := config.Load(a.ConfigPath)
			if err != nil {
				cfg = config.Default()
			}

			switch key {
			case "log_file":
				cfg.LogFile = value
			case "scripts_dir":
				cfg.ScriptsDir = value
			case "inventory_file":
				cfg.InventoryFile = value
			default:
				return fmt.Errorf("unsupported key: %s", key)
			}

			b, err := yaml.Marshal(&cfg)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(dirOf(a.ConfigPath), 0o755); err != nil {
				return err
			}
			return os.WriteFile(a.ConfigPath, b, 0o644)
		},
	}
	setCmd.Flags().StringVar(&key, "key", "", "Config key (log_file, scripts_dir, inventory_file)")
	setCmd.Flags().StringVar(&value, "value", "", "Config value")
	cmd.AddCommand(setCmd)

	return cmd
}

func dirOf(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			if i == 0 {
				return "/"
			}
			return p[:i]
		}
	}
	return "."
}
