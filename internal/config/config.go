package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address string `env:"HTTP_ADDRESS" env-default:"localhost:8082"`
}

type Config struct {
	ENV    string `env:"ENV" env-default:"dev"`
	DBPath string `env:"DB_PATH" env-default:"sqlite/dev"`
	DBName string `env:"DB_NAME" env-default:"api.db"`
	HTTPServer
	JWTKey string `env:"JWT_KEY" env-default:"supersecretjwtbdckey"`
}

func LoadConfig() *Config {
	var cfg Config
	var envPath string

	flag.StringVar(&envPath, "config", "", "Path to .env file")
	flag.Parse()

	//this is if after deploying server its difficult to access config file, check if we loaded config file where we deployed server. ex. aws, cloud vps, etc...
	if envPath == "" {
		envPath = os.Getenv("CONFIG_PATH")
	}

	//if its not on cloud hosting site, hardcode it to file in project
	if envPath == "" {
		envPath = "./config/env"
	}

	err := cleanenv.ReadConfig(envPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read .env file from %s: %v", envPath, err)
	}

	return &cfg
}
