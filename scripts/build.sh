#!/usr/bin/env bash

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

mkdir -p dist

REPO_VERSION=$(./scripts/get-version.sh)

# Define output names to match GitHub workflow
ARTIFACT_NAME="clu-$(go env GOOS)-$(go env GOARCH)"
ARTIFACT_DIR="dist/${ARTIFACT_NAME}"
mkdir -p "${ARTIFACT_DIR}"

CGO_ENABLED=0 go build -ldflags "-X github.com/huntermatthews/clu/pkg/global.Version=${REPO_VERSION}" -o "${ARTIFACT_DIR}/clu" ./cmd/main.go
