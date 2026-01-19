#!/usr/bin/env bash

set -euo pipefail

# Script to build RPM package locally from source tarball
# Usage: ./scripts/build-rpm.sh <version> [architecture]

if [[ ${#} -lt 1 || ${#} -gt 2 ]]; then
    echo "Error: Version argument is required, architecture is optional"
    echo "Usage: ${0} <version> [architecture]"
    echo "Example: ${0} 1.2.3"
    echo "Example: ${0} 1.2.3 aarch64"
    exit 1
fi

VERSION="${1}"
ARCH="${2:-x86_64}"  # Default to x86_64 if not specified

# Convert architecture naming conventions
if [[ "${ARCH}" == "arm64" || "${ARCH}" == "aarch64" ]]; then
    RPM_ARCH="aarch64"
    GOARCH="arm64"
elif [[ "${ARCH}" == "x86_64" || "${ARCH}" == "amd64" ]]; then
    RPM_ARCH="x86_64"
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

echo "Building RPM for clu version: ${VERSION}, architecture: ${RPM_ARCH}"

# Create source tarball if it doesn't exist
if [[ ! -f "${TARBALL_PATH}" ]]; then
    echo "Creating source tarball from git..."
    # Use same approach as GitHub archives (common pattern)
    git archive --format=tar.gz --prefix=clu-${VERSION}/ HEAD > "${TARBALL_PATH}"
    echo "Created ${TARBALL_PATH}"
fi

# Set up build tree
echo "Setting up RPM build tree..."
if command -v rpmdev-setuptree >/dev/null 2>&1; then
    rpmdev-setuptree
else
    # Manual setup if rpmdevtools not available
    mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
    echo "Created RPM build directories manually (rpmdevtools not found)"
fi

# Copy tarball and spec
echo "Copying source tarball and spec file..."
cp "${TARBALL_PATH}" ~/rpmbuild/SOURCES/
cp redhat/clu.spec ~/rpmbuild/SPECS/

# Check for Go compiler (flexible approach)
if ! command -v go >/dev/null 2>&1 && [[ ! -x /usr/lib/go-1.20/bin/go ]]; then
    echo "Error: Go compiler not found"
    echo "Install Go with one of these options:"
    echo "  Ubuntu 20.04: sudo apt-get install golang-1.20-go"
    echo "  Ubuntu 22.04+: sudo apt-get install golang-go"
    echo "  Or use official installer from https://golang.org/dl/"
    exit 1
fi

# Build RPM
echo "Building RPM from source..."
export GOARCH="${GOARCH}"
export CGO_ENABLED=0
rpmbuild --define "_version ${VERSION}" \
         --define "_buildhost $(hostname -f)" \
         --target "${RPM_ARCH}" \
         -bb ~/rpmbuild/SPECS/clu.spec

# Show results
echo "Generated RPM:"
ls -la ~/rpmbuild/RPMS/${RPM_ARCH}/clu-*.rpm

echo "To install: sudo rpm -i ~/rpmbuild/RPMS/${RPM_ARCH}/clu-${VERSION}-1.*.rpm"
