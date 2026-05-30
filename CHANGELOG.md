# Changelog

## [0.2.0] - 2026-05-30

### Added

- `--version` flag now displays version from `git describe` in local builds
- VM integration tests via Python test runner

### Changed

- Config key renamed from `basedir` to `project`; `basedir` is now deprecated
  and produces a warning
- Unit tests overhauled and simpler to understand
- Logging moved away from `slog` to a nicer user format (retaining timestamps)
- `--help` and invalid flags now use consistent single-error output; bad flags
  exit non-zero

### Fixed

- VM wait loop now has a timeout (300 iterations) instead of unbounded polling
- `lp1878225Quirk` `lsb_release` check now has a 30-second timeout
- `UnmarshalYAML` now rejects empty system maps and multiple system keys
- Propagate `lsb_release` failures instead of silently ignoring

### Removed

- Direct LXD Go API client dependency; all LXD interaction uses the `lxc` CLI
- Mockery-generated mocks (~20k lines)

### Infrastructure

- CI pipeline via GitHub Actions with pre-commit, golangci-lint, gosec, staticcheck
- Extracted `internal/omnienv` package for reuse and clean separation

[0.2.0]: https://github.com/dbungert/omnienv/releases/tag/v0.2.0
