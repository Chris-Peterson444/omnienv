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

## usage

1. Identify a project directory you would like to associate with a container or
   vm.
2. Add a file `.omnienv.yaml` to this directory with the following contents:
```yaml
system: noble
```
3. Run `oe --launch`. The container will be created, the project directory
   mounted at `/project` in that environment, and an interactive shell will
   start in the same directory you are in right now (but in the environment).
4. Standard LXD management commands can be used with the container. For
   instance, this container can be deleted with `lxc delete myproject-noble`,
   where `myproject` is the basename of the directory containing
   `.omnienv.yaml`.
5. Return to this same instance later by running `oe` from the directory with
   `.omnienv.yaml` or lower.
6. To run a non-interactive command inside the environment, pass it after `--`:
   `oe -- make build`. Positional arguments after a flag terminator or after
   non-option args are treated as a command to execute.

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
* `system`: the OS version of environment to use. At this time only Ubuntu is
  supported, and only using the series names, so `jammy` for Ubuntu 22.04 and
  `noble` for Ubuntu 24.04 and so on. Defaults to the `DEFAULT_SERIES`
  environment variable if set.
* `system` (map form): specify a custom launch image. For example:
  ```yaml
  system:
    jammy:
      image: ubuntu:j
  ```
* `virtualization` (optional): use a `container` (default) or `vm`.
* `label` (optional): the prefix for the environment name, this is inferred
  from the basename of the `rootdir` config. The full LXD instance name is
  `<label>-<system>`.
* `rootdir` (optional): which directory to mount read-write in the environment.
  If unspecified, this is set to the parent directory of `.omnienv.yaml`.
* `backend` (optional): which backend to use. Only `lxd` is implemented.

The deprecated keys `project` and `series` are accepted but produce a warning.

## expected project direction

* The config file format is under active work, and the terms used may change.
* Support expected for Ubuntu based on version numbers, and other Linux
  distributions handled by the LXD `images` remote.
