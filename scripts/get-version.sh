#! /usr/bin/env bash

# Get the current version of the repo using git describe
# common script to be used by other build scripts
git describe --dirty --always --match "v[0-9]*"
