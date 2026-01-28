#!/usr/bin/env bash

set -euo pipefail

# Script to build RPM package directly from checkout

if [[ ${#} -lt 1 || ${#} -gt 2 ]]; then
    echo "Error: Version argument is required, architecture is optional"
    echo "Usage: ${0} <version> [<amd64|arm64>]"
    echo "Example: ${0} 1.2.3"
    echo "Example: ${0} 1.2.3 arm64"
    exit 1
fi

export VERSION="${1}"
ARCH="${2:-amd64}"  # Default to amd64 if not specified

# Map architecture names
case "${ARCH}" in
    amd64) RPM_ARCH="x86_64"; GOARCH="amd64" ;;
    arm64) RPM_ARCH="aarch64"; GOARCH="arm64" ;;
    *) echo "Error: Unsupported architecture '${ARCH}'"; exit 1 ;;
esac

# Validate version format (RPM/DEB compatible)
if ! [[ "${VERSION}" =~ ^[0-9]+(\.[0-9]+)*(-[a-zA-Z0-9]+)*$ ]]; then
    echo "Error: Invalid version format '${VERSION}'"
    echo "Use compatible patterns: 1.2.3, 1.2.3-alpha, 1.2.3-rc1, etc."
    echo "Avoid: tildes (~), colons (:), or plus signs (+)"
    exit 1
fi

echo "Building RPM for clu version: ${VERSION}, architecture: ${RPM_ARCH}"

# Set up build tree
echo "Setting up RPM build tree..."
mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

# Symlink source to BUILD directory (no tarball, no copy)
echo "Symlinking source to build directory..."
rm -rf ~/rpmbuild/BUILD/clu-${VERSION}
ln -s "$PWD" ~/rpmbuild/BUILD/clu-${VERSION}

# Copy spec file to SPECS directory
echo "Copying spec file..."
cp redhat/clu.spec ~/rpmbuild/SPECS/

# Check for Go compiler
if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go compiler not found"
    echo "Install Go with one of these options:"
    echo "  Ubuntu 22.04+: sudo apt-get install golang-go"
    echo "  Or use official installer from https://golang.org/dl/"
    exit 1
fi

# Export GOARCH for cross-compilation
export GOARCH

# Build RPM
echo "Building RPM..."
rpmbuild --define "_version ${VERSION}" \
         --define "_buildhost $(hostname -f)" \
         --target "${RPM_ARCH}" \
         -bb ~/rpmbuild/SPECS/clu.spec

# Move the final binary RPM to the top of the repo and delete other RPMs
REPO_TOP="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
RPM_FILE="$HOME/rpmbuild/RPMS/${RPM_ARCH}/clu-${VERSION}-1.${RPM_ARCH}.rpm"
if [ -f "$RPM_FILE" ]; then
    mv "$RPM_FILE" "$REPO_TOP/"
fi


# Show results
echo ""
echo "RPM build complete!"
echo "Generated RPM:"
ls -la "$REPO_TOP/clu-${VERSION}-1.${RPM_ARCH}.rpm" 2>/dev/null || true
echo ""
echo "To install: sudo rpm -i $REPO_TOP/clu-${VERSION}-1.${RPM_ARCH}.rpm"
