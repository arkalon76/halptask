#!/usr/bin/env bash
set -e

# HalpTask Installer Script
# Usage: curl -fsSL https://raw.githubusercontent.com/kenth/halptask/main/scripts/install.sh | bash

REPO="kenth/halptask"
BINARY_NAME="halptask"

echo "🚀 Installing HalpTask..."

# 1. Detect OS
OS_RAW=$(uname -s)
case "${OS_RAW}" in
    Linux*)     OS="Linux";;
    Darwin*)    OS="Darwin";;
    *)          echo "❌ Unsupported operating system: ${OS_RAW}. HalpTask install script supports Linux and macOS."; exit 1;;
esac

# 2. Detect Architecture
ARCH_RAW=$(uname -m)
case "${ARCH_RAW}" in
    x86_64|amd64)   ARCH="x86_64";;
    arm64|aarch64)  ARCH="arm64";;
    *)              echo "❌ Unsupported architecture: ${ARCH_RAW}. Supported architectures: x86_64, arm64."; exit 1;;
esac

# 3. Determine Version
if [ -z "${VERSION}" ]; then
    echo "🔍 Finding latest release..."
    LATEST_RELEASE=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null || true)
    if [ -n "${LATEST_RELEASE}" ]; then
        TAG=$(echo "${LATEST_RELEASE}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi
    if [ -z "${TAG}" ]; then
        echo "⚠️ Could not automatically detect latest tag from GitHub API, falling back to latest release download..."
        DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}_${OS}_${ARCH}"
    else
        DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${BINARY_NAME}_${OS}_${ARCH}"
    fi
else
    TAG="${VERSION}"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${BINARY_NAME}_${OS}_${ARCH}"
fi

# 4. Determine Destination Directory
INSTALL_DIR="/usr/local/bin"
USE_SUDO=false

if [ ! -w "${INSTALL_DIR}" ]; then
    if [ -d "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
    elif command -v sudo >/dev/null 2>&1; then
        USE_SUDO=true
    else
        INSTALL_DIR="$HOME/bin"
        mkdir -p "${INSTALL_DIR}"
    fi
fi

# 5. Download Binary
TMP_DIR=$(mktemp -d)
trap 'rm -rf "${TMP_DIR}"' EXIT

TMP_BINARY="${TMP_DIR}/${BINARY_NAME}"
echo "📥 Downloading binary for ${OS}/${ARCH}..."
if ! curl -fsSL "${DOWNLOAD_URL}" -o "${TMP_BINARY}"; then
    echo "❌ Download failed! Could not download from ${DOWNLOAD_URL}"
    exit 1
fi

chmod +x "${TMP_BINARY}"

# 6. Install Binary
echo "📦 Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
if [ "${USE_SUDO}" = true ]; then
    sudo mv "${TMP_BINARY}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    mv "${TMP_BINARY}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo "✅ HalpTask successfully installed!"
echo "🎉 Run 'halptask' to get started."
