package omnienv

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var nameTests = []struct {
	summary string
	config  Config
	opts    Opts

	system      string
	name        string
	launchImage string
}{{
	summary:     "basic name",
	config:      Config{Label: "l", System: NewSystem("s")},
	system:      "s",
	name:        "l-s",
	launchImage: "ubuntu-daily:s",
}, {
	summary:     "foo-bar",
	config:      Config{Label: "foo", System: NewSystem("bar")},
	system:      "bar",
	name:        "foo-bar",
	launchImage: "ubuntu-daily:bar",
}, {
	summary:     "opts override",
	config:      Config{Label: "l", System: NewSystem("sys-from-config")},
	opts:        Opts{System: "sys-from-opts"},
	system:      "sys-from-opts",
	name:        "l-sys-from-opts",
	launchImage: "ubuntu-daily:sys-from-opts",
}}

func TestName(t *testing.T) {
	for _, test := range nameTests {
		app := App{Config: test.config, Opts: test.opts}
		assert.Equal(t, test.system, app.system(), test.summary)
		assert.Equal(t, test.name, app.name(), test.summary)
		assert.Equal(t, test.launchImage, app.launchImage(), test.summary)
	}
}

var lxcInstanceStatusTests = []struct {
	summary string
	cmd     *exec.Cmd
	status  string
	errMsg  string
}{{
	summary: "ok",
	cmd:     exec.Command("/bin/echo", "Status: RUNNING"),
	status:  "RUNNING",
}, {
	summary: "info fails",
	cmd:     exec.Command("/bin/false"),
	errMsg:  "failed to get instance info",
}, {
	summary: "no status",
	cmd:     exec.Command("/bin/echo", "just some output"),
	errMsg:  "could not determine status",
}}

func TestLxcInstanceStatus(t *testing.T) {
	for _, test := range lxcInstanceStatusTests {
		restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
			return test.cmd
		})
		status, err := App{Config: Config{Label: "l", System: NewSystem("s")}}.lxcInstanceStatus()
		restoreCmd()
		if test.errMsg != "" {
			assert.ErrorContains(t, err, test.errMsg, test.summary)
		} else {
			assert.Nil(t, err, test.summary)
			assert.Equal(t, test.status, status, test.summary)
		}
	}
}

var startIfNeededTests = []struct {
	summary string
	cmd     *exec.Cmd
	errMsg  string
}{{
	summary: "running",
	cmd:     exec.Command("/bin/echo", "Status: RUNNING"),
}, {
	summary: "unknown",
	cmd:     exec.Command("/bin/echo", "Status: UNKNOWN"),
	errMsg:  "no handler for Status UNKNOWN",
}}

func TestStartIfNeeded(t *testing.T) {
	for _, test := range startIfNeededTests {
		restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
			return test.cmd
		})
		err := App{Config: Config{Label: "l", System: NewSystem("s")}}.StartIfNeeded()
		restoreCmd()
		if test.errMsg != "" {
			assert.ErrorContains(t, err, test.errMsg, test.summary)
		} else {
			assert.Nil(t, err, test.summary)
		}
	}
}

func TestStartIfNeededStopped(t *testing.T) {
	callCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		callCount++
		if callCount == 1 {
			return exec.Command("/bin/echo", "Status: STOPPED")
		}
		return exec.Command("/bin/true")
	})
	defer restoreCmd()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	assert.Nil(t, app.StartIfNeeded())
}

func TestStartFails(t *testing.T) {
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/false")
	})
	defer restoreCmd()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.start()
	assert.ErrorContains(t, err, "failed to start instance")
}

var isVMTests = []struct {
	summary string
	cmd     *exec.Cmd
	vm      bool
	errMsg  string
}{{
	summary: "virtual-machine",
	cmd:     exec.Command("/bin/echo", "Type: virtual-machine"),
	vm:      true,
}, {
	summary: "container",
	cmd:     exec.Command("/bin/echo", "Type: container"),
	vm:      false,
}, {
	summary: "info fails",
	cmd:     exec.Command("/bin/false"),
	errMsg:  "failed to get instance info",
}, {
	summary: "no type",
	cmd:     exec.Command("/bin/echo", "Status: RUNNING"),
	errMsg:  "could not determine type",
}}

func TestIsVM(t *testing.T) {
	for _, test := range isVMTests {
		restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
			return test.cmd
		})
		vm, err := App{Config: Config{Label: "l", System: NewSystem("s")}}.isVM()
		restoreCmd()
		if test.errMsg != "" {
			assert.ErrorContains(t, err, test.errMsg, test.summary)
		} else {
			assert.Nil(t, err, test.summary)
			assert.Equal(t, test.vm, vm, test.summary)
		}
	}
}

func TestWaitNotVM(t *testing.T) {
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/echo", "Type: container")
	})
	defer restoreCmd()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	assert.Nil(t, app.Wait())
}

func TestWaitVMExecOk(t *testing.T) {
	callCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		callCount++
		if callCount == 1 {
			return exec.Command("/bin/echo", "Type: virtual-machine")
		}
		return exec.Command("/bin/true")
	})
	defer restoreCmd()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	assert.Nil(t, app.Wait())
}

func TestWaitVMExecFailsNonExitError(t *testing.T) {
	callCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		callCount++
		if callCount == 1 {
			return exec.Command("/bin/echo", "Type: virtual-machine")
		}
		return exec.Command("/nonexistent-binary")
	})
	defer restoreCmd()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.Wait()
	assert.Error(t, err)
}

func TestWaitVMStrangeExitCode(t *testing.T) {
	callCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		callCount++
		if callCount == 1 {
			return exec.Command("/bin/echo", "Type: virtual-machine")
		}
		return exec.Command("/bin/false")
	})
	defer restoreCmd()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.Wait()
	assert.ErrorContains(t, err, "strange lxc exec exit code 1")
}

func TestWaitVMTimeout(t *testing.T) {
	callCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		callCount++
		if callCount == 1 {
			return exec.Command("/bin/echo", "Type: virtual-machine")
		}
		return exec.Command("/bin/sh", "-c", "exit 255")
	})
	defer restoreCmd()
	restoreSleep := Patch(&timeSleep, func(_ time.Duration) {})
	defer restoreSleep()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.Wait()
	assert.ErrorContains(t, err, "timed out waiting")
}

func TestWaitVMEventualSuccess(t *testing.T) {
	callCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		callCount++
		if callCount == 1 {
			return exec.Command("/bin/echo", "Type: virtual-machine")
		}
		if callCount < 5 {
			return exec.Command("/bin/sh", "-c", "exit 255")
		}
		return exec.Command("/bin/true")
	})
	defer restoreCmd()
	restoreSleep := Patch(&timeSleep, func(_ time.Duration) {})
	defer restoreSleep()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	assert.Nil(t, app.Wait())
}

var isUbuntuJammyTests = []struct {
	summary string
	cmd     *exec.Cmd
	want    bool
}{{
	summary: "true",
	cmd:     exec.Command("/bin/printf", "Distributor ID: Ubuntu\nRelease: 22.04"),
	want:    true,
}, {
	summary: "not Ubuntu",
	cmd:     exec.Command("/bin/echo", "Debian"),
	want:    false,
}, {
	summary: "command fails",
	cmd:     exec.Command("/bin/false"),
	want:    false,
}, {
	summary: "wrong version",
	cmd:     exec.Command("/bin/printf", "Distributor ID: Ubuntu\nRelease: 24.04"),
	want:    false,
}, {
	summary: "no release",
	cmd:     exec.Command("/bin/printf", "Distributor ID: Ubuntu"),
	want:    false,
}}

func TestIsUbuntuJammy(t *testing.T) {
	for _, test := range isUbuntuJammyTests {
		restoreCmdCtx := Patch(&commandContext, func(_ context.Context, _ string, _ ...string) *exec.Cmd {
			return test.cmd
		})
		jammy, err := App{Config: Config{Label: "l", System: NewSystem("s")}}.isUbuntuJammy()
		restoreCmdCtx()
		assert.Nil(t, err, test.summary)
		assert.Equal(t, test.want, jammy, test.summary)
	}
}

func TestShellContainerOk(t *testing.T) {
	cmdCallCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		cmdCallCount++
		switch cmdCallCount {
		case 1:
			return exec.Command("/bin/echo", "Status: RUNNING") // StartIfNeeded
		case 2:
			return exec.Command("/bin/echo", "Type: container") // Wait → isVM
		case 3:
			return exec.Command("/bin/true") // lxcExec
		default:
			return exec.Command("/bin/true")
		}
	})
	defer restoreCmd()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	assert.Nil(t, app.Shell())
}

func TestShellStartIfNeededFails(t *testing.T) {
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/false")
	})
	defer restoreCmd()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.Shell()
	assert.ErrorContains(t, err, "failed to start instance")
}

func TestShellWaitFails(t *testing.T) {
	cmdCallCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		cmdCallCount++
		if cmdCallCount == 1 {
			return exec.Command("/bin/echo", "Status: RUNNING") // StartIfNeeded
		}
		return exec.Command("/bin/false") // Wait → isVM
	})
	defer restoreCmd()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.Shell()
	assert.ErrorContains(t, err, "failed to wait for instance")
}

func TestShellLxcExecFails(t *testing.T) {
	cmdCallCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		cmdCallCount++
		switch cmdCallCount {
		case 1:
			return exec.Command("/bin/echo", "Status: RUNNING") // StartIfNeeded
		case 2:
			return exec.Command("/bin/echo", "Type: container") // Wait → isVM
		case 3:
			return exec.Command("/bin/false") // lxcExec
		default:
			return exec.Command("/bin/true")
		}
	})
	defer restoreCmd()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.Shell()
	assert.ErrorContains(t, err, "failed to lxc exec")
}

func TestShellWithParams(t *testing.T) {
	cmdCallCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		cmdCallCount++
		switch cmdCallCount {
		case 1:
			return exec.Command("/bin/echo", "Status: RUNNING") // StartIfNeeded
		case 2:
			return exec.Command("/bin/echo", "Type: container") // Wait → isVM
		case 3:
			return exec.Command("/bin/true") // lxcExec
		default:
			return exec.Command("/bin/true")
		}
	})
	defer restoreCmd()
	app := App{
		Config: Config{Label: "l", System: NewSystem("s")},
		Opts:   Opts{Params: []string{"echo", "hi"}},
	}
	assert.Nil(t, app.Shell())
}

func TestLaunchContainerOk(t *testing.T) {
	cmdCallCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		cmdCallCount++
		switch cmdCallCount {
		case 1:
			return exec.Command("/bin/true") // lxc launch
		case 2:
			return exec.Command("/bin/echo", "Type: container") // lxc info
		case 3:
			return exec.Command("/bin/true") // lxc exec use_pty
		case 4:
			return exec.Command("/bin/true") // lxc exec cloud-init
		default:
			return exec.Command("/bin/true")
		}
	})
	defer restoreCmd()

	restoreCmdCtx := Patch(&commandContext, func(_ context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/echo", "Debian")
	})
	defer restoreCmdCtx()

	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	assert.Nil(t, app.Launch())
}

func TestLxcLaunchFails(t *testing.T) {
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/false")
	})
	defer restoreCmd()

	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.lxcLaunch()
	assert.ErrorContains(t, err, "failed to create instance")
}

func TestLxcLaunchVM(t *testing.T) {
	var args []string
	restoreCmd := Patch(&command, func(name string, arg ...string) *exec.Cmd {
		args = append([]string{name}, arg...)
		return exec.Command("/bin/true")
	})
	defer restoreCmd()

	app := App{Config: Config{Label: "l", System: NewSystem("s"), Virtualization: "vm"}}
	assert.Nil(t, app.lxcLaunch())
	assert.Contains(t, args, "--vm")
}

func TestLaunchWaitFails(t *testing.T) {
	cmdCallCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		cmdCallCount++
		if cmdCallCount == 1 {
			return exec.Command("/bin/true") // lxc launch
		}
		return exec.Command("/bin/false") // lxc info
	})
	defer restoreCmd()

	restoreCmdCtx := Patch(&commandContext, func(_ context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/echo", "Debian")
	})
	defer restoreCmdCtx()

	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.Launch()
	assert.ErrorContains(t, err, "failed to wait for instance")
}

func TestLaunchUsePtyFails(t *testing.T) {
	cmdCallCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		cmdCallCount++
		switch cmdCallCount {
		case 1:
			return exec.Command("/bin/true") // lxc launch
		case 2:
			return exec.Command("/bin/echo", "Type: container") // lxc info
		case 3:
			return exec.Command("/bin/false") // lxc exec use_pty
		default:
			return exec.Command("/bin/true")
		}
	})
	defer restoreCmd()

	restoreCmdCtx := Patch(&commandContext, func(_ context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/echo", "Debian")
	})
	defer restoreCmdCtx()

	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.Launch()
	assert.ErrorContains(t, err, "use_pty setup failure")
}

func TestLaunchCloudInitFails(t *testing.T) {
	cmdCallCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		cmdCallCount++
		switch cmdCallCount {
		case 1:
			return exec.Command("/bin/true") // lxc launch
		case 2:
			return exec.Command("/bin/echo", "Type: container") // lxc info
		case 3:
			return exec.Command("/bin/true") // lxc exec use_pty
		case 4:
			return exec.Command("/bin/false") // lxc exec cloud-init
		default:
			return exec.Command("/bin/true")
		}
	})
	defer restoreCmd()

	restoreCmdCtx := Patch(&commandContext, func(_ context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/echo", "Debian")
	})
	defer restoreCmdCtx()

	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.Launch()
	assert.ErrorContains(t, err, "cloud-init failure")
}

func TestLaunchQuirkFails(t *testing.T) {
	restoreCmdCtx := Patch(&commandContext, func(_ context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/printf", "Distributor ID: Ubuntu\nRelease: 22.04")
	})
	defer restoreCmdCtx()

	cmdCallCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		cmdCallCount++
		switch cmdCallCount {
		case 1:
			return exec.Command("/bin/true") // lxc launch
		case 2:
			return exec.Command("/bin/echo", "Type: container") // lxc info
		case 3:
			return exec.Command("/bin/true") // lxc exec use_pty
		case 4:
			return exec.Command("/bin/false") // lxc exec bus wait (quirk)
		default:
			return exec.Command("/bin/true")
		}
	})
	defer restoreCmd()

	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.Launch()
	assert.ErrorContains(t, err, "LP: #1878225 workaround failure")
}

var sudoLoginTests = []struct {
	summary string
	script  string
	app     App

	expected []string
}{{
	summary:  "simple shell",
	script:   "echo hi",
	app:      App{Config: Config{Label: "l", System: NewSystem("s")}},
	expected: []string{"sudo", "--login", "--user", "user", "sh", "-c", "echo hi"},
}, {
	summary:  "cd command",
	script:   `cd "/project" && exec $SHELL`,
	app:      App{Config: Config{Label: "l", System: NewSystem("s")}},
	expected: []string{"sudo", "--login", "--user", "user", "sh", "-c", `cd "/project" && exec $SHELL`},
}}

func TestSudoLogin(t *testing.T) {
	for _, test := range sudoLoginTests {
		assert.Equal(t, test.expected, test.app.sudoLogin(test.script), test.summary)
	}
}

func TestLp1878225QuirkNotJammy(t *testing.T) {
	restoreCmdCtx := Patch(&commandContext, func(_ context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/echo", "Debian")
	})
	defer restoreCmdCtx()
	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	assert.Nil(t, app.lp1878225Quirk())
}

func TestLp1878225QuirkBusWaitFails(t *testing.T) {
	restoreCmdCtx := Patch(&commandContext, func(_ context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/printf", "Distributor ID: Ubuntu\nRelease: 22.04")
	})
	defer restoreCmdCtx()

	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/false")
	})
	defer restoreCmd()

	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.lp1878225Quirk()
	assert.ErrorContains(t, err, "bus wait failure")
}

func TestLp1878225QuirkSeededStopFails(t *testing.T) {
	restoreCmdCtx := Patch(&commandContext, func(_ context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/printf", "Distributor ID: Ubuntu\nRelease: 22.04")
	})
	defer restoreCmdCtx()

	cmdCallCount := 0
	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		cmdCallCount++
		if cmdCallCount == 1 {
			return exec.Command("/bin/true")
		}
		return exec.Command("/bin/false")
	})
	defer restoreCmd()

	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	err := app.lp1878225Quirk()
	assert.ErrorContains(t, err, "seeded stop failure")
}

func TestLp1878225QuirkJammyOk(t *testing.T) {
	restoreCmdCtx := Patch(&commandContext, func(_ context.Context, _ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/printf", "Distributor ID: Ubuntu\nRelease: 22.04")
	})
	defer restoreCmdCtx()

	restoreCmd := Patch(&command, func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("/bin/true")
	})
	defer restoreCmd()

	app := App{Config: Config{Label: "l", System: NewSystem("s")}}
	assert.Nil(t, app.lp1878225Quirk())
}
