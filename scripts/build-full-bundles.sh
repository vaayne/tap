#!/bin/sh
# Add pinned agent-browser binaries to existing GoReleaser archives.
# Usage: scripts/build-full-bundles.sh VERSION [DIST_DIR] [OUTPUT_DIR]

set -eu

if [ $# -lt 1 ] || [ $# -gt 3 ]; then
  echo "Usage: $0 VERSION [DIST_DIR] [OUTPUT_DIR]" >&2
  exit 1
fi

VERSION=$1
DIST_DIR=${2:-dist}
OUTPUT_DIR=${3:-${DIST_DIR}/full}
ROOT=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)

# shellcheck source=../deps/agent-browser.sh
. "${ROOT}/deps/agent-browser.sh"

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
mkdir -p "$OUTPUT_DIR" "$WORK/downloads"
DIST_DIR=$(CDPATH= cd -- "$DIST_DIR" && pwd)
OUTPUT_DIR=$(CDPATH= cd -- "$OUTPUT_DIR" && pwd)

license="${WORK}/downloads/agent-browser.LICENSE"
curl -fsSL "$AGENT_BROWSER_LICENSE_URL" -o "$license"
actual_license=$(sha256_file "$license")
if [ "$actual_license" != "$AGENT_BROWSER_LICENSE_SHA256" ]; then
  echo "Checksum mismatch for agent-browser LICENSE" >&2
  exit 1
fi

bundle() {
  os=$1
  arch=$2
  libc=$3
  extension=$4

  agent_browser_resolve "$os" "$arch" "$libc" || {
    echo "No agent-browser binary for ${os}/${arch}/${libc}" >&2
    return 1
  }

  thin_base="tap_${VERSION}_${os}_${arch}"
  thin_archive="${DIST_DIR}/${thin_base}.${extension}"
  if [ ! -f "$thin_archive" ]; then
    echo "Missing GoReleaser archive: ${thin_archive}" >&2
    return 1
  fi

  binary="${WORK}/downloads/${AGENT_BROWSER_ASSET}"
  if [ ! -f "$binary" ]; then
    echo "Downloading agent-browser ${AGENT_BROWSER_VERSION} (${os}/${arch}/${libc})..."
    curl -fsSL "${AGENT_BROWSER_RELEASE_BASE}/${AGENT_BROWSER_ASSET}" -o "$binary"
  fi
  actual=$(sha256_file "$binary")
  if [ "$actual" != "$AGENT_BROWSER_SHA256" ]; then
    echo "Checksum mismatch for ${AGENT_BROWSER_ASSET}" >&2
    echo "  want: ${AGENT_BROWSER_SHA256}" >&2
    echo "  got:  ${actual}" >&2
    return 1
  fi

  root="${WORK}/${os}-${arch}-${libc}"
  rm -rf "$root"
  mkdir -p "$root/licenses"
  case "$extension" in
    zip) unzip -q "$thin_archive" -d "$root" ;;
    tar.gz) tar -xzf "$thin_archive" -C "$root" ;;
    *) echo "Unsupported archive format: ${extension}" >&2; return 1 ;;
  esac

  target="agent-browser"
  [ "$os" = "windows" ] && target="agent-browser.exe"
  cp "$binary" "${root}/${target}"
  chmod +x "${root}/${target}"
  cp "$license" "${root}/licenses/agent-browser.LICENSE"

  full_os=$os
  [ "$libc" = "musl" ] && full_os="linux_musl"
  full_base="tap_${VERSION}_${full_os}_${arch}_full"
  output="${OUTPUT_DIR}/${full_base}.${extension}"
  case "$extension" in
    zip) (cd "$root" && zip -qr "$output" .) ;;
    tar.gz) tar -czf "$output" -C "$root" . ;;
  esac
  echo "Built ${output}"
}

bundle darwin amd64 default tar.gz
bundle darwin arm64 default tar.gz
bundle linux amd64 glibc tar.gz
bundle linux arm64 glibc tar.gz
bundle linux amd64 musl tar.gz
bundle linux arm64 musl tar.gz
bundle windows amd64 default zip

(
  cd "$OUTPUT_DIR"
  : > full-checksums.txt
  for archive in *_full.tar.gz *_full.zip; do
    [ -f "$archive" ] || continue
    printf '%s  %s\n' "$(sha256_file "$archive")" "$archive" >> full-checksums.txt
  done
)
echo "Wrote ${OUTPUT_DIR}/full-checksums.txt"
