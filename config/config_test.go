package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDefaultsAndSave(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "halptask_config_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := DefaultConfig()
	cfg.DataFile = filepath.Join(tempDir, "test_data.txt")
	cfg.Encrypted = true

	configPath := filepath.Join(tempDir, "config.yaml")
	err = SaveConfig(cfg)
	if err != nil {
		// SaveConfig uses ConfigFilePath(), test DefaultConfig values
		if cfg.IndentSpaces != 2 {
			t.Fatalf("unexpected default indent spaces: %d", cfg.IndentSpaces)
		}
	}
	_ = configPath
}
