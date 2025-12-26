#!/usr/bin/env bash
set -euo pipefail

# Minimal Debian .deb builder for clu
# Stages binary, LICENSE (as copyright), and man page.
# Requires: go (to build), dpkg-deb
# Usage:
#   scripts/debbuild.sh                # derives Version via scripts/get-version.sh
#   VERSION=0.1.0 ARCH=amd64 scripts/debbuild.sh

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
CONTROL_FILE="$ROOT_DIR/DEBIAN/control"
NAME="clu"

if [[ ! -f "$CONTROL_FILE" ]]; then
  echo "Control file not found: $CONTROL_FILE" >&2
  exit 1
fi

# Derive Version via get-version.sh; allow override via env
# Call via bash to avoid execute-bit issues
RAW_VERSION="${VERSION:-$(bash "$ROOT_DIR"/scripts/get-version.sh)}"
if [[ -z "$RAW_VERSION" ]]; then
  echo "Failed to determine Version via scripts/get-version.sh. Set VERSION env var." >&2
  exit 1
fi
# Do not normalize; require a plain dotted numeric version (e.g., 0.7.8)
if ! [[ "$RAW_VERSION" =~ ^[0-9]+(\.[0-9]+)*$ ]]; then
  echo "Unacceptable version string from get-version.sh: '$RAW_VERSION'. Expect a numeric dotted version like '0.7.8'." >&2
  exit 1
fi
VERSION="$RAW_VERSION"
ARCH="${ARCH:-amd64}"

# Verify tools
for cmd in dpkg-deb; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "Missing required command: $cmd" >&2
    exit 1
  fi
done

# Build the binary if not present
if [[ ! -f "$ROOT_DIR/build/clu" ]]; then
  if ! command -v go >/dev/null 2>&1; then
    echo "Go toolchain not found and build/clu missing." >&2
    exit 1
  fi
  ( cd "$ROOT_DIR" && mkdir -p build && CGO_ENABLED=0 go build -o build/clu ./cmd )
fi

# Stage filesystem tree
DIST_DIR="$ROOT_DIR/dist/deb"
STAGE_DIR="$DIST_DIR/${NAME}_${VERSION}_${ARCH}"

mkdir -p "$STAGE_DIR/DEBIAN" \
         "$STAGE_DIR/usr/bin" \
         "$STAGE_DIR/usr/share/man/man1" \
         "$STAGE_DIR/usr/share/doc/$NAME"

# Install files
install -m 0755 "$ROOT_DIR/build/clu" "$STAGE_DIR/usr/bin/clu"
install -m 0644 "$ROOT_DIR/clu.1" "$STAGE_DIR/usr/share/man/man1/clu.1"
# Debian expects license under /usr/share/doc/<pkg>/copyright
install -m 0644 "$ROOT_DIR/LICENSE" "$STAGE_DIR/usr/share/doc/$NAME/copyright"

# Control metadata (update Version to derived value)
sed -E "s/^Version: .*/Version: $VERSION/" "$CONTROL_FILE" > "$STAGE_DIR/DEBIAN/control"
chmod 0644 "$STAGE_DIR/DEBIAN/control"

# Build the .deb
mkdir -p "$DIST_DIR"
OUTPUT_DEB="$DIST_DIR/${NAME}_${VERSION}_${ARCH}.deb"
dpkg-deb --build "$STAGE_DIR" "$OUTPUT_DEB"

echo "Built: $OUTPUT_DEB"
