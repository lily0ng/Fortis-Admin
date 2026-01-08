package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fortis-admin/internal/config"
	"fortis-admin/internal/logging"
	"fortis-admin/internal/ui"
)

type App struct {
	ConfigPath string

	Debug   bool
	Quiet   bool
	Verbose bool

	ColorMode ui.ColorMode

	Config config.Config
	Log    *logging.Logger

	scriptsDir string
}

func New() *App {
	return &App{
		ColorMode: ui.ColorAuto,
		Config:    config.Default(),
		Log:       logging.New(os.Stdout, os.Stderr, logging.LevelInfo, false, false),
	}
}

func (a *App) Init() error {
	if a.ConfigPath == "" {
		return errors.New("config path is required")
	}

	cfg, err := config.Load(a.ConfigPath)
	if err == nil {
		a.Config = cfg
	}

	level := logging.LevelInfo
	if a.Debug {
		level = logging.LevelDebug
	}
	a.Log = logging.New(os.Stdout, os.Stderr, level, a.Quiet, a.Verbose)

	if a.Config.ScriptsDir != "" {
		a.scriptsDir = a.Config.ScriptsDir
	} else {
		a.scriptsDir = a.detectScriptsDir()
	}
	return nil
}

func (a *App) detectScriptsDir() string {
	exe, err := os.Executable()
	if err == nil {
		base := filepath.Dir(exe)
		candidates := []string{
			filepath.Join(base, "scripts"),
			filepath.Join(base, "..", "scripts"),
		}
		for _, c := range candidates {
			if fi, err := os.Stat(c); err == nil && fi.IsDir() {
				return c
			}
		}
	}
	return filepath.Join(".", "scripts")
}

func (a *App) ScriptsDir() string { return a.scriptsDir }

func (a *App) RunScript(ctx context.Context, scriptName string, args ...string) error {
	p := scriptName
	if !filepath.IsAbs(p) {
		p = filepath.Join(a.scriptsDir, scriptName)
	}
	cmdArgs := append([]string{p}, args...)
	cmd := exec.CommandContext(ctx, "bash", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	a.Log.Debugf("running script: %s", strings.Join(append([]string{"bash"}, cmdArgs...), " "))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("script failed (%s): %w", scriptName, err)
	}
	return nil
}
