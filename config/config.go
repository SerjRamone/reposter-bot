// Package config ...
package config

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Posts contents range of posts id for repost
type Posts struct {
	From int `mapstructure:"from"`
	To   int `mapstructure:"to"`
}

// Delay contents range of min and max delay between reposts in minutes
type Delay struct {
	From int `mapstructure:"from"`
	To   int `mapstructure:"to"`
}

// Channel channel data for reposts
type Channel struct {
	ID          int   `mapstructure:"id"`
	Posts       Posts `mapstructure:"posts"`
	Delay       Delay `mapstructure:"delay"`
	PostedAt    time.Time
	IsInProcess bool
}

// Config is a config
type Config struct {
	BotToken   string    `mapstructure:"bot_token"`
	RetryDelay int       `mapstructure:"retry_delay"`
	Channels   []Channel `mapstructure:"channels"`
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment. Once.
func Get(pathToConfig string) *Config {
	once.Do(func() {
		viper.SetConfigFile(pathToConfig)
		viper.SetConfigType("yml")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("error reading env file", err)
		}

		// Viper unmarshals the loaded env varialbes into the struct
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal(err)
		}

		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Configuration:", string(configBytes))
	})

	return &config
}
