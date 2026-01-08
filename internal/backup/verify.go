package backup

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func Verify(opts VerifyOptions) (VerifyResult, error) {
	if opts.BackupPath == "" {
		return VerifyResult{}, errors.New("--backup is required")
	}
	fi, err := os.Stat(opts.BackupPath)
	if err != nil {
		return VerifyResult{}, err
	}
	if fi.Size() == 0 {
		return VerifyResult{BackupPath: opts.BackupPath, OK: false, Reason: "backup file is empty"}, nil
	}
	sum, _, err := sha256File(opts.BackupPath)
	if err != nil {
		return VerifyResult{}, err
	}

	res := VerifyResult{BackupPath: opts.BackupPath, OK: true, SHA256: sum}

	// If there is a sidecar meta file, validate checksum.
	metaPath := strings.TrimSuffix(opts.BackupPath, filepath.Ext(opts.BackupPath)) + ".meta.json"
	if strings.HasSuffix(opts.BackupPath, ".tar.gz") {
		metaPath = strings.TrimSuffix(opts.BackupPath, ".tar.gz") + ".meta.json"
	}
	if b, err := os.ReadFile(metaPath); err == nil {
		var raw map[string]any
		if err := json.Unmarshal(b, &raw); err == nil {
			want, _ := raw["sha256"].(string)
			if want != "" && want != sum {
				res.OK = false
				res.Reason = "checksum mismatch vs meta"
				return res, nil
			}
		}
	}

	if opts.Full {
		// Restore simulation to detect archive corruption.
		tmp, err := os.MkdirTemp("", "fortis-restore-verify-*")
		if err != nil {
			return VerifyResult{}, err
		}
		defer os.RemoveAll(tmp)
		if err := Restore(RestoreOptions{BackupPath: opts.BackupPath, TargetDir: tmp, DryRun: false}); err != nil {
			res.OK = false
			res.Reason = "restore simulation failed: " + err.Error()
			return res, nil
		}
	}

	if opts.Quick {
		res.Reason = "checksum computed"
	} else {
		res.Reason = "checksum validated"
	}
	return res, nil
}
