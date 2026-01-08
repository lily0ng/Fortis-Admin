package cluster

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type InitOptions struct {
	InventoryPath string
	SSHKeyPath    string
	Force         bool
}

type InitResult struct {
	InventoryPath string   `json:"inventory_path"`
	SSHKeyPath    string   `json:"ssh_key_path"`
	CreatedFiles  []string `json:"created_files"`
	Notes         []string `json:"notes"`
}

func Init(opts InitOptions) (InitResult, error) {
	if strings.TrimSpace(opts.InventoryPath) == "" {
		return InitResult{}, errors.New("inventory path is required")
	}
	res := InitResult{InventoryPath: opts.InventoryPath, SSHKeyPath: opts.SSHKeyPath}

	if opts.SSHKeyPath != "" {
		if _, err := os.Stat(opts.SSHKeyPath); err != nil {
			res.Notes = append(res.Notes, "ssh key not found at provided path")
		} else {
			res.Notes = append(res.Notes, "ssh key found")
		}
	} else {
		res.Notes = append(res.Notes, "ssh key path not provided; use --ssh-key to target a specific key")
	}

	if err := os.MkdirAll(filepath.Dir(opts.InventoryPath), 0o755); err != nil {
		return res, err
	}
	if !opts.Force {
		if _, err := os.Stat(opts.InventoryPath); err == nil {
			return res, fmt.Errorf("inventory already exists at %s (use --force to overwrite)", opts.InventoryPath)
		}
	}

	tpl := "servers:\n" +
		"  - hostname: web01\n" +
		"    ip: 192.168.1.101\n" +
		"    os: Ubuntu 22.04\n" +
		"    status: online\n" +
		"    groups: [webservers, production]\n" +
		"    tags: [linux]\n" +
		"    ssh_user: root\n" +
		"    ssh_port: 22\n" +
		"  - hostname: db01\n" +
		"    ip: 192.168.1.102\n" +
		"    os: CentOS 8\n" +
		"    status: online\n" +
		"    groups: [databases, production]\n" +
		"    tags: [linux]\n" +
		"    ssh_user: root\n" +
		"    ssh_port: 22\n"

	if err := os.WriteFile(opts.InventoryPath, []byte(tpl), 0o644); err != nil {
		return res, err
	}
	res.CreatedFiles = append(res.CreatedFiles, opts.InventoryPath)
	res.Notes = append(res.Notes,
		"Next steps:",
		"1) Validate inventory: fortis cluster inventory --inventory-file "+opts.InventoryPath+" --output json",
		"2) Test SSH connectivity: fortis cluster exec --inventory-file "+opts.InventoryPath+" --hosts web01 --command \"uname -a\"",
		"3) Distribute SSH keys safely: ssh-copy-id (manual) or your config management tool",
	)

	return res, nil
}
