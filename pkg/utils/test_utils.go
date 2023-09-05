package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Roll-play/roll-play-backend/pkg/config"
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
			os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	fmt.Println(connectionString, "AAAAAAAAAAAAAAA", isDocker)
	db, err := sqlx.Open("pgx", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func ExecSchema(db *sqlx.DB, schema string) error {
	statements := strings.Split(schema, ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		_, err := db.Exec(statement)
		if err != nil {
			return fmt.Errorf("error executing SQL statement: %v", err)
		}
	}

	return nil
}

func TeardownTestDB(db *sqlx.DB) {
	// Close the test database connections
	if db != nil {
		db.Close()
	}
}
