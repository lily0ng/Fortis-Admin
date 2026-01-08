package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogFile       string `yaml:"log_file"`
	ScriptsDir    string `yaml:"scripts_dir"`
	InventoryFile string `yaml:"inventory_file"`
}

func Default() Config {
	return Config{
		LogFile:       "/var/log/fortis/fortis.log",
		ScriptsDir:    "",
		InventoryFile: "/etc/fortis/inventory.yaml",
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, errors.New("config path is empty")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
