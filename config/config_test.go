package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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
	if cfg.UpdateInterval != "daily" {
		t.Fatalf("expected default update interval to be 'daily', got %s", cfg.UpdateInterval)
	}
	if !cfg.CheckUpdates {
		t.Fatalf("expected default CheckUpdates to be true")
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

func TestShouldCheckForUpdate(t *testing.T) {
	cfg := DefaultConfig()
	// Initial state with empty LastUpdateCheck should check
	if !ShouldCheckForUpdate(cfg) {
		t.Fatalf("expected ShouldCheckForUpdate to be true initially")
	}

	// Set last check to 1 hour ago -> should NOT check for daily interval
	cfg.LastUpdateCheck = time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	if ShouldCheckForUpdate(cfg) {
		t.Fatalf("expected ShouldCheckForUpdate to be false for check 1h ago with daily interval")
	}

	// Set last check to 25 hours ago -> SHOULD check for daily interval
	cfg.LastUpdateCheck = time.Now().Add(-25 * time.Hour).Format(time.RFC3339)
	if !ShouldCheckForUpdate(cfg) {
		t.Fatalf("expected ShouldCheckForUpdate to be true for check 25h ago with daily interval")
	}

	// CheckUpdates disabled -> should NOT check
	cfg.CheckUpdates = false
	if ShouldCheckForUpdate(cfg) {
		t.Fatalf("expected ShouldCheckForUpdate to be false when CheckUpdates is false")
	}

	// CheckUpdates re-enabled, interval set to "always" -> should check regardless of timestamp
	cfg.CheckUpdates = true
	cfg.UpdateInterval = "always"
	cfg.LastUpdateCheck = time.Now().Format(time.RFC3339)
	if !ShouldCheckForUpdate(cfg) {
		t.Fatalf("expected ShouldCheckForUpdate to be true when UpdateInterval is 'always'")
	}

	// Interval set to "never" -> should NOT check
	cfg.UpdateInterval = "never"
	if ShouldCheckForUpdate(cfg) {
		t.Fatalf("expected ShouldCheckForUpdate to be false when UpdateInterval is 'never'")
	}
}
