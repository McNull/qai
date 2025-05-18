#!/bin/bash

set -e

VERSION=$(grep 'APP_VERSION' ./app/constants.go | awk '{ print $3 }' | tr -d '"')

echo "Building version: ${VERSION}"

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
    GOOS=${os} GOARCH=${arch} go build -ldflags="-s -w" -trimpath -o "${output_dir}/${output_name}"
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
    echo "Zipping ${os} build..."
    zip -r -9 "release/qai-${VERSION}-${os}.zip" "release/${VERSION}/${os}"
done

tree -h release

# Prompt user for confirmation before uploading to GitHub
echo "The following files will be uploaded to GitHub:"
ls -lh release/qai-${VERSION}-*.zip
echo "Do you want to continue? (y/n)"
read -r confirm
if [[ ! $confirm =~ ^[yY]$ ]]; then
    echo "Aborting release process."
    exit 1
fi

# Create GitHub release
echo "Creating GitHub release v${VERSION}..."

# Check if GitHub CLI is installed
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI (gh) is not installed. Please install it first."
    echo "Visit https://cli.github.com/ for installation instructions."
    exit 1
fi

# Create the release and upload assets
gh release create "v${VERSION}" \
    --title "QAI v${VERSION}" \
    --notes "Release version ${VERSION}" \
    "release/qai-${VERSION}-linux.zip" \
    "release/qai-${VERSION}-macos.zip" \
    "release/qai-${VERSION}-windows.zip"

echo "GitHub release created successfully!"