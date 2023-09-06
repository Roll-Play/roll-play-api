package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Roll-play/roll-play-backend/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func SetupTestDB(envPath string) (*sqlx.DB, error) {
	isDocker, err := strconv.ParseBool(os.Getenv("DOCKER"))

	if err != nil {
		isDocker = false
	}

	err = config.Config(isDocker, envPath)

	if err != nil {
		return nil, err
	}

	connectionString :=
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")+"_test")

	db, err := sqlx.Open("pgx", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(migrationSource string) error {
	m, err := migrate.New(migrationSource, os.Getenv("DB_URL"))

	if err != nil {
		return err
	}

	m.Down()

	err = m.Up()

	if err != nil {
		return err
	}

	return nil
}

func TeardownTestDB(db *sqlx.DB) {
	// Close the test database connections
	if db != nil {
		db.Close()
	}
}
