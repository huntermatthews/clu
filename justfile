# Justfile for Python projects (manage common tasks)


# Justfile settings
_default: help

# Print the list of targets and their descriptions.
help:
    @just --list


# build it all
[group('build')]
build: clean _dirs wheel zipapp

# Create the build / dist dirs
[group('build')]
_dirs:
    mkdir build dist

# Build a proper python package (wheel)
[group('build')]
wheel: _dirs
    uv build --wheel

# build the program as a single file executable (zipapp)
[group('build')]
zipapp: _dirs wheel
    uv pip install --target=build dist/clu-*-any.whl
    uv run python -m zipapp build -m "clu.cli:main" -o dist/clu -p "/usr/bin/env python3"

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

# uv build
# uv pip install dist/clu-0.1.0-py3-none-any.whl --target=build/
# uv run python -m zipapp build -m "clu.cli:main" -o dist/clu -p "/usr/bin/env python3"
