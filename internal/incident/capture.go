package incident

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func Capture(ctx context.Context, opts CaptureOptions) (CaptureManifest, error) {
	if opts.CaseID == "" {
		return CaptureManifest{}, errors.New("case id is required")
	}
	caseID := sanitizeCaseID(opts.CaseID)

	outDir := opts.OutputDir
	if outDir == "" {
		outDir = filepath.Join(".", "evidence", caseID)
	} else {
		outDir = filepath.Join(outDir, caseID)
	}

	if err := ensureDir(outDir); err != nil {
		return CaptureManifest{}, err
	}

	host, _ := os.Hostname()
	m := CaptureManifest{
		CaseID:      caseID,
		Timestamp:   time.Now(),
		OutputDir:   outDir,
		CollectedBy: opts.CollectedBy,
		Notes: []string{
			"OS=" + runtime.GOOS,
			"HOST=" + host,
		},
	}

	want := normalizeCaptureTypes(opts.Types)

	// Always collect basic system artifacts.
	if want.all || want.system {
		_ = writeCmdOutput(filepath.Join(outDir, "uname.txt"), "uname", "-a")
		_ = writeCmdOutput(filepath.Join(outDir, "date.txt"), "date")
	}

	if runtime.GOOS == "linux" {
		if want.all || want.system {
			_ = writeCmdOutput(filepath.Join(outDir, "os-release.txt"), "bash", "-lc", "cat /etc/os-release 2>/dev/null || true")
			_ = writeCmdOutput(filepath.Join(outDir, "users.txt"), "bash", "-lc", "who 2>/dev/null || true")
			_ = writeCmdOutput(filepath.Join(outDir, "kernel-modules.txt"), "bash", "-lc", "lsmod 2>/dev/null | head -n 500 || true")
			_ = writeCmdOutput(filepath.Join(outDir, "systemctl-failed.txt"), "bash", "-lc", "systemctl --failed 2>/dev/null || true")
		}

		if want.all || want.processes {
			_ = writeCmdOutput(filepath.Join(outDir, "processes.txt"), "bash", "-lc", "ps auxww 2>/dev/null | head -n 500 || true")
		}

		if want.all || want.network {
			_ = writeCmdOutput(filepath.Join(outDir, "network.txt"), "bash", "-lc", "ss -tulpn 2>/dev/null | head -n 500 || netstat -tulpn 2>/dev/null | head -n 500 || true")
			_ = writeCmdOutput(filepath.Join(outDir, "listening.txt"), "bash", "-lc", "lsof -i -P -n 2>/dev/null | head -n 500 || true")
		}

		if want.all || want.logs {
			_ = writeCmdOutput(filepath.Join(outDir, "authlog-tail.txt"), "bash", "-lc", "(tail -n 400 /var/log/auth.log 2>/dev/null || tail -n 400 /var/log/secure 2>/dev/null || true)")
			_ = writeCmdOutput(filepath.Join(outDir, "syslog-tail.txt"), "bash", "-lc", "tail -n 400 /var/log/syslog 2>/dev/null || tail -n 400 /var/log/messages 2>/dev/null || true")
			_ = writeCmdOutput(filepath.Join(outDir, "journal-tail.txt"), "bash", "-lc", "journalctl -n 400 --no-pager 2>/dev/null || true")
		}
	}

	if want.disk {
		m.Notes = append(m.Notes, "disk evidence requested: imaging not implemented (use external tooling)")
	}
	if want.memory {
		m.Notes = append(m.Notes, "memory evidence requested: acquisition not implemented (LiME integration not included)")
	}

	if opts.Integrity {
		files, _ := os.ReadDir(outDir)
		for _, e := range files {
			if e.IsDir() {
				continue
			}
			p := filepath.Join(outDir, e.Name())
			sum, size, err := sha256File(p)
			if err != nil {
				continue
			}
			m.Files = append(m.Files, FileHash{Path: e.Name(), SHA256: sum, Size: size})
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		bestEffortWrite(filepath.Join(outDir, "manifest.json"), b)
	}

	if opts.Compress {
		tarPath := filepath.Join(filepath.Dir(outDir), fmt.Sprintf("%s-%s.tar.gz", caseID, m.Timestamp.Format("20060102-150405")))
		if err := tarGzDir(ctx, outDir, tarPath); err != nil {
			m.Notes = append(m.Notes, "compression_failed="+err.Error())
		} else {
			m.Notes = append(m.Notes, "compressed="+tarPath)
		}
	}

	if opts.Encrypt {
		m.Notes = append(m.Notes, "encrypt requested: not implemented (use disk encryption / external tooling)")
	}
	if opts.Chain {
		m.Notes = append(m.Notes, "chain-of-custody requested: manifest.json written when --integrity is enabled")
	}

	return m, nil
}

type captureTypeSet struct {
	all       bool
	memory    bool
	disk      bool
	network   bool
	logs      bool
	system    bool
	processes bool
}

func normalizeCaptureTypes(types []string) captureTypeSet {
	if len(types) == 0 {
		return captureTypeSet{all: true}
	}

	set := captureTypeSet{}
	for _, t := range types {
		t = strings.ToLower(strings.TrimSpace(t))
		switch t {
		case "", "all":
			set.all = true
		case "memory":
			set.memory = true
		case "disk":
			set.disk = true
		case "network":
			set.network = true
		case "logs":
			set.logs = true
		case "system":
			set.system = true
		case "processes":
			set.processes = true
		}
	}

	if !(set.all || set.memory || set.disk || set.network || set.logs || set.system || set.processes) {
		set.all = true
	}
	return set
}

func tarGzDir(ctx context.Context, srcDir, outPath string) error {
	_ = ctx
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()
	w := tar.NewWriter(gz)
	defer w.Close()

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(filepath.Dir(srcDir), path)
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return nil
		}
		hdr.Name = rel
		if err := w.WriteHeader(hdr); err != nil {
			return nil
		}
		r, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer r.Close()
		_, _ = io.Copy(w, r)
		return nil
	})
}
