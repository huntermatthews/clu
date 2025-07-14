
_default: help

# Clean up build / dist artifacts
clean:
  rm -rf build dist

# Print out the list of valid targets
help:
  @just --list

# zip the program as a single file executable (zipapp)
zip:
  mkdir -p dist
  uv run python -m zipapp src -m "clu.cli:main" -o dist/clu -p "/usr/bin/env python3"
