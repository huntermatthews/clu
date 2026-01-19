#!/usr/bin/env bash

set -euo pipefail

# Script to build DEB package locally from source tarball
# Usage: ./scripts/build-deb.sh <version> [architecture]

if [[ ${#} -lt 1 || ${#} -gt 2 ]]; then
    echo "Error: Version argument is required, architecture is optional"
    echo "Usage: ${0} <version> [architecture]"
    echo "Example: ${0} 1.2.3"
    echo "Example: ${0} 1.2.3 arm64"
    exit 1
fi

VERSION="${1}"
ARCH="${2:-amd64}"  # Default to amd64 if not specified

# Convert architecture naming conventions
if [[ "${ARCH}" == "arm64" || "${ARCH}" == "aarch64" ]]; then
    DEB_ARCH="arm64"
    GOARCH="arm64"
elif [[ "${ARCH}" == "x86_64" || "${ARCH}" == "amd64" ]]; then
    DEB_ARCH="amd64"
    GOARCH="amd64"
else
    echo "Error: Unsupported architecture '${ARCH}'"
    echo "Supported: x86_64, amd64, aarch64, arm64"
    exit 1
fi

# Validate version format (RPM/DEB compatible)
if ! [[ "${VERSION}" =~ ^[0-9]+(\.[0-9]+)*(-[a-zA-Z0-9]+)*$ ]]; then
    echo "Error: Invalid version format '${VERSION}'"
    echo "Use compatible patterns: 1.2.3, 1.2.3-alpha, 1.2.3-rc1, etc."
    echo "Avoid: tildes (~), colons (:), or plus signs (+)"
    exit 1
fi

TARBALL_PATH="clu-${VERSION}.tar.gz"
ORIG_TARBALL="clu_${VERSION}.orig.tar.gz"

echo "Building DEB for clu version: ${VERSION}, architecture: ${DEB_ARCH}"

# Create source tarball if it doesn't exist
if [[ ! -f "${TARBALL_PATH}" ]]; then
    echo "Creating source tarball..."
    git archive --format=tar.gz --prefix=clu-${VERSION}/ HEAD > "${TARBALL_PATH}"
    echo "Created ${TARBALL_PATH}"
fi

# Create Debian orig tarball (required naming convention)
if [[ ! -f "${ORIG_TARBALL}" ]]; then
    echo "Creating Debian orig tarball..."
    cp "${TARBALL_PATH}" "${ORIG_TARBALL}"
fi

# Extract source
echo "Extracting source..."
rm -rf "clu-${VERSION}"
tar -xzf "${ORIG_TARBALL}"

# Copy debian packaging files
echo "Setting up Debian package files..."
if [[ ! -d "debian" ]]; then
    echo "Error: debian/ directory not found"
    echo "Please create debian/ directory with control files first"
    exit 1
fi
cp -r debian "clu-${VERSION}/"

# Update changelog with correct version
echo "Updating changelog version to ${VERSION}..."
sed -i "1s/([^)]*)/(${VERSION}-1)/" "clu-${VERSION}/debian/changelog"

# Build package
echo "Building DEB package..."

# Check for required build tools
if ! command -v dpkg-buildpackage >/dev/null 2>&1; then
    echo "Error: dpkg-buildpackage not found"
    echo "Install build dependencies with: sudo apt-get install dpkg-dev build-essential golang-go-1.20"
    exit 1
fi

# Check for Go compiler (flexible approach)
if ! command -v go >/dev/null 2>&1 && [[ ! -x /usr/lib/go-1.20/bin/go ]]; then
    echo "Error: Go compiler not found"
    echo "Install Go with one of these options:"
    echo "  Ubuntu 20.04: sudo apt-get install golang-1.20-go"
    echo "  Ubuntu 22.04+: sudo apt-get install golang-go"
    echo "  Or use official installer from https://golang.org/dl/"
    exit 1
fi

# Build the package (unsigned for local builds)
# Use -d flag to skip build dependency checks when cross-compiling
# Use dpkg-buildpackage directly to avoid debuild wrapper issues with architecture detection
export GOARCH="${GOARCH}"
export CGO_ENABLED=0
(cd "clu-${VERSION}" && dpkg-buildpackage -us -uc -ui -i -b -d -a"${DEB_ARCH}")

# Show results
echo ""
echo "DEB build complete!"
echo "Generated packages:"
ls -la clu_*.deb clu_*.changes 2>/dev/null || true
echo ""
echo "To install: sudo dpkg -i clu_${VERSION}-1_${DEB_ARCH}.deb"
echo "If dependencies missing: sudo apt-get install -f"

# DEBUG: Alternative Go 1.20 check method (package-based)
# if ! dpkg -l | grep -qE "^ii.*(golang-1\.20|golang-1\.20-go)"; then
#     echo "Error: golang-1.20 not found"
#     echo "Install Go 1.20 with: sudo apt-get install golang-1.20"
#     echo "Or install the compiler directly: sudo apt-get install golang-1.20-go"
#     echo "Note: golang-1.20 is available in Ubuntu 20.04 main repositories"
#     exit 1
# fi
