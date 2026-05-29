# TODO

Items from a code review of the current codebase, ordered by rough priority.

## High Priority

- [X] **Add `context.Context` with timeout to VM wait loop** (`app.go:107-128`). The
  `wait()` method polls `lxc exec … /bin/true` in an unbounded loop. If a VM never
  becomes reachable, this hangs forever. Add a context with a deadline or a max-retry
  count.

- [X] **Add timeout to `lp1878225Quirk` bus-wait loop** (`app.go:172-184`). The `for i
  in $(seq 10)` fallback is already bounded, but the initial `lsb_release` check has
  no timeout on the outer `lxc exec`. Not critical, but adds defense in depth.

- [ ] **Replace global variable indirection with dependency injection**
  (`globals.go`). `command`, `timeSleep`
  are reassignable package vars mutated by tests. This makes production reasoning hard
  and tests leaky. Move them into `App` struct fields or an interface.

- [X] **Add `Run() error` pattern to `main.go`**. Currently `main()` does real work
  and calls `fatal()` inline. Idiomatic Go extracts a `func (app App) Run() error` and
  keeps `main()` as:

  ```go
  func main() {
      if err := run(); err != nil {
          log.Fatal(err)
      }
  }
  ```

- [X] **Validate `UnmarshalYAML` map-parsing edge case** (`cfg.go:33-57`). If the map
  branch succeeds with an empty map, `System` gets zero-valued silently. Multiple map
  keys also silently overwrite (last wins). Add validation.

- [X] **Fall back from `os.Getenv("PWD")` to `os.Getwd()`** (`app.go:265`). If the
  shell doesn't set `PWD`, the subdirectory-relative `cd` silently falls through and
  `dest` stays `/project`. Log a warning or fall back to `os.Getwd()`.

## Medium Priority

- [X] **Extract an internal package** (e.g. `pkg/omnienv`). Currently everything is in
  `package main`. Extracting config parsing, instance management, and cloud-init
  logic into a separate package prevents import cycles, enables reuse, and improves
  testability.

- [X] **Consolidate LXD interaction approach** (`app.go`). Some operations use the Go
  LXD client (`start`, `isVM`), others shell out to the `lxc` CLI (`lxcRun`,
  `lxcExec`, `wait`). Pick one approach and migrate consistently.

- [X] **Set up CI pipeline**. The repo has no `.github/workflows/` or equivalent CI
  config. Add at minimum: `go build`, `go test`, `go vet`, and the pre-commit hooks.

## Low Priority

- [X] **Refactor `lp1878225Quirk` into detection + remediation** (`app.go:138-192`).
  The sentinel-value pattern (exit code 225) works but is opaque. Split the function
  into `isAffected()` and `applyWorkaround()`.

- [X] **Un-pin or pin-finalize pre-commit hooks** (`.pre-commit-config.yaml`). The
  `tekwizely/pre-commit-golang` rev is `v1.0.0-rc.1` (a release candidate). Consider
  moving to a stable release. Also `go-sec-repo-mod` and `golangci-lint-mod` are
  commented out -- decide whether to enable or remove.

## Test Coverage Gaps

- [X] **`startIfNeeded`**: only tests connect-fail path. No coverage for the state
  machine (Stopped → start, Running → no-op, unexpected status → error).

- [X] **`wait()`**: untested.

- [X] **`launch()`**: untested (involves real LXD, but could integration-test the
  config-piping / args construction).

- [X] **`lp1878225Quirk()`**: untested.

- [X] **`lxcRun()`**: untested.

## Test Hygiene

- [X] **Use `t.Cleanup` instead of manual `defer` in `TestGetConfig`**
  (`cfg_test.go:194-208`). The `os.Chdir` restore happens in a `defer`; if the test
  panics between `Chdir` and the defer, the working directory is corrupted.
  `t.Cleanup` is the idiomatic fix.

- [X] **Make `UserInfo` fields exported** (`user.go:6-7`). `uid` and `gid` are
  unexported. Fine while tests are in `package main`, but will break if extracted to a
  separate package.
