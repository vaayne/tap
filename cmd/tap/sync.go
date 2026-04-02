package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	apiBaseURL = "https://tap.vaayne.com"
	syncAPI    = apiBaseURL + "/api/sync"
	searchAPI  = apiBaseURL + "/api/search"
	contentAPI = apiBaseURL + "/api/scripts"

	syncTTL      = 24 * time.Hour
	lastSyncFile = ".last_sync"
)

// defaultSitesDir returns the default cache directory for site scripts.
// Uses $XDG_CACHE_HOME/tap/sites, falling back to ~/.cache/tap/sites.
func defaultSitesDir() string {
	if dir := os.Getenv("XDG_CACHE_HOME"); dir != "" {
		return filepath.Join(dir, "tap", "sites")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".cache", "tap", "sites")
}

// syncManifestItem represents a single script in the remote manifest.
type syncManifestItem struct {
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	UpdatedAt string `json:"updatedAt"`
}

// syncManifest is the response from the /api/sync endpoint.
type syncManifest struct {
	Scripts []syncManifestItem `json:"scripts"`
}

// searchResult represents a script returned by the search API.
type searchResult struct {
	Name        string                  `json:"name"`
	Site        string                  `json:"site"`
	Description string                  `json:"description"`
	Domain      string                  `json:"domain"`
	ReadOnly    bool                    `json:"readOnly"`
	Example     string                  `json:"example"`
	Args        map[string]searchArgDef `json:"args"`
	UsageCount  int                     `json:"usageCount"`
}

type searchArgDef struct {
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

type searchResponse struct {
	Scripts []searchResult `json:"scripts"`
	Sites   []string       `json:"sites"`
}

// ensureScripts checks if the sites directory has any scripts and if the
// cache is fresh. Syncs if no scripts exist or if the last sync is older
// than syncTTL (24h).
func ensureScripts(dir string, verbose bool) error {
	if hasScripts(dir) && !isSyncStale(dir) {
		return nil
	}
	if !hasScripts(dir) {
		if verbose {
			log.Printf("No local scripts found, syncing from %s...", apiBaseURL)
		} else {
			fmt.Fprintf(os.Stderr, "Syncing scripts from %s...\n", apiBaseURL)
		}
	} else {
		if verbose {
			log.Printf("Cache expired, syncing from %s...", apiBaseURL)
		} else {
			fmt.Fprintf(os.Stderr, "Updating scripts from %s...\n", apiBaseURL)
		}
	}
	return syncScripts(dir, verbose)
}

// isSyncStale returns true if the last sync was more than syncTTL ago
// or if the timestamp file doesn't exist.
func isSyncStale(dir string) bool {
	path := filepath.Join(dir, lastSyncFile)
	info, err := os.Stat(path)
	if err != nil {
		return true
	}
	return time.Since(info.ModTime()) > syncTTL
}

// touchLastSync writes/updates the .last_sync file.
func touchLastSync(dir string) {
	path := filepath.Join(dir, lastSyncFile)
	_ = os.MkdirAll(dir, 0o755)
	_ = os.WriteFile(path, []byte(time.Now().Format(time.RFC3339)), 0o644)
}

// hasScripts returns true if the directory contains at least one .js file.
func hasScripts(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			subEntries, err := os.ReadDir(filepath.Join(dir, e.Name()))
			if err != nil {
				continue
			}
			for _, se := range subEntries {
				if !se.IsDir() && strings.HasSuffix(se.Name(), ".js") {
					return true
				}
			}
		}
	}
	return false
}

// syncScripts fetches the remote manifest and downloads new/changed scripts.
// It also removes scripts that no longer exist remotely.
func syncScripts(dir string, verbose bool) error {
	client := &http.Client{Timeout: 30 * time.Second}

	// 1. Fetch remote manifest
	if verbose {
		log.Printf("Fetching manifest from %s", syncAPI)
	}
	resp, err := client.Get(syncAPI)
	if err != nil {
		return fmt.Errorf("fetch manifest: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch manifest: HTTP %d", resp.StatusCode)
	}

	var manifest syncManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return fmt.Errorf("decode manifest: %w", err)
	}

	// 2. Build local hash index
	localHashes := buildLocalHashes(dir)

	// 3. Determine what to download and what to delete
	remoteNames := make(map[string]bool, len(manifest.Scripts))
	var toDownload []syncManifestItem
	for _, s := range manifest.Scripts {
		remoteNames[s.Name] = true
		if localHash, ok := localHashes[s.Name]; !ok || localHash != s.Hash {
			toDownload = append(toDownload, s)
		}
	}

	// Find local scripts to delete
	var toDelete []string
	for name := range localHashes {
		if !remoteNames[name] {
			toDelete = append(toDelete, name)
		}
	}

	if len(toDownload) == 0 && len(toDelete) == 0 {
		fmt.Fprintf(os.Stderr, "Already up to date. (%d scripts)\n", len(manifest.Scripts))
		return nil
	}

	// 4. Download new/changed scripts
	downloaded := 0
	for _, s := range toDownload {
		if err := downloadScript(client, dir, s.Name); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to download %s: %v\n", s.Name, err)
			continue
		}
		downloaded++
		if verbose {
			log.Printf("Downloaded: %s", s.Name)
		}
	}

	// 5. Delete removed scripts
	deleted := 0
	for _, name := range toDelete {
		path := scriptPath(dir, name)
		if err := os.Remove(path); err == nil {
			deleted++
			if verbose {
				log.Printf("Deleted: %s", name)
			}
			// Clean up empty parent directory
			parentDir := filepath.Dir(path)
			entries, _ := os.ReadDir(parentDir)
			if len(entries) == 0 {
				_ = os.Remove(parentDir)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Synced: %d downloaded, %d deleted, %d total\n",
		downloaded, deleted, len(manifest.Scripts))

	touchLastSync(dir)
	return nil
}

// downloadScript fetches a single script's content and writes it to disk.
func downloadScript(client *http.Client, dir, name string) error {
	url := fmt.Sprintf("%s/%s/content", contentAPI, name)
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("fetch %s: %w", name, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch %s: HTTP %d", name, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read %s: %w", name, err)
	}

	path := scriptPath(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}

	return os.WriteFile(path, body, 0o644)
}

// scriptPath converts a script name like "google/search" to a file path.
func scriptPath(dir, name string) string {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) == 2 {
		return filepath.Join(dir, parts[0], parts[1]+".js")
	}
	return filepath.Join(dir, name+".js")
}

// buildLocalHashes scans the sites directory and computes SHA-256 hashes.
func buildLocalHashes(dir string) map[string]string {
	hashes := make(map[string]string)
	sites, err := os.ReadDir(dir)
	if err != nil {
		return hashes
	}
	for _, site := range sites {
		if !site.IsDir() {
			continue
		}
		sitePath := filepath.Join(dir, site.Name())
		files, err := os.ReadDir(sitePath)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".js") {
				continue
			}
			content, err := os.ReadFile(filepath.Join(sitePath, f.Name()))
			if err != nil {
				continue
			}
			h := sha256.Sum256(content)
			name := site.Name() + "/" + strings.TrimSuffix(f.Name(), ".js")
			hashes[name] = hex.EncodeToString(h[:])
		}
	}
	return hashes
}

// searchOnline queries the remote search API.
func searchOnline(query string) (*searchResponse, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	url := fmt.Sprintf("%s?q=%s", searchAPI, query)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search: HTTP %d", resp.StatusCode)
	}

	var result searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode search: %w", err)
	}
	return &result, nil
}
