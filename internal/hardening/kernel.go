package hardening

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type KernelOptions struct {
	ApplyPreset string
	Param       string
	Value       string
	Persist     bool
	Yes         bool
	DryRun      bool
}

type KernelResult struct {
	Plan []string
}

func ApplyKernel(ctx context.Context, opts KernelOptions) (KernelResult, error) {
	if runtime.GOOS != "linux" {
		return KernelResult{Plan: []string{"not supported on this OS"}}, nil
	}

	pairs := map[string]string{}
	if opts.Param != "" {
		pairs[opts.Param] = opts.Value
	}
	for k, v := range presetSysctls(opts.ApplyPreset) {
		pairs[k] = v
	}
	if len(pairs) == 0 {
		return KernelResult{}, errors.New("no kernel parameters specified")
	}

	res := KernelResult{}
	for k, v := range pairs {
		res.Plan = append(res.Plan, fmt.Sprintf("sysctl -w %s=%s", k, v))
		if opts.Persist {
			res.Plan = append(res.Plan, fmt.Sprintf("persist %s=%s", k, v))
		}
	}

	if opts.DryRun {
		return res, nil
	}
	if !opts.Yes {
		return res, errors.New("refusing to apply kernel changes without --yes")
	}

	for k, v := range pairs {
		_ = exec.CommandContext(ctx, "sysctl", "-w", fmt.Sprintf("%s=%s", k, v)).Run()
		if opts.Persist {
			if err := persistSysctl(k, v); err != nil {
				return res, err
			}
		}
	}
	return res, nil
}

func presetSysctls(name string) map[string]string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "network":
		return map[string]string{
			"net.ipv4.conf.all.accept_redirects": "0",
			"net.ipv4.conf.all.send_redirects":   "0",
		}
	case "security":
		return map[string]string{
			"kernel.kptr_restrict":  "2",
			"kernel.dmesg_restrict": "1",
		}
	case "memory":
		return map[string]string{
			"vm.mmap_min_addr": "65536",
		}
	default:
		return map[string]string{}
	}
}
