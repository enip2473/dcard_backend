package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL" yaml:"database_url,omitempty"`
	Port        string `mapstructure:"PORT" yaml:"port,omitempty"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	RedisAddr   string `mapstructure:"REDIS_ADDR" yaml:"redis_addr,omitempty"`
	RedisPass   string `mapstructure:"REDIS_PASS" yaml:"redis_pass,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.BindEnv("DATABASE_URL"); err == nil {
		viper.SetDefault("DATABASE_URL", "")
	}

	if err := viper.BindEnv("PORT"); err == nil {
		viper.SetDefault("PORT", 4444)
	}

	if err := viper.BindEnv("REDIS_ADDR"); err == nil {
		viper.SetDefault("REDIS_ADDR", "")
	}

	if err := viper.BindEnv("REDIS_PASS"); err == nil {
		viper.SetDefault("REDIS_PASS", "")
	}

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: No configuration file found or could not be read. Relying on environment variables. %s\n", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %s", err)
	}

	fmt.Printf("Using configuration:\n%+v\n", cfg)
	return &cfg, nil
}
