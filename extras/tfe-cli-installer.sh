#!/bin/bash
set -euo pipefail

# Prepare variables.
PROJECT=tfe-cli
REPO=rgreinho/${PROJECT}
LATEST_TAG=$(git ls-remote --tags --refs --sort="v:refname" https://github.com/${REPO}.git | tail -n1 | sed 's/.*\///')
VERSION=${VERSION:-$LATEST_TAG}
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
OPT_DIR="/usr/local/opt/${REPO}/${VERSION}"
ARCH=amd64
BINARY="${PROJECT}-${VERSION}-${PLATFORM}-${ARCH}"


# Download the binaries.
mkdir -p "${OPT_DIR}"
pushd "${OPT_DIR}"
curl -LO "https://github.com/${REPO}/releases/download/${VERSION}/${PROJECT}-${VERSION}-${PLATFORM}-${ARCH}"
popd

# Create the simlink
SRC="${OPT_DIR}/${BINARY}"
TARGET="/usr/local/bin/${PROJECT}"
echo "Updating permissions..."
chmod +x "${SRC}"
echo "Linking ${SRC} to ${TARGET}..."
ln -sf "${SRC}" "${TARGET}"
