package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Configuration struct {
	DatabaseFileName string
}

var Config Configuration

func LoadEnv() error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	err = godotenv.Load(fmt.Sprintf("%s/config/env/local.env", cwd))

	if err != nil {
		return err
	}

	Config = Configuration{
		DatabaseFileName: os.Getenv("DB_FILE_NAME"),
	}

	return nil
}

func FetchConfig() Configuration {
	return Config
}
