package omnienv

import "log"

var Verbose bool

func debugLog(format string, args ...any) {
	if Verbose {
		log.Printf("DEBUG: "+format, args...)
	}
}
