package incident

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GenerateReport(opts ReportOptions) (string, error) {
	format := strings.TrimSpace(opts.Format)
	if format == "" {
		format = "markdown"
	}
	if format != "markdown" {
		return "", errors.New("only markdown report output is implemented")
	}

	out := opts.Output
	if out == "" {
		out = filepath.Join(".", fmt.Sprintf("incident-report-%s.md", time.Now().Format("20060102-150405")))
	}
	_ = os.MkdirAll(filepath.Dir(out), 0o755)

	b := &strings.Builder{}
	b.WriteString("# Incident Report\n\n")
	b.WriteString("Generated: " + time.Now().Format(time.RFC3339) + "\n\n")

	if opts.Executive {
		b.WriteString("## Executive Summary\n\n")
		b.WriteString("- Summary: (fill in)\n")
		b.WriteString("- Impact: (fill in)\n")
		b.WriteString("- Current status: (fill in)\n\n")
	}

	if opts.Technical {
		b.WriteString("## Technical Details\n\n")
		b.WriteString("- Timeline: (fill in)\n")
		b.WriteString("- Indicators observed: (fill in)\n")
		b.WriteString("- Affected systems: (fill in)\n\n")
	}

	if opts.Evidence {
		b.WriteString("## Evidence\n\n")
		b.WriteString("- Evidence directory: (fill in)\n")
		b.WriteString("- Hash manifest: (fill in)\n\n")
	}

	b.WriteString("## Actions Taken\n\n")
	b.WriteString("- Capture performed\n")
	b.WriteString("- Triage performed\n")
	b.WriteString("- Containment (if any)\n\n")

	if err := os.WriteFile(out, []byte(b.String()), 0o600); err != nil {
		return "", err
	}
	return out, nil
}
