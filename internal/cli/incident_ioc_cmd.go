package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/incident"
)

func newIncidentIOCCmd(a *app.App) *cobra.Command {
	var storePath string
	var format string

	cmd := &cobra.Command{
		Use:   "ioc",
		Short: "Indicator of Compromise management",
	}
	cmd.Flags().StringVar(&storePath, "store", incident.DefaultIOCStorePath(), "IOC store path")
	cmd.Flags().StringVar(&format, "format", "json", "Output format (json, text)")

	cmd.AddCommand(newIncidentIOCListCmd(&storePath, &format))
	cmd.AddCommand(newIncidentIOCAddCmd(&storePath, &format))
	cmd.AddCommand(newIncidentIOCRemoveCmd(&storePath, &format))
	cmd.AddCommand(newIncidentIOCImportCmd(&storePath, &format))
	cmd.AddCommand(newIncidentIOCExportCmd(&storePath))

	_ = a
	return cmd
}

func newIncidentIOCListCmd(storePath *string, format *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List IOCs",
		RunE: func(cmd *cobra.Command, args []string) error {
			store := incident.IOCStore{Path: *storePath}
			iocs, err := store.Load()
			if err != nil {
				return err
			}

			switch strings.ToLower(strings.TrimSpace(*format)) {
			case "text":
				for _, it := range iocs {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n", it.Value, it.Type, it.Source)
				}
				return nil
			default:
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(iocs)
			}
		},
	}
}

func newIncidentIOCAddCmd(storePath *string, format *string) *cobra.Command {
	var value string
	var typ string
	var source string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an IOC",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(value) == "" {
				return errors.New("--value is required")
			}
			store := incident.IOCStore{Path: *storePath}
			iocs, err := store.Load()
			if err != nil {
				return err
			}
			iocs = incident.AddIOC(iocs, value, typ, source)
			if err := store.Save(iocs); err != nil {
				return err
			}
			return newIncidentIOCListCmd(storePath, format).RunE(cmd, args)
		},
	}
	cmd.Flags().StringVar(&value, "value", "", "IOC value")
	cmd.Flags().StringVar(&typ, "type", "generic", "IOC type")
	cmd.Flags().StringVar(&source, "source", "manual", "IOC source")
	return cmd
}

func newIncidentIOCRemoveCmd(storePath *string, format *string) *cobra.Command {
	var value string
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an IOC",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(value) == "" {
				return errors.New("--value is required")
			}
			store := incident.IOCStore{Path: *storePath}
			iocs, err := store.Load()
			if err != nil {
				return err
			}
			iocs = incident.RemoveIOC(iocs, value)
			if err := store.Save(iocs); err != nil {
				return err
			}
			return newIncidentIOCListCmd(storePath, format).RunE(cmd, args)
		},
	}
	cmd.Flags().StringVar(&value, "value", "", "IOC value")
	return cmd
}

func newIncidentIOCImportCmd(storePath *string, format *string) *cobra.Command {
	var file string
	var typ string
	var source string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import IOCs from a text file (one per line)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(file) == "" {
				return errors.New("--file is required")
			}
			imported, err := incident.ImportIOCsFromTextFile(file, typ, source)
			if err != nil {
				return err
			}
			store := incident.IOCStore{Path: *storePath}
			iocs, err := store.Load()
			if err != nil {
				return err
			}
			for _, it := range imported {
				iocs = incident.AddIOC(iocs, it.Value, it.Type, it.Source)
			}
			if err := store.Save(iocs); err != nil {
				return err
			}
			return newIncidentIOCListCmd(storePath, format).RunE(cmd, args)
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Input file")
	cmd.Flags().StringVar(&typ, "type", "generic", "IOC type")
	cmd.Flags().StringVar(&source, "source", "import", "IOC source")
	return cmd
}

func newIncidentIOCExportCmd(storePath *string) *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export IOCs to a text file (one per line)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(out) == "" {
				return errors.New("--output is required")
			}
			store := incident.IOCStore{Path: *storePath}
			iocs, err := store.Load()
			if err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
				return err
			}
			lines := make([]string, 0, len(iocs))
			for _, it := range iocs {
				lines = append(lines, it.Value)
			}
			return os.WriteFile(out, []byte(strings.Join(lines, "\n")+"\n"), 0o600)
		},
	}
	cmd.Flags().StringVar(&out, "output", "", "Output file")
	return cmd
}
