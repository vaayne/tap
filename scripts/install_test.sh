#!/bin/sh
# Hermetic smoke test for online bootstrap and full-bundle installation.

set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
FIXTURES="${WORK}/fixtures"
mkdir -p "$FIXTURES/bin" "$FIXTURES/thin" "$FIXTURES/full/licenses"

sh "${ROOT}/scripts/install.sh" --help | grep -q -- '--full'

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo "unsupported test architecture: $arch" >&2; exit 1 ;;
esac
case "$os" in
  darwin|linux) ;;
  *) echo "installer smoke test skipped on $os"; exit 0 ;;
esac

version="9.9.9"
thin_archive="tap_${version}_${os}_${arch}.tar.gz"
full_archive="tap_${version}_${os}_${arch}_full.tar.gz"

cat > "${FIXTURES}/thin/tap" <<'EOF'
#!/bin/sh
echo 'tap fixture'
EOF
chmod +x "${FIXTURES}/thin/tap"
tar -czf "${FIXTURES}/${thin_archive}" -C "${FIXTURES}/thin" .

cat > "${FIXTURES}/fixture-agent-browser" <<'EOF'
#!/bin/sh
case "${1:-}" in
  --version) echo 'agent-browser 9.9.9' ;;
  install)
    echo 'fixture Chrome installed'
    [ -z "${FIXTURE_CHROME_LOG:-}" ] || : > "$FIXTURE_CHROME_LOG"
    ;;
  *) exit 0 ;;
esac
EOF
chmod +x "${FIXTURES}/fixture-agent-browser"
printf '%s\n' 'fixture license' > "${FIXTURES}/agent-browser.LICENSE"
ab_sha=$(sha256_file "${FIXTURES}/fixture-agent-browser")
license_sha=$(sha256_file "${FIXTURES}/agent-browser.LICENSE")

cat > "${FIXTURES}/agent-browser.sh" <<EOF
AGENT_BROWSER_VERSION="v9.9.9"
AGENT_BROWSER_RELEASE_BASE="https://fixtures.invalid"
AGENT_BROWSER_LICENSE_URL="https://fixtures.invalid/agent-browser.LICENSE"
AGENT_BROWSER_LICENSE_SHA256="${license_sha}"
agent_browser_resolve() {
  AGENT_BROWSER_ASSET="fixture-agent-browser"
  AGENT_BROWSER_SHA256="${ab_sha}"
  return 0
}
EOF

cp "${FIXTURES}/thin/tap" "${FIXTURES}/full/tap"
cp "${FIXTURES}/fixture-agent-browser" "${FIXTURES}/full/agent-browser"
cp "${FIXTURES}/agent-browser.LICENSE" "${FIXTURES}/full/licenses/agent-browser.LICENSE"
tar -czf "${FIXTURES}/${full_archive}" -C "${FIXTURES}/full" .
printf '%s  %s\n' "$(sha256_file "${FIXTURES}/${thin_archive}")" "$thin_archive" > "${FIXTURES}/checksums.txt"
printf '%s  %s\n' "$(sha256_file "${FIXTURES}/${full_archive}")" "$full_archive" > "${FIXTURES}/full-checksums.txt"
printf '{"tag_name":"v%s"}\n' "$version" > "${FIXTURES}/latest"

cat > "${FIXTURES}/bin/curl" <<'EOF'
#!/bin/sh
set -eu
output=""
url=""
writeout=""
while [ $# -gt 0 ]; do
  case "$1" in
    -o) output=$2; shift 2 ;;
    -w) writeout=$2; shift 2 ;;
    -*) shift ;;
    *) url=$1; shift ;;
  esac
done
case "$url" in
  */api.github.com/*/releases/latest|https://api.github.com/*/releases/latest) source="${FIXTURES}/latest" ;;
  */deps/agent-browser.sh) source="${FIXTURES}/agent-browser.sh" ;;
  */fixture-agent-browser) source="${FIXTURES}/fixture-agent-browser" ;;
  */agent-browser.LICENSE) source="${FIXTURES}/agent-browser.LICENSE" ;;
  */checksums.txt) source="${FIXTURES}/checksums.txt" ;;
  */full-checksums.txt) source="${FIXTURES}/full-checksums.txt" ;;
  *.tar.gz) source="${FIXTURES}/$(basename "$url")" ;;
  *) echo "unexpected URL: $url" >&2; exit 1 ;;
esac
if [ -n "$output" ]; then
  cp "$source" "$output"
else
  cat "$source"
fi
if [ -n "$writeout" ]; then
  printf '%s' '200'
fi
EOF
chmod +x "${FIXTURES}/bin/curl"

run_install() {
  destination=$1
  shift
  FIXTURES="$FIXTURES" FIXTURE_CHROME_LOG="${WORK}/chrome-installed" \
    HOME="$WORK/home" PATH="${FIXTURES}/bin:/usr/bin:/bin" \
    sh "${ROOT}/scripts/install.sh" --dir "$destination" "$@"
  test -x "${destination}/tap"
  test -x "${destination}/agent-browser"
  "${destination}/agent-browser" --version | grep -q '9.9.9'
  test -f "${destination}/../share/tap/licenses/agent-browser.LICENSE"
}

run_install "${WORK}/online-bin"
test -f "${WORK}/chrome-installed"
rm "${WORK}/chrome-installed"
run_install "${WORK}/full-bin" --full --skip-chrome
test ! -e "${WORK}/chrome-installed"

cp "${FIXTURES}/fixture-agent-browser" "${FIXTURES}/bin/agent-browser"
FIXTURES="$FIXTURES" HOME="$WORK/home" PATH="${FIXTURES}/bin:/usr/bin:/bin" \
  sh "${ROOT}/scripts/install.sh" --dir "${WORK}/existing-bin" --skip-chrome >/dev/null
test -x "${WORK}/existing-bin/tap"
test ! -e "${WORK}/existing-bin/agent-browser"
rm "${FIXTURES}/bin/agent-browser"

printf 'tampered' >> "${FIXTURES}/${thin_archive}"
if FIXTURES="$FIXTURES" HOME="$WORK/home" PATH="${FIXTURES}/bin:/usr/bin:/bin" \
  sh "${ROOT}/scripts/install.sh" --dir "${WORK}/rejected-bin" \
    --without-agent-browser >"${WORK}/rejected.out" 2>"${WORK}/rejected.err"; then
  echo "tampered archive was accepted" >&2
  exit 1
fi
grep -q 'checksum mismatch' "${WORK}/rejected.err"
test ! -e "${WORK}/rejected-bin/tap"
echo "installer bootstrap and full bundle: ok"
