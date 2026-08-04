package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type TagConfig struct {
	Name  string `yaml:"name"`
	Emoji string `yaml:"emoji"`
	Color string `yaml:"color"` // LipGloss hex color string or name
}

type Config struct {
	AutoSave     bool        `yaml:"auto_save"`
	CheckUpdates bool        `yaml:"check_updates"`
	GithubRepo   string      `yaml:"github_repo"`
	DataFile     string      `yaml:"data_file"`
	Encrypted    bool        `yaml:"encrypted"`
	IndentSpaces int         `yaml:"indent_spaces"`
	LeaderKey    string      `yaml:"leader_key"`
	ShowWhichKey    bool        `yaml:"show_which_key"`
	Theme           string      `yaml:"theme"`
	DefaultItemType string      `yaml:"default_item_type"` // "bullet" or "task"
	Tags            []TagConfig `yaml:"tags,omitempty"`
}

func GetDefaultTagConfigs() []TagConfig {
	return []TagConfig{
		{Name: "bug", Emoji: "🐛", Color: "#f7768e"},     // Red
		{Name: "urgent", Emoji: "🔥", Color: "#ff9e64"},  // Orange
		{Name: "feature", Emoji: "✨", Color: "#7aa2f7"}, // Blue
		{Name: "work", Emoji: "💼", Color: "#bb9af7"},    // Purple
		{Name: "home", Emoji: "🏠", Color: "#9ece6a"},    // Green
		{Name: "idea", Emoji: "💡", Color: "#e0af68"},    // Yellow
		{Name: "pin", Emoji: "📌", Color: "#7dcfff"},     // Cyan
		{Name: "review", Emoji: "👀", Color: "#f7768e"},  // Pink/Red
	}
}

func DefaultConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	defaultData := filepath.Join(home, ".config", "halptask", "data.txt")
	return &Config{
		AutoSave:        true,
		CheckUpdates:    true,
		GithubRepo:      "arkalon76/halptask",
		DataFile:        defaultData,
		Encrypted:       false,
		IndentSpaces:    2,
		LeaderKey:       " ",
		ShowWhichKey:    true,
		Theme:           "default",
		DefaultItemType: "bullet",
		Tags:            GetDefaultTagConfigs(),
	}
}

func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "halptask")
}

func ConfigFilePath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return cfg, err
	}

	path := ConfigFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Save default config if not existing
			_ = SaveConfig(cfg)
			return cfg, nil
		}
		return cfg, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return cfg, err
	}
	if cfg.IndentSpaces <= 0 {
		cfg.IndentSpaces = 2
	}
	if cfg.LeaderKey == "" {
		cfg.LeaderKey = " "
	}
	if cfg.GithubRepo == "" {
		cfg.GithubRepo = "arkalon76/halptask"
	}
	if cfg.DefaultItemType == "" {
		cfg.DefaultItemType = "bullet"
	} else {
		if cfg.DefaultItemType == "bulletpoint" || cfg.DefaultItemType == "bullet_point" {
			cfg.DefaultItemType = "bullet"
		}
		if cfg.DefaultItemType != "bullet" && cfg.DefaultItemType != "task" {
			cfg.DefaultItemType = "bullet"
		}
	}
	if len(cfg.Tags) == 0 {
		cfg.Tags = GetDefaultTagConfigs()
	}
	return cfg, nil
}

func SaveConfig(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFilePath(), data, 0644)
}
