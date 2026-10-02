# go-common.mk
# Shared targets for building a single Go program.
#
# Tunables (set in the sub-Makefile BEFORE the include to customize):
#   PREFIX           install prefix
#   CGO_ENABLED      cgo flag
#   PLATFORMS        cross-build target list
#   SRCS             source deps for the binary
#   TEST_UNITS_PKGS  packages passed to test-units
#   CLEAN_EXTRA      extra paths for `clean` to rm



# 'make install PREFIX=/usr/local' for a system-wide install
# This default ASSUMES that you have XDG Base Directory Specification compliant directories set up in your environment.
# (or at least something similar to ~/.local/bin and ~/.local/share/man/man1)
PREFIX            ?= $(HOME)/.local
CGO_ENABLED       ?= 0
PLATFORMS         ?= darwin-arm64 linux-amd64 linux-arm64 windows-amd64
GO_BUILD_PACKAGE  ?= ./cmd/$(BINARY)
SRCS              ?= $(wildcard *.go) go.mod $(wildcard go.sum)
TEST_UNITS_PKGS   ?= ./...
TEST_GROUP        ?=


# Non Tunables -- Do not override these in the sub-Makefile
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
	go test -v $(TEST_UNITS_PKGS) $(if $(TEST_GROUP),-run $(TEST_GROUP))

.PHONY: test-list
test-list: ## List available test groups (use as TEST_GROUP= value with make test)
	@go test -list '.*' ./... | grep '^Test' | awk -F_ '{print $$1}' | sort -u

.PHONY: coverage
coverage: ## Run tests with coverage and show summary
go test -coverprofile=coverage.out ./...  ${all_packages}

.PHONY: coverage-html
coverage-html: coverage ## Open HTML coverage report
	go tool cover -html=coverage.out


##
##@ Linting and Formatting
##

.PHONY: tidy-default
tidy-default: ## Run go mod tidy
	go mod tidy

.PHONY: fmt-check
fmt-check: ## Check formatting without modifying files
	@gofmt -d $(shell find . -name '*.go' -not -path "./vendor/*")


.PHONY: fmt
fmt: ## Format code
	go fmt $(shell find . -name '*.go' -not -path "./vendor/*") 
	
.PHONY: lint-default
lint-default: require-golangci-lint ## Run linter
	golangci-lint run

.PHONY: vet
vet: ## Run go vet
	go vet ./...


##
##@ Setup / Cleanup
##

.PHONY: setup-default
setup-default: ## Create dist/ directory
	@mkdir -p dist
	go mod download

.PHONY: clean-default
clean-default: ## Clean up build artifacts
	@rm -rf dist/* coverage.out

#
# wildcard target for making sure default rules are used if no extension or overide is
# provided,
#
%: %-default
	@true
