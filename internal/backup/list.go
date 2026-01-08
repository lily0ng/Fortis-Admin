package backup

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type ListedBackup struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Archive   string    `json:"archive"`
	SizeBytes int64     `json:"size_bytes"`
	SHA256    string    `json:"sha256"`
}

func List(opts ListOptions) ([]ListedBackup, error) {
	if opts.TargetDir == "" {
		return nil, errors.New("target dir is required")
	}
	entries, err := os.ReadDir(opts.TargetDir)
	if err != nil {
		return nil, err
	}
	out := []ListedBackup{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".meta.json") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(opts.TargetDir, name))
		if err != nil {
			continue
		}
		var raw map[string]any
		if err := json.Unmarshal(b, &raw); err != nil {
			continue
		}
		id, _ := raw["id"].(string)
		archive, _ := raw["archive_path"].(string)
		sha, _ := raw["sha256"].(string)
		createdStr, _ := raw["created_at"].(string)
		created, _ := time.Parse(time.RFC3339, createdStr)
		sizeF, _ := raw["size_bytes"].(float64)

		if opts.Filter != "" && !strings.Contains(id, opts.Filter) {
			continue
		}
		out = append(out, ListedBackup{ID: id, CreatedAt: created, Archive: archive, SizeBytes: int64(sizeF), SHA256: sha})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}
