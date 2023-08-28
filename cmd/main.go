package main

import (
	"fmt"
	"log"

	"github.com/Roll-play/roll-play-backend/pkg/api"
	"github.com/Roll-play/roll-play-backend/pkg/config"
)


func main() {
	envVars, err := config.Config()

	if err != nil {
		log.Fatal(err)
	}
	
	connectionString :=
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", envVars["DB_HOST"], envVars["DB_USER"], envVars["DB_PASSWORD"], envVars["DB_NAME"])

	app, err := api.NewApp(connectionString)

	if err != nil {
		log.Fatal(err)
	}
	
	app.Server.Logger.Fatal(app.Server.Start(envVars["PORT"]))
}