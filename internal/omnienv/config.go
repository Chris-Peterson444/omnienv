package omnienv

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

var cfgName = ".omnienv.yaml"

var ErrCfgNotFound = errors.New("Config not found")

type System struct {
	Name  string
	Image string
}

func NewSystem(val string) System {
	return System{Name: val}
}

func (sys System) LaunchImage() string {
	if sys.Image == "" {
		return "ubuntu-daily:" + sys.Name
	}
	return sys.Image
}

func (sys *System) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// if string, unmarshal to Name and done
	var err error
	var name string
	if err = unmarshal(&name); err == nil {
		*sys = NewSystem(name)
		return nil
	}

	type systemImage struct {
		Image string
	}

	// if map with single key, unmarshal key to Name and set Image
	var dict map[string]systemImage
	if err = unmarshal(&dict); err != nil {
		return err
	}

	switch len(dict) {
	case 0:
		return errors.New("empty system map")
	case 1:
		for name, img := range dict {
			*sys = System{Name: name, Image: img.Image}
		}
		return nil
	default:
		return errors.New("multiple system keys, expected one")
	}
}

type Config struct {
	// System indicates what distribution and version to base this upon.
	// Specifying distribution not yet implemented.
	// When distribution is omitted, Ubuntu is used.
	System System
	// Label is set to the basename of RootDir by default, or may be
	// overwritten with the supplied value.  Used as part of the name of
	// the instance, along with System.
	Label string
	// RootDir is the writable base directory that the container has access
	// to.  This field is optional, and if unsupplied uses the parent
	// directory of the omnienv.yaml config.
	RootDir string `yaml:"basedir"`
	// Backend indicates upon what we are running the instance.
	// Only "lxd" is implemented.
	Backend string
	// Virtualization chooses between "container" (default) and "vm".
	Virtualization string

	// unsupported keys that are unmarshalled for warning purposes
	Project string
	Series  string
}

func (cfg Config) isVM() bool {
	return cfg.Virtualization == "vm"
}

func (cfg Config) lxdLaunchConfig(user UserInfo) string {
	tmap := map[string]string{
		"WORKDIR":  cfg.RootDir,
		"HOST_UID": strconv.Itoa(user.UID),
		"HOST_GID": strconv.Itoa(user.GID),
	}

	template := `
config:
  raw.idmap: |-
    uid ${HOST_UID} 1000
    gid ${HOST_GID} 1000
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
    source: ${WORKDIR}
`
	return os.Expand(template, func(key string) string {
		return tmap[key]
	})
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func findConfig(dir string) (string, error) {
	for {
		cur := dir + "/" + cfgName
		if exists(cur) {
			return cur, nil
		}
		if dir == "/" {
			break
		}
		dir = filepath.Dir(dir)
	}

	return "", ErrCfgNotFound
}

func loadConfig(path string) (Config, error) {
	// directory traversal is an intended feature
	data, err := os.ReadFile(path) //gosec:disable G304
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.RootDir == "" {
		cfg.RootDir = filepath.Dir(path)
	}

	if cfg.Label == "" {
		cfg.Label = filepath.Base(cfg.RootDir)
	}

	if cfg.System.Name == "" {
		cfg.System = NewSystem(os.Getenv("DEFAULT_SERIES"))
	}

	if cfg.Virtualization == "" {
		cfg.Virtualization = "container"
	}

	if cfg.Project != "" {
		log.Printf("WARN: unsupported key project=%s", cfg.Project)
	}
	if cfg.Series != "" {
		log.Printf("WARN: unsupported key series=%s", cfg.Series)
	}
	debugLog("loadConfig config=%+v", cfg)
	return cfg, nil
}

func GetConfig() (Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	cfgPath, err := findConfig(dir)
	if err != nil {
		return Config{}, err
	}

	return loadConfig(cfgPath)
}
