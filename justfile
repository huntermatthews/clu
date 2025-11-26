# Justfile for Python projects (manage common tasks)
# -*- mode: makefile -*-

# Justfile settings
_default: help

# Print the list of targets and their descriptions.
help:
    @just --list


# build it all
[group('build')]
build: zipapp


# Create the required build / dist dirs
[group('build')]
_dirs:
    mkdir -p build dist


# Build a proper python package (wheel)
[group('build')]
wheel: clean _dirs
    uv build --wheel


# build the program as a single file executable (zipapp)
[group('build')]
zipapp: wheel
    uv pip install --target=build dist/clu-*-any.whl
    uv run python -m zipapp build -m "clu.cli:main" -o dist/clu -p "/usr/bin/env python3" -c


# Clean up build artifacts
[group('build')]
clean:
    rm -rf build dist
    find src -type d -name '*.egg-info' -exec rm -r {} + -depth
    find src -type f -name '*.py[co]' -delete -o -type d -name __pycache__ -delete


# Run tests using pytest
[group('dev')]
test *args:
    uv run pytest tests {{args}}


# Type check the code using mypy
[group('dev')]
types:
    uv run mypy src tests


# Format the code (--diff only for now)
[group('dev')]
format:
    uv run ruff format --diff src tests


# Format the code and write it
[group('dev')]
reformat:
    uv run ruff format src tests


# Run linter using ruff
[group('dev')]
lint *args:
    uv run ruff check {{args}}

# Synchronize dependencies (update the clu version string)
[group('dev')]
sync:
    uv sync --reinstall-package clu

# [group('dev')]
# tag *args:
#     git tag --annotate --sign {{args}}


###
# go build -ldflags \"-X github.com/huntermatthews/clu/pkg.Version=1.2.3\"
