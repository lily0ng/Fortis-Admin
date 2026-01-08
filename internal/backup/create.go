package backup

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Create(opts CreateOptions) (BackupMeta, error) {
	if opts.TargetDir == "" {
		return BackupMeta{}, errors.New("target dir is required")
	}
	if len(opts.Sources) == 0 {
		return BackupMeta{}, errors.New("at least one source is required")
	}
	if opts.Type == "" {
		opts.Type = BackupFull
	}
	if opts.Compress == "" {
		opts.Compress = CompressionGzip
	}

	id := fmt.Sprintf("backup-%s", time.Now().Format("20060102-150405"))
	if err := ensureDir(opts.TargetDir); err != nil {
		return BackupMeta{}, err
	}

	ext := "tar"
	var notes []string
	switch opts.Compress {
	case CompressionGzip:
		ext = "tar.gz"
	case CompressionNone:
		ext = "tar"
		notes = append(notes, "compression=none")
	default:
		notes = append(notes, "compression requested but not implemented, using gzip")
		ext = "tar.gz"
		opts.Compress = CompressionGzip
	}

	archivePath := filepath.Join(opts.TargetDir, id+"."+ext)
	f, err := os.Create(archivePath)
	if err != nil {
		return BackupMeta{}, err
	}
	defer f.Close()

	var tw *tar.Writer
	var gw *gzip.Writer
	var out io.Writer = f
	if opts.Compress == CompressionGzip {
		gw = gzip.NewWriter(f)
		defer gw.Close()
		out = gw
	}
	w := tar.NewWriter(out)
	defer w.Close()
	tw = w
	_ = tw

	for _, src := range opts.Sources {
		src = filepath.Clean(src)
		_ = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			rel := strings.TrimPrefix(path, string(filepath.Separator))
			if rel == "" {
				rel = path
			}
			if matchAnyGlob(path, opts.Exclude) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if info.IsDir() {
				return nil
			}
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

	sum, size, err := sha256File(archivePath)
	if err != nil {
		return BackupMeta{}, err
	}

	meta := BackupMeta{
		ID:             id,
		CreatedAt:      time.Now(),
		Type:           opts.Type,
		Sources:        opts.Sources,
		ArchivePath:    archivePath,
		SizeBytes:      size,
		ChecksumSHA256: sum,
		Encrypted:      opts.Encrypt,
		Compression:    string(opts.Compress),
		Notes:          notes,
	}

	if opts.Encrypt {
		meta.Notes = append(meta.Notes, "encrypt requested: not implemented (hook point)")
	}

	// Write a sidecar metadata JSON for listing.
	metaPath := filepath.Join(opts.TargetDir, id+".meta.json")
	metaJSON := fmt.Sprintf("{\n  \"id\": %q,\n  \"created_at\": %q,\n  \"type\": %q,\n  \"archive_path\": %q,\n  \"size_bytes\": %d,\n  \"sha256\": %q,\n  \"encrypted\": %v,\n  \"compression\": %q,\n  \"sources\": %q,\n  \"notes\": %q\n}\n", meta.ID, meta.CreatedAt.Format(time.RFC3339), meta.Type, meta.ArchivePath, meta.SizeBytes, meta.ChecksumSHA256, meta.Encrypted, meta.Compression, strings.Join(meta.Sources, ","), strings.Join(meta.Notes, ";"))
	_ = os.WriteFile(metaPath, []byte(metaJSON), 0o600)

	return meta, nil
}
