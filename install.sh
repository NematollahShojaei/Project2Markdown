#!/bin/bash

# Copyright 2026 R3D HILLS. All Rights Reserved.
# Enterprise Zero-Install Script for Project2Markdown (Linux & macOS)

set -e

# Colors for UI
GREEN='\033[0;32m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${CYAN}==================================================${NC}"
echo -e "${CYAN}  Installing Project2Markdown (P2M) Enterprise... ${NC}"
echo -e "${CYAN}==================================================${NC}"

# 1. Detect Operating System
OS="$(uname -s)"
case "${OS}" in
    Linux*)     TARGET_OS="linux";;
    Darwin*)    TARGET_OS="mac";;
    *)          echo -e "${RED}Error: Unsupported operating system: ${OS}${NC}"; exit 1;;
esac

# 2. Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64*)
        if [ "$TARGET_OS" = "mac" ]; then
            TARGET_ARCH="intel"
        else
            TARGET_ARCH="amd64"
        fi
        ;;
    arm64*|aarch64*)
        if [ "$TARGET_OS" = "mac" ]; then
            TARGET_ARCH="arm"
        else
            TARGET_ARCH="arm64"
        fi
        ;;
    *)
        echo -e "${RED}Error: Unsupported architecture: ${ARCH}${NC}"; exit 1;;
esac

BINARY_NAME="p2m-${TARGET_OS}-${TARGET_ARCH}"
DOWNLOAD_URL="https://github.com/nematollahshojaei/project2markdown/releases/latest/download/${BINARY_NAME}"
INSTALL_DIR="/usr/local/bin"
DEST_FILE="${INSTALL_DIR}/p2m"

echo -e "Detected System: ${GREEN}${OS} (${ARCH})${NC}"
echo -e "Downloading latest release..."

# 3. Download the binary to a temporary location
TMP_FILE=$(mktemp)
if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"; then
    echo -e "${RED}Error: Failed to download the binary. Please check your internet connection or ensure a release exists on GitHub.${NC}"
    rm -f "$TMP_FILE"
    exit 1
fi

# 4. Make it executable
chmod +x "$TMP_FILE"

# 5. Move to installation directory (requires sudo)
echo -e "Installing to ${GREEN}${INSTALL_DIR}${NC} (You may be prompted for your password)"
if ! sudo mv "$TMP_FILE" "$DEST_FILE"; then
    echo -e "${RED}Error: Failed to move binary to ${INSTALL_DIR}. Do you have sudo privileges?${NC}"
    rm -f "$TMP_FILE"
    exit 1
fi

echo -e "${GREEN}==================================================${NC}"
echo -e "${GREEN}  SUCCESS! P2M has been installed successfully.   ${NC}"
echo -e "${GREEN}==================================================${NC}"
echo -e "You can now run it from anywhere using the command:"
echo -e "  ${CYAN}p2m --cli${NC}"
echo ""