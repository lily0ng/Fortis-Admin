package backup

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type SnapshotBackend string

const (
	SnapshotBackendNone  SnapshotBackend = "none"
	SnapshotBackendLVM   SnapshotBackend = "lvm"
	SnapshotBackendZFS   SnapshotBackend = "zfs"
	SnapshotBackendBtrfs SnapshotBackend = "btrfs"
)

type SnapshotOptions struct {
	Backend SnapshotBackend
	Volume  string
	Name    string
	Keep    int
	Remote  string
	DryRun  bool
	Apply   bool
	Yes     bool
}

type SnapshotResult struct {
	Timestamp time.Time       `json:"timestamp"`
	Backend   SnapshotBackend `json:"backend"`
	Volume    string          `json:"volume"`
	Name      string          `json:"name"`
	DryRun    bool            `json:"dry_run"`
	Planned   []string        `json:"planned"`
	Notes     []string        `json:"notes"`
}

func DetectSnapshotBackend(ctx context.Context) SnapshotBackend {
	_ = ctx
	if exec.Command("bash", "-lc", "command -v zfs >/dev/null 2>&1").Run() == nil {
		return SnapshotBackendZFS
	}
	if exec.Command("bash", "-lc", "command -v btrfs >/dev/null 2>&1").Run() == nil {
		return SnapshotBackendBtrfs
	}
	if exec.Command("bash", "-lc", "command -v lvcreate >/dev/null 2>&1").Run() == nil {
		return SnapshotBackendLVM
	}
	return SnapshotBackendNone
}

func ManageSnapshots(ctx context.Context, opts SnapshotOptions) (SnapshotResult, error) {
	_ = ctx
	if opts.Backend == "" {
		opts.Backend = DetectSnapshotBackend(ctx)
	}
	if strings.TrimSpace(opts.Volume) == "" {
		return SnapshotResult{}, errors.New("--volume is required")
	}
	if opts.Keep == 0 {
		opts.Keep = 7
	}

	res := SnapshotResult{Timestamp: time.Now(), Backend: opts.Backend, Volume: opts.Volume, Name: opts.Name, DryRun: opts.DryRun || !opts.Apply}

	if opts.Backend == SnapshotBackendNone {
		res.Notes = append(res.Notes, "no snapshot backend detected")
		return res, nil
	}

	// Plan only. Apply is intentionally a stub (safe-by-default).
	res.Planned = append(res.Planned, "create snapshot")
	res.Planned = append(res.Planned, "rotate snapshots (keep="+itoa(opts.Keep)+")")
	res.Planned = append(res.Planned, "space check")
	if opts.Remote != "" {
		res.Planned = append(res.Planned, "remote sync to "+opts.Remote+" (stub)")
	}

	if opts.Apply {
		if !opts.Yes {
			return res, errors.New("refusing to apply snapshots without --yes")
		}
		res.Notes = append(res.Notes, "apply not implemented; use native tooling (zfs/btrfs/lvm) or scripts")
	}

	return res, nil
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	b := make([]byte, 0, 32)
	for i > 0 {
		b = append(b, byte('0'+(i%10)))
		i /= 10
	}
	for l, r := 0, len(b)-1; l < r; l, r = l+1, r-1 {
		b[l], b[r] = b[r], b[l]
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}
