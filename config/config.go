package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseUrl string
	RdbAddr string
	RdbPassword string
}

var (
	once sync.Once
	config *Config
)

func GetConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		config = &Config{
			BaseUrl: os.Getenv("BASE_URL"),
			RdbAddr: os.Getenv("RDB_ADDR"),
			RdbPassword: os.Getenv("RDB_PASSWORD"),
		}
	})

	return config
}