package backup

import (
	"fmt"
	"time"
)

type TestDROptions struct {
	BackupPath string
	TargetDir  string
	DryRun     bool
}

type TestDRResult struct {
	Timestamp time.Time `json:"timestamp"`
	Backup    string    `json:"backup"`
	Target    string    `json:"target"`
	DryRun    bool      `json:"dry_run"`
	OK        bool      `json:"ok"`
	Note      string    `json:"note"`
}

func TestDR(opts TestDROptions) (TestDRResult, error) {
	res := TestDRResult{Timestamp: time.Now(), Backup: opts.BackupPath, Target: opts.TargetDir, DryRun: opts.DryRun}
	if err := Restore(RestoreOptions{BackupPath: opts.BackupPath, TargetDir: opts.TargetDir, DryRun: opts.DryRun}); err != nil {
		res.OK = false
		res.Note = err.Error()
		return res, err
	}
	res.OK = true
	if opts.DryRun {
		res.Note = "dry-run restore simulation completed"
	} else {
		res.Note = fmt.Sprintf("restore completed to %s", opts.TargetDir)
	}
	return res, nil
}
