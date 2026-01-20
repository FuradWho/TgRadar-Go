package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Telegram struct {
		AppID       int    `mapstructure:"app_id"`
		AppHash     string `mapstructure:"app_hash"`
		SessionFile string `mapstructure:"session_file"`
		Phone       string `mapstructure:"phone"`
		Password    string `mapstructure:"password"`
		Proxy       string `mapstructure:"proxy"`
	} `mapstructure:"telegram"`

	Monitor struct {
		WindowSeconds int  `mapstructure:"window_seconds"`
		Debug         bool `mapstructure:"debug"`
	} `mapstructure:"monitor"`

	AI struct {
		APIKey   string `mapstructure:"api_key"`
		BaseURL  string `mapstructure:"base_url"`
		Model    string `mapstructure:"model"`
		Language string `mapstructure:"language"`
	} `mapstructure:"ai"`
}

// LoadConfig parses config.yml from current directory
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // Config filename (without extension)
	viper.SetConfigType("yml")    // Config file type
	viper.AddConfigPath(".")      // Search path: current directory

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w\nensure config.yaml exists in current directory", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Simple validation
	if cfg.Telegram.AppID == 0 || cfg.Telegram.AppHash == "" {
		return nil, fmt.Errorf("config error: app_id and app_hash cannot be empty")
	}

	return &cfg, nil
}
