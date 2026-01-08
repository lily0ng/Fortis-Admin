package incident

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

type AnalysisReport struct {
	Timestamp time.Time `json:"timestamp"`
	Input     string    `json:"input"`
	IOCs      []string  `json:"iocs"`
	Matches   []string  `json:"matches"`
	Notes     []string  `json:"notes"`
}

func Analyze(ctx context.Context, opts AnalyzeOptions) (string, error) {
	_ = ctx
	if opts.Input == "" {
		return "", errors.New("--input is required")
	}

	rep := AnalysisReport{Timestamp: time.Now(), Input: opts.Input}

	needles := []string{}
	if opts.IOCFile != "" {
		b, err := os.ReadFile(opts.IOCFile)
		if err != nil {
			return "", err
		}
		for _, ln := range strings.Split(string(b), "\n") {
			ln = strings.TrimSpace(ln)
			if ln == "" || strings.HasPrefix(ln, "#") {
				continue
			}
			needles = append(needles, ln)
		}
		rep.IOCs = needles
	}

	// Scan common files inside capture output for IOC matches.
	candidates := []string{"authlog-tail.txt", "syslog-tail.txt", "journal-tail.txt", "processes.txt", "network.txt"}
	for _, c := range candidates {
		p := filepath.Join(opts.Input, c)
		if _, err := os.Stat(p); err != nil {
			continue
		}
		if len(needles) == 0 {
			continue
		}
		hits, err := grepFile(p, needles, 100)
		if err != nil {
			continue
		}
		for _, h := range hits {
			rep.Matches = append(rep.Matches, fmt.Sprintf("%s: %s", c, h))
		}
	}

	if opts.Timeline {
		rep.Notes = append(rep.Notes, "timeline requested: use `incident timeline` command")
	}
	if opts.Correlate {
		rep.Notes = append(rep.Notes, "correlation requested: basic IOC matching applied")
	}

	out := opts.Report
	if out == "" {
		out = filepath.Join(".", fmt.Sprintf("analysis-%s.json", time.Now().Format("20060102-150405")))
	}
	_ = os.MkdirAll(filepath.Dir(out), 0o755)

	f, err := os.Create(out)
	if err != nil {
		return "", err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rep); err != nil {
		return "", err
	}
	return out, nil
}
