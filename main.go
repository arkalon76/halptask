package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kenth/halptask/config"
	"github.com/kenth/halptask/model"
	"github.com/kenth/halptask/ui"
	"github.com/kenth/halptask/updater"
)

const Version = "0.0.5"

func main() {
	filePathFlag := flag.String("f", "", "Path to halptask data file")
	filePathLongFlag := flag.String("file", "", "Path to halptask data file")
	encryptFlag := flag.Bool("encrypt", false, "Force enable encryption")
	versionFlag := flag.Bool("version", false, "Print version")
	updateFlag := flag.Bool("update", false, "Check for updates and perform auto-update if a new version is available")
	checkUpdateFlag := flag.Bool("check-update", false, "Check if a new version is available")
	repoFlag := flag.String("repo", "", "Override target GitHub repository (e.g. owner/repo)")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("halptask v%s\n", Version)
		os.Exit(0)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
	}

	targetRepo := cfg.GithubRepo
	if *repoFlag != "" {
		targetRepo = *repoFlag
	}

	if *checkUpdateFlag || *updateFlag {
		fmt.Printf("Checking for updates from %s...\n", targetRepo)
		rel, err := updater.CheckForUpdate(Version, targetRepo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
			os.Exit(1)
		}

		if rel.NewRepo != "" && rel.NewRepo != cfg.GithubRepo {
			fmt.Printf("Repository has moved to %s. Updating config...\n", rel.NewRepo)
			cfg.GithubRepo = rel.NewRepo
			_ = config.SaveConfig(cfg)
		}

		if !rel.HasUpdate {
			fmt.Printf("halptask v%s is up to date! (Latest: v%s)\n", Version, rel.Version)
			os.Exit(0)
		}

		fmt.Printf("A new version of halptask is available: v%s (current: v%s)\n", rel.Version, Version)

		if *updateFlag {
			canUpdate, realPath, reason := updater.CanUpdate()
			if !canUpdate {
				fmt.Fprintf(os.Stderr, "Cannot update executable at %s: %s\n", realPath, reason)
				os.Exit(1)
			}

			fmt.Printf("Updating halptask binary at %s...\n", realPath)
			if err := updater.DoUpdate(rel); err != nil {
				fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully updated halptask to v%s!\n", rel.Version)
			os.Exit(0)
		}
		os.Exit(0)
	}

	targetFilePath := cfg.DataFile
	if *filePathFlag != "" {
		targetFilePath = *filePathFlag
	} else if *filePathLongFlag != "" {
		targetFilePath = *filePathLongFlag
	}

	isEncrypted := cfg.Encrypted || *encryptFlag

	storage := model.NewStorage(targetFilePath, isEncrypted)

	appModel, initCmd := ui.InitialModel(cfg, storage)

	p := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if initCmd != nil {
		go func() {
			// Exec initial command if any
		}()
	}

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running halptask: %v\n", err)
		os.Exit(1)
	}
}
