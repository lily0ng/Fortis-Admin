package backup

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RestoreWizardOptions struct {
	BackupPath  string
	TargetDir   string
	Items       []string
	DryRun      bool
	Interactive bool
}

type RestoreWizardResult struct {
	Timestamp  time.Time `json:"timestamp"`
	BackupPath string    `json:"backup_path"`
	TargetDir  string    `json:"target_dir"`
	Items      []string  `json:"items"`
	DryRun     bool      `json:"dry_run"`
	OK         bool      `json:"ok"`
	Note       string    `json:"note"`
}

func RunRestoreWizard(ctx context.Context, opts RestoreWizardOptions) (RestoreWizardResult, error) {
	_ = ctx
	res := RestoreWizardResult{Timestamp: time.Now(), BackupPath: opts.BackupPath, TargetDir: opts.TargetDir, Items: opts.Items, DryRun: opts.DryRun}
	if strings.TrimSpace(opts.BackupPath) == "" {
		return res, errors.New("--backup is required")
	}

	if opts.Interactive {
		r := bufio.NewReader(os.Stdin)
		if strings.TrimSpace(opts.TargetDir) == "" {
			fmt.Fprint(os.Stdout, "Restore target directory: ")
			v, _ := r.ReadString('\n')
			opts.TargetDir = strings.TrimSpace(v)
			res.TargetDir = opts.TargetDir
		}
		if opts.TargetDir == "" {
			return res, errors.New("--target is required")
		}
		fmt.Fprintf(os.Stdout, "Proceed with restore to %s? (yes/no): ", opts.TargetDir)
		v, _ := r.ReadString('\n')
		v = strings.ToLower(strings.TrimSpace(v))
		if v != "yes" {
			res.OK = false
			res.Note = "cancelled"
			return res, errors.New("restore cancelled")
		}
	}

	if strings.TrimSpace(opts.TargetDir) == "" {
		// default safe location
		opts.TargetDir = filepath.Join(".", "restore")
		res.TargetDir = opts.TargetDir
	}

	if err := Restore(RestoreOptions{BackupPath: opts.BackupPath, TargetDir: opts.TargetDir, Items: opts.Items, DryRun: opts.DryRun}); err != nil {
		res.OK = false
		res.Note = err.Error()
		return res, err
	}
	res.OK = true
	if opts.DryRun {
		res.Note = "dry-run restore plan completed"
	} else {
		res.Note = "restore completed"
	}
	return res, nil
}
