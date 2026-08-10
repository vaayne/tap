# Runtime dependency manifest

`agent-browser.sh` is the canonical version/checksum manifest used by the
online installer and release bundler. Tap does not embed these executables.

To upgrade agent-browser:

1. Change `AGENT_BROWSER_VERSION`.
2. Download every upstream platform asset and replace all SHA-256 values.
3. Replace the pinned license checksum if it changed.
4. Run `sh -n deps/agent-browser.sh scripts/*.sh` and
   `scripts/install_test.sh`.
5. Run a GoReleaser snapshot followed by `scripts/build-full-bundles.sh`, then
   verify `dist/full/full-checksums.txt` and inspect each archive.

The upstream release currently has no Windows arm64 binary. Tap still ships a
thin Windows arm64 archive, but bootstrap/full-bundle installation must fail
clearly rather than substituting an x64 executable.
