# go-common.mk
# Shared targets for building a single Go program.
#
# Tunables (set in the primary Makefile BEFORE the include of this file to customize):
#   BINARY            binary name
#   PREFIX            install prefix
#   CGO_ENABLED       cgo flag
#   PLATFORMS         cross-build target list
#   GO_BUILD_PACKAGE  package passed to go build
#   SRCS              source deps for the binary
#   GO_PKGS           packages passed to testing, formatting, linting, and vetting
#   TEST_GROUP        optional test-name pattern passed to go test -run



# 'make install PREFIX=/usr/local' for a system-wide install
# This default ASSUMES that you have XDG Base Directory Specification compliant directories set up in your environment.
# (or at least something similar to ~/.local/bin and ~/.local/share/man/man1)
PREFIX            ?= $(HOME)/.local
CGO_ENABLED       ?= 0
PLATFORMS         ?= darwin-arm64 linux-amd64 linux-arm64 windows-amd64
GO_BUILD_PACKAGE  ?= ./cmd/$(BINARY)
SRCS              ?= $(wildcard *.go) go.mod $(wildcard go.sum)
GO_PKGS           ?= ./...
TEST_GROUP        ?=


# Non Tunables -- Do not override these in the primary Makefile
PLATFORM_TARGETS := $(patsubst %,dist/$(BINARY)-%,$(PLATFORMS))

##
##@ Build
##

.PHONY: build-default
build-default: setup dist/$(BINARY) ## Setup and build $(BINARY) for the current platform

dist/$(BINARY): $(SRCS)
	CGO_ENABLED=$(CGO_ENABLED) go build -o dist/$(BINARY) $(GO_BUILD_PACKAGE)

.PHONY: build-all
build-all: setup $(PLATFORM_TARGETS) ## Setup and build $(BINARY) for all PLATFORMS

dist/$(BINARY)-%: setup $(SRCS)
	@echo "Building $(BINARY) for $*..."
	STEM=$*
	OS=$${STEM%-*}
	ARCH=$${STEM#*-}
	GOOS=$$OS GOARCH=$$ARCH CGO_ENABLED=$(CGO_ENABLED) go build -o $@ $(GO_BUILD_PACKAGE)

.PHONY: install-default
install-default: build ## Install $(BINARY) binary to $PREFIX
	install -d $(PREFIX)/bin
	install -m 755 dist/$(BINARY) $(PREFIX)/bin/$(BINARY)


##
##@ Testing
##

.PHONY: test
test: test-units ## Run all tests

.PHONY: test-units
test-units: ## Run unit tests with full output (optional: TEST_GROUP=TestFoo)
	go test -v $(GO_PKGS) $(if $(TEST_GROUP),-run $(TEST_GROUP))

.PHONY: test-list
test-list: ## List available test groups (use as TEST_GROUP= value with make test)
	@go test -list '.*' $(GO_PKGS) | grep '^Test' | awk -F_ '{print $$1}' | sort -u

.PHONY: coverage
coverage: ## Run tests with coverage and show summary
	go test -coverprofile=coverage.out $(GO_PKGS)

.PHONY: coverage-html
coverage-html: coverage ## Open HTML coverage report
	go tool cover -html=coverage.out


##
##@ Linting and Formatting
##

.PHONY: tidy-default
tidy-default: ## Run go mod tidy
	go mod tidy

.PHONY: fmt
fmt: ## Re-format all go code
	go fmt $(GO_PKGS)

.PHONY: lint-default
lint-default: require-golangci-lint ## Run linter
	golangci-lint run $(GO_PKGS)

.PHONY: vet
vet: ## Run go vet
	go vet $(GO_PKGS)


##
##@ Setup / Cleanup
##

.PHONY: setup-default
setup-default: ## Create dist/ directory
	@mkdir -p dist

.PHONY: clean-default
clean-default: ## Clean up build artifacts
	@rm -rf dist/* coverage.out

#
# wildcard target for making sure default rules are used if no extension or overide is
# provided,
#
%: %-default
	@true


# ----------------------------------------------------------------------------------------
# Language agnostic Makefile settings and rules.

# EL7 ships with 3.82, so this is the current minimum.
# macOS ships with an ancient version of Make (3.81), so you will need to install a newer version via Homebrew.
# 3.81 is missing the .ONESHELL feature which is critical for simple target happiness.

MINIMUM_GNU_MAKE := 3.82
ifneq "$(MINIMUM_GNU_MAKE)" "$(firstword $(sort $(MINIMUM_GNU_MAKE) ${MAKE_VERSION}))"
  $(info ERROR: This Makefile requires GNU Make $(MINIMUM_GNU_MAKE) or greater.)
  $(info ERROR: You used $(MAKE_COMMAND), which is GNU Make $(MAKE_VERSION).)
  $(error Use gmake instead (macOS: brew install make))
endif

# Makefile settings
MAKEFLAGS += --no-builtin-rules --no-builtin-variables --no-print-directory --warn-undefined-variables
# if we want --warn-undefined-variables, we have to declare GNUMAKEFLAGS as empty to prevent a warning
GNUMAKEFLAGS ?=
SHELL := bash
.DEFAULT_GOAL := help
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

##
##@ Help
##
# I hate how complex (and long) this awk script is but I love what it does - so here it lives.
.PHONY: help
help: ## Display this help
	@awk -v prog='$(BINARY)' ' \
	BEGIN { \
		FS = ":.*##"; \
		printf "\nUsage:\n  make <target> \033[36m\033[0m\n"; \
	} \
	/^##@/ { \
		section = substr($$0, 5); \
		if (!(section in seen)) { \
			seen[section] = 1; \
			order[++count] = section; \
		} \
		next; \
	} \
	/^[a-zA-Z0-9_()$$%-]+:.*?##/ { \
		target = $$1; \
		desc = $$2; \
		gsub(/\$$\(BINARY\)/, prog, target); \
		gsub(/\$$\(BINARY\)/, prog, desc); \
		is_default = target ~ /-default$$/; \
		if (is_default) { \
			target = substr(target, 1, length(target) - 8); \
		} \
		if (section == "") { \
			section = "Other"; \
			if (!(section in seen)) { \
				seen[section] = 1; \
				order[++count] = section; \
			} \
		} \
		entry_section[++entry_count] = section; \
		entry_target[entry_count] = target; \
		entry_desc[entry_count] = desc; \
		entry_default[entry_count] = is_default; \
		if (!is_default) { \
			explicit[target] = 1; \
		} \
	} \
	END { \
		for (entry = 1; entry <= entry_count; entry++) { \
			target = entry_target[entry]; \
			if (entry_default[entry] && explicit[target]) { \
				continue; \
			} \
			if (displayed[target]) { \
				continue; \
			} \
			displayed[target] = 1; \
			section = entry_section[entry]; \
			items[section] = items[section] sprintf("  \033[36m%-15s\033[0m %s\n", target, entry_desc[entry]); \
		} \
		for (i = 1; i <= count; i++) { \
			section = order[i]; \
			if (items[section] != "") { \
				printf "\n\033[1m%s\033[0m\n%s", section, items[section]; \
			} \
		} \
	} \
	' $(MAKEFILE_LIST)


# Verify that a required external tool is installed.
# Usage: add require-<toolname> as a prerequisite to any target that needs it.
# If a tool is missing, it is likely available via Homebrew (macOS) or EPEL (Linux).
.PHONY: require-%
require-%:
	@which $* > /dev/null 2>&1 || { echo "Error: '$*' is not installed (try: brew install $* or epel-release)"; exit 1; }


.PHONY: mk-debug
mk-debug: ## Test that the Makefile is functional
	@echo "Makefile is functional."
	echo "Using Make version: ${MAKE_VERSION}"
