# TODO

Items from a code review of the current codebase, ordered by rough priority.

## High Priority

- [ ] **Replace global variable indirection with dependency injection**
  (`globals.go`). `command`, `commandContext`, `timeSleep`
  are reassignable package vars mutated by tests. This makes production reasoning hard
  and tests leaky. Move them into `App` struct fields or an interface.

- [X] **Fix `Run()` error swallowing** (`cmd/oe/main.go:14-15`). When `GetOpts`
  returns a parse error, `Run()` returns `nil` instead of propagating the error.
  This means bad flags (e.g. `--invalid`) print help text but exit 0. Change to
  `return err`.

- [X] **Fix `runDevNull` misleading name/behavior** (`internal/omnienv/exec.go`).
  The function does not suppress stdout/stderr — it leaves `cmd.Stdout` and
  `cmd.Stderr` nil, inheriting the parent process's fds. Either redirect to
  `io.Discard` or rename to `runQuiet`.

- [X] **Embed version string for `--version`** (`cmd/oe/main.go`). The `--version`
  flag uses `debug.ReadBuildInfo()` which returns `"unknown"` for locally-built
  binaries. Add a `var Version = "dev"` in a `version.go` and wire `-ldflags` into
  the Makefile so tagged releases display the correct version.

## Medium Priority

- [X] **Create CHANGELOG.md**. No release notes exist beyond git tags. Add a
  changelog summarizing the v0.1 → v0.2 delta: new features, bug fixes, test
  improvements, CI setup, and notable refactors.

- [X] **Update README project status**. Currently says "Pre-alpha". At v0.2,
  consider updating to reflect a more stable state.

- [ ] Skip bionic-vm integration test
