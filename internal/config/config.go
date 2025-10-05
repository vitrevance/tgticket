package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TelegramBotToken string `yaml:"telegram_bot_token"`
	TelegramChatID   int64  `yaml:"telegram_chat_id"`
	AdminUser        string `yaml:"admin_user"`
	AdminPassword    string `yaml:"admin_password"`
	ServerAddr       string `yaml:"server_addr"`
	PublicAddr       string `yaml:"public_addr"`
}

func LoadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	if cfg.ServerAddr == "" {
		cfg.ServerAddr = ":8080"
	}
	return cfg, nil
}
