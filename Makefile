build:
	@go fmt ./...
	@go build -o bin/anokha-backend

run: build
	@./bin/anokha-backend

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
