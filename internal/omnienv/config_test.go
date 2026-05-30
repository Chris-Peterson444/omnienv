package omnienv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var cfgTests = []struct {
	summary string
	config  Config

	vm bool
}{{
	summary: "basic name",
	config:  Config{Label: "l", System: NewSystem("s")},
	vm:      false,
}, {
	summary: "vm",
	config:  Config{Label: "foo", System: NewSystem("bar"), Virtualization: "vm"},
	vm:      true,
}}

func TestCfg(t *testing.T) {
	for _, test := range cfgTests {
		assert.Equal(t, test.vm, test.config.isVM(), test.summary)
	}
}

func TestFindCfg(t *testing.T) {
	tempdir := t.TempDir()
	assert.Nil(t, os.WriteFile(tempdir+"/"+cfgName, []byte{}, 0644))
	_, err := findConfig(tempdir)
	assert.Nil(t, err)
}

func TestFindParentCfg(t *testing.T) {
	tempdir := t.TempDir()
	dirname := tempdir + "/foo"
	assert.Nil(t, os.Mkdir(dirname, 0750))
	assert.Nil(t, os.WriteFile(tempdir+"/"+cfgName, []byte{}, 0644))
	_, err := findConfig(dirname)
	assert.Nil(t, err)
}

func patchEnv(key string, mock string) func() {
	original, ok := os.LookupEnv(key)
	_ = os.Setenv(key, mock)
	return func() {
		if ok {
			_ = os.Setenv(key, original)
		} else {
			_ = os.Unsetenv(key)
		}
	}
}

var loadCfgTests = []struct {
	summary   string
	data      string
	seriesEnv string

	config Config
	warn   string
	image  string
}{{
	summary: "system",
	data:    "system: plucky",
	config: Config{
		System:         NewSystem("plucky"),
		Virtualization: "container",
	},
}, {
	summary: "system / label",
	data: `
system: warty
label: ubiquity
`,
	config: Config{
		System:         NewSystem("warty"),
		Label:          "ubiquity",
		Virtualization: "container",
	},
}, {
	summary: "system environ / lxd / default virt",
	data:    "backend: lxd",
	config: Config{
		System:         NewSystem("zesty"),
		Backend:        "lxd",
		Virtualization: "container",
	},
}, {
	summary: "system environ / lxd / vm",
	data: `
backend: lxd
virtualization: vm
`,
	config: Config{
		System:         NewSystem("zesty"),
		Backend:        "lxd",
		Virtualization: "vm",
	},
	image: "ubuntu-daily:zesty",
}, {
	summary: "project dir",
	data:    "project: proj",
	config: Config{
		Label:          "proj",
		ProjectDir:     "proj",
		System:         NewSystem("zesty"),
		Virtualization: "container",
	},
}, {
	summary: "deprecated basedir",
	data:    "basedir: /tmp/foo",
	config: Config{
		Basedir:        "/tmp/foo",
		System:         NewSystem("zesty"),
		Virtualization: "container",
	},
	warn: `msg="unsupported key"`,
}, {
	summary: "series",
	data:    "series: warty",
	config: Config{
		Series:         "warty",
		System:         NewSystem("zesty"),
		Virtualization: "container",
	},
	warn: `msg="unsupported key"`,
}, {
	summary: "manual remote for image",
	data: `
system:
    jammy:
        image: ubuntu:j
`,

	config: Config{
		System:         System{"jammy", "ubuntu:j"},
		Virtualization: "container",
	},
	image: "ubuntu:j",
}}

func TestLoadCfg(t *testing.T) {
	tempdir := t.TempDir()
	dirname := tempdir + "/foo"
	assert.Nil(t, os.Mkdir(dirname, 0750))
	filename := dirname + "/" + cfgName

	restoreEnv := patchEnv("DEFAULT_SERIES", "zesty")
	defer restoreEnv()

	for _, test := range loadCfgTests {
		err := os.WriteFile(filename, []byte(test.data), 0644)
		assert.Nil(t, err, test.summary)
		actual, err := loadConfig(filename)
		assert.Nil(t, err, test.summary)
		if test.config.ProjectDir == "" {
			test.config.ProjectDir = dirname
		}
		if test.config.Label == "" {
			test.config.Label = "foo"
		}
		assert.Equal(t, test.config, actual, test.summary)

		if test.image != "" {
			assert.Equal(
				t, test.config.System.LaunchImage(), test.image,
				test.summary,
			)
		}
	}
}

func TestNotLoadCfgUnreadable(t *testing.T) {
	tempdir := t.TempDir()
	filename := tempdir + "/" + cfgName
	assert.Nil(t, os.WriteFile(filename, []byte{}, 000))
	_, err := loadConfig(filename)
	assert.NotNil(t, err)
}

func TestNotLoadCfgUnmarshalable(t *testing.T) {
	tempdir := t.TempDir()
	data := []byte(`{`)
	filename := tempdir + "/" + cfgName
	assert.Nil(t, os.WriteFile(filename, data, 0644))
	_, err := loadConfig(filename)
	assert.NotNil(t, err)
}

func TestGetConfig(t *testing.T) {
	tempdir := t.TempDir()

	curdir, err := os.Getwd()
	assert.Nil(t, err)
	assert.Nil(t, os.Chdir(tempdir))
	t.Cleanup(func() { _ = os.Chdir(curdir) })

	data := []byte("system: warty")
	filename := tempdir + "/" + cfgName
	assert.Nil(t, os.WriteFile(filename, data, 0644))
	actual, err := GetConfig()
	assert.Nil(t, err)
	assert.Equal(t, "warty", actual.System.Name)
}

func TestNotGetConfig(t *testing.T) {
	curdir, err := os.Getwd()
	assert.Nil(t, err)
	assert.Nil(t, os.Chdir("/"))
	t.Cleanup(func() { _ = os.Chdir(curdir) })

	_, err = GetConfig()
	assert.Equal(t, ErrCfgNotFound, err)
}

func TestLXDLaunchConfigWorkdir(t *testing.T) {
	cfg := Config{ProjectDir: "/tmp/b"}
	expected := `
config:
  raw.idmap: |-
    uid 1234 1000
    gid 5678 1000
  user.vendor-data: |
    #cloud-config
    users:
      - name: user
        sudo: ALL=(ALL) NOPASSWD:ALL
        groups: users,admin
        shell: /bin/bash
devices:
  workdir:
    type: disk
    readonly: false
    shift: false
    path: /project
    source: /tmp/b
`

	assert.Equal(t, expected, cfg.lxdLaunchConfig(UserInfo{1234, 5678}))
}

func TestUnmarshalSystemEmptyMap(t *testing.T) {
	data := []byte("system: {}")
	var cfg Config
	err := yaml.Unmarshal(data, &cfg)
	assert.ErrorContains(t, err, "empty system map")
}

func TestUnmarshalSystemMultipleKeys(t *testing.T) {
	data := []byte(`
system:
  jammy:
    image: ubuntu:j
  noble:
    image: ubuntu:n
`)
	var cfg Config
	err := yaml.Unmarshal(data, &cfg)
	assert.ErrorContains(t, err, "multiple system keys")
}

func TestUnmarshalSystemDictFails(t *testing.T) {
	data := []byte(`system: [a, b]`)
	var cfg Config
	err := yaml.Unmarshal(data, &cfg)
	assert.ErrorContains(t, err, "cannot unmarshal")
}
