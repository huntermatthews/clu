include common.mk

# Variables
PREFIX ?= /usr/local
BINARY := clu
GO_SOURCES := $(shell find cmd pkg -name '*.go') go.mod go.sum Makefile common.mk
PLATFORMS := darwin-arm64 windows-amd64 linux-amd64 linux-arm64


##
##@ Build
##

.PHONY: build
build: setup dist/$(BINARY) ## setup and build $(BINARY)

dist/$(BINARY): $(GO_SOURCES) ## Build $(BINARY) for the current platform
	CGO_ENABLED=0 go build -o dist/$(BINARY) ./cmd/clu

.PHONY: clean
clean: ## Clean up build artifacts
	@rm -rf dist/* coverage.out .go-md2man-installed

.PHONY: all
all: setup ## setup and build for all platforms
	@for plat in $(PLATFORMS); do
		export GOOS=$${plat%-*}
		export GOARCH=$${plat#*-}
		echo "Building $(BINARY) for $$GOOS-$$GOARCH..."
		$(MAKE) dist/$(BINARY)
		mv dist/$(BINARY) dist/$(BINARY)-$$GOOS-$$GOARCH
	done

.PHONY: install
install: build manpage ## Install $(BINARY) binary, manpage, and documentation
	install -d $(PREFIX)/bin
	install -d $(PREFIX)/share/man/man1
	install -d $(PREFIX)/share/doc/$(BINARY)
	install -m 755 dist/$(BINARY) $(PREFIX)/bin/$(BINARY)
	install -m 644 $(BINARY).1 $(PREFIX)/share/man/man1/$(BINARY).1
	install -m 644 README.md $(PREFIX)/share/doc/$(BINARY)/README.md


##
##@ Documentation
##

.PHONY: man
man: .go-md2man-installed setup ## Generate man page from markdown using go-md2man
	go-md2man -in $(BINARY).1.md -out $(BINARY).1
	# embed is limited to same or sub-pkgs
	cp $(BINARY).1 pkg/subcmd/$(BINARY).1

.go-md2man-installed:
	go install github.com/cpuguy83/go-md2man/v2@latest
	@touch .go-md2man-installed


##
##@ Testing
##

TEST_GROUP ?=

.PHONY: test
test: test-units ## Run all tests

.PHONY: test-units
test-units: ## Run unit tests with full output (optional: TEST_GROUP=TestFoo)
	go test -v ./pkg/... ./cmd/... $(if $(TEST_GROUP),-run $(TEST_GROUP))

.PHONY: test-list
test-list: ## List available test groups (use as TEST_GROUP= value with make test)
	@go test -list '.*' ./... | grep '^Test' | awk -F_ '{print $$1}' | sort -u

.PHONY: coverage
coverage: ## Run tests with coverage and show summary
	go test -coverpkg=./pkg/...,./cmd/... -coverprofile=coverage.out ./pkg/... ./cmd/...
	go tool cover -func=coverage.out

.PHONY: coverage-html
coverage-html: coverage ## Open HTML coverage report
	go tool cover -html=coverage.out


##
##@ Linting and Formatting
##

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: fmt-check
fmt-check: ## Check formatting without modifying files
	@gofmt -d $(shell find . -name '*.go' -not -path "./vendor/*")

# .PHONY: lint
# lint: ## Run linter
# 	golangci-lint run

.PHONY: vet
vet: ## Run go vet
	go vet ./...


##
##@ Setup and Tools
##

.PHONY: setup
setup: ## Setup Go modules
	@mkdir -p dist
	go mod download
