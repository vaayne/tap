package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/vaayne/tap"
	"gopkg.in/yaml.v3"
)

const (
	embeddedSkillRoot   = tap.EmbeddedTapWebSkillRoot
	embeddedSkillConfig = tap.EmbeddedTapWebSkillConfig
)

type skillMetadata struct {
	Author  string `yaml:"author"`
	Version string `yaml:"version"`
}

type skillConfig struct {
	Metadata skillMetadata `yaml:"metadata"`
}

func skillCmd() *cli.Command {
	return &cli.Command{
		Name:  "skill",
		Usage: "Manage the embedded tap-web skill",
		Description: `The tap-web skill provides AI agent capabilities for web access.
It is embedded in the tap binary and can be installed to the skills directory.`,
		Commands: []*cli.Command{
			{
				Name:   "install",
				Usage:  "Install or update the embedded skill",
				Action: skillInstallAction,
				Flags: append([]cli.Flag{
					&cli.BoolFlag{
						Name:  "force",
						Usage: "Reinstall even if already up to date",
					},
				}, skillPathFlags("Custom installation directory (default: ~/.config/tap/skills/tap-web/)")...),
			},
			{
				Name:   "version",
				Usage:  "Show embedded skill version",
				Action: skillVersionAction,
				Flags:  skillPathFlags("Custom installation directory to check"),
			},
			{
				Name:   "path",
				Usage:  "Show skill installation path",
				Action: skillPathAction,
				Flags:  skillPathFlags("Custom installation directory"),
			},
		},
	}
}

func skillPathFlags(usage string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"dir"},
			Usage:   usage,
			Sources: cli.EnvVars("TAP_SKILL_DIR"),
		},
	}
}

// skillInstallAction extracts the embedded skill to the skills directory
func skillInstallAction(ctx context.Context, cmd *cli.Command) error {
	force := cmd.Bool("force")
	skillDir := resolveSkillDir(cmd)
	targetVersion := targetSkillVersion()

	if !force {
		if installed, current := isSkillInstalled(skillDir); installed {
			if current == targetVersion {
				fmt.Printf("Skill already up to date (%s) at %s\n", current, skillDir)
				return nil
			}
			fmt.Printf("Updating skill: %s -> %s\n", current, targetVersion)
		}
	}

	if err := extractEmbeddedSkill(skillDir, cmd.Root().Bool("verbose")); err != nil {
		return fmt.Errorf("failed to extract skill: %w", err)
	}

	fmt.Printf("✓ Skill installed (%s) at %s\n", getEmbeddedSkillVersion(), skillDir)
	return nil
}

// skillVersionAction shows the installed skill version relative to the current tap CLI version.
func skillVersionAction(ctx context.Context, cmd *cli.Command) error {
	cliVersion := currentCLIVersion()
	targetVersion := targetSkillVersion()
	skillDir := resolveSkillDir(cmd)

	fmt.Printf("tap CLI version: %s\n", cliVersion)
	if targetVersion != cliVersion {
		fmt.Printf("Bundled skill version: %s\n", targetVersion)
	}

	if installed, installedVersion := isSkillInstalled(skillDir); installed {
		fmt.Printf("Installed skill version: %s\n", installedVersion)
		if installedVersion == targetVersion {
			fmt.Println("Status: in sync")
		} else {
			fmt.Println("Status: mismatch (run 'tap skill install')")
		}
	} else {
		fmt.Println("Installed skill version: not installed")
		fmt.Println("Status: missing (run 'tap skill install')")
	}
	return nil
}

// skillPathAction shows the skill installation path
func skillPathAction(ctx context.Context, cmd *cli.Command) error {
	skillDir := resolveSkillDir(cmd)
	fmt.Println(skillDir)
	return nil
}

// defaultSkillDir returns the default skill installation directory
func defaultSkillDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "tap", "skills", "tap-web")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "tap", "skills", "tap-web")
}

// resolveSkillDir returns the skill directory, using custom path if provided
func resolveSkillDir(cmd *cli.Command) string {
	if customPath := cmd.String("path"); customPath != "" {
		return customPath
	}
	return defaultSkillDir()
}

func getEmbeddedSkillVersion() string {
	content, err := fs.ReadFile(tap.EmbeddedTapWebSkillFS(), embeddedSkillConfig)
	if err != nil {
		return "unknown"
	}
	return parseSkillVersion(content)
}

func currentCLIVersion() string {
	if version == "" {
		return "unknown"
	}
	if version == "dev" || strings.HasPrefix(version, "v") {
		return version
	}
	return "v" + version
}

func targetSkillVersion() string {
	cliVersion := currentCLIVersion()
	if cliVersion == "dev" || cliVersion == "unknown" {
		return getEmbeddedSkillVersion()
	}
	return cliVersion
}

func isSkillInstalled(skillDir string) (bool, string) {
	content, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if err != nil {
		return false, ""
	}
	return true, parseSkillVersion(content)
}

func parseSkillVersion(content []byte) string {
	frontmatter, ok := extractFrontmatter(string(content))
	if !ok {
		return "unknown"
	}

	var config skillConfig
	if err := yaml.Unmarshal([]byte(frontmatter), &config); err != nil {
		return "unknown"
	}
	if config.Metadata.Version == "" {
		return "unknown"
	}
	return config.Metadata.Version
}

func extractFrontmatter(content string) (string, bool) {
	if !strings.HasPrefix(content, "---") {
		return "", false
	}

	end := strings.Index(content[3:], "---")
	if end == -1 {
		return "", false
	}

	return content[3 : end+3], true
}

// extractEmbeddedSkill extracts all embedded skill files to the destination
func extractEmbeddedSkill(destDir string, verbose bool) error {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("create skill directory: %w", err)
	}

	embeddedFS := tap.EmbeddedTapWebSkillFS()

	return fs.WalkDir(embeddedFS, embeddedSkillRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(embeddedSkillRoot, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0o755); err != nil {
				return fmt.Errorf("create directory %s: %w", destPath, err)
			}
			return nil
		}

		// Copy file
		content, err := fs.ReadFile(embeddedFS, path)
		if err != nil {
			return fmt.Errorf("read embedded file %s: %w", path, err)
		}

		if err := os.WriteFile(destPath, content, 0o644); err != nil {
			return fmt.Errorf("write file %s: %w", destPath, err)
		}

		if verbose {
			log.Printf("Extracted: %s", relPath)
		}

		return nil
	})
}

func installEmbeddedSkillIfNeeded(skillDir string) error {
	if installed, _ := isSkillInstalled(skillDir); installed {
		return nil
	}
	return extractEmbeddedSkill(skillDir, false)
}

func autoInstallEmbeddedSkill() {
	skillDir := os.Getenv("TAP_SKILL_DIR")
	if skillDir == "" {
		skillDir = defaultSkillDir()
	}
	if err := installEmbeddedSkillIfNeeded(skillDir); err != nil {
		// Silently ignore; user can manually install with `tap skill install`.
		_ = err
	}
}
