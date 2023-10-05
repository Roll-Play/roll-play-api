package utils

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/Roll-play/roll-play-backend/pkg/config"
	"github.com/Roll-play/roll-play-backend/pkg/entities"
	repository "github.com/Roll-play/roll-play-backend/pkg/repositories"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
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

func SetupUser(db *sqlx.DB, username string, t *testing.T) (uuid.UUID, error) {
	userRepository := repository.NewUserRepository(db)

	savedUser, err := userRepository.Create(entities.User{
		Username: username,
		Email:    username + "@test",
		Password: "test",
	})

	if err != nil {
		return uuid.New(), err
	}

	return savedUser.Id, nil
}

func CreateSheetAndSheetUser(db *sqlx.DB, sheet *entities.Sheet, userId uuid.UUID, permission int, owner bool) error {
	err := db.Get(sheet, `INSERT INTO sheets (name, description, properties, background) VALUES ($1, $2, $3, $4)
	RETURNING id, name, description, properties, background`,
		sheet.Name, sheet.Description, sheet.Properties, sheet.Background)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO sheet_user (sheet_id, user_id, permission, owner) VALUES ($1, $2, $3, $4)`, sheet.Id, userId, permission, owner)
	return err
}
