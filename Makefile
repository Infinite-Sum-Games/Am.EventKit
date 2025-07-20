GO_BIN := $(shell go env GOPATH)/bin

setup:
	go install github.com/air-verse/air@latest

	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

	go install github.com/evilmartians/lefthook@latest

	go install github.com/pressly/goose/v3/cmd/goose@latest

	go mod tidy

	@lefthook install

	@echo "All CLI tools installed successfully in $(GO_BIN)"

build:
	@go fmt ./...
	@go build -o bin/anokha-backend

run: build
	@./bin/anokha-backend

test:
	@go test -v ./...

up:
	@goose -dir ./db/migrations/ -no-versioning up

seed: build
	@go run ./seed/seeder.go

down:
	@goose -dir ./db/migrations/ -no-versioning down

truncate: build
	go run ./truncate/truncate.go

# For docker users
doc:
	@docker compose up -d

# For podman users
pod:
	@podman compose down
	@podman compose up -d
