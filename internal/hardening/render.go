package hardening

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"

	"gopkg.in/yaml.v3"
)

type OutputFormat string

const (
	FormatJSON OutputFormat = "json"
	FormatYAML OutputFormat = "yaml"
	FormatHTML OutputFormat = "html"
)

func DetectFormat(output string) OutputFormat {
	// if output is a file path, infer by extension
	switch {
	case hasSuffix(output, ".yaml") || hasSuffix(output, ".yml"):
		return FormatYAML
	case hasSuffix(output, ".html") || hasSuffix(output, ".htm"):
		return FormatHTML
	case hasSuffix(output, ".json"):
		return FormatJSON
	default:
		return FormatJSON
	}
}

func Render(w io.Writer, rep Report, format OutputFormat) error {
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(rep)
	case FormatYAML:
		b, err := yaml.Marshal(rep)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		return err
	case FormatHTML:
		return renderHTML(w, rep)
	default:
		return errors.New("unsupported format")
	}
}

func renderHTML(w io.Writer, rep Report) error {
	const tpl = `<!doctype html>
<html>
<head>
<meta charset="utf-8"/>
<title>Fortis Audit Report</title>
<style>
body{font-family:system-ui,-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;margin:24px}
code{background:#f2f2f2;padding:2px 6px;border-radius:4px}
.badge{display:inline-block;padding:2px 8px;border-radius:999px;font-size:12px}
.pass{background:#e8fff0;color:#116a2c}
.fail{background:#ffe8e8;color:#7d1b1b}
.warn{background:#fff8e0;color:#6b4b00}
.skip{background:#eef2ff;color:#2a3a7a}
</style>
</head>
<body>
<h1>Fortis Audit Report</h1>
<p><b>Timestamp:</b> {{ .Timestamp }}</p>
<p><b>Host:</b> {{ .Hostname }} <b>OS:</b> {{ .Platform }}</p>
<p><b>Profile:</b> {{ .Profile }} <b>Level:</b> {{ .Level }}</p>
<p><b>Score:</b> {{ .Score }}/100 ({{ .ScoreLabel }})</p>
<p><b>Stats:</b> Passed {{ .Passed }} | Failed {{ .Failed }} | Warnings {{ .Warnings }} | Skipped {{ .Skipped }}</p>
<hr/>
<table cellpadding="8" cellspacing="0" border="0">
<thead><tr><th align="left">ID</th><th align="left">Result</th><th align="left">Title</th><th align="left">Recommendation</th></tr></thead>
<tbody>
{{ range .Findings }}
<tr>
<td><code>{{ .ID }}</code></td>
<td>
<span class="badge {{ .Result }}">{{ .Result }}</span>
</td>
<td>{{ .Title }}</td>
<td>{{ .Recommendation }}</td>
</tr>
{{ end }}
</tbody>
</table>
</body>
</html>`

	t, err := template.New("rep").Parse(tpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, rep); err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func hasSuffix(s, suf string) bool {
	if len(s) < len(suf) {
		return false
	}
	return s[len(s)-len(suf):] == suf
}
