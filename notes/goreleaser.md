# GoReleaser migration plan

## Objective

Replace release-specific GitHub Actions logic and native DEB/RPM build workflows with GoReleaser, while retaining Make targets for local development, testing, manpage generation, and local installation.

The migration must preserve these published artifacts:

| Artifact | Platforms |
|---|---|
| Raw `clu` binary | Linux amd64/arm64, Darwin arm64 |
| Raw `clu.exe` binary | Windows amd64 |
| DEB package | Linux amd64/arm64 |
| RPM package | Linux amd64/arm64 |
| SHA-256 checksum manifest | All release artifacts |
| GitHub build attestations | All release artifacts |

## Current responsibilities

- `Makefile` builds `./cmd/clu` with `CGO_ENABLED=0`, generates the embedded manpage, and installs the binary, manpage, and README.
- `release.yaml` cross-builds binaries in a four-entry matrix, renames artifacts, creates checksums, creates a draft release, and finally publishes it.
- `build-deb.yaml` uses `dpkg-buildpackage` and `debian/rules`.
- `build-rpm.yaml` uses `rpmbuild` and `redhat/clu.spec`.
- Each workflow creates GitHub build provenance attestations.

## Target design

A tag push starts one release workflow:

1. Check out the tagged commit with complete Git history.
2. Set up the pinned Go toolchain and `go-md2man`.
3. Generate `pkg/subcmd/clu.1` before building.
4. Run GoReleaser once to build binaries, packages, checksums, and the GitHub Release.
5. Attest the GoReleaser checksum manifest with GitHub Actions.

Keep Make targets for `build`, `build-all`, `man`, `install`, `test`, `coverage`, `fmt`, `lint`, and `vet`. GoReleaser is a release orchestrator, not a replacement for local developer tooling.

## GoReleaser configuration

Add a version-2 `.goreleaser.yaml` with these components.

### Build

Configure one Go build with:

- `main: ./cmd/clu`
- `binary: clu`
- `env: [CGO_ENABLED=0]`
- targets: `linux/amd64`, `linux/arm64`, `darwin/arm64`, and `windows/amd64`
- explicit artifact naming to preserve current release names:
  - `clu-linux-amd64`
  - `clu-linux-arm64`
  - `clu-darwin-arm64`
  - `clu-windows-amd64.exe`

Use a pre-build hook, or a preceding workflow step, to generate the manpage. The generated file must exist before GoReleaser invokes `go build`, because it is embedded by `pkg/subcmd/help.go`.

### Raw binaries and checksums

- Configure binary-format archives, rather than the default compressed archives, if the existing bare executable release assets must be retained.
- Configure one predictable checksum filename, such as `checksums.txt` or `SHA256SUMS`.
- Ensure the checksum manifest includes raw binaries and DEB/RPM packages.

### DEB and RPM packages

Use an nFPM configuration that packages Linux build artifacts only.

Required package contents:

```text
/usr/bin/clu
/usr/share/man/man1/clu.1
/usr/share/doc/clu/README.md
```

Set package metadata deliberately:

- package name: `clu`
- license: `LGPL-2.1-only`
- homepage: `https://github.com/huntermatthews/clu`
- summary/description: use the canonical README program summary
- release number: `1`
- DEB architecture names: `amd64`, `arm64`
- RPM architecture names: `x86_64`, `aarch64`

Configure filenames to preserve the current conventions:

```text
clu_<version>-1_<arch>.deb
clu-<version>-1.<arch>.rpm
```

nFPM removes the need for `dpkg-buildpackage`, `debhelper`, `rpmbuild`, `debian/rules`, and `redhat/clu.spec`. It packages the already-built static Linux artifacts directly.

### GitHub Release

Configure GoReleaser to:

- use the pushed `vX.Y.Z` tag as the release source;
- upload all configured binaries, packages, and checksums;
- generate a changelog or use explicitly supplied release notes;
- publish only after all configured artifacts are complete.

Decide whether releases should remain drafts for human review. Unlike the current workflow, package generation occurs before the release is published, so no separate draft-and-later-upload coordination is necessary.

### Attestations

Keep GitHub Actions provenance generation after GoReleaser.

- Grant `contents: write`, `attestations: write`, and `id-token: write`.
- Use `actions/attest@v4` with the GoReleaser checksum manifest.
- Verify that checksum-subject attestations support the existing `gh attestation verify` commands and SLSA claims before removing the current per-artifact attestation steps.

## CI workflow plan

Replace the three release workflows with one tag-triggered workflow after validation.

Required job steps:

1. `actions/checkout` with `fetch-depth: 0`.
2. `actions/setup-go` using the project-supported Go version.
3. Install a pinned `go-md2man` version for the host architecture.
4. Generate the manpage, or let the GoReleaser pre-build hook do so.
5. Run `goreleaser release --clean` through the GoReleaser GitHub Action.
6. Attest the checksum manifest.

Add a separate pre-release quality job if release gating is desired:

```text
make test
make vet
make lint
```

The current release workflows do not run these checks.

## Migration phases

### Phase 1: establish artifact expectations

Record expected artifacts from a current tagged release:

- exact file names;
- package metadata and installed file paths;
- checksum format;
- GitHub release body behavior;
- attestation verification commands.

Decide whether raw binaries must remain bare files or may become archives.

### Phase 2: add configuration without publishing

Add `.goreleaser.yaml` and run locally:

```text
goreleaser release --snapshot --clean
```

Validate:

- all four binary targets build;
- Windows artifact ends in `.exe`;
- manpage is embedded and packaged;
- DEB and RPM packages contain the expected files;
- version and architecture names match expectations;
- checksums cover every intended release asset.

### Phase 3: parallel CI validation

Add a non-publishing workflow for tags or manual dispatch.

Compare GoReleaser artifacts with the existing release pipeline. Do not remove current workflows until the following are equivalent:

- executable behavior;
- package installation paths;
- package version and release fields;
- release asset names;
- checksum verification;
- provenance verification.

### Phase 4: release cutover

1. Replace `release.yaml`, `build-deb.yaml`, and `build-rpm.yaml` with the single GoReleaser release workflow.
2. Perform one monitored tagged release.
3. Verify published binaries, packages, checksums, and attestations.
4. Retire `debian/` and `redhat/` packaging definitions only after that release is accepted.

### Phase 5: documentation cleanup

Update `README.md` release download, checksum, and attestation instructions if artifact names or checksum names changed.

## Decisions required before implementation

1. Should GitHub releases remain drafts, or publish automatically after a successful GoReleaser run?
2. Must raw binaries retain their current names and unarchived format?
3. Should release notes remain the current hand-written text, or use generated changelogs?
4. Is checksum-manifest attestation acceptable in place of individual artifact attestation steps?
5. Should package builds retain native `dpkg`/`rpmbuild` tooling for policy reasons, or is nFPM output acceptable?
6. How should the CLI version be injected?

## Version caveat

The current workflows pass `VERSION`, but the Makefile and Go sources do not consume it. `pkg/global/buildinfo.go` reads `runtime/debug.ReadBuildInfo()` instead.

GoReleaser can supply the tag version through `ldflags`, but `clu` needs an explicit Go variable to receive it. Resolve this separately so a release binary reports the intended version rather than the default development module version.

## Rollback

Keep the existing workflows and package definitions through the first successful GoReleaser-produced release. If artifact, package, release, or attestation validation fails, restore the current workflows and correct the configuration before attempting another cutover.
