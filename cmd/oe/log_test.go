package main

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/dbungert/omnienv/internal/omnienv"
	"github.com/stretchr/testify/assert"
)

var logTests = []struct {
	summary string
	verbose bool
}{{
	summary: "verbose false",
	verbose: false,
}, {
	summary: "verbose true",
	verbose: true,
}}

func patchLogger() (func(), *bytes.Buffer) {
	buf := &bytes.Buffer{}
	original := stderr
	stderr = buf
	return func() { stderr = original }, buf
}

func TestSetupLogging(t *testing.T) {
	restore, buf := patchLogger()
	defer restore()

	for _, test := range logTests {
		buf.Reset()
		omnienv.Verbose = test.verbose
		log.SetOutput(stderr)

		log.Printf("INFO: info")
		if omnienv.Verbose {
			log.Printf("DEBUG: debug")
		}

		lines := strings.Split(buf.String(), "\n")
		assert.Contains(t, lines[0], "INFO: info")
		if test.verbose {
			assert.Contains(t, lines[1], "DEBUG: debug")
			assert.Len(t, lines, 3)
		} else {
			assert.Len(t, lines, 2)
		}
	}
}
