package hardening

import (
	"bufio"
	"context"
	"os/exec"
	"runtime"
)

type PackageAuditOptions struct {
	List bool
}

type PackageAuditReport struct {
	Manager       string   `json:"manager" yaml:"manager"`
	TotalPackages int      `json:"total_packages" yaml:"total_packages"`
	Packages      []string `json:"packages,omitempty" yaml:"packages,omitempty"`
	Note          string   `json:"note" yaml:"note"`
}

func AuditPackages(ctx context.Context, opts PackageAuditOptions) (PackageAuditReport, error) {
	rep := PackageAuditReport{Note: "CVE integration not implemented yet"}
	if runtime.GOOS != "linux" {
		rep.Manager = "unsupported"
		return rep, nil
	}

	if _, err := exec.LookPath("dpkg-query"); err == nil {
		rep.Manager = "dpkg"
		cmd := exec.CommandContext(ctx, "dpkg-query", "-W", "-f", "${binary:Package}\n")
		out, err := cmd.StdoutPipe()
		if err != nil {
			return rep, err
		}
		if err := cmd.Start(); err != nil {
			return rep, err
		}
		s := bufio.NewScanner(out)
		for s.Scan() {
			rep.TotalPackages++
			if opts.List {
				rep.Packages = append(rep.Packages, s.Text())
			}
			if rep.TotalPackages > 5000 && opts.List {
				// avoid huge output
				break
			}
		}
		_ = cmd.Wait()
		return rep, nil
	}

	if _, err := exec.LookPath("rpm"); err == nil {
		rep.Manager = "rpm"
		cmd := exec.CommandContext(ctx, "rpm", "-qa")
		out, err := cmd.StdoutPipe()
		if err != nil {
			return rep, err
		}
		if err := cmd.Start(); err != nil {
			return rep, err
		}
		s := bufio.NewScanner(out)
		for s.Scan() {
			rep.TotalPackages++
			if opts.List {
				rep.Packages = append(rep.Packages, s.Text())
			}
			if rep.TotalPackages > 5000 && opts.List {
				break
			}
		}
		_ = cmd.Wait()
		return rep, nil
	}

	rep.Manager = "unknown"
	return rep, nil
}
