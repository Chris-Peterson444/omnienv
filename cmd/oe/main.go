package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/dbungert/omnienv/internal/omnienv"
)

func Run() error {
	opts, err := GetOpts(os.Args[1:])
	if err != nil {
		return nil
	}

	if opts.Version {
		ver := "unknown"
		if bi, ok := debug.ReadBuildInfo(); ok {
			ver = bi.Main.Version
		}
		fmt.Printf("omnienv version: %v\n", ver)
		return nil
	}

	setupLogging(opts.Verbose)
	slog.Debug("cmdline", "opts", opts) // #nosec G706 — debug-only local CLI tool, negligible injection risk

	cfg, err := omnienv.GetConfig()
	if err != nil {
		return fmt.Errorf("fatal error: %w", err)
	}

	app := omnienv.App{Config: cfg, Opts: opts}

	if opts.Launch {
		if err := app.Launch(); err != nil {
			return fmt.Errorf("failed to launch: %w", err)
		}
	}

	if err := app.Shell(); err != nil {
		return fmt.Errorf("failed to create shell: %w", err)
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
}
