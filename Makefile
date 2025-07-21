ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

GO_BIN := $(shell go env GOPATH)/bin
GOOSE_DRIVER := postgres
GOOSE_DBSTRING := $(DB_URL)
GOOSE_MIGRATION_DIR := ./db/migrations/

# TODO: Right now, we are writing down the name of every package that needs to 
# be tested explicitly. However, there should be a way to ignore these:
# db/gen, tests/, seed/ via some pattern-matching regex
TEST_PACKAGES := \
	github.com/Thanus-Kumaar/anokha-2025-backend \
	github.com/Thanus-Kumaar/anokha-2025-backend/api/auth \
	github.com/Thanus-Kumaar/anokha-2025-backend/api/event \
	github.com/Thanus-Kumaar/anokha-2025-backend/api/profile \
	github.com/Thanus-Kumaar/anokha-2025-backend/api/staff \
	github.com/Thanus-Kumaar/anokha-2025-backend/cmd \
	github.com/Thanus-Kumaar/anokha-2025-backend/models \
	github.com/Thanus-Kumaar/anokha-2025-backend/mail \
	github.com/Thanus-Kumaar/anokha-2025-backend/middleware \
	github.com/Thanus-Kumaar/anokha-2025-backend/pkg

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
	# TODO: The following tool needs to be downloaded and installed into GO_BIN
	# curl https://gotest-release.s3.amazonaws.com/gotest_linux > gotest && chmod +x gotest
	@echo "All CLI tools installed successfully in $(GO_BIN)"

dev:
	@air

build:
	@go mod tidy
	@go test $(TEST_PACKAGES)
	@go fmt ./...
	@go build -o $(BIN_NAME) main.go

run: build
	@./$(BIN_NAME)

# Requires the gotest tool for colored outputs
test:
	@gotest $(TEST_PACKAGES)

up:
	@goose -dir $(GOOSE_MIGRATION_DIR) -no-versioning $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

seed: build
	@go run seed/seed.go seed/truncate.go seed/main.go -s

down:
	@goose -dir $(GOOSE_MIGRATION_DIR) -no-versioning $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down

clean:
	@go run seed/seed.go seed/truncate.go seed/main.go -c

# For docker users
doc:
	@docker compose up -d

# For podman users
pod:
	@podman compose down
	@podman compose up -d
