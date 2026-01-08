package hardening

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ComplianceGap struct {
	Control string `json:"control"`
	Title   string `json:"title"`
	IssueID string `json:"issue_id"`
	Details string `json:"details,omitempty"`
}

type ComplianceReport struct {
	Timestamp time.Time       `json:"timestamp"`
	Standard  string          `json:"standard"`
	Host      string          `json:"host"`
	Score     int             `json:"score"`
	Gaps      []ComplianceGap `json:"gaps"`
	Evidence  []string        `json:"evidence"`
	Notes     []string        `json:"notes"`
}

type ComplianceOptions struct {
	Standard        string
	CollectEvidence bool
	GapAnalysis     bool
	Format          OutputFormat
}

func GenerateComplianceReport(ctx context.Context, opts ComplianceOptions) (ComplianceReport, error) {
	if strings.TrimSpace(opts.Standard) == "" {
		opts.Standard = "pci-dss"
	}
	if opts.Format == "" {
		opts.Format = FormatJSON
	}

	auditProfile := "cis"
	switch strings.ToLower(strings.TrimSpace(opts.Standard)) {
	case "pci", "pci-dss":
		auditProfile = "pci"
	case "hipaa":
		auditProfile = "hipaa"
	case "gdpr":
		auditProfile = "cis"
	case "iso27001", "iso-27001":
		auditProfile = "cis"
	}

	rep, err := RunAudit(ctx, AuditOptions{Profile: auditProfile, Level: "basic", Fix: false, Yes: false, Verbose: false})
	if err != nil {
		return ComplianceReport{}, err
	}

	out := ComplianceReport{
		Timestamp: time.Now(),
		Standard:  opts.Standard,
		Host:      rep.Hostname,
		Score:     rep.Score,
	}

	// Very lightweight mapping: failed findings are gaps.
	if opts.GapAnalysis {
		for _, f := range rep.Findings {
			if f.Result != ResultFail {
				continue
			}
			out.Gaps = append(out.Gaps, ComplianceGap{Control: controlForStandard(opts.Standard, f.ID), Title: f.Title, IssueID: f.ID, Details: f.Details})
		}
	}

	if opts.CollectEvidence {
		out.Evidence = append(out.Evidence, fmt.Sprintf("audit_profile=%s", rep.Profile))
		out.Evidence = append(out.Evidence, fmt.Sprintf("report_hash=%s", rep.ReportHash))
		out.Notes = append(out.Notes, "evidence collection is lightweight; attach artifacts manually if required")
	}

	return out, nil
}

func controlForStandard(standard, findingID string) string {
	s := strings.ToLower(strings.TrimSpace(standard))
	if s == "pci" || s == "pci-dss" {
		return "PCI-DSS" + "::" + findingID
	}
	if s == "hipaa" {
		return "HIPAA" + "::" + findingID
	}
	if s == "gdpr" {
		return "GDPR" + "::" + findingID
	}
	if s == "iso27001" || s == "iso-27001" {
		return "ISO27001" + "::" + findingID
	}
	return strings.ToUpper(s) + "::" + findingID
}

func RenderCompliance(w *os.File, rep ComplianceReport, format OutputFormat) error {
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(rep)
	case FormatHTML:
		// Simple HTML wrapper; keep in Go for portability.
		html := "<!doctype html><html><head><meta charset=\"utf-8\"/><title>Compliance Report</title></head><body>" +
			fmt.Sprintf("<h1>Compliance Report (%s)</h1>", rep.Standard) +
			fmt.Sprintf("<p><b>Timestamp:</b> %s</p>", rep.Timestamp.Format(time.RFC3339)) +
			fmt.Sprintf("<p><b>Host:</b> %s</p>", rep.Host) +
			fmt.Sprintf("<p><b>Score:</b> %d/100</p>", rep.Score) +
			"<h2>Gaps</h2><ul>"
		for _, g := range rep.Gaps {
			html += fmt.Sprintf("<li><code>%s</code> %s</li>", g.Control, g.Title)
		}
		html += "</ul></body></html>"
		_, err := w.Write([]byte(html))
		return err
	default:
		return errors.New("unsupported compliance export format")
	}
}

func ResolveComplianceOutputPath(exportFmt string, ts time.Time) (string, OutputFormat, error) {
	f := DetectFormat(exportFmt)
	if exportFmt != "" {
		if strings.Contains(exportFmt, "/") || strings.Contains(exportFmt, "\\") || strings.Contains(exportFmt, ".") {
			return exportFmt, f, nil
		}
	}
	name := fmt.Sprintf("compliance-%s.%s", ts.Format("20060102-150405"), string(f))
	preferred := filepath.Join("/var/log/fortis", name)
	if canWriteDir("/var/log/fortis") {
		return preferred, f, nil
	}
	return filepath.Join(".", name), f, nil
}

func canWriteDir(dir string) bool {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false
	}
	f, err := os.CreateTemp(dir, ".perm")
	if err != nil {
		return false
	}
	name := f.Name()
	_ = f.Close()
	_ = os.Remove(name)
	return true
}
