#!/bin/sh
# Install or update Tap and its pinned agent-browser runtime dependency.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh -s -- --full
#   curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh -s -- --dir /usr/local/bin

set -eu

REPO="vaayne/tap"
INSTALL_DIR="${HOME}/.local/bin"
FULL=0
WITH_AGENT_BROWSER=1
INSTALL_CHROME=1

usage() {
  cat <<'EOF'
Usage: install.sh [options]

Options:
  --dir PATH                 Install binaries under PATH (default ~/.local/bin)
  --full                     Use the Tap + agent-browser full bundle
  --skip-chrome              Do not run agent-browser install
  --without-agent-browser    Install only Tap
  -h, --help                 Show this help
EOF
}

while [ $# -gt 0 ]; do
  case "$1" in
    --dir) INSTALL_DIR="$2"; shift 2 ;;
    --dir=*) INSTALL_DIR="${1#--dir=}"; shift ;;
    --full) FULL=1; shift ;;
    --without-agent-browser) WITH_AGENT_BROWSER=0; shift ;;
    --skip-chrome) INSTALL_CHROME=0; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; exit 1 ;;
  esac
done

if [ "$FULL" -eq 1 ] && [ "$WITH_AGENT_BROWSER" -eq 0 ]; then
  echo "Error: --full and --without-agent-browser cannot be combined" >&2
  exit 1
fi

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

LIBC="default"
if [ "$OS" = "linux" ]; then
  LIBC="glibc"
  if (ldd --version 2>&1 || true) | grep -qi musl || ls /lib/ld-musl-*.so.1 >/dev/null 2>&1; then
    LIBC="musl"
  fi
fi

LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": *"//;s/".*//')
if [ -z "$LATEST" ]; then
  echo "Error: could not determine latest Tap version" >&2
  exit 1
fi
VERSION="${LATEST#v}"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

needs_agent_browser=0
if [ "$WITH_AGENT_BROWSER" -eq 1 ]; then
  manifest="${TMPDIR}/agent-browser.sh"
  if [ "$FULL" -eq 1 ] || ! command -v agent-browser >/dev/null 2>&1; then
    manifest_status=$(curl -sSL -o "$manifest" -w '%{http_code}' "https://raw.githubusercontent.com/${REPO}/${LATEST}/deps/agent-browser.sh") || {
      echo "Error: could not download the agent-browser dependency manifest" >&2
      exit 1
    }
    if [ "$manifest_status" = "404" ]; then
      if [ "$FULL" -eq 1 ]; then
        echo "Error: Tap ${LATEST} does not publish full agent-browser bundles" >&2
        exit 1
      fi
      echo "Warning: Tap ${LATEST} predates agent-browser bootstrap; installing Tap only." >&2
      WITH_AGENT_BROWSER=0
    elif [ "$manifest_status" != "200" ]; then
      echo "Error: dependency manifest returned HTTP ${manifest_status}" >&2
      exit 1
    else
      # shellcheck source=/dev/null
      . "$manifest"
    fi
  fi

  if [ "$WITH_AGENT_BROWSER" -eq 1 ] && [ "$FULL" -eq 1 ]; then
    if [ "$OS" = "windows" ] && [ "$ARCH" = "arm64" ]; then
      echo "Error: agent-browser has no native Windows arm64 release" >&2
      echo "Use --without-agent-browser or install a supported build manually." >&2
      exit 1
    fi
  elif [ "$WITH_AGENT_BROWSER" -eq 1 ] && ! command -v agent-browser >/dev/null 2>&1; then
    if ! agent_browser_resolve "$OS" "$ARCH" "$LIBC"; then
      echo "Error: agent-browser ${AGENT_BROWSER_VERSION} has no native binary for ${OS}/${ARCH}/${LIBC}" >&2
      echo "Use --without-agent-browser or install a supported build manually." >&2
      exit 1
    fi
    needs_agent_browser=1
  fi
fi

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

verify_sha256() {
  actual=$(sha256_file "$1")
  if [ "$actual" != "$2" ]; then
    echo "Error: checksum mismatch for $1" >&2
    echo "  want: $2" >&2
    echo "  got:  $actual" >&2
    exit 1
  fi
}

archive_os="$OS"
suffix=""
if [ "$FULL" -eq 1 ]; then
  suffix="_full"
  [ "$LIBC" = "musl" ] && archive_os="linux_musl"
fi
if [ "$OS" = "windows" ]; then
  extension="zip"
else
  extension="tar.gz"
fi
ARCHIVE="tap_${VERSION}_${archive_os}_${ARCH}${suffix}.${extension}"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${ARCHIVE}"

echo "Installing Tap ${VERSION} (${OS}/${ARCH}, ${LIBC})..."
echo "Downloading ${URL}..."
curl -fsSL -o "${TMPDIR}/${ARCHIVE}" "$URL"
if [ "$FULL" -eq 1 ]; then
  checksum_asset="full-checksums.txt"
else
  checksum_asset="checksums.txt"
fi
curl -fsSL "https://github.com/${REPO}/releases/download/${LATEST}/${checksum_asset}" -o "${TMPDIR}/${checksum_asset}"
expected=$(awk -v archive="$ARCHIVE" '$2 == archive { print $1 }' "${TMPDIR}/${checksum_asset}")
if [ -z "$expected" ]; then
  echo "Error: ${ARCHIVE} is absent from ${checksum_asset}" >&2
  exit 1
fi
verify_sha256 "${TMPDIR}/${ARCHIVE}" "$expected"
case "$extension" in
  zip) unzip -q "${TMPDIR}/${ARCHIVE}" -d "$TMPDIR" ;;
  tar.gz) tar -xzf "${TMPDIR}/${ARCHIVE}" -C "$TMPDIR" ;;
esac

mkdir -p "$INSTALL_DIR"
tap_source="${TMPDIR}/tap"
tap_target="${INSTALL_DIR}/tap"
if [ "$OS" = "windows" ]; then
  tap_source="${tap_source}.exe"
  tap_target="${tap_target}.exe"
fi
cp "$tap_source" "$tap_target"
chmod +x "$tap_target"
echo "Installed Tap ${VERSION} to ${tap_target}"

installed_agent_browser=0
agent_browser_target="${INSTALL_DIR}/agent-browser"
[ "$OS" = "windows" ] && agent_browser_target="${agent_browser_target}.exe"

if [ "$WITH_AGENT_BROWSER" -eq 1 ]; then
  if [ "$FULL" -eq 1 ]; then
    agent_browser_source="${TMPDIR}/agent-browser"
    [ "$OS" = "windows" ] && agent_browser_source="${agent_browser_source}.exe"
    cp "$agent_browser_source" "$agent_browser_target"
    chmod +x "$agent_browser_target"
    installed_agent_browser=1
  elif [ "$needs_agent_browser" -eq 1 ]; then
    echo "Downloading agent-browser ${AGENT_BROWSER_VERSION}..."
    curl -fsSL "${AGENT_BROWSER_RELEASE_BASE}/${AGENT_BROWSER_ASSET}" -o "$agent_browser_target"
    verify_sha256 "$agent_browser_target" "$AGENT_BROWSER_SHA256"
    chmod +x "$agent_browser_target"
    installed_agent_browser=1

    license_dir="${INSTALL_DIR}/../share/tap/licenses"
    mkdir -p "$license_dir"
    curl -fsSL "$AGENT_BROWSER_LICENSE_URL" -o "${license_dir}/agent-browser.LICENSE"
    verify_sha256 "${license_dir}/agent-browser.LICENSE" "$AGENT_BROWSER_LICENSE_SHA256"
  fi

  if [ "$FULL" -eq 1 ]; then
    license_dir="${INSTALL_DIR}/../share/tap/licenses"
    mkdir -p "$license_dir"
    cp "${TMPDIR}/licenses/agent-browser.LICENSE" "${license_dir}/agent-browser.LICENSE"
  fi

  if [ "$installed_agent_browser" -eq 1 ]; then
    echo "Installed $("$agent_browser_target" --version) to ${agent_browser_target}"
    if [ "$INSTALL_CHROME" -eq 1 ]; then
      echo "Installing agent-browser Chrome runtime..."
      "$agent_browser_target" install
    else
      echo "Skipped Chrome installation; run '${agent_browser_target} install' when ready."
    fi
  else
    echo "Using existing $(agent-browser --version)"
  fi
fi

case ":${PATH}:" in
  *":${INSTALL_DIR}:"*) ;;
  *)
    echo ""
    echo "Add ${INSTALL_DIR} to your PATH:"
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    ;;
esac
