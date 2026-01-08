package incident

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type IOC struct {
	Value   string    `json:"value"`
	Type    string    `json:"type"`
	Source  string    `json:"source"`
	AddedAt time.Time `json:"added_at"`
}

type IOCStore struct {
	Path string
}

func DefaultIOCStorePath() string {
	h, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".fortis", "iocs.json")
	}
	return filepath.Join(h, ".fortis", "iocs.json")
}

func (s IOCStore) Load() ([]IOC, error) {
	b, err := os.ReadFile(s.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []IOC{}, nil
		}
		return nil, err
	}
	var iocs []IOC
	if err := json.Unmarshal(b, &iocs); err != nil {
		return nil, err
	}
	return iocs, nil
}

func (s IOCStore) Save(iocs []IOC) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return err
	}
	sort.Slice(iocs, func(i, j int) bool { return iocs[i].Value < iocs[j].Value })
	b, err := json.MarshalIndent(iocs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.Path, b, 0o600)
}

func AddIOC(iocs []IOC, value, typ, source string) []IOC {
	value = strings.TrimSpace(value)
	if value == "" {
		return iocs
	}
	for _, it := range iocs {
		if it.Value == value {
			return iocs
		}
	}
	iocs = append(iocs, IOC{Value: value, Type: typ, Source: source, AddedAt: time.Now()})
	return iocs
}

func RemoveIOC(iocs []IOC, value string) []IOC {
	value = strings.TrimSpace(value)
	out := make([]IOC, 0, len(iocs))
	for _, it := range iocs {
		if it.Value == value {
			continue
		}
		out = append(out, it)
	}
	return out
}

func ImportIOCsFromTextFile(path string, typ string, source string) ([]IOC, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	iocs := []IOC{}
	for _, ln := range strings.Split(string(b), "\n") {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}
		iocs = AddIOC(iocs, ln, typ, source)
	}
	return iocs, nil
}
