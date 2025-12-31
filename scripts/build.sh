#!/usr/bin/env zsh

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

mkdir -p dist

REPO_VERSION=$(git describe --dirty --always --match "v[0-9]*")

# Define output names to match GitHub workflow
ARTIFACT_NAME="clu-$(go env GOOS)-$(go env GOARCH)"
ARTIFACT_DIR="dist/${ARTIFACT_NAME}"
mkdir -p "${ARTIFACT_DIR}"

# Build stripped production binary
CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/huntermatthews/clu/pkg.Version=${REPO_VERSION}" -o "${ARTIFACT_DIR}/clu" ./cmd/main.go
# Build binary with debug symbols
CGO_ENABLED=0 go build -ldflags "-X github.com/huntermatthews/clu/pkg.Version=${REPO_VERSION}" -o "${ARTIFACT_DIR}/clu-debug" ./cmd/main.go
