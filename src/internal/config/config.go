package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string        `mapstructure:"database_url"`
	Port        string        `mapstructure:"port"`
	LogLevel    string        `mapstructure:"log_level"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	var cfg Config

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %s", err)
	}

	return &cfg, nil
}
