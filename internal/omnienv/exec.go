package omnienv

import (
	"io"
	"os"
)

func (app App) run(args ...string) error {
	cmd := app.command(args[0], args[1:]...)
	app.debugLog("run command=%v", args)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (app App) runDevNull(args ...string) error {
	cmd := app.command(args[0], args[1:]...)
	app.debugLog("run command=%v", args)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}
