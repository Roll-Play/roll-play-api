-include .env

migratecreate:
	migrate create -ext sql -dir migrations -seq $(name)

migrateup:
	migrate -path migrations/ -database ${DB_URL} -verbose up

migratedown:
	migrate -path migrations/ -database ${DB_URL} -verbose down 1

run:
	go run cmd/main.go

test:
	go test -v ./...

postgres-up:
	# Start a PostgreSQL container in detached mode with environment variables
	sudo docker run --name test-postgres -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -e POSTGRES_DB="${DB_NAME}_test" -p 5432:5432 -d postgres:latest

postgres-down:
	sudo docker stop test-postgres
	sudo docker rm test-postgres

coverage:
	go test -cover ./...

build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/app cmd/main.go


.PHONY: migratecreate migrateup migratedown run test coverage build postgres-up postgres-down
