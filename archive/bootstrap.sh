#!/usr/bin/env bash
# bootstrap.sh - Download and run dotsetup binary
# Usage: git clone <repo> && cd dotfile && bash bootstrap.sh
set -euo pipefail

REPO="aikenhong/dotsetup"
BINARY="dotsetup"
INSTALL_DIR="${HOME}/.local/bin"

# Detect platform
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "${ARCH}" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: ${ARCH}"; exit 1 ;;
esac

echo "🔧 dotsetup bootstrap"
echo "   Platform: ${OS}/${ARCH}"

# Check if dotsetup is already built locally
if [ -f "./dotsetup" ]; then
  echo "   Found local binary, running..."
  exec ./dotsetup "$@"
fi

# Check if Go is available → build from source
if command -v go &>/dev/null; then
  echo "   Go detected, building from source..."
  go build -o ./dotsetup ./cmd/dotsetup/
  exec ./dotsetup "$@"
fi

# Fallback: download pre-built binary
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY}_${OS}_${ARCH}"

echo "   Downloading pre-built binary..."
mkdir -p "${INSTALL_DIR}"

if command -v curl &>/dev/null; then
  curl -fsSL "${DOWNLOAD_URL}" -o "${INSTALL_DIR}/${BINARY}"
elif command -v wget &>/dev/null; then
  wget -q "${DOWNLOAD_URL}" -O "${INSTALL_DIR}/${BINARY}"
else
  echo "Error: neither curl nor wget found"
  exit 1
fi

chmod +x "${INSTALL_DIR}/${BINARY}"
echo "   Installed to ${INSTALL_DIR}/${BINARY}"

# Ensure ~/.local/bin is in PATH
if [[ ":${PATH}:" != *":${INSTALL_DIR}:"* ]]; then
  echo "   Note: Add ${INSTALL_DIR} to your PATH"
  export PATH="${INSTALL_DIR}:${PATH}"
fi

exec "${INSTALL_DIR}/${BINARY}" "$@"
