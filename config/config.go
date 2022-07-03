package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Configuration struct {
	DatabaseFileName string
	Port             string
}

var Config Configuration

//LoadEnv : Load The Environment Variables Needed
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
		Port:             os.Getenv("PORT"),
	}

	return nil
}

func FetchConfig() Configuration {
	return Config
}
