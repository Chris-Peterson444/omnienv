package main

import (
	"log"

	"github.com/dbungert/omnienv/internal/omnienv"
)

func setupLogging(verbose bool) {
	log.SetOutput(stderr)
	omnienv.Verbose = verbose
}
