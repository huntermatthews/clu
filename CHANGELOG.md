# Changelog

All notable *user visible* changes to Clu will be documented in this file.

This changelog follows [keepachangelog](https://keepachangelog.com/en/1.0.0/) format, and is intended
for human consumption.

This project versioning [Semantic Versioning](https://semver.org) but formatted in python's
[PEP440](https://peps.python.org/pep-0440/) so it can work in python.

<!-- Added Changed Removed Fixed Bugs --- Please leave this comment before unreleased section.-->

## [Unreleased]

## [v1.3] - 2025-09-18

### Added

- Added verbosity/tiers of facts. Allows easy selection of how much detail you want to see.
- Added a full man page as our user documentation.
- Better --help text.
- Source for /System/Library/CoreServices/SystemVersion.plist on Darwin.

### Fixed

- Simplified program logic in number of areas, including more comments.
- Improved testing code to be more pylance compliant.


## [v1.2] - 2025-09-04

### Added

- Noted that man page was incomplete - its a WIP.
- Added a `--net` command line option to the report command - required now for any network queries.
  `dnf checkupdate` was too slow.
- Implement a first cut of a parser for the AWS IMDSv2 service.
  IMDS has a slightly complicated data model and this needs improvements, but this works for now.

### Fixed

- Due to dependency handling, some parsers were being called twice.
- Fixed one of the tests - strange editor cut-paste error maybe.
- Fix the default argument handling when ./clu is run without the `report` command being explicitly
  given.


## [v1.1] - 2025-09-01

### Added

- Implement the [versioningit](https://github.com/jwodder/versioningit) plugin to replace the old
  setuptools_scm code. Its simpler, cleaner and I like the default versions it creates better.
- Implement parsing of `lsmem` to hopefully, **finally** get a reliable way to determine how much RAM
  a system has. Note that ram is not memory, and /proc/meminfo doesn't yield total ram.

### Bugs

- `dnf check-update` is extremely slow due to checking for the new updates. A future version will gate
  this (and other) network dependant checks behind a --net command line flag to enable networked sources.


## [v1.0] - 2025-08-30

### Added

- GitHub workflow (actions) to build the zip file when we tag a new release.
- The beginnings of a man page. No unix command is complete without a man page.

### Changed

- Changed error handling in the input functions to be consistent. More consistent maybe.
- Parsing of sys_dmi was substantially improved.

### Removed

- Removed the custom logging levels of trace and verbose - mypy didn't like it and it was just
  more trouble than it was worth.

### Fixed

- Fix the fact name for the /no_salt parser - not sure how this got messed up.


## [v0.100] - 2025-08-21

### Added

- `dnf check-update` is now reported - note that this fails on ubuntu.
- Implemented basic network information from `ip addr` command.
- Started thinking about keeping a changelog.
- Got the testing code to pull all the expected results from 1 large host specific json file
  to avoid a ton of duplication. (Makes future changes much easier).

### Changed

- Updated some of the tests.
- Reworked the --debug and --verbosity command line flags.


<!-- markdownlint-disable-file MD024 -->
