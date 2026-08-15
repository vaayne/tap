# Installation

Tap executes site programs and Defuddle through the external `agent-browser`
CLI. The default installer sets up both tools while keeping them as separate
executables.

## Recommended: online bootstrap

```bash
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh | sh
```

The installer:

1. Downloads the latest Tap release archive.
2. Verifies it against the release `checksums.txt`.
3. Preserves an existing `agent-browser` on `PATH`; otherwise downloads Tap's
   pinned native version and verifies its SHA-256.
4. Installs the agent-browser license alongside Tap's shared data.
5. When it installs a new agent-browser binary, runs `agent-browser install` to
   install Chrome for Testing.

Both executables default to `~/.local/bin`. Add that directory to `PATH` when
the installer asks, then verify the installation:

```bash
tap doctor
```

### Installer options

```text
--dir PATH                 Install binaries under PATH
--full                     Download the Tap + agent-browser full archive
--skip-chrome              Do not run agent-browser install
--without-agent-browser    Install only Tap
-h, --help                 Show installer help
```

For example, use an existing system Chrome without downloading Chrome for
Testing:

```bash
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh \
  | sh -s -- --skip-chrome
```

## Full bundles and offline transfer

Each release also publishes `_full` archives containing:

```text
tap
agent-browser[.exe]
licenses/agent-browser.LICENSE
```

Full bundles are available for:

| Platform | Architectures |
|---|---|
| macOS | amd64, arm64 |
| Linux (glibc) | amd64, arm64 |
| Linux (musl) | amd64, arm64 |
| Windows | amd64 |

The shell installer on Windows requires Git Bash, MSYS2, or Cygwin. Native
PowerShell users should download and extract the Windows zip directly.

Download the matching archive and verify it with `full-checksums.txt` from the
same [GitHub Release](https://github.com/vaayne/tap/releases/latest). After
extraction, Tap automatically discovers an `agent-browser` executable beside
it.

The full bundle does **not** contain Chrome. Fully offline use therefore
requires an existing compatible Chrome installation. Otherwise run this once
when online:

```bash
tap browser install
```

Upstream does not currently publish a Windows arm64 agent-browser binary. Tap
still publishes its thin Windows arm64 archive, but bootstrap and full-bundle
installation intentionally refuse to substitute an x64 runtime.

## Tap-only installation

These methods install only Tap and require `agent-browser` to be installed
separately:

```bash
curl -fsSL https://raw.githubusercontent.com/vaayne/tap/main/scripts/install.sh \
  | sh -s -- --without-agent-browser

go install github.com/vaayne/tap/cmd/tap@latest
```

## Runtime resolution

Tap resolves the agent-browser executable in this order:

1. An explicit override (`TAP_AGENT_BROWSER` or the Go library option).
2. `agent-browser[.exe]` beside the running Tap binary.
3. `agent-browser` on `PATH`.

Tap inherits `AGENT_BROWSER_SESSION` and all other agent-browser environment
settings unchanged. It never creates, names, persists, or closes sessions.

## Upgrades and repair

Tap and agent-browser remain independently executable and upgradeable:

```bash
tap upgrade
tap browser upgrade
tap doctor
```

`tap doctor --fix` delegates runtime repair to an already installed
`agent-browser doctor --fix`; it does not silently install a missing executable.
To update both tools to the versions shipped together, rerun the Tap installer
or install the latest full bundle.
