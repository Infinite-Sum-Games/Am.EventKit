ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

DB_URL := "postgresql://postgres:1234@localhost:5432/postgres"
GO_BIN := $(shell go env GOPATH)/bin
GOOSE_DRIVER := postgres
GOOSE_DBSTRING := $(DB_URL)
GOOSE_MIGRATION_DIR := ./db/migrations/

ifeq ($(OS),Windows_NT)
	BIN_NAME := bin/anokha-backend.exe
else
	BIN_NAME := bin/anokha-backend
endif

setup:
	go install github.com/air-verse/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/evilmartians/lefthook@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@lefthook install
	@echo "All CLI tools installed successfully in $(GO_BIN)"

dev:
	@air

build:
	@go mod tidy
	@go fmt ./...
	@go build -o $(BIN_NAME) main.go

run: build
	@./$(BIN_NAME)

up:
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

upone:
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up-by-one

seed: build
	@go run seed/seed.go seed/truncate.go seed/main.go -s

down:
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) reset

downto:
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down-to $(v)

status:
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) status

clean:
	@go run seed/seed.go seed/truncate.go seed/main.go -c

doc:
	@docker compose up -d

pod:
	@podman compose down
	@podman compose up -d
