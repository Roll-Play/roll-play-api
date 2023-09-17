package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Roll-play/roll-play-backend/pkg/api"
	"github.com/Roll-play/roll-play-backend/pkg/config"
)

func main() {
	isDocker, err := strconv.ParseBool(os.Getenv("DOCKER"))

	if err != nil {
		isDocker = false
	}

	err = config.Config(isDocker, "./.env")

	if err != nil {
		log.Fatal(err)
	}

	connectionString :=
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_SLL_MODE"))

	app, err := api.NewApp(connectionString)

	if err != nil {
		log.Fatal(err)
	}

	app.Server.Logger.Fatal(app.Server.Start(os.Getenv("PORT")))
}
