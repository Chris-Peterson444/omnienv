# TODO

## Before v0.2

- [ ] **Fix README/config key mismatch**. README documents `rootdir:` but the
  struct tag in `config.go:77` says `yaml:"basedir"`. Fix whichever is wrong so
  users can successfully configure the root directory.

- [ ] **Change `var cfgName` to `const`** (`config.go:13`). Never mutated.

## v0.3+

- [ ] **Test `NewApp()`**. Trivial constructor (0% coverage) — only matters if
  it grows logic.

- [ ] **Add empty-args guard to `run()`/`runDevNull()`** (`exec.go`). Currently
  panics on `args[0]` if called with no arguments. All existing callers are
  safe, but it's a footgun.

- [ ] **Review `test_runner.py` system list**. `resolute` references Ubuntu
  25.04; verify it's intentional and works.

- [ ] **Use `t.Helper()` in `patchEnv`** (`config_test.go`). Minor test hygiene.
