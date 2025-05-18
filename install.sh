#!/bin/bash

set -e

# Function to log messages
log() {
  echo "[QAI Installer] $1"
}

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [[ "$OS" == "darwin" ]]; then
  OS="macos"
elif [[ "$OS" != "linux" ]]; then
  log "Unsupported OS: $OS. This installer only supports macOS and Linux."
  exit 1
fi

# Detect architecture
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
  ARCH="arm64"
else
  log "Unsupported architecture: $ARCH. This installer only supports amd64 and arm64."
  exit 1
fi

log "Detected OS: $OS, Architecture: $ARCH"

# Determine installation directory
INSTALL_DIR=""
if [[ -d "$HOME/.local/bin" && ":$PATH:" == *":$HOME/.local/bin:"* ]]; then
  INSTALL_DIR="$HOME/.local/bin"
else
  INSTALL_DIR="/usr/local/bin"
fi

# Check if we need sudo for the install directory
SUDO=""
if [[ ! -w "$INSTALL_DIR" ]]; then
  if [[ -z $(command -v sudo) ]]; then
    log "Error: Need write permission for $INSTALL_DIR"
    exit 1
  fi
  log "Note: Using sudo to install to $INSTALL_DIR"
  SUDO="sudo"
fi

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Get latest release information
log "Fetching latest release information from GitHub..."
if [[ -z $(command -v curl) ]]; then
  log "Error: curl is required but not installed."
  exit 1
fi

GITHUB_REPO="McNull/qai"
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | 
                grep '"tag_name":' | 
                sed -E 's/.*"([^"]+)".*/\1/' | 
                sed 's/^v//')

if [[ -z "$LATEST_VERSION" ]]; then
  log "Error: Could not determine the latest version."
  exit 1
fi

log "Latest version: $LATEST_VERSION"

# Download the appropriate zip file
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/v$LATEST_VERSION/qai-$LATEST_VERSION-$OS.zip"
DOWNLOAD_PATH="$TMP_DIR/qai-$LATEST_VERSION-$OS.zip"

log "Downloading from $DOWNLOAD_URL..."
curl -L -o "$DOWNLOAD_PATH" "$DOWNLOAD_URL"

if [[ ! -f "$DOWNLOAD_PATH" ]]; then
  log "Error: Failed to download the zip file."
  exit 1
fi

# Extract the zip file
log "Extracting..."
unzip -q "$DOWNLOAD_PATH" -d "$TMP_DIR"

# Find and copy the executable
BINARY_PATH="$TMP_DIR/release/$LATEST_VERSION/$OS/$ARCH/qai"
if [[ "$OS" == "windows" ]]; then
  BINARY_PATH="$BINARY_PATH.exe"
fi

if [[ ! -f "$BINARY_PATH" ]]; then
  log "Error: Binary not found at expected path: $BINARY_PATH"
  exit 1
fi

DEST_PATH="$INSTALL_DIR/qai"
log "Installing to $DEST_PATH..."
# ensure the destination directory exists
$SUDO mkdir -p "$INSTALL_DIR"
$SUDO cp "$BINARY_PATH" "$DEST_PATH"
$SUDO chmod +x "$DEST_PATH"


# Verify installation
if command -v qai >/dev/null 2>&1; then
  
  # Check if the config file exists
  CONFIG_DIR="$HOME/.config/qai"
  CONFIG_FILE="$CONFIG_DIR/config.json"

  if [[ ! -d "$CONFIG_DIR" ]]; then
    log "Creating config directory at $CONFIG_DIR"
    mkdir -p "$CONFIG_DIR"
  fi

  if [[ ! -f "$CONFIG_FILE" ]]; then
    log "Creating default config file at $CONFIG_FILE"
    qai --create-config
  fi

  log "Installation completed successfully!"
  log "You can now run 'qai' from your terminal."

  log "Adjust your config file at $CONFIG_FILE to customize your settings."
  log "A new config file can be created with 'qai --create-config'."
  log "For more information, visit: https://github.com/McNull/qai"

else
  log "Could not find 'qai' in your PATH."
  log "Note: You may need to add $INSTALL_DIR to your PATH or restart your terminal."
fi

