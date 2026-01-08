package incident

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TimelineEvent struct {
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"mtime"`
}

func BuildTimeline(ctx context.Context, opts TimelineOptions) (string, error) {
	_ = ctx
	if opts.Source == "" {
		return "", errors.New("--source is required")
	}

	from, _ := parseTime(opts.From)
	to, _ := parseTime(opts.To)

	events := []TimelineEvent{}
	limit := 5000
	_ = filepath.Walk(opts.Source, func(path string, info os.FileInfo, err error) error {
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
		mt := info.ModTime()
		if !from.IsZero() && mt.Before(from) {
			return nil
		}
		if !to.IsZero() && mt.After(to) {
			return nil
		}
		events = append(events, TimelineEvent{Path: path, Size: info.Size(), Mode: info.Mode().String(), ModTime: mt})
		if len(events) >= limit {
			return errors.New("timeline limit reached")
		}
		return nil
	})

	ext := strings.TrimSpace(opts.Export)
	if ext == "" {
		ext = "csv"
	}

	out := filepath.Join(".", fmt.Sprintf("timeline-%s.%s", time.Now().Format("20060102-150405"), ext))
	f, err := os.Create(out)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if opts.Visualize || ext == "html" {
		return out, renderTimelineHTML(f, events)
	}
	if ext == "json" {
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return out, enc.Encode(events)
	}
	// default CSV
	w := csv.NewWriter(f)
	_ = w.Write([]string{"path", "size", "mode", "mtime"})
	for _, e := range events {
		_ = w.Write([]string{e.Path, fmt.Sprintf("%d", e.Size), e.Mode, e.ModTime.Format(time.RFC3339)})
	}
	w.Flush()
	return out, w.Error()
}

func parseTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, nil
	}
	// supports "YYYY-MM-DD HH:MM"
	return time.ParseInLocation("2006-01-02 15:04", s, time.Local)
}

func renderTimelineHTML(w io.Writer, events []TimelineEvent) error {
	const tpl = `<!doctype html><html><head><meta charset="utf-8"/>
<title>Fortis Timeline</title>
<style>body{font-family:system-ui;margin:20px}table{border-collapse:collapse;width:100%}th,td{border-bottom:1px solid #eee;padding:8px;text-align:left}code{background:#f2f2f2;padding:2px 6px;border-radius:4px}</style>
</head><body>
<h1>Forensic Timeline</h1>
<table><thead><tr><th>Time</th><th>Path</th><th>Mode</th><th>Size</th></tr></thead><tbody>
{{ range . }}<tr><td>{{ .ModTime }}</td><td><code>{{ .Path }}</code></td><td>{{ .Mode }}</td><td>{{ .Size }}</td></tr>{{ end }}
</tbody></table>
</body></html>`
	t, err := template.New("t").Parse(tpl)
	if err != nil {
		return err
	}
	return t.Execute(w, events)
}
