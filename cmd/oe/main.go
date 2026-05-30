package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/dbungert/omnienv/internal/omnienv"
	"github.com/jessevdk/go-flags"
)

func Run() error {
	opts, err := GetOpts(os.Args[1:])
	if err != nil {
		return err
	}

	if opts.Version {
		ver := Version
		if ver == "dev" {
			if bi, ok := debug.ReadBuildInfo(); ok {
				ver = bi.Main.Version
			}
		}
		fmt.Printf("omnienv version: %v\n", ver)
		return nil
	}

	log.SetOutput(os.Stderr)
	omnienv.Verbose = opts.Verbose
	if omnienv.Verbose {
		// G706 regards log injection, but we log to stderr
		log.Printf("DEBUG: cmdline opts=%+v", opts) // #nosec G706
	}

	cfg, err := omnienv.GetConfig()
	if err != nil {
		return fmt.Errorf("fatal error: %w", err)
	}

	app := omnienv.NewApp(cfg, opts)

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
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) {
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		}
		log.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
