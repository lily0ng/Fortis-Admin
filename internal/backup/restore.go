package backup

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Restore(opts RestoreOptions) error {
	if opts.BackupPath == "" {
		return errors.New("--backup is required")
	}
	if opts.TargetDir == "" {
		return errors.New("--target is required")
	}
	if err := ensureDir(opts.TargetDir); err != nil {
		return err
	}

	f, err := os.Open(opts.BackupPath)
	if err != nil {
		return err
	}
	defer f.Close()

	var r io.Reader = f
	if strings.HasSuffix(opts.BackupPath, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer gz.Close()
		r = gz
	}

	tr := tar.NewReader(r)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		name := filepath.Clean(h.Name)
		if name == "." || name == string(filepath.Separator) {
			continue
		}
		if len(opts.Items) > 0 {
			match := false
			for _, it := range opts.Items {
				it = strings.TrimPrefix(filepath.Clean(it), string(filepath.Separator))
				if it == "" {
					continue
				}
				if strings.HasPrefix(name, it) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		dest := filepath.Join(opts.TargetDir, name)
		if opts.DryRun {
			continue
		}
		if h.FileInfo().IsDir() {
			_ = ensureDir(dest)
			continue
		}
		if err := ensureDir(filepath.Dir(dest)); err != nil {
			return err
		}
		out, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, h.FileInfo().Mode())
		if err != nil {
			return err
		}
		_, _ = io.Copy(out, tr)
		_ = out.Close()
	}
	return nil
}
