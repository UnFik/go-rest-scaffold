package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// NewViper is a function to load config from config.json
// You can change the implementation, for example load from env file, consul, etcd, etc
func NewViper() *viper.Viper {
	err := godotenv.Load()
	if err != nil {
		// handle error if .env file not found or other error,
		// but maybe just log it or ignore if we expect env vars from system
		// for now we can ignore as it might be prod
	}

	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")
	err = config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	// Enable reading from environment variables
	config.AutomaticEnv()

	// Bind specific env vars to config keys
	config.BindEnv("web.port", "APP_PORT")
	config.BindEnv("app.name", "APP_NAME")

	// Set defaults (fallback jika env tidak ada dan config.json tidak ada)
	config.SetDefault("web.port", 3000)

	return config
}
