package omnienv

import (
	"log"
	"os/exec"
	"time"
)

var command = exec.Command
var commandContext = exec.CommandContext
var timeSleep = time.Sleep
var Verbose bool

func debugLog(format string, args ...any) {
	if Verbose {
		log.Printf("DEBUG: "+format, args...)
	}
}
