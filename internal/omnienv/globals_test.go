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

	Verbose = true
	defer func() { Verbose = false }()

	debugLog("hello %s", "world")
	assert.Contains(t, buf.String(), "DEBUG: hello world")
}

func TestDebugLogNotVerbose(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	Verbose = false
	debugLog("should not appear")

	assert.Empty(t, buf)
}
