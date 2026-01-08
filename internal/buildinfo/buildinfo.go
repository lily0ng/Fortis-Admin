package buildinfo

import (
	"fmt"
	"runtime"
)

var (
	Version   = "1.0.0"
	Commit    = "a1b2c3d4e5f6"
	BuildDate = "2024-01-15"
)

func GoVersion() string { return runtime.Version() }

func Platform() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}
