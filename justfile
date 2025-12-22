# Justfile for Python projects (manage common tasks)

# Justfile settings
_default: help

# Print the list of targets and their descriptions.
help:
    @just --list

# build it all
[group('build')]
build: clu

# Clean up build artifacts
[group('build')]
clean:
    @rm -rf build/*

# Build the Go CLI tool
[group('build')]
clu: clean
    #!/bin/zsh
    CLU_VERSION=$(git describe --dirty --always --match "v[0-9]*")
    mkdir -p build
    platforms=("darwin:arm64" "linux:amd64")

    for plat in "${platforms[@]}"; do
        os=${plat%:*}
        arch=${plat#*:}
        echo "Building for $os on $arch..."
        GOOS=$os GOARCH=$arch go build -ldflags "-X github.com/huntermatthews/clu/pkg.Version=${CLU_VERSION}" -o build/clu-$os-$arch ./cmd/main.go
    done

    LOCAL_OS=$(go env GOOS)
    LOCAL_ARCH=$(go env GOARCH)
    echo "Creating symlink for local build: $LOCAL_OS on $LOCAL_ARCH..."
    ln -s clu-${LOCAL_OS}-${LOCAL_ARCH}  build/clu

[group('info')]
report-build-info:
    @echo "This is a {{os()}} system running on {{arch()}} with {{num_cpus()}} logical CPUs."
