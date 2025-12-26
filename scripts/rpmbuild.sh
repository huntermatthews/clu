#!/usr/bin/env bash
set -euo pipefail

# Minimal rpmbuild wrapper for clu
# - Packages required files into a Source0 tarball: cmd/, README.md, LICENSE, clu.1
# - Builds RPM using redhat/clu.spec
# Usage:
#   scripts/rpmbuild.sh            # derives Version via scripts/get-version.sh
#   VERSION=0.1.0 scripts/rpmbuild.sh  # override version

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
SPEC_FILE="$ROOT_DIR/redhat/clu.spec"
NAME="clu"

if [[ ! -f "$SPEC_FILE" ]]; then
  echo "Spec file not found: $SPEC_FILE" >&2
  exit 1
fi

# Derive Version via get-version.sh; allow override via env
# Call via bash to avoid execute-bit issues
RAW_VERSION="${VERSION:-$(bash "$ROOT_DIR"/scripts/get-version.sh)}"
if [[ -z "$RAW_VERSION" ]]; then
  echo "Failed to determine version via scripts/get-version.sh. Set VERSION env var." >&2
  exit 1
fi
# Do not normalize; require a plain dotted numeric version (e.g., 0.7.8)
if ! [[ "$RAW_VERSION" =~ ^[0-9]+(\.[0-9]+)*$ ]]; then
  echo "Unacceptable version string from get-version.sh: '$RAW_VERSION'. Expect a numeric dotted version like '0.7.8'." >&2
  exit 1
fi
VERSION="$RAW_VERSION"

if ! command -v rpmbuild >/dev/null 2>&1; then
  echo "rpmbuild is not installed. Please install rpm-build." >&2
  exit 1
fi

TOPDIR="$ROOT_DIR/rpmbuild"
SOURCES="$TOPDIR/SOURCES"
SPECS="$TOPDIR/SPECS"
BUILD="$TOPDIR/BUILD"
RPMS="$TOPDIR/RPMS"
SRPMS="$TOPDIR/SRPMs"
TMPDIR="$TOPDIR/tmp"
STAGE="$TMPDIR/$NAME-$VERSION"

mkdir -p "$SOURCES" "$SPECS" "$BUILD" "$RPMS" "$SRPMS" "$STAGE"

# Verify required sources exist
for f in README.md LICENSE clu.1; do
  if [[ ! -f "$ROOT_DIR/$f" ]]; then
    echo "Missing required file: $ROOT_DIR/$f" >&2
    exit 1
  fi
done
if [[ ! -d "$ROOT_DIR/cmd" ]]; then
  echo "Missing required directory: $ROOT_DIR/cmd" >&2
  exit 1
fi

# Stage files under clu-<version>/
cp -a "$ROOT_DIR/cmd" "$STAGE/"
cp -a "$ROOT_DIR/README.md" "$ROOT_DIR/LICENSE" "$ROOT_DIR/clu.1" "$STAGE/"
# Include Go module files if present
for f in go.mod go.sum; do
  [[ -f "$ROOT_DIR/$f" ]] && cp -a "$ROOT_DIR/$f" "$STAGE/"
done

# Create Source0 tarball
TARBALL="$SOURCES/$NAME-$VERSION.tar.gz"
( cd "$TMPDIR" && tar -czf "$TARBALL" "$NAME-$VERSION" )

# Copy spec and update Version to the derived value
sed -E "s/^(Version:[[:space:]]*).*/\1$VERSION/" "$SPEC_FILE" > "$SPECS/$(basename "$SPEC_FILE")"

# Build RPM
rpmbuild -ba "$SPECS/$(basename "$SPEC_FILE")" --define "_topdir $TOPDIR"

# Output built RPM paths
echo "\nBuilt RPMs:"
find "$RPMS" -type f -name "*.rpm" -print
