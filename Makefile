.PHONY: all clean test
include .env
export $(shell sed 's/=.*//' .env)

init:
	uvx pre-commit install
	make auth.newkey

dev:
	env
	wgo run ./cmd/api/main.go

start:
	go run ./cmd/api/main.go

infra.up:
	docker compose up -d

i install:
	@echo "Installing dependencies..."
	go mod download

tidy:
	go mod tidy -v

c clean:
	rm -f covprofile*.out covprofile*.xml covprofile*.html
	rm -rf tmp

lint:
	# Run go vet to examine Go source code and report suspicious constructs
	go vet ./...
	# Format Go code using go fmt to ensure standard style
	go fmt ./...
	# Run gosec to check for security issues in Go code
	go tool gosec ./...
	# Run golangci-lint for comprehensive linting (style, bugs, etc.)
	go tool golangci-lint run

f fmt:
	go fmt ./...

g generate:
	@go generate ./...
	@echo 'Generating sqlc code...'
	@go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
	@echo 'Generating mocks with mockery...'
	@go tool mockery --config .mockery.yaml
	@go tool swag init -g cmd/api/main.go

b build:
	go build -o apiserver ./api/cmd

t test:
	# for clear cache `-count=1`
	@GIN_MODE=test go test -short $$(go list ./... | grep -v '/mock' | grep -v '/tests/integration')

it integration.test:
	@GIN_MODE=test go test ./tests/integration/...

# Run both unit and integration tests
test.all:
	@GIN_MODE=test make test
	@GIN_MODE=test make integration.test

# Unit test coverage
tc.unit test.cov.unit:
	@go test -covermode=atomic -coverprofile=covprofile-unit.out -short $$(go list ./... | grep -v '/mock' | grep -v '/tests/integration') || true
	@if [ -f covprofile-unit.out ]; then go tool cover -html=covprofile-unit.out -o covprofile-unit.html; fi

# Integration test coverage
tc.integration test.cov.integration:
	@go test -covermode=atomic -coverprofile=covprofile-integration.out ./tests/integration/... || true
	@if [ -f covprofile-integration.out ]; then go tool cover -html=covprofile-integration.out -o covprofile-integration.html; fi

# Combined coverage (legacy)
tc test.cov:
	@go test -covermode=atomic -coverprofile=covprofile.out $$(go list ./... | grep -v '/mock')
	@if [ -f covprofile.out ]; then go tool cover -html=covprofile.out; fi

# auth
auth.newkey:
	openssl ecparam -name prime256v1 -genkey -noout -out ecdsa_private_key.pem
	openssl ec -in ecdsa_private_key.pem -pubout -out ecdsa_public_key.pem

## ============ Database Migrations ============

# Database connection string for migrate CLI
DB_URL := postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DBNAME)?sslmode=$(DB_SSLMODE)

# Create a new migration file
# Usage: make migrate.create name=create_users_table
migrate.create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required"; \
		echo "Usage: make migrate.create name=create_users_table"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir db/migrations -seq $(name)

# Run all pending migrations
migrate.up:
	@echo "Running all pending migrations..."
	podman run --rm --network host -v $$(pwd)/db/migrations:/migrations migrate/migrate \
		-path=/migrations -database "$(DB_URL)" up

# Rollback migrations
# Usage: make migrate.down [steps=1]
migrate.down:
	@if [ -z "$(steps)" ]; then \
		echo "Rolling back 1 migration..."; \
		go run github.com/golang-migrate/migrate/v4/cmd/migrate -database "$(DB_URL)" -path db/migrations down 1; \
	else \
		echo "Rolling back $(steps) migrations..."; \
		go run github.com/golang-migrate/migrate/v4/cmd/migrate -database "$(DB_URL)" -path db/migrations down $(steps); \
	fi

# Show current migration version
migrate.version:
	@echo "Current migration version:"
	go run github.com/golang-migrate/migrate/v4/cmd/migrate -database "$(DB_URL)" -path db/migrations version

# Force migration to specific version (use with caution)
# Usage: make migrate.force version=20231201000001
migrate.force:
	@if [ -z "$(version)" ]; then \
		echo "Error: version parameter is required"; \
		echo "Usage: make migrate.force version=20231201000001"; \
		exit 1; \
	fi
	@echo "Forcing migration to version $(version)..."
	go run github.com/golang-migrate/migrate/v4/cmd/migrate -database "$(DB_URL)" -path db/migrations force $(version)

# Validate migration files
migrate.validate:
	@echo "Validating migration files..."
	go run github.com/golang-migrate/migrate/v4/cmd/migrate -path db/migrations validate

## ============ End Database Migrations ============
