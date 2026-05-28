package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ForegroundColor string `json:"foreground_color"`
	BackgroundColor string `json:"background_color"`
	Bold            bool   `json:"bold"`
	Italic          bool   `json:"italic"`
	Underline       bool   `json:"underline"`
	TabWidth        int    `json:"tab_width"`
	SideMargin      int    `json:"side_margin"`
}

func Default() *Config {
	return &Config{
		ForegroundColor: "default",
		BackgroundColor: "default",
		Bold:            false,
		Italic:          false,
		Underline:       false,
		TabWidth:        4,
		SideMargin:      2,
	}
}

func configDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "goread")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "goread")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			cfg := Default()
			_ = cfg.Save()
			return cfg, nil
		}
		return nil, err
	}

	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Save() error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}
