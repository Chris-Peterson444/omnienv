# TODO

## Before v0.2

- [X] **Rename config key to `project`**. Was `basedir`, then `rootdir`, now
  `project`. Added `basedir` as a deprecated (warning) key.

- [X] **Change `var cfgName` to `const`** (`config.go:13`). Never mutated.

## v0.3+

- [X] **Test `NewApp()`**. Trivial constructor (0% coverage) — only matters if
  it grows logic.

- [X] **Add empty-args guard to `run()`/`runDevNull()`** (`exec.go`). Currently
  panics on `args[0]` if called with no arguments. All existing callers are
  safe, but it's a footgun.

- [X] **Review `test_runner.py` system list**. `resolute` is Ubuntu 26.04
  (intentional, this container runs it).

- [X] **Use `t.Helper()` in `patchEnv`** (`config_test.go`). Minor test hygiene.

- [X] test coverage for more of ./cmd ?
