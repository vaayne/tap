#!/bin/sh
# Install or update tap — https://github.com/vaayne/tap
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh -s -- --dir /usr/local/bin

set -eu

REPO="vaayne/tap"
INSTALL_DIR="${HOME}/.local/bin"

# Parse flags
while [ $# -gt 0 ]; do
  case "$1" in
    --dir) INSTALL_DIR="$2"; shift 2 ;;
    --dir=*) INSTALL_DIR="${1#--dir=}"; shift ;;
    *) echo "Unknown option: $1" >&2; exit 1 ;;
  esac
done

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

case "$OS" in
  linux|darwin) ;;
  mingw*|msys*|cygwin*) OS="windows" ;;
  *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

# Fetch latest version from GitHub API
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": *"//;s/".*//')
if [ -z "$LATEST" ]; then
  echo "Error: could not determine latest version" >&2
  exit 1
fi

VERSION="${LATEST#v}"
echo "Installing tap ${VERSION} (${OS}/${ARCH})..."

# Build download URL
if [ "$OS" = "windows" ]; then
  ARCHIVE="tap_${VERSION}_${OS}_${ARCH}.zip"
else
  ARCHIVE="tap_${VERSION}_${OS}_${ARCH}.tar.gz"
fi
URL="https://github.com/${REPO}/releases/download/${LATEST}/${ARCHIVE}"

# Download and extract
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${URL}..."
curl -fsSL -o "${TMPDIR}/${ARCHIVE}" "$URL"

if [ "$OS" = "windows" ]; then
  unzip -q "${TMPDIR}/${ARCHIVE}" -d "$TMPDIR"
else
  tar -xzf "${TMPDIR}/${ARCHIVE}" -C "$TMPDIR"
fi

# Install binary
mkdir -p "$INSTALL_DIR"
cp "${TMPDIR}/tap" "${INSTALL_DIR}/tap"
chmod +x "${INSTALL_DIR}/tap"

echo "Installed tap ${VERSION} to ${INSTALL_DIR}/tap"

# agent-browser is Tap's runtime dependency. Keep ownership with its own
# installer rather than downloading or versioning it inside Tap.
if ! command -v agent-browser >/dev/null 2>&1; then
  echo ""
  echo "Tap requires agent-browser:"
  echo "  npm install -g agent-browser"
  echo "  agent-browser install"
fi

# Check if install dir is in PATH
case ":${PATH}:" in
  *":${INSTALL_DIR}:"*) ;;
  *)
    echo ""
    echo "Add ${INSTALL_DIR} to your PATH:"
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    ;;
esac
