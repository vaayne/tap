# Pinned agent-browser release used by Tap's installer and full bundles.
# Keep checksums immutable: upgrade the version and all platform entries together.
AGENT_BROWSER_VERSION="v0.33.2"
AGENT_BROWSER_RELEASE_BASE="https://github.com/vercel-labs/agent-browser/releases/download/${AGENT_BROWSER_VERSION}"
AGENT_BROWSER_LICENSE_URL="https://raw.githubusercontent.com/vercel-labs/agent-browser/${AGENT_BROWSER_VERSION}/LICENSE"
AGENT_BROWSER_LICENSE_SHA256="014bb31e83d5c2e76aea1cc6e82217346ab41362f32cb355ad0f5c10aa0aeaff"

# agent_browser_resolve OS ARCH LIBC sets AGENT_BROWSER_ASSET and
# AGENT_BROWSER_SHA256. OS uses darwin/linux/windows; ARCH uses amd64/arm64.
agent_browser_resolve() {
  AGENT_BROWSER_ASSET=""
  AGENT_BROWSER_SHA256=""

  case "$1/$2/$3" in
    darwin/arm64/*)
      AGENT_BROWSER_ASSET="agent-browser-darwin-arm64"
      AGENT_BROWSER_SHA256="cbb517902bcaa3b7a6384fd9f25dd274da3df2bb6a3ba9c3e85806d78213c26b"
      ;;
    darwin/amd64/*)
      AGENT_BROWSER_ASSET="agent-browser-darwin-x64"
      AGENT_BROWSER_SHA256="a6bb1c10124f624a9b1fd0eecabf774477cdb710e3552fb843f1f7f664b8f326"
      ;;
    linux/arm64/musl)
      AGENT_BROWSER_ASSET="agent-browser-linux-musl-arm64"
      AGENT_BROWSER_SHA256="eec7d0a27e32b96a4f9b9fbdd0c070d058e5b4eaa1bd6be1fffe926321c5d01c"
      ;;
    linux/amd64/musl)
      AGENT_BROWSER_ASSET="agent-browser-linux-musl-x64"
      AGENT_BROWSER_SHA256="ca7e6589158fd9276897ec66367105704a215f95b1df4c4abb193244d0260eda"
      ;;
    linux/arm64/*)
      AGENT_BROWSER_ASSET="agent-browser-linux-arm64"
      AGENT_BROWSER_SHA256="6ccaba1eb26a0e6f5c23c59d2c63e6e0237fde82713cfdb543ba506490cac9c1"
      ;;
    linux/amd64/*)
      AGENT_BROWSER_ASSET="agent-browser-linux-x64"
      AGENT_BROWSER_SHA256="b7bc3dfcf0a7326c1f5a60423163259ba2349eebfa5bd2e70e111af743da4a49"
      ;;
    windows/amd64/*)
      AGENT_BROWSER_ASSET="agent-browser-win32-x64.exe"
      AGENT_BROWSER_SHA256="291f0c33c2fbcbf159b5868065ab412dfd8722d6299821e010cf0715964f2cba"
      ;;
  esac

  [ -n "$AGENT_BROWSER_ASSET" ]
}
