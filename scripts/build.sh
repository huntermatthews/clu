#!/usr/bin/env zsh

# This silly script is so we don't need to install `just` to build the project
# Sadly its not available on all CI systems by default or our ancient C7/Ubuntu systems.

cd "$(git rev-parse --show-toplevel)"

mkdir -p dist

REPO_VERSION=$(git describe --dirty --always --match "v[0-9]*")

go build -ldflags "-X github.com/huntermatthews/clu/pkg.Version=${REPO_VERSION}" -o ./dist/clu-$(go env GOOS)-$(go env GOARCH) ./cmd/main.go
