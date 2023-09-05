package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Config(docker bool, envPath string) error {

	// ugly hack to accommodate config with .env files
	// or docker compose
	if !docker {
		if err := godotenv.Load(envPath); err != nil {
			return err
		}
	}

	port := os.Getenv("PORT")

	if []byte(port)[0] == ':' {
		return nil
	}

	os.Setenv("PORT", string(append([]byte(":"), []byte(port)...)))

	return nil
}
