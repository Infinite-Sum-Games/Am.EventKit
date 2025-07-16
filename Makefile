GO_BIN := $(shell go env GOPATH)/bin

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
	@goose -dir ./db/migrations/ -no-versioning up

seed:
	@goose -dir ./db/seed/ -no-versioning up

down:
	@goose -dir ./db/migrations/ -no-versioning down

# For docker users
doc:
	@docker compose up -d

# For podman users
pod:
	@podman compose down
	@podman compose up -d
