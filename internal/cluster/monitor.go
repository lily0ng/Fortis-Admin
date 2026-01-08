package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type MonitorOptions struct {
	InventoryPath string
	Group         string
	Hosts         []string
	HostsFile     string

	SSHUser    string
	SSHPort    int
	SSHKey     string
	SSHTimeout time.Duration
	Parallel   int

	Metrics []string
	Output  string // text|json
}

type HostMetrics struct {
	Host       string `json:"host"`
	OK         bool   `json:"ok"`
	LoadAvg    string `json:"load_avg,omitempty"`
	Uptime     string `json:"uptime,omitempty"`
	MemSummary string `json:"mem_summary,omitempty"`
	DiskRoot   string `json:"disk_root,omitempty"`
	Error      string `json:"error,omitempty"`
	Health     int    `json:"health"`
}

type MonitorReport struct {
	Timestamp time.Time     `json:"timestamp"`
	Hosts     []HostMetrics `json:"hosts"`
	Health    int           `json:"health"`
}

func Monitor(ctx context.Context, opts MonitorOptions) (MonitorReport, error) {
	metrics := map[string]bool{}
	if len(opts.Metrics) == 0 {
		metrics["cpu"] = true
		metrics["memory"] = true
		metrics["disk"] = true
		metrics["uptime"] = true
	} else {
		for _, m := range opts.Metrics {
			metrics[strings.ToLower(strings.TrimSpace(m))] = true
		}
	}

	script := []string{"bash", "-lc"}
	cmdParts := []string{}
	if metrics["cpu"] {
		cmdParts = append(cmdParts, "echo LOADAVG=$(cat /proc/loadavg 2>/dev/null | awk '{print $1\" \"$2\" \"$3}' || true)")
	}
	if metrics["uptime"] {
		cmdParts = append(cmdParts, "echo UPTIME=$(uptime -p 2>/dev/null || true)")
	}
	if metrics["memory"] {
		cmdParts = append(cmdParts, "echo MEM=$(free -m 2>/dev/null | awk '/Mem:/ {print $3\"/\"$2\"MB\"}' || true)")
	}
	if metrics["disk"] {
		cmdParts = append(cmdParts, "echo DISK=$(df -h / 2>/dev/null | awk 'NR==2 {print $5\" used\"}' || true)")
	}
	if len(cmdParts) == 0 {
		cmdParts = append(cmdParts, "echo OK=1")
	}
	sshCommand := strings.Join(cmdParts, "; ")
	_ = script

	res, err := Exec(ctx, ExecOptions{
		Command:       strings.Join(append(script, sshCommand), " "),
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
		return MonitorReport{}, err
	}

	rep := MonitorReport{Timestamp: time.Now()}
	total := 0
	count := 0
	for _, r := range res {
		hm := HostMetrics{Host: r.Host, OK: r.OK, Health: 0}
		if !r.OK {
			hm.Error = r.Error
			hm.Health = 0
			rep.Hosts = append(rep.Hosts, hm)
			count++
			continue
		}
		// parse key=value lines
		for _, ln := range strings.Split(r.Output, "\n") {
			ln = strings.TrimSpace(ln)
			if strings.HasPrefix(ln, "LOADAVG=") {
				hm.LoadAvg = strings.TrimPrefix(ln, "LOADAVG=")
			}
			if strings.HasPrefix(ln, "UPTIME=") {
				hm.Uptime = strings.TrimPrefix(ln, "UPTIME=")
			}
			if strings.HasPrefix(ln, "MEM=") {
				hm.MemSummary = strings.TrimPrefix(ln, "MEM=")
			}
			if strings.HasPrefix(ln, "DISK=") {
				hm.DiskRoot = strings.TrimPrefix(ln, "DISK=")
			}
		}
		// crude health heuristic
		hm.Health = 80
		if hm.DiskRoot != "" && strings.Contains(hm.DiskRoot, "%") {
			// leave as-is; detailed parsing can come later
		}
		if hm.LoadAvg == "" {
			hm.Health -= 10
		}
		if hm.MemSummary == "" {
			hm.Health -= 10
		}
		if hm.Uptime == "" {
			hm.Health -= 5
		}
		total += hm.Health
		count++
		rep.Hosts = append(rep.Hosts, hm)
	}
	if count > 0 {
		rep.Health = total / count
	}

	return rep, nil
}

func EncodeMonitorReportJSON(rep MonitorReport) ([]byte, error) {
	b, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func MonitorReportToText(rep MonitorReport) string {
	b := fmt.Sprintf("Cluster Health: %d/100\n", rep.Health)
	for _, h := range rep.Hosts {
		if !h.OK {
			b += fmt.Sprintf("%s\tERROR\t%s\n", h.Host, h.Error)
			continue
		}
		b += fmt.Sprintf("%s\tOK\tHealth=%d\tLoad=%s\tMem=%s\tDisk=%s\tUptime=%s\n", h.Host, h.Health, h.LoadAvg, h.MemSummary, h.DiskRoot, h.Uptime)
	}
	return b
}
