package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v3"
)

func upgradeCmd() *cli.Command {
	return &cli.Command{
		Name:  "upgrade",
		Usage: "Upgrade tap to the latest version",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Upgrade even if already on the latest version",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runUpgrade(ctx, cmd.Bool("force"))
		},
	}
}

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func runUpgrade(ctx context.Context, force bool) error {
	current := version
	fmt.Printf("Current version: %s\n", current)

	// Fetch latest release
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/vaayne/tap/releases/latest", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch latest release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var release ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("decode release: %w", err)
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	fmt.Printf("Latest version:  %s\n", latest)

	if !force && current == latest {
		fmt.Println("Already up to date.")
		return nil
	}

	// Find matching asset
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	wantName := fmt.Sprintf("tap_%s_%s_%s.tar.gz", latest, goos, goarch)
	if goos == "windows" {
		wantName = fmt.Sprintf("tap_%s_%s_%s.zip", latest, goos, goarch)
	}

	var downloadURL string
	for _, a := range release.Assets {
		if a.Name == wantName {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no release asset found for %s/%s (%s)", goos, goarch, wantName)
	}

	// Download
	fmt.Printf("Downloading %s...\n", wantName)
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("create download request: %w", err)
	}
	dlResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer func() { _ = dlResp.Body.Close() }()

	if dlResp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %s", dlResp.Status)
	}

	// Extract binary from tar.gz
	newBinary, err := extractTapBinary(dlResp.Body)
	if err != nil {
		return fmt.Errorf("extract binary: %w", err)
	}

	// Replace current binary
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate current binary: %w", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return fmt.Errorf("resolve symlink: %w", err)
	}

	// Write to temp file next to the binary, then rename (atomic on same fs)
	tmpPath := exe + ".new"
	if err := os.WriteFile(tmpPath, newBinary, 0o755); err != nil {
		return fmt.Errorf("write new binary: %w", err)
	}

	if err := os.Rename(tmpPath, exe); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("replace binary: %w", err)
	}

	fmt.Printf("Upgraded tap: %s -> %s\n", current, latest)
	fmt.Println("Run 'tap skill install' to update your local skill installation.")
	return nil
}

func extractTapBinary(r io.Reader) ([]byte, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("gzip: %w", err)
	}
	defer func() { _ = gz.Close() }()

	tr := tar.NewReader(gz)
	binName := "tap"
	if runtime.GOOS == "windows" {
		binName = "tap.exe"
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("binary %q not found in archive", binName)
		}
		if err != nil {
			return nil, err
		}
		if filepath.Base(hdr.Name) == binName {
			return io.ReadAll(tr)
		}
	}
}
