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

type CatalogEntry struct {
	Path string
	Size int64
}

func Catalog(opts CatalogOptions) ([]CatalogEntry, error) {
	if opts.BackupPath == "" {
		return nil, errors.New("--backup is required")
	}
	f, err := os.Open(opts.BackupPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var r io.Reader = f
	if strings.HasSuffix(opts.BackupPath, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		defer gz.Close()
		r = gz
	}

	entries := []CatalogEntry{}
	tr := tar.NewReader(r)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		name := filepath.Clean(h.Name)
		if opts.Search != "" && !strings.Contains(name, opts.Search) {
			continue
		}
		entries = append(entries, CatalogEntry{Path: name, Size: h.Size})
		if len(entries) > 5000 {
			break
		}
	}
	return entries, nil
}
