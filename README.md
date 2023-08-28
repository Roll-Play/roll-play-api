# Roll & Play backend

## Minimal requirements to run this application:

- docker compose
- GNU Make
- go-migrate:

  ```bash
  $ sudo apt-get update
  $ curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
  $ sudo apt-get install migrate
  ```

## Running the application with docker compose:

- create and fill the .env file using .env.example as a guide
- run `$ sudo docker compose up --build`

## Running the application manually:

- create and fill the .env file using .env.example as a guide
- spin up a local instance of postgres
- run migrations with `make migrateup`
- run go run ./cmd/main.go

## Migrations:

### Creating migrations

- `make migrate create name=migration_name`

### Running migrations

- `make migrateup`

### Undoing migrations

- `make migratedown`

## Project structure

### cmd directory

- The `cmd` directory contains the main applications for the project.

### pkg directory

- The pkg directory is used to store shared packages or libraries that can be imported and used across different parts of the project. These packages are meant to be reusable components for the project.
