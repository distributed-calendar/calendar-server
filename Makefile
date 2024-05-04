export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=''
export POSTGRES_DBNAME=postgres
export CONFIG_PATH=${PWD}/config.yaml

.PHONY: build run test compose-up test-run compose-down test-setup test-teardown run-go

test: test-setup test-run test-teardown

test-setup: compose-up migrate-up

test-teardown: migrate-down compose-down

migrate-up:
	migrate -path migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DBNAME}?sslmode=disable up

migrate-down:
	migrate -path migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DBNAME}?sslmode=disable down -all

test-run:
	go test ./...

build:
	go build .

run: compose-up migrate-up run-go migrate-down compose-down

run-go:
	go run .

compose-up:
	docker-compose -f docker-compose.yaml up -d && sleep 2

compose-down:
	docker-compose -f docker-compose.yaml down -v
