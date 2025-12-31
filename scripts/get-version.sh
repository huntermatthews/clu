#!/usr/bin/env bash

# This script outputs the current version of the repo based on git tags.
# for more complicated versioning schemes, modify this script - IE, its a monorepo with multiple go modules
# Having it as an external script allows just and github actions to both use the same logic.

set -euo pipefail

git describe --dirty --always --match "v[0-9]*"
