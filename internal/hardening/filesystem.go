package hardening

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type FilesystemOptions struct {
	Root string
}

type FilesystemReport struct {
	Root             string   `json:"root" yaml:"root"`
	SUIDFiles        []string `json:"suid_files" yaml:"suid_files"`
	SGIDFiles        []string `json:"sgid_files" yaml:"sgid_files"`
	WorldWritable    []string `json:"world_writable" yaml:"world_writable"`
	MountFSTabIssues []string `json:"fstab_issues" yaml:"fstab_issues"`
}

func ScanFilesystem(ctx context.Context, opts FilesystemOptions) (FilesystemReport, error) {
	_ = ctx
	if opts.Root == "" {
		opts.Root = "/"
	}

	rep := FilesystemReport{Root: opts.Root}
	if runtime.GOOS != "linux" {
		return rep, nil
	}

	limit := 200
	_ = filepath.WalkDir(opts.Root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		// Skip pseudo filesystems
		if d.IsDir() {
			base := filepath.Base(path)
			if path == "/proc" || path == "/sys" || path == "/dev" || base == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		mode := info.Mode()
		if mode&os.ModeSetuid != 0 {
			if len(rep.SUIDFiles) < limit {
				rep.SUIDFiles = append(rep.SUIDFiles, path)
			}
		}
		if mode&os.ModeSetgid != 0 {
			if len(rep.SGIDFiles) < limit {
				rep.SGIDFiles = append(rep.SGIDFiles, path)
			}
		}
		if mode.Perm()&0o002 != 0 {
			if len(rep.WorldWritable) < limit {
				rep.WorldWritable = append(rep.WorldWritable, path)
			}
		}
		return nil
	})

	// Basic fstab checks
	if b, err := os.ReadFile("/etc/fstab"); err == nil {
		lines := strings.Split(string(b), "\n")
		for _, ln := range lines {
			ln = strings.TrimSpace(ln)
			if ln == "" || strings.HasPrefix(ln, "#") {
				continue
			}
			fields := strings.Fields(ln)
			if len(fields) < 4 {
				continue
			}
			mount := fields[1]
			opts := fields[3]
			if mount == "/tmp" && !strings.Contains(opts, "noexec") {
				rep.MountFSTabIssues = append(rep.MountFSTabIssues, "/tmp missing noexec")
			}
			if mount == "/var/tmp" && !strings.Contains(opts, "noexec") {
				rep.MountFSTabIssues = append(rep.MountFSTabIssues, "/var/tmp missing noexec")
			}
		}
	}

	return rep, nil
}
