#!/bin/bash
set -euo pipefail

# Script to build RPM package locally from source tarball
# Usage: ./scripts/build-rpm.sh <version>

if [[ ${#} -ne 1 ]]; then
    echo "Error: Version argument is required"
    echo "Usage: ${0} <version>"
    echo "Example: ${0} 1.2.3"
    exit 1
fi

VERSION="${1}"
TARBALL_PATH="clu-${VERSION}.tar.gz"

echo "Building RPM for clu version: ${VERSION}"

# Create source tarball if it doesn't exist
if [[ ! -f "${TARBALL_PATH}" ]]; then
    echo "Creating source tarball..."
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

# Build RPM
echo "Building RPM from source..."
rpmbuild --define "_version ${VERSION}" -bb ~/rpmbuild/SPECS/clu.spec

# Show results
echo "Generated RPM:"
ls -la ~/rpmbuild/RPMS/x86_64/clu-*.rpm

echo "To install: sudo rpm -i ~/rpmbuild/RPMS/x86_64/clu-${VERSION}-1.*.rpm"
