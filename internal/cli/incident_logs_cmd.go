package cli

import (
	"io"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/incident"
)

func newIncidentLogsCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "logs",
		Short:  "Log analyzer",
		Hidden: true,
	}
	cmd.AddCommand(newIncidentLogsAnalyzeCmd(a))
	return cmd
}

func newIncidentLogsAnalyzeCmd(a *app.App) *cobra.Command {
	var input string
	var iocFile string
	var output string

	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze logs and search for IOC matches",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = args
			out, err := incident.AnalyzeLogs(cmd.Context(), incident.LogAnalyzeOptions{
				InputPath: input,
				IOCFile:   iocFile,
				Output:    output,
			})
			if err != nil {
				return err
			}
			_, _ = io.WriteString(cmd.OutOrStdout(), out+"\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&input, "input", "", "Log file or directory")
	cmd.Flags().StringVar(&iocFile, "ioc", "", "IOC file")
	cmd.Flags().StringVar(&output, "output", "", "Output report file (json)")

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("ioc")
	_ = a
	return cmd
}
