#!/bin/bash
set -euo pipefail

# Prepare variables.
REPO=rgreinho/tfe-cli
LATEST_TAG=$(git ls-remote --tags --refs --sort="v:refname" git://github.com/${REPO}.git | tail -n1 | sed 's/.*\///')
VERSION=${VERSION:-$LATEST_TAG}
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
OPT_DIR="/usr/local/opt/${REPO}/${VERSION}"
BINARY="tfe-client-${VERSION}-${PLATFORM}-amd64"
export GITHUB_OAUTH_TOKEN=${GITHUB_TOKEN}

# Download the binaries.
mkdir -p "${OPT_DIR}"
fetch --repo="https://github.com/${REPO}" \
  --tag="${VERSION}" \
  --release-asset="${BINARY}" \
  ${OPT_DIR}

# Create the simlink
SRC="${OPT_DIR}/${BINARY}"
TARGET="/usr/local/bin/tfe-client"
echo "Updating permissions..."
chmod +x "${SRC}"
echo "Linking ${SRC} to ${TARGET}..."
ln -sf "${SRC}" "${TARGET}"
