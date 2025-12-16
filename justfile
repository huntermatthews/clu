# Justfile for Python projects (manage common tasks)
# -*- mode: makefile -*-

# Justfile settings
_default: help

# Print the list of targets and their descriptions.
help:
    @just --list

# build it all
[group('build')]
build: clu

# Clean up build artifacts
[group('build')]
clean:
    rm -rf build/*

# Build the Go CLI tool
[group('build')]
@clu:
    mkdir -p build
    # go build -ldflags "-X github.com/huntermatthews/clu/pkg.Version=1.2.3" -o build/clu ./cmd/main.go
    GOOS=darwin GOARCH=arm64 go build -o build/clu-macos ./cmd/main.go
    GOOS=linux GOARCH=amd64 go build -o build/clu-linux-amd64 ./cmd/main.go

# @echo "This is a {{os()}} system running on {{arch()}} with {{num_cpus()}} logical CPUs."
