package cluster

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Hostname string   `yaml:"hostname" json:"hostname"`
	IP       string   `yaml:"ip" json:"ip"`
	OS       string   `yaml:"os" json:"os"`
	Status   string   `yaml:"status" json:"status"`
	Groups   []string `yaml:"groups" json:"groups"`
	Tags     []string `yaml:"tags" json:"tags"`
	SSHUser  string   `yaml:"ssh_user" json:"ssh_user"`
	SSHPort  int      `yaml:"ssh_port" json:"ssh_port"`
}

type Inventory struct {
	Servers []Server `yaml:"servers" json:"servers"`
}

func LoadInventory(path string) (Inventory, error) {
	var inv Inventory
	b, err := os.ReadFile(path)
	if err != nil {
		return inv, err
	}
	if err := yaml.Unmarshal(b, &inv); err != nil {
		return inv, err
	}
	return inv, nil
}

func HostsFromFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		ln := strings.TrimSpace(s.Text())
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}
		out = append(out, ln)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, errors.New("no hosts found in file")
	}
	return out, nil
}

func FilterByGroup(inv Inventory, group string) []Server {
	group = strings.TrimSpace(group)
	if group == "" {
		return nil
	}
	out := []Server{}
	for _, s := range inv.Servers {
		for _, g := range s.Groups {
			if g == group {
				out = append(out, s)
				break
			}
		}
	}
	return out
}

func FindByHostnameOrIP(inv Inventory, host string) *Server {
	for i := range inv.Servers {
		s := &inv.Servers[i]
		if s.Hostname == host || s.IP == host {
			return s
		}
	}
	return nil
}
