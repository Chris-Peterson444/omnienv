package omnienv

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugLogVerbose(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	app := App{Opts: Opts{Verbose: true}}
	app.debugLog("hello %s", "world")
	assert.Contains(t, buf.String(), "DEBUG: hello world")
}

func TestDebugLogNotVerbose(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	app := App{Opts: Opts{Verbose: false}}
	app.debugLog("should not appear")

	assert.Empty(t, buf)
}
