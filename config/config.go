package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/miun173/rss-reader/helper"
)

// Config :nodoc:
type Config struct {
	port         string
	cronInterval uint64
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}

// GetConfig :nodoc:
func GetConfig() Config {
	port := "8080"
	if val, ok := os.LookupEnv("PORT"); ok {
		port = val
	}

	cronInterval := 1
	if val, ok := os.LookupEnv("CRON_INTERVAL"); ok {
		cronInterval = helper.StringToInt(val)
	}

	return Config{
		port:         port,
		cronInterval: uint64(cronInterval),
	}
}

// Port :nodoc:
func (cfg *Config) Port() string {
	return cfg.port
}

// CronInterval :nodoc:
func (cfg *Config) CronInterval() uint64 {
	return cfg.cronInterval
}
