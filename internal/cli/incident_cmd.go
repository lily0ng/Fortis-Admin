package cli

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"

	"fortis-admin/internal/app"
	"fortis-admin/internal/incident"
)

func newIncidentCmd(a *app.App) *cobra.Command {
	var (
		chainOfCustody bool
		encrypt        bool
		verbose        bool
		quiet          bool
	)

	cmd := &cobra.Command{
		Use:   "incident",
		Short: "Incident response toolkit",
	}
	cmd.GroupID = "incident"
	cmd.Flags().BoolVar(&chainOfCustody, "chain-of-custody", false, "Maintain chain of custody")
	cmd.Flags().BoolVar(&encrypt, "encrypt", false, "Encrypt sensitive data")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Show detailed output")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "Minimal output")

	cmd.AddCommand(newIncidentCaptureCmd(a))
	cmd.AddCommand(newIncidentTriageCmd(a))
	cmd.AddCommand(newIncidentAnalyzeCmd(a))
	cmd.AddCommand(newIncidentHuntCmd(a))
	cmd.AddCommand(newIncidentReportCmd(a))
	cmd.AddCommand(newIncidentTimelineCmd(a))
	cmd.AddCommand(newIncidentContainCmd(a))
	cmd.AddCommand(newIncidentEradicateCmd(a))
	cmd.AddCommand(newIncidentNetworkForensicsCmd(a))
	cmd.AddCommand(newIncidentIOCCmd(a))
	cmd.AddCommand(newIncidentLogsCmd(a))
	setGroupHelp(cmd, "INCIDENT RESPONSE COMMANDS", "fortis incident [command] [flags]", func(w io.Writer) {
		_ = chainOfCustody
		_ = encrypt
		_ = verbose
		_ = quiet

		io.WriteString(w, "COMMANDS:\n")
		io.WriteString(w, "  capture [flags]                  Capture forensic evidence\n")
		io.WriteString(w, "    --case string                  Case identifier (required)\n")
		io.WriteString(w, "    --type strings                 Evidence types (memory, disk, network, logs, all)\n")
		io.WriteString(w, "    --output string                Output directory\n")
		io.WriteString(w, "    --compress                     Compress captured data\n")
		io.WriteString(w, "    --integrity                    Generate integrity checksums\n\n")

		io.WriteString(w, "  triage [flags]                   Perform system triage\n")
		io.WriteString(w, "    --quick                        Quick triage (basic system info)\n")
		io.WriteString(w, "    --full                         Full triage (comprehensive)\n")
		io.WriteString(w, "    --processes                    Analyze running processes\n")
		io.WriteString(w, "    --network                      Analyze network connections\n")
		io.WriteString(w, "    --persistence                  Check for persistence mechanisms\n\n")
		io.WriteString(w, "    --output string                Output file path\n\n")

		io.WriteString(w, "  analyze [flags]                  Analyze captured data\n")
		io.WriteString(w, "    --input string                 Input directory or file\n")
		io.WriteString(w, "    --ioc string                   IOC definition file\n")
		io.WriteString(w, "    --timeline                     Create timeline of events\n")
		io.WriteString(w, "    --correlate                    Correlate multiple evidence sources\n")
		io.WriteString(w, "    --report string                Generate analysis report\n\n")

		io.WriteString(w, "  hunt [flags]                     Hunt for threats and IOCs\n")
		io.WriteString(w, "    --yara string                  YARA rules file/directory\n")
		io.WriteString(w, "    --sigma string                 Sigma rules for detection\n")
		io.WriteString(w, "    --memory                       Scan process memory\n")
		io.WriteString(w, "    --filesystem                   Scan filesystem\n")
		io.WriteString(w, "    --registry                     Scan Windows registry (Wine/Cross-platform)\n\n")

		io.WriteString(w, "  report [flags]                   Generate incident report\n")
		io.WriteString(w, "    --template string              Report template\n")
		io.WriteString(w, "    --format string                Output format (pdf, html, docx, markdown)\n")
		io.WriteString(w, "    --executive                    Include executive summary\n")
		io.WriteString(w, "    --technical                    Include technical details\n")
		io.WriteString(w, "    --evidence                     Include evidence references\n\n")

		io.WriteString(w, "  timeline [flags]                 Create forensic timeline\n")
		io.WriteString(w, "    --source string                Data source directory\n")
		io.WriteString(w, "    --from string                  Start time (YYYY-MM-DD HH:MM)\n")
		io.WriteString(w, "    --to string                    End time (YYYY-MM-DD HH:MM)\n")
		io.WriteString(w, "    --visualize                    Generate visual timeline\n")
		io.WriteString(w, "    --export string                Export format (csv, json, html)\n\n")

		io.WriteString(w, "  contain [flags]                  Execute containment procedures\n")
		io.WriteString(w, "    --isolate                      Network isolation\n")
		io.WriteString(w, "    --quarantine string            Quarantine suspicious files\n")
		io.WriteString(w, "    --accounts strings             Suspend compromised accounts\n")
		io.WriteString(w, "    --services strings             Stop suspicious services\n")
		io.WriteString(w, "    --revert                       Revert to last known good state\n\n")

		io.WriteString(w, "  eradicate [flags]                Remove threats from system\n")
		io.WriteString(w, "    --malware                      Remove detected malware\n")
		io.WriteString(w, "    --persistence                  Remove persistence mechanisms\n")
		io.WriteString(w, "    --artifacts                    Clean attack artifacts\n")
		io.WriteString(w, "    --validate                     Verify removal success\n\n")

		io.WriteString(w, "FLAGS:\n")
		io.WriteString(w, "  --chain-of-custody           Maintain chain of custody\n")
		io.WriteString(w, "  --encrypt                    Encrypt sensitive data\n")
		io.WriteString(w, "  --verbose                    Show detailed output\n")
		io.WriteString(w, "  --quiet                      Minimal output\n\n")

		io.WriteString(w, "EXAMPLES:\n")
		io.WriteString(w, "  fortis incident capture --case incident-001 --type all\n")
		io.WriteString(w, "  fortis incident triage --full --output /evidence/triage.txt\n")
		io.WriteString(w, "  fortis incident hunt --yara ./rules/malware.yara\n")
		io.WriteString(w, "  fortis incident report --format pdf --executive --technical\n")
	})

	return cmd
}

func newIncidentNetworkForensicsCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "net",
		Short:  "Network forensics helper (safe)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = args
			return a.RunScript(cmd.Context(), "network-forensics.sh")
		},
	}
	cmd.Flags().StringVarP(new(string), "output", "o", "", "")
	// keep command hidden; placeholder flag kept hidden intentionally
	_ = cmd.Flags().MarkHidden("output")
	return cmd
}

func newIncidentCaptureCmd(a *app.App) *cobra.Command {
	var (
		caseID    string
		types     []string
		outputDir string
		compress  bool
		integrity bool
	)
	cmd := &cobra.Command{
		Use:   "capture",
		Short: "Collect forensic evidence",
		RunE: func(cmd *cobra.Command, args []string) error {
			if caseID == "" {
				return errors.New("--case is required")
			}
			m, err := incident.Capture(cmd.Context(), incident.CaptureOptions{
				CaseID:      caseID,
				Types:       types,
				OutputDir:   outputDir,
				Compress:    compress,
				Integrity:   integrity,
				Encrypt:     getBoolFlag(cmd, "encrypt"),
				Chain:       getBoolFlag(cmd, "chain-of-custody"),
				Verbose:     a.Verbose,
				CollectedBy: os.Getenv("USER"),
			})
			if err != nil {
				return err
			}
			// Write a manifest summary to stdout for automation.
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(m)
		},
	}
	cmd.Flags().StringVar(&caseID, "case", "", "Case identifier (required)")
	cmd.Flags().StringSliceVar(&types, "type", nil, "Evidence types (memory, disk, network, logs, all)")
	cmd.Flags().StringVar(&outputDir, "output", "", "Output directory")
	cmd.Flags().BoolVar(&compress, "compress", false, "Compress captured data")
	cmd.Flags().BoolVar(&integrity, "integrity", false, "Generate integrity checksums")
	return cmd
}

func newIncidentTriageCmd(a *app.App) *cobra.Command {
	var (
		quick       bool
		full        bool
		processes   bool
		network     bool
		persistence bool
		output      string
	)
	cmd := &cobra.Command{
		Use:   "triage",
		Short: "Perform rapid system triage",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := incident.Triage(cmd.Context(), incident.TriageOptions{
				Quick:       quick,
				Full:        full,
				Processes:   processes,
				Network:     network,
				Persistence: persistence,
				Output:      output,
			})
			if err != nil {
				return err
			}
			_, _ = io.WriteString(cmd.OutOrStdout(), out+"\n")
			return nil
		},
	}
	cmd.Flags().BoolVar(&quick, "quick", false, "Quick triage (basic system info)")
	cmd.Flags().BoolVar(&full, "full", false, "Full triage (comprehensive)")
	cmd.Flags().BoolVar(&processes, "processes", false, "Analyze running processes")
	cmd.Flags().BoolVar(&network, "network", false, "Analyze network connections")
	cmd.Flags().BoolVar(&persistence, "persistence", false, "Check for persistence mechanisms")
	cmd.Flags().StringVar(&output, "output", "", "Output file path")
	return cmd
}

func newIncidentAnalyzeCmd(a *app.App) *cobra.Command {
	var (
		input     string
		ioc       string
		timeline  bool
		correlate bool
		report    string
	)
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze captured data",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := incident.Analyze(cmd.Context(), incident.AnalyzeOptions{
				Input:     input,
				IOCFile:   ioc,
				Timeline:  timeline,
				Correlate: correlate,
				Report:    report,
			})
			if err != nil {
				return err
			}
			_, _ = io.WriteString(cmd.OutOrStdout(), out+"\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&input, "input", "", "Input directory or file")
	cmd.Flags().StringVar(&ioc, "ioc", "", "IOC definition file")
	cmd.Flags().BoolVar(&timeline, "timeline", false, "Create timeline of events")
	cmd.Flags().BoolVar(&correlate, "correlate", false, "Correlate multiple evidence sources")
	cmd.Flags().StringVar(&report, "report", "", "Generate analysis report")
	return cmd
}

func newIncidentHuntCmd(a *app.App) *cobra.Command {
	var (
		yara       string
		sigma      string
		memory     bool
		filesystem bool
		registry   bool
	)
	cmd := &cobra.Command{
		Use:   "hunt",
		Short: "Hunt for threats/IOCs",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{"--yara", yara, "--sigma", sigma}
			if memory {
				argv = append(argv, "--memory")
			}
			if filesystem {
				argv = append(argv, "--filesystem")
			}
			if registry {
				argv = append(argv, "--registry")
			}
			return a.RunScript(cmd.Context(), "malware-hunt.sh", argv...)
		},
	}
	cmd.Flags().StringVar(&yara, "yara", "", "YARA rules file/directory")
	cmd.Flags().StringVar(&sigma, "sigma", "", "Sigma rules for detection")
	cmd.Flags().BoolVar(&memory, "memory", false, "Scan process memory")
	cmd.Flags().BoolVar(&filesystem, "filesystem", false, "Scan filesystem")
	cmd.Flags().BoolVar(&registry, "registry", false, "Scan Windows registry")
	return cmd
}

func newIncidentReportCmd(a *app.App) *cobra.Command {
	var (
		template  string
		format    string
		execSum   bool
		technical bool
		evidence  bool
		output    string
	)
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate incident reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := incident.GenerateReport(incident.ReportOptions{
				Template:  template,
				Format:    format,
				Executive: execSum,
				Technical: technical,
				Evidence:  evidence,
				Output:    output,
			})
			if err != nil {
				return err
			}
			_, _ = io.WriteString(cmd.OutOrStdout(), out+"\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&template, "template", "", "Report template")
	cmd.Flags().StringVar(&format, "format", "markdown", "Output format (pdf, html, docx, markdown)")
	cmd.Flags().BoolVar(&execSum, "executive", false, "Include executive summary")
	cmd.Flags().BoolVar(&technical, "technical", false, "Include technical details")
	cmd.Flags().BoolVar(&evidence, "evidence", false, "Include evidence references")
	cmd.Flags().StringVar(&output, "output", "", "Output file path")
	return cmd
}

func newIncidentTimelineCmd(a *app.App) *cobra.Command {
	var (
		source    string
		from      string
		to        string
		visualize bool
		exportFmt string
	)
	cmd := &cobra.Command{
		Use:   "timeline",
		Short: "Create forensic timeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := incident.BuildTimeline(cmd.Context(), incident.TimelineOptions{
				Source:    source,
				From:      from,
				To:        to,
				Visualize: visualize,
				Export:    exportFmt,
			})
			if err != nil {
				return err
			}
			_, _ = io.WriteString(cmd.OutOrStdout(), out+"\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&source, "source", "", "Data source directory")
	cmd.Flags().StringVar(&from, "from", "", "Start time (YYYY-MM-DD HH:MM)")
	cmd.Flags().StringVar(&to, "to", "", "End time (YYYY-MM-DD HH:MM)")
	cmd.Flags().BoolVar(&visualize, "visualize", false, "Generate visual timeline")
	cmd.Flags().StringVar(&exportFmt, "export", "", "Export format (csv, json, html)")
	return cmd
}

func newIncidentContainCmd(a *app.App) *cobra.Command {
	var (
		isolate    bool
		quarantine string
		accounts   []string
		services   []string
		revert     bool
	)
	cmd := &cobra.Command{
		Use:   "contain",
		Short: "Execute containment procedures",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{}
			if isolate {
				argv = append(argv, "--isolate")
			}
			if quarantine != "" {
				argv = append(argv, "--quarantine", quarantine)
			}
			for _, a1 := range accounts {
				argv = append(argv, "--accounts", a1)
			}
			for _, s := range services {
				argv = append(argv, "--services", s)
			}
			if revert {
				argv = append(argv, "--revert")
			}
			return a.RunScript(cmd.Context(), "containment-procedures.sh", argv...)
		},
	}
	cmd.Flags().BoolVar(&isolate, "isolate", false, "Network isolation")
	cmd.Flags().StringVar(&quarantine, "quarantine", "", "Quarantine suspicious files")
	cmd.Flags().StringSliceVar(&accounts, "accounts", nil, "Suspend compromised accounts")
	cmd.Flags().StringSliceVar(&services, "services", nil, "Stop suspicious services")
	cmd.Flags().BoolVar(&revert, "revert", false, "Revert to last known good state")
	return cmd
}

func newIncidentEradicateCmd(a *app.App) *cobra.Command {
	var (
		malware     bool
		persistence bool
		artifacts   bool
		validate    bool
	)
	cmd := &cobra.Command{
		Use:   "eradicate",
		Short: "Remove threats from system",
		RunE: func(cmd *cobra.Command, args []string) error {
			argv := []string{}
			if malware {
				argv = append(argv, "--malware")
			}
			if persistence {
				argv = append(argv, "--persistence")
			}
			if artifacts {
				argv = append(argv, "--artifacts")
			}
			if validate {
				argv = append(argv, "--validate")
			}
			return a.RunScript(cmd.Context(), "eradication-tools.sh", argv...)
		},
	}
	cmd.Flags().BoolVar(&malware, "malware", false, "Remove detected malware")
	cmd.Flags().BoolVar(&persistence, "persistence", false, "Remove persistence mechanisms")
	cmd.Flags().BoolVar(&artifacts, "artifacts", false, "Clean attack artifacts")
	cmd.Flags().BoolVar(&validate, "validate", false, "Verify removal success")
	return cmd
}
