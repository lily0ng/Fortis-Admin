package cli

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/cluster"
)

func newClusterInitCmd(a *app.App) *cobra.Command {
	var out string
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a cluster inventory and SSH setup guidance (safe)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = args
			invPath := getStringFlag(cmd, "inventory-file")
			sshKey := getStringFlag(cmd, "ssh-key")

			res, err := cluster.Init(cluster.InitOptions{InventoryPath: invPath, SSHKeyPath: sshKey, Force: force})
			if err != nil {
				return err
			}
			if out != "" {
				// also write JSON output file
				b, _ := json.MarshalIndent(res, "", "  ")
				_ = writeFile0600(out, b)
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(res)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing inventory file")
	cmd.Flags().StringVar(&out, "output", "", "Write init result JSON to file")
	_ = a
	return cmd
}
