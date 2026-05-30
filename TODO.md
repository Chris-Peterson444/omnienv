# TODO

## Before v0.2

- [X] **Rename config key to `project`**. Was `basedir`, then `rootdir`, now
  `project`. Added `basedir` as a deprecated (warning) key.

- [X] **Change `var cfgName` to `const`** (`config.go:13`). Never mutated.

## v0.3+

- [X] **Test `NewApp()`**. Trivial constructor (0% coverage) — only matters if
  it grows logic.

- [ ] **Add empty-args guard to `run()`/`runDevNull()`** (`exec.go`). Currently
  panics on `args[0]` if called with no arguments. All existing callers are
  safe, but it's a footgun.

- [ ] **Review `test_runner.py` system list**. `resolute` references Ubuntu
  25.04; verify it's intentional and works.

- [ ] **Use `t.Helper()` in `patchEnv`** (`config_test.go`). Minor test hygiene.

- [ ] test coverage for more of ./cmd ?
