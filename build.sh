#!/bin/bash

# Extract version from version.go
VERSION=$(grep 'var VERSION' version.go | awk '{ print $4 }' | tr -d '"')

# Define build function
build() {
    local os=$1
    local arch=$2

    local output_os="${os}"

    if [ "${os}" == "darwin" ]; then
        output_os="macos"
    fi

    local output_dir="release/${VERSION}/${output_os}/${arch}"
    local output_name="qai"

    if [ "${os}" == "windows" ]; then
        output_name="qai.exe"
    fi

    echo "Building for OS: ${os}, ARCH: ${arch}"
    mkdir -p "${output_dir}"
    GOOS=${os} GOARCH=${arch} go build -ldflags="-s" -ldflags="-w" -o "${output_dir}/${output_name}"
    echo "Build complete: ${output_dir}/${output_name}"
}

# Clean up previous builds
rm -rf release

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

# zip each ./release/${VERSION}/${os} directory into ./release/${os}-${VERSION}.zip
for os in linux macos windows; do
    zip -r -9 "release/${os}-${VERSION}.zip" "release/${VERSION}/${os}"
done