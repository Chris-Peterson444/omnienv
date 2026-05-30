package omnienv

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunEmptyArgs(t *testing.T) {
	app := App{}
	err := app.run()
	assert.ErrorIs(t, err, errMissingArgs)
}

func TestRunDevNullEmptyArgs(t *testing.T) {
	app := App{}
	err := app.runDevNull()
	assert.ErrorIs(t, err, errMissingArgs)
}

func TestRun(t *testing.T) {
	app := App{
		command: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("/bin/true")
		},
	}
	assert.Nil(t, app.run("/bin/true"))
}

func TestRunDevNull(t *testing.T) {
	app := App{
		command: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("/bin/true")
		},
	}
	assert.Nil(t, app.runDevNull("/bin/true"))
}

func TestLxcExec(t *testing.T) {
	app := App{
		command: func(arg0 string, argv ...string) *exec.Cmd {
			assert.Equal(t, "lxc", arg0)
			assert.Equal(t, []string{"exec", "-", "--", "bar"}, argv)
			cmd := exec.Command("/bin/true")
			cmd.Args = append([]string{arg0}, argv...)
			return cmd
		},
	}
	assert.Nil(t, app.lxcExec("bar"))
}
