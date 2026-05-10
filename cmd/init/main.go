// cmd/init creates ~/.mangahub/ directory structure and writes a default config.yaml.
// Run once before anything else: go run ./cmd/init
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"mangahub/pkg/utils"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot find home dir:", err)
		os.Exit(1)
	}

	dirs := []string{
		filepath.Join(home, ".mangahub"),
		filepath.Join(home, ".mangahub", "logs"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			fmt.Fprintln(os.Stderr, "mkdir:", err)
			os.Exit(1)
		}
	}

	cfgPath := filepath.Join(home, ".mangahub", "config.yaml")
	if _, err := os.Stat(cfgPath); err == nil {
		fmt.Println("Config already exists:", cfgPath)
		fmt.Println("Use --force to overwrite (not implemented here, just delete and re-run).")
		return
	}

	cfg := utils.DefaultConfig()
	if err := utils.SaveConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "save config:", err)
		os.Exit(1)
	}

	fmt.Println("MangaHub initialized!")
	fmt.Println()
	fmt.Println("  Config :", cfgPath)
	fmt.Println("  Database:", cfg.Database.Path)
	fmt.Println("  Logs    :", cfg.Logging.Path)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Add manga data : go run ./cmd/seed")
	fmt.Println("  2. Start servers  : make run-all   (or run each separately)")
}
