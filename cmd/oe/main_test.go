package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
)

func withArgs(t *testing.T, args []string) {
	saved := os.Args
	os.Args = args
	t.Cleanup(func() { os.Args = saved })
}

func TestRunVersionDev(t *testing.T) {
	withArgs(t, []string{"oe", "--version"})
	saved := Version
	Version = "dev"
	t.Cleanup(func() { Version = saved })

	r, w, err := os.Pipe()
	assert.Nil(t, err)
	savedStdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = savedStdout })

	assert.Nil(t, Run())

	assert.Nil(t, w.Close())
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	assert.Contains(t, buf.String(), "omnienv version:")
}

func TestRunVersionSet(t *testing.T) {
	withArgs(t, []string{"oe", "--version"})
	saved := Version
	Version = "v0.2"
	t.Cleanup(func() { Version = saved })

	r, w, err := os.Pipe()
	assert.Nil(t, err)
	savedStdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = savedStdout })

	assert.Nil(t, Run())

	assert.Nil(t, w.Close())
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	assert.Equal(t, "omnienv version: v0.2\n", buf.String())
}

func TestRunHelp(t *testing.T) {
	withArgs(t, []string{"oe", "--help"})
	err := Run()
	var flagsErr *flags.Error
	assert.ErrorAs(t, err, &flagsErr)
	assert.Equal(t, flags.ErrHelp, flagsErr.Type)
}

func TestRunBadFlag(t *testing.T) {
	withArgs(t, []string{"oe", "--invalid"})
	err := Run()
	assert.Error(t, err)
}
