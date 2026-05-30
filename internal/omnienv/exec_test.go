package omnienv

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
