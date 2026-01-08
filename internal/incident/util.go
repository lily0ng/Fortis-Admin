package incident

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func writeCmdOutput(path string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	b := out
	if err != nil {
		b = append(b, []byte("\nERROR: "+err.Error()+"\n")...)
	}
	return os.WriteFile(path, b, 0o600)
}

func sha256File(path string) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}

func ensureDir(p string) error { return os.MkdirAll(p, 0o755) }

func sanitizeCaseID(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "\\", "-")
	if s == "" {
		return "incident"
	}
	return s
}

func bestEffortWrite(path string, b []byte) {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, b, 0o600)
}

func grepFile(path string, needles []string, maxHits int) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(b, []byte("\n"))
	hits := []string{}
	for _, ln := range lines {
		l := string(ln)
		for _, n := range needles {
			if n == "" {
				continue
			}
			if strings.Contains(l, n) {
				hits = append(hits, l)
				if maxHits > 0 && len(hits) >= maxHits {
					return hits, nil
				}
				break
			}
		}
	}
	return hits, nil
}
