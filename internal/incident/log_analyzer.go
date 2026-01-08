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

type LogAnalyzeOptions struct {
	InputPath string
	IOCFile   string
	Output    string
}

type LogMatch struct {
	File string `json:"file"`
	Line string `json:"line"`
}

type LogAnalyzeReport struct {
	Timestamp time.Time  `json:"timestamp"`
	InputPath string     `json:"input"`
	IOCFile   string     `json:"ioc_file"`
	IOCs      []string   `json:"iocs"`
	Matches   []LogMatch `json:"matches"`
	Total     int        `json:"total"`
}

func AnalyzeLogs(ctx context.Context, opts LogAnalyzeOptions) (string, error) {
	_ = ctx
	if opts.InputPath == "" {
		return "", errors.New("--input is required")
	}

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
	}

	rep := LogAnalyzeReport{
		Timestamp: time.Now(),
		InputPath: opts.InputPath,
		IOCFile:   opts.IOCFile,
		IOCs:      needles,
	}

	files := []string{}
	fi, err := os.Stat(opts.InputPath)
	if err != nil {
		return "", err
	}
	if fi.IsDir() {
		_ = filepath.Walk(opts.InputPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				base := filepath.Base(path)
				if base == ".git" {
					return filepath.SkipDir
				}
				return nil
			}
			if info.Size() > 10*1024*1024 {
				return nil
			}
			files = append(files, path)
			return nil
		})
	} else {
		files = append(files, opts.InputPath)
	}

	if len(needles) == 0 {
		// if no IOC file, just return inventory.
		rep.Total = 0
		out := opts.Output
		if out == "" {
			out = filepath.Join(".", fmt.Sprintf("log-analysis-%s.json", time.Now().Format("20060102-150405")))
		}
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return "", err
		}
		f, err := os.Create(out)
		if err != nil {
			return "", err
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return out, enc.Encode(rep)
	}

	for _, p := range files {
		hits, err := grepFile(p, needles, 50)
		if err != nil {
			continue
		}
		for _, h := range hits {
			rep.Matches = append(rep.Matches, LogMatch{File: p, Line: h})
		}
	}
	rep.Total = len(rep.Matches)

	out := opts.Output
	if out == "" {
		out = filepath.Join(".", fmt.Sprintf("log-analysis-%s.json", time.Now().Format("20060102-150405")))
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return "", err
	}
	f, err := os.Create(out)
	if err != nil {
		return "", err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return out, enc.Encode(rep)
}
