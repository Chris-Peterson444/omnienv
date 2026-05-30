package omnienv

import (
	"errors"
	"io"
	"os"
)

var errMissingArgs = errors.New("missing command arguments")

func (app App) run(args ...string) error {
	if len(args) == 0 {
		return errMissingArgs
	}
	cmd := app.command(args[0], args[1:]...)
	app.debugLog("run command=%v", args)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (app App) runDevNull(args ...string) error {
	if len(args) == 0 {
		return errMissingArgs
	}
	cmd := app.command(args[0], args[1:]...)
	app.debugLog("run command=%v", args)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}
