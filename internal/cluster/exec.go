package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type ExecOptions struct {
	Command       string
	Group         string
	Hosts         []string
	InventoryPath string
	HostsFile     string

	SSHUser    string
	SSHPort    int
	SSHKey     string
	SSHTimeout time.Duration

	Parallel int
	Output   string // combined|json
}

type ExecResult struct {
	Host   string `json:"host"`
	OK     bool   `json:"ok"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func Exec(ctx context.Context, opts ExecOptions) ([]ExecResult, error) {
	if strings.TrimSpace(opts.Command) == "" {
		return nil, errors.New("command is required")
	}
	if opts.SSHPort == 0 {
		opts.SSHPort = 22
	}
	if opts.SSHTimeout == 0 {
		opts.SSHTimeout = 30 * time.Second
	}
	if opts.Parallel <= 0 {
		opts.Parallel = 4
	}

	targets := []string{}
	if len(opts.Hosts) > 0 {
		targets = append(targets, opts.Hosts...)
	}
	if opts.HostsFile != "" {
		hs, err := HostsFromFile(opts.HostsFile)
		if err != nil {
			return nil, err
		}
		targets = append(targets, hs...)
	}

	var inv Inventory
	if opts.InventoryPath != "" {
		loaded, err := LoadInventory(opts.InventoryPath)
		if err == nil {
			inv = loaded
		}
	}
	if opts.Group != "" {
		for _, s := range FilterByGroup(inv, opts.Group) {
			if s.Hostname != "" {
				targets = append(targets, s.Hostname)
			} else if s.IP != "" {
				targets = append(targets, s.IP)
			}
		}
	}

	// de-dup
	seen := map[string]struct{}{}
	uniq := []string{}
	for _, h := range targets {
		h = strings.TrimSpace(h)
		if h == "" {
			continue
		}
		if _, ok := seen[h]; ok {
			continue
		}
		seen[h] = struct{}{}
		uniq = append(uniq, h)
	}
	if len(uniq) == 0 {
		return nil, errors.New("no target hosts")
	}

	type job struct{ host string }
	jobs := make(chan job)
	results := make(chan ExecResult)

	worker := func() {
		for j := range jobs {
			res := runSSH(ctx, inv, j.host, opts)
			results <- res
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < opts.Parallel; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); worker() }()
	}

	go func() {
		for _, h := range uniq {
			jobs <- job{host: h}
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	out := []ExecResult{}
	for r := range results {
		out = append(out, r)
	}
	return out, nil
}

func runSSH(ctx context.Context, inv Inventory, host string, opts ExecOptions) ExecResult {
	user := opts.SSHUser
	port := opts.SSHPort

	if s := FindByHostnameOrIP(inv, host); s != nil {
		if user == "" && s.SSHUser != "" {
			user = s.SSHUser
		}
		if s.SSHPort != 0 {
			port = s.SSHPort
		}
	}

	target := host
	if user != "" {
		target = fmt.Sprintf("%s@%s", user, host)
	}

	args := []string{"-p", fmt.Sprintf("%d", port), "-o", "BatchMode=yes", "-o", "StrictHostKeyChecking=accept-new"}
	if opts.SSHKey != "" {
		args = append(args, "-i", opts.SSHKey)
	}
	args = append(args, target, opts.Command)

	ctx2, cancel := context.WithTimeout(ctx, opts.SSHTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx2, "ssh", args...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return ExecResult{Host: host, OK: false, Output: string(b), Error: err.Error()}
	}
	return ExecResult{Host: host, OK: true, Output: string(b)}
}

func EncodeExecResultsJSON(res []ExecResult) ([]byte, error) {
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return nil, err
	}
	return b, nil
}
