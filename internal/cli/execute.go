package cli

import "os"

func Execute() {
	cmd := NewRootCmd(os.Stdout, os.Stderr)
	_ = cmd.Execute()
}
