#!/bin/bash

# Extract version from version.go
VERSION=$(grep 'var VERSION' version.go | awk '{ print $4 }' | tr -d '"')

# Define build function
build() {
    local os=$1
    local arch=$2
    local output_dir="release/${VERSION}/${os}/${arch}"
    local output_name="qai"

    echo "Building for OS: ${os}, ARCH: ${arch}"
    mkdir -p "${output_dir}"
    GOOS=${os} GOARCH=${arch} go build -ldflags="-s" -ldflags="-w" -o "${output_dir}/${output_name}"
    echo "Build complete: ${output_dir}/${output_name}"
}

# Build for Linux x64 and ARM
build linux amd64
build linux arm64

# Build for Windows x64 and ARM
build windows amd64
build windows arm64

# Build for macOS x64 and ARM
build darwin amd64
build darwin arm64

echo "All builds are complete."