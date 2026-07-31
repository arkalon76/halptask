package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kenth/halptask/config"
	"github.com/kenth/halptask/model"
	"github.com/kenth/halptask/ui"
)

const Version = "0.1.0"

func main() {
	filePathFlag := flag.String("f", "", "Path to halptask data file")
	filePathLongFlag := flag.String("file", "", "Path to halptask data file")
	encryptFlag := flag.Bool("encrypt", false, "Force enable encryption")
	versionFlag := flag.Bool("version", false, "Print version")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("halptask v%s\n", Version)
		os.Exit(0)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
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
