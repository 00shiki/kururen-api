include .env

migration:
	psql -d "dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT}" -f ./migrations/ddl.sql
	psql -d "dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT}" -f ./migrations/dml.sql

db:
	psql -d "dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT}"

build:
	go build -o build/app main.go

clean:
	rm -rf ./build

run: build
	./build/app