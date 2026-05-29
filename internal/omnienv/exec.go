package omnienv

import (
	"io"
	"os"
)

func run(args ...string) error {
	cmd := command(args[0], args[1:]...)
	debugLog("run command=%v", args)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDevNull(args ...string) error {
	cmd := command(args[0], args[1:]...)
	debugLog("run command=%v", args)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}
