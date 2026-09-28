# go-common.mk adoption analysis

`go-common.mk` has been added from another repository and overlaps heavily with this repo’s existing `Makefile`. This note compares the two without proposing immediate code changes.

## High-level summary

`go-common.mk` duplicates many generic Go project operations currently defined in `Makefile`:

- build
- cross-build
- install
- test
- coverage
- tidy
- format checking
- vet
- setup
- clean

The current `Makefile`, however, contains important `clu`-specific behavior that `go-common.mk` does not currently model:

- `clu` builds from `./cmd/clu`, not from the repository root.
- `clu` generates and embeds a man page at `pkg/subcmd/clu.1`.
- `install` installs the binary, man page, and README.
- coverage is configured with explicit `-coverpkg` values.
- `setup` currently runs `go mod download`.
- `fmt` exists locally, but not in `go-common.mk`.

As-is, `go-common.mk` is a useful base for standardization, but it assumes a simpler Go layout than this repository currently has.

## Important direct-adoption issue

`go-common.mk` builds with:

```make
go build -o dist/$(BINARY) .
```

This repository has no root `.go` files. The main package is under:

```text
cmd/clu/main.go
```

The current `Makefile` correctly uses:

```make
CGO_ENABLED=0 go build -o dist/$(BINARY) ./cmd/clu
```

Replacing the build target with `go-common.mk` as-is would therefore break `make build` for `clu`.

A reusable shared makefile probably needs a tunable such as:

```make
GO_BUILD_PACKAGE ?= .
```

or similar.

## Target comparison

| Area | Current `Makefile` | `go-common.mk` | Difference |
|---|---|---|---|
| `build` | Builds `./cmd/clu` | Builds `.` | Not equivalent for this repo |
| `dist/$(BINARY)` | Depends on Go sources and generated man page | Depends on `$(SRCS)` | `go-common.mk` needs custom `SRCS` and build package support |
| cross-build | Target is `all` | Target is `build-all` | Same purpose, different target name |
| cross-build implementation | Recursive `$(MAKE)` loop, then `mv` | Pattern targets `dist/$(BINARY)-%` | `go-common.mk` is cleaner and more make-native |
| `install` | Installs binary, man page, README | Installs binary only | Not equivalent |
| man page generation | Present | Absent | Repo-specific behavior remains local |
| `test` | Delegates to `test-units` | Same | Equivalent |
| `test-units` | `./pkg/... ./cmd/...` | `$(TEST_UNITS_PKGS)`, default `./...` | Configurable in `go-common.mk` |
| `test-list` | Same recipe | Same recipe | Equivalent |
| `coverage` | Uses `-coverpkg=./pkg/...,./cmd/...` | Plain `go test -coverprofile ./...` | Coverage semantics differ |
| `coverage-html` | Same | Same | Equivalent once coverage target is acceptable |
| `tidy` | Same | Same | Equivalent |
| `fmt` | Present | Missing | Would be lost unless kept locally or added |
| `fmt-check` | Runs `gofmt -d`, does not fail on differences | Uses `test -z "$$(gofmt -l .)"`, fails if unformatted | `go-common.mk` is stricter/better as a check |
| `lint` | Commented out | Enabled, depends on `require-golangci-lint` | `require-golangci-lint` is not defined in current `common.mk` |
| `vet` | Same | Same | Equivalent |
| `setup` | Creates `dist`, runs `go mod download` | Creates `dist` only | Behavior differs |
| `clean` | Removes `dist/*`, `coverage.out`, `.go-md2man-installed` | Removes `dist/*`, `coverage.out`, `$(CLEAN_EXTRA)` | Can be made equivalent with `CLEAN_EXTRA := .go-md2man-installed` |

## Variable/default differences

| Variable | Current `Makefile` | `go-common.mk` | Notes |
|---|---|---|---|
| `BINARY` | `clu` | Expected to be set before include | Compatible |
| `PREFIX` | `/usr/local` | `$(HOME)/.local` | Install default changes significantly |
| `CGO_ENABLED` | Inline in build recipe | Default/exported as `0` | `go-common.mk` is more consistent |
| `PLATFORMS` | Same values, different order | Same values, different order | Output set is equivalent |
| source deps | `GO_SOURCES` includes `cmd`, `pkg`, man markdown, makefiles | `SRCS` defaults to root `*.go`, `go.mod`, `go.sum` | Default is insufficient for `clu` |
| test packages | Hard-coded `./pkg/... ./cmd/...` | `TEST_UNITS_PKGS ?= ./...` | Tunable |
| clean extras | Hard-coded | `CLEAN_EXTRA` | Better abstraction in `go-common.mk` |

## `common.mk` interaction

`go-common.mk` explicitly expects to be included after `common.mk`, and that matters.

In particular, `go-common.mk`’s cross-build rule relies on `.ONESHELL` so these shell variables survive across recipe lines:

```make
STEM=$*
OS=$${STEM%-*}
ARCH=$${STEM#*-}
GOOS=$$OS GOARCH=$$ARCH go build -o $@ .
```

That works with this repository’s `common.mk`, which sets:

```make
.ONESHELL:
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
```

So the dependency on `common.mk` is real, not just organizational.

## Things `go-common.mk` improves

- Cleaner cross-build implementation using real pattern targets.
- Centralized defaults for common Go tasks.
- Exported `CGO_ENABLED`.
- Configurable `TEST_UNITS_PKGS`.
- Configurable `CLEAN_EXTRA`.
- `fmt-check` actually fails when files need formatting, unlike the current `gofmt -d` recipe.

## Things that need attention before standardizing

For `go-common.mk` to work well across multiple Go repositories, consider adding first-class tunables/hooks instead of expecting target overrides that produce “overriding recipe” warnings.

Useful additions may include:

```make
GO_BUILD_PACKAGE ?= .
GO_BUILD_FLAGS ?=
BUILD_DEPS ?=
COVER_PKGS ?= ./...
COVERPKG ?=
INSTALL_EXTRA_TARGETS ?=
SETUP_EXTRA_TARGETS ?=
GOFMT_PATHS ?= .
```

Possibly also add:

```make
.PHONY: all
all: build-all
```

if existing repositories commonly use `make all`.

Also, either `common.mk` or `go-common.mk` should define `require-golangci-lint`, because the current `lint` target depends on it but this repository does not currently provide it.

## Likely `clu`-specific leftovers after adoption

Even with `go-common.mk`, this repository probably still needs local targets or overrides for:

- `man`
- `$(MANPAGE)`
- `.go-md2man-installed`
- installing the man page and README
- using `./cmd/clu` as the build package
- preserving current coverage behavior, if desired
- possibly `fmt`
- possibly `all` as an alias for `build-all`

So the split would likely be:

- `go-common.mk`: shared Go standards
- `Makefile`: repo identity, docs/manpage generation, install extras, and package-layout-specific settings
