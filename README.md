# omnienv

omnienv aims to make it simpler to create and use container-based (or vm)
virtual environments, and to make a directory tree available to that virtual
environment.

omnienv is implemented using LXD and built with golang.

omnienv mounts the working directory of the project being managed at `/project`.

omnienv arranges for a `user` account to be created in the environment with
passwordless sudo access, mapped to the host user's uid/gid. The resulting
shell has a pty for compatibility with various terminal applications.

omnienv shells can handle waiting for startup of containers and
virtual-machines, including handling the non-instant wait time for a LXD vm to
go from stopped to shell-ready. If the environment is stopped when a shell is
requested, omnienv first transparently starts that environment and initiates
the shell when possible.

## project status

Usable but expect breaking changes. Config file format under active development.

## installation

    go install github.com/dbungert/omnienv/cmd/oe@latest

## quick start

1. Identify a project directory you would like to associate with a container or
   vm.
2. Add a file `.omnienv.yaml` to this directory with the following contents:
```yaml
system: noble
```
3. Run `oe --launch`. The container will be created and the project directory
   will be mounted at `/project` in that environment

## usage examples

```
oe --launch
```

> Create the environment instance, and shell into the instance when ready.

```
oe
```

> Shell into an existing instance, first starting if needed.

```
oe make
```

> Shell into instance, run `make`, and return the exit code.

```
lxc remove foo-resolute
```

> Standard LXD management commands can be used with the container. In this
> example, delete the container for the foo project of series resolute.

## options

* `--launch`: Create the LXD environment (container or VM) before opening a
  shell.
* `-s`, `--system`: Override the `system` value from the config file.
* `-v`, `--verbose`: Increase logging verbosity to DEBUG level.
* `--version`: Print the version and exit.

## config file format

An omnienv project is defined by the `.omnienv.yaml` config file and location.
The parent directory of that config is the working directory, and that working
directory is mounted read-write at `/project` in the environment.

The config file is discovered by walking up the directory tree from the current
working directory.

These fields are supported:

| Field | Description | Accepted values | Default |
|-------|-------------|----------------|---------|
| `system` | the OS version of environment to use. At this time only Ubuntu is supported, and only using the series names, so `bionic` for Ubuntu 18.04, `jammy` for Ubuntu 22.04, `noble` for Ubuntu 24.04 and so on. Only bionic and newer are supported. | bionic or newer Ubuntu series name (e.g. `bionic`, `jammy`, `noble`, `plucky`). | `DEFAULT_SERIES` environment variable, empty if unset. |
| `virtualization` | use a `container` or `vm`. | `container`, `vm`. | `container`. |
| `label` | the prefix for the environment name, this is inferred from the basename of the `project` config. The full LXD instance name is `<label>-<system>`. | Any string. | basename of the `project` directory. |
| `project` | which directory to mount read-write in the environment. | A filesystem directory path. | Parent directory of `.omnienv.yaml`. |

The `system` field's map form can specify a custom launch image. For example:

```yaml
system:
  jammy:
    image: ubuntu:j
```

The deprecated keys `basedir` and `series` are accepted but produce a warning.

## known issues

* Bionic (18.04) does not work as a VM.

## expected project direction

* The config file format is under active work, and the terms used may change.
* Support expected for Ubuntu based on version numbers, and other Linux
  distributions handled by the LXD `images` remote.
