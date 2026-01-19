# Justfile for Python projects (manage common tasks)

# Get the current version from git
REPO_VERSION := `scripts/get-version.sh`

# Required programs
go-md2man := require("go-md2man")

# Justfile settings
_default: help

# Print the list of targets and their descriptions.
help:
    @just --list

# Build the program and all its varients.
[group('build')]
build: clu clu-symlink
    @mkdir -p dist

# Clean up build artifacts
[group('build')]
clean:
    @rm -rf dist/* coverage.out

# Run Go tests
[group('build')]
test:
    go test ./pkg/... ./cmd/...

# Run tests with coverage and show summary
[group('build')]
coverage:
    go test -coverpkg=./pkg/...,./cmd/... -coverprofile=coverage.out ./pkg/... ./cmd/...
    go tool cover -func=coverage.out

# Open HTML coverage report
[group('build')]
coverage-html: coverage
    go tool cover -html=coverage.out

# Build the Go CLI tool
[group('build')]
clu:
    #!/usr/bin/env bash

    set -euo pipefail

    platforms=("darwin:arm64")

    for plat in "${platforms[@]}"; do
        export GOOS=${plat%:*}
        export GOARCH=${plat#*:}
        echo "Building for $GOOS on $GOARCH..."
        go build -ldflags "-X github.com/huntermatthews/clu/pkg/global.Version={{REPO_VERSION}}" -o ./dist/clu-$(go env GOOS)-$(go env GOARCH) ./cmd
    done

# Build the Go CLI tool
[group('build')]
clu-symlink:
    rm -f dist/clu
    ln -s clu-$(go env GOOS)-$(go env GOARCH)  dist/clu

# Generate man page from markdown using go-md2man
[group('build')]
manpage:
    {{go-md2man}} -in clu.1.md -out clu.1

# Show the current version from git
[group('package')]
version:
    echo "Current repo version is: {{REPO_VERSION}}"
