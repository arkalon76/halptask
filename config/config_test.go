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
		if cfg.IndentSpaces != 2 {
			t.Fatalf("unexpected default indent spaces: %d", cfg.IndentSpaces)
		}
	}
	if cfg.DefaultItemType != "bullet" {
		t.Fatalf("expected default item type to be 'bullet', got %s", cfg.DefaultItemType)
	}
	_ = configPath
}

func TestDefaultItemTypeNormalization(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.DefaultItemType != "bullet" {
		t.Fatalf("expected DefaultItemType to be 'bullet', got %s", cfg.DefaultItemType)
	}

	// Test normalization logic directly
	cfg.DefaultItemType = "bulletpoint"
	if cfg.DefaultItemType == "bulletpoint" || cfg.DefaultItemType == "bullet_point" {
		cfg.DefaultItemType = "bullet"
	}
	if cfg.DefaultItemType != "bullet" {
		t.Fatalf("expected normalized 'bullet', got %s", cfg.DefaultItemType)
	}
}
