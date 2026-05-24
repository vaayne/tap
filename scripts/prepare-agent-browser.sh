#!/usr/bin/env bash
set -euo pipefail

VERSION="${AGENT_BROWSER_VERSION:-0.27.0}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="$ROOT/browser/bin"
BASE_URL="https://github.com/vercel-labs/agent-browser/releases/download/v${VERSION}"

mkdir -p "$OUT_DIR"

fetch() {
  local asset="$1"
  local path="$OUT_DIR/$asset"
  if [ -s "$path" ]; then
    return 0
  fi
  local tmp="$path.tmp"
  echo "Downloading $asset"
  curl -fsSL -o "$tmp" "$BASE_URL/$asset"
  chmod 0755 "$tmp"
  mv "$tmp" "$path"
}

fetch agent-browser-darwin-arm64
fetch agent-browser-darwin-x64
fetch agent-browser-linux-arm64
fetch agent-browser-linux-x64
fetch agent-browser-win32-x64.exe
