# Justfile for Python projects (manage common tasks)

# Justfile settings
_default: help

# Print the list of targets and their descriptions.
help:
    @just --list

# build it all
[group('build')]
build: clu clu-symlink
    @mkdir -p dist

# Clean up build artifacts
[group('build')]
clean:
    @rm -rf dist/*

# Run Go tests
[group('build')]
test:
    go test ./pkg/... ./cmd/...

# Build the Go CLI tool
[group('build')]
clu:
    #!/usr/bin/env zsh

    platforms=("darwin:arm64")# "linux:amd64")

    for plat in "${platforms[@]}"; do
        export GOOS=${plat%:*}
        export GOARCH=${plat#*:}
        echo "Building for $GOOS on $GOARCH..."
        scripts/build.sh
    done

# Build the Go CLI tool
[group('build')]
build-manylinux:
    # biowebdev05 is our build server for ancient Linux compatibility
    ssh biowebdev05 "cd clu && git pull && ./scripts/build.sh"
    scp biowebdev05:clu/dist/clu-linux-amd64 dist/clu-manylinux2014
    scp dist/clu-manylinux2014 itbrepo02.nhgri.nih.gov:/srv/webroot/matthewsht/clu-manylinux2014
    echo 'curl -fsSL -o clu https://itbrepo02.nhgri.nih.gov/matthewsht/clu-manylinux2014 && chmod +x clu'

# Build the Go CLI tool
[group('build')]
clu-symlink:
    rm -f dist/clu
    ln -s clu-$(go env GOOS)-$(go env GOARCH)  dist/clu

[group('info')]
report-build-info:
    @echo "This is a {{os()}} system running on {{arch()}} with {{num_cpus()}} logical CPUs."
