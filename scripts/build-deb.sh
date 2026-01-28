#!/usr/bin/env bash

set -euo pipefail

# Usage: ./scripts/build-deb.sh <version> [<amd64|arm64>]

if [[ ${#} -lt 1 || ${#} -gt 2 ]]; then
    echo "Error: Version argument is required, architecture is optional"
    echo "Usage: ${0} <version> [<amd64|arm64>]"
    echo "Example: ${0} 1.2.3"
    echo "Example: ${0} 1.2.3 arm64"
    exit 1
fi

export VERSION="${1}"
GOARCH="${2:-amd64}"  # Default to amd64 if not specified
DEB_ARCH="${GOARCH}"   # we don't need a mapping here - deb/ubuntu use amd64/arm64

# Validate version format (RPM/DEB compatible)
if ! [[ "${VERSION}" =~ ^[0-9]+(\.[0-9]+)*(-[a-zA-Z0-9]+)*$ ]]; then
    echo "Error: Invalid version format '${VERSION}'"
    echo "Use compatible patterns: 1.2.3, 1.2.3-alpha, 1.2.3-rc1, etc."
    echo "Avoid: tildes (~), colons (:), or plus signs (+)"
    exit 1
fi

echo "Building DEB for clu version: ${VERSION}, architecture: ${DEB_ARCH}"

# Update changelog with correct version
echo "Updating changelog version to ${VERSION}..."
sed -i "1s/([^)]*)/(${VERSION}-1)/" debian/changelog

# Check for required build tools
if ! command -v dpkg-buildpackage >/dev/null 2>&1; then
    echo "Error: dpkg-buildpackage not found"
    echo "Install build dependencies with: sudo apt-get install dpkg-dev build-essential debhelper"
    exit 1
fi

# Check for Go compiler
if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go compiler not found"
    echo "Install Go with one of these options:"
    echo "  Ubuntu 22.04: sudo apt-get install golang-1.22-go"
    echo "  Ubuntu 24.04+: sudo apt-get install golang-go"
    exit 1
fi

# Build package
echo "Building DEB package..."

# Export GOARCH for cross-compilation
export GOARCH

# Build the package (unsigned for local builds, binary-only)
dpkg-buildpackage \
    --unsigned-source \
    --unsigned-changes \
    --unsigned-buildinfo \
    --diff-ignore=.* \
    --build=binary \
    --no-check-builddeps \
    --host-arch="${DEB_ARCH}"

# Move the .deb file to the top of the repo and delete other build artifacts
# I hate that we have to do this, but dpkg-buildpackage insists on putting them
# in the parent directory of the source tree.
REPO_TOP="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
DEB_FILE="../clu_${VERSION}-1_${DEB_ARCH}.deb"
if [ -f "$DEB_FILE" ]; then
    mv "$DEB_FILE" "$REPO_TOP/"
fi
# Remove other build artifacts
rm -f ../clu_*_${DEB_ARCH}.changes ../clu_*_${DEB_ARCH}.buildinfo ../clu_*_${DEB_ARCH}.dsc 2>/dev/null || true

# Show results
echo ""
echo "DEB build complete!"
echo "Generated package:"
ls -la "$REPO_TOP/clu_${VERSION}-1_${DEB_ARCH}.deb" 2>/dev/null || true
echo ""
echo "To install: sudo dpkg -i $REPO_TOP/clu_${VERSION}-1_${DEB_ARCH}.deb"
