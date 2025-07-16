ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

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
	@go test -v ./...
	@go fmt ./...
	@go build -o $(BIN_NAME) main.go

run: build
	@./$(BIN_NAME)

test:
	@go test -v ./...

up:
	@goose -dir $(GOOSE_MIGRATION_DIR) -no-versioning $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

seed:
	@goose -dir ./db/seed/ -no-versioning $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

down:
	@goose -dir $(GOOSE_MIGRATION_DIR) -no-versioning $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down

# For docker users
doc:
	@docker compose up -d

# For podman users
pod:
	@podman compose down
	@podman compose up -d
