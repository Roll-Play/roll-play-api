package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var keys = []string{	
	"PORT",
	"DB_HOST",
	"DB_USER",
	"DB_PASSWORD",
	"DB_NAME",
}

func Config() (map[string]string, error) {
	isDocker, err := strconv.ParseBool(os.Getenv("DOCKER"))

	if err != nil {
		isDocker = false
	}

	// ugly hack to accommodate config with .env files
	// or docker compose
	if !isDocker {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	envVars := make(map[string]string, len(keys))

	for _, key := range keys {
		if key == "PORT" {
			port := os.Getenv(key)

			if []byte(port)[0] == ':' {
				envVars[key] = port
				continue	
			}

			envVars[key] = string(append([]byte(":"), []byte(port)...))
			continue
		}
		envVars[key] = os.Getenv(key)	
	}

	return envVars, nil
}