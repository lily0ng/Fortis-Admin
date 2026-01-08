package cli

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
)

func newLogsCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Log utilities",
	}
	cmd.GroupID = "utilities"

	var file string
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show application logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := file
			if p == "" {
				p = a.Config.LogFile
			}
			f, err := os.Open(p)
			if err != nil {
				return err
			}
			defer f.Close()

			s := bufio.NewScanner(f)
			for s.Scan() {
				_, _ = cmd.OutOrStdout().Write(append(s.Bytes(), '\n'))
			}
			return s.Err()
		},
	}
	showCmd.Flags().StringVar(&file, "file", "", "Log file path")
	cmd.AddCommand(showCmd)

	var follow bool
	tailCmd := &cobra.Command{
		Use:   "tail",
		Short: "Tail application logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := file
			if p == "" {
				p = a.Config.LogFile
			}
			if p == "" {
				return errors.New("log file is empty")
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tailArgs := []string{"-n", "200"}
			if follow {
				tailArgs = append(tailArgs, "-f")
			}
			tailArgs = append(tailArgs, p)

			c := exec.CommandContext(ctx, "tail", tailArgs...)
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.ErrOrStderr()
			return c.Run()
		},
	}
	tailCmd.Flags().StringVar(&file, "file", "", "Log file path")
	tailCmd.Flags().BoolVar(&follow, "follow", true, "Follow log output")
	_ = tailCmd.Flags().MarkHidden("follow")

	cmd.AddCommand(tailCmd)
	return cmd
}
