# Justfile for Python projects (manage common tasks)


# Justfile settings
_default: help


# Clean up build artifacts
[group('build')]
clean:
    rm -rf build dist
    find src -type d -name '*.egg-info' -exec rm -r {} + -depth
    find src -type f -name '*.py[co]' -delete -o -type d -name __pycache__ -delete


# Print the list of targets and their descriptions.
help:
    @just --list


# Run tests using pytest
[group('test')]
test *args:
    uv run pytest tests {{args}}


# Type check the code using mypy
[group('test')]
mypy:
    uv run mypy src tests


# build the program as a single file executable (zipapp)
[group('build')]
build:
    mkdir -p dist
    uv run python -m zipapp src -m "clu.cli:main" -o dist/clu -p "/usr/bin/env python3"
