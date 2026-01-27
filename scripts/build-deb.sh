#!/usr/bin/env bash

set -euo pipefail

# Usage: ./scripts/build-deb.sh <version> [architecture]

if [[ ${#} -lt 1 || ${#} -gt 2 ]]; then
    echo "Error: Version argument is required, architecture is optional"
    echo "Usage: ${0} <version> [architecture]"
    echo "Example: ${0} 1.2.3"
    echo "Example: ${0} 1.2.3 arm64"
    exit 1
fi

export VERSION="${1}"
ARCH="${2:-amd64}"  # Default to amd64 if not specified

# Map architecture names (matching build-deb.yaml workflow)
case "${ARCH}" in
    amd64|x86_64) DEB_ARCH="amd64"; GOARCH="amd64" ;;
    arm64|aarch64) DEB_ARCH="arm64"; GOARCH="arm64" ;;
    *) echo "Error: Unsupported architecture '${ARCH}'"; exit 1 ;;
esac

# Validate version format (RPM/DEB compatible)
if ! [[ "${VERSION}" =~ ^[0-9]+(\.[0-9]+)*(-[a-zA-Z0-9]+)*$ ]]; then
    echo "Error: Invalid version format '${VERSION}'"
    echo "Use compatible patterns: 1.2.3, 1.2.3-alpha, 1.2.3-rc1, etc."
    echo "Avoid: tildes (~), colons (:), or plus signs (+)"
    exit 1
fi

echo "Building DEB for clu version: ${VERSION}, architecture: ${DEB_ARCH}"

# Verify debian directory exists
if [[ ! -d "debian" ]]; then
    echo "Error: debian/ directory not found"
    echo "Please run this script from the repository root"
    exit 1
fi

# Update changelog with correct version
echo "Updating changelog version to ${VERSION}..."
sed -i "1s/([^)]*)/(${VERSION}-1)/" debian/changelog

# Check for required build tools
if ! command -v dpkg-buildpackage >/dev/null 2>&1; then
    echo "Error: dpkg-buildpackage not found"
    echo "Install build dependencies with: sudo apt-get install dpkg-dev build-essential"
    exit 1
fi

# Check for Go compiler
if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go compiler not found"
    echo "Install Go with one of these options:"
    echo "  Ubuntu 20.04/22.04: sudo apt-get install golang-1.22-go"
    echo "  Ubuntu 24.04+: sudo apt-get install golang-go"
    exit 1
fi

# Build package
echo "Building DEB package..."

# Export GOARCH for cross-compilation
export GOARCH

# Build the package (unsigned for local builds, binary-only)
# Use -d flag to skip build dependency checks when cross-compiling
# Use -b for binary-only build (no source package needed)
dpkg-buildpackage -us -uc -ui -i -b -d -a"${DEB_ARCH}"

# Show results
echo ""
echo "DEB build complete!"
echo "Generated packages:"
ls -la ../clu_*.deb ../clu_*.changes 2>/dev/null || true
echo ""
echo "To install: sudo dpkg -i ../clu_${VERSION}-1_${DEB_ARCH}.deb"
