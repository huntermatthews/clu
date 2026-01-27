include common.mk

# Variables
PREFIX ?= /usr/local
PKG := $(shell go list -m)
PLATFORMS := darwin-arm64 windows-amd64 linux-amd64 linux-arm64


##
##@ Build
##

.PHONY: build
build: setup ## Build the program for the current platform
	CGO_ENABLED=0 go build -ldflags "-X $(PKG)/pkg/global.Version=$(VERSION)" -o dist/clu ./cmd/clu

.PHONY: clean
clean: ## Clean up build artifacts
	@rm -rf dist/* coverage.out .go-md2man-installed

.PHONY: clu-all
clu-all: ## Build the Go CLI tool
	@for plat in $(PLATFORMS); do
		$(MAKE) clu-$$plat
	done

.PHONY: clu-%
clu-%: ## Build clu for a specific platform (e.g., make clu-darwin-amd64)
	@GOOS=$(word 1,$(subst -, ,$*))
	GOARCH=$(word 2,$(subst -, ,$*))
	GOOS=$$GOOS GOARCH=$$GOARCH $(MAKE) build
	mv dist/clu dist/clu-$$GOOS-$$GOARCH


.PHONY: install
install: build ## Install clu binary, manpage, and documentation
	install -d $(PREFIX)/bin
	install -d $(PREFIX)/share/man/man1
	install -d $(PREFIX)/share/doc/clu
	install -m 755 dist/clu $(PREFIX)/bin/clu
	install -m 644 clu.1 $(PREFIX)/share/man/man1/clu.1
	install -m 644 README.md $(PREFIX)/share/doc/clu/README.md


.PHONY: manpage
manpage: .go-md2man-installed ## Generate man page from markdown using go-md2man
	go-md2man -in clu.1.md -out clu.1

.go-md2man-installed:
	go install github.com/cpuguy83/go-md2man/v2@latest
	@touch .go-md2man-installed

.PHONY: version
version: ## Print the current version
	@echo $(VERSION)


##
##@ Testing
##

.PHONY: test
test: ## Run Go tests
	go test ./pkg/... ./cmd/...

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
