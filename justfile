# Justfile for Python projects (manage common tasks)


# Justfile settings
_default: help

# Print the list of targets and their descriptions.
help:
    @just --list


# build the program as a single file executable (zipapp)
[group('build')]
build: clean
    mkdir -p dist
    uv run python -m zipapp src -m "clu.cli:main" -o dist/clu -p "/usr/bin/env python3"


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
mypy:
    uv run mypy src tests


# Run linter using ruff
[group('dev')]
lint *args:
    uv run ruff check {{args}}
