# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [Unreleased]

### Changed

- Fixed Makefile: binary now correctly rebuilds when any `.go` source file changes

## [1.2.1] - 2026-05-18

### Changed

- Refactored `internal/facts/sources`: adopted defer-add pattern, cleaned up `FactsDB`, fixed `Tier` type across all source files
- Simplified `GetTier` to a loop-based implementation
- Backported cleanup, fixes, and improvements from related project
- Bumped `github.com/alecthomas/kong` 1.14.0 → 1.15.0
- Makefile: force standard shell flags on targets; general cleanup
- Fixed secrets passing to RPM and DEB build steps in CI
- Bumped `actions/upload-artifact` v6 → v7, `actions/download-artifact` v7 → v8, `actions/attest-build-provenance` v3 → v4, `softprops/action-gh-release` v2 → v3


## [1.1.1] - 2026-02-18

### Added

- `help` subcommand: prints the embedded man page
- Multi-distro Linux and macOS package update checking
- SLSA attestation now also covers the SHA256SUMS checksum file

### Fixed

- `.exe` extension handling for Windows binaries simplified and corrected
- Merge regression in a renamed source file
- Windows `systeminfo` parsing and testdata alignment

### Changed

- Bumped `github.com/alecthomas/kong` 1.13.0 → 1.14.0


## [1.1.0] - 2026-02-18

### Added

- `check` subcommand
- `version` subcommand
- Salt fact source (replaces the old `no_salt` stub)
- `FileAgeReader` input type: stat a file's age
- Mock input type for testing

### Changed

- Windows: handle `.exe` suffix on tool paths
- Simplified Makefile build targets
- Bumped minimum `make` version requirement to 3.82
- Release workflow improvements in CI


## [1.0.0] - 2026-01-29

### Added

- First Official Release


[Unreleased]: https://github.com/huntermatthews/clu/compare/v1.2.1...HEAD
[1.2.1]: https://github.com/huntermatthews/clu/compare/v1.1.1...v1.2.1
[1.1.1]: https://github.com/huntermatthews/clu/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/huntermatthews/clu/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/huntermatthews/clu/releases/tag/v1.0.0
