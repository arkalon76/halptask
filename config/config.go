package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AutoSave     bool   `yaml:"auto_save"`
	DataFile     string `yaml:"data_file"`
	Encrypted    bool   `yaml:"encrypted"`
	IndentSpaces int    `yaml:"indent_spaces"`
	LeaderKey    string `yaml:"leader_key"`
	ShowWhichKey bool   `yaml:"show_which_key"`
	Theme        string `yaml:"theme"`
}

func DefaultConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	defaultData := filepath.Join(home, ".config", "halptask", "data.txt")
	return &Config{
		AutoSave:     true,
		DataFile:     defaultData,
		Encrypted:    false,
		IndentSpaces: 2,
		LeaderKey:    " ",
		ShowWhichKey: true,
		Theme:        "default",
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
