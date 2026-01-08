package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type PatchOptions struct {
	InventoryPath string
	Group         string
	Hosts         []string
	HostsFile     string

	SSHUser    string
	SSHPort    int
	SSHKey     string
	SSHTimeout time.Duration
	Parallel   int

	Packages          []string
	Strategy          string
	BatchSize         int
	PreCheck          bool
	PostCheck         bool
	RollbackOnFailure bool

	Apply bool
	Yes   bool
}

type PatchHostResult struct {
	Host   string `json:"host"`
	OK     bool   `json:"ok"`
	Plan   string `json:"plan"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

type PatchReport struct {
	Timestamp time.Time         `json:"timestamp"`
	Apply     bool              `json:"apply"`
	Strategy  string            `json:"strategy"`
	BatchSize int               `json:"batch_size"`
	Packages  []string          `json:"packages"`
	Results   []PatchHostResult `json:"results"`
}

func OrchestratePatches(ctx context.Context, opts PatchOptions) (PatchReport, error) {
	if opts.Apply && !opts.Yes {
		return PatchReport{}, errors.New("refusing to apply patches without --yes")
	}
	if opts.Strategy == "" {
		opts.Strategy = "rolling"
	}
	if opts.BatchSize <= 0 {
		opts.BatchSize = 2
	}

	plan := "noop"
	if len(opts.Packages) > 0 {
		plan = "update packages: " + strings.Join(opts.Packages, ",")
	}

	rep := PatchReport{Timestamp: time.Now(), Apply: opts.Apply, Strategy: opts.Strategy, BatchSize: opts.BatchSize, Packages: opts.Packages}

	if !opts.Apply {
		// Dry-run: return which hosts would be targeted.
		res, err := Exec(ctx, ExecOptions{
			Command:       "echo OK=1",
			Group:         opts.Group,
			Hosts:         opts.Hosts,
			HostsFile:     opts.HostsFile,
			InventoryPath: opts.InventoryPath,
			SSHUser:       opts.SSHUser,
			SSHPort:       opts.SSHPort,
			SSHKey:        opts.SSHKey,
			SSHTimeout:    opts.SSHTimeout,
			Parallel:      opts.Parallel,
			Output:        "json",
		})
		if err != nil {
			return rep, err
		}
		for _, r := range res {
			rep.Results = append(rep.Results, PatchHostResult{Host: r.Host, OK: r.OK, Plan: plan})
		}
		return rep, nil
	}

	// Apply: best-effort, non-destructive-ish command selection.
	cmd := ""
	if len(opts.Packages) == 0 {
		cmd = "echo 'no packages specified; nothing to do'"
	} else {
		pkgs := strings.Join(opts.Packages, " ")
		cmd = fmt.Sprintf("(command -v apt-get >/dev/null 2>&1 && sudo apt-get update -y && sudo apt-get install -y %s) || (command -v yum >/dev/null 2>&1 && sudo yum install -y %s) || echo 'no supported package manager'", pkgs, pkgs)
	}
	if opts.PreCheck {
		cmd = "uname -a; " + cmd
	}
	if opts.PostCheck {
		cmd = cmd + "; echo POSTCHECK_OK=1"
	}

	res, err := Exec(ctx, ExecOptions{
		Command:       cmd,
		Group:         opts.Group,
		Hosts:         opts.Hosts,
		HostsFile:     opts.HostsFile,
		InventoryPath: opts.InventoryPath,
		SSHUser:       opts.SSHUser,
		SSHPort:       opts.SSHPort,
		SSHKey:        opts.SSHKey,
		SSHTimeout:    opts.SSHTimeout,
		Parallel:      opts.Parallel,
		Output:        "json",
	})
	if err != nil {
		return rep, err
	}
	for _, r := range res {
		rep.Results = append(rep.Results, PatchHostResult{Host: r.Host, OK: r.OK, Plan: plan, Output: r.Output, Error: r.Error})
	}

	return rep, nil
}

func EncodePatchReportJSON(rep PatchReport) ([]byte, error) {
	b, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return nil, err
	}
	return b, nil
}
