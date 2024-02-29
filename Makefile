export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=''
export POSTGRES_DBNAME=postgres
export CONFIG_PATH=${PWD}/config.yaml

.PHONY: build run test compose-up test-run compose-down test-setup test-teardown run-go

test: test-setup test-run test-teardown

test-setup: compose-up

test-teardown: compose-down

test-run:
	go test ./...

build:
	go build .

run: compose-up run-go compose-down

run-go:
	go run .

compose-up:
	docker-compose -f docker-compose.yaml up -d

compose-down:
	docker-compose -f docker-compose.yaml down -v
