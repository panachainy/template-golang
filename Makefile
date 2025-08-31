include .env
export $(shell sed 's/=.*//' .env)

dev:
	wgo run ./cmd/api/main.go

start:
	go run ./cmd/api/main.go

infra.up:
	docker-compose -f ./docker-compose.yml up -d

i: install
install:
	@echo "Installing dependencies..."
	go mod download

setup:
	make auth.newkey
	@echo "Installing go tools..."
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

tidy:
	go mod tidy -v

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
	go run github.com/golang-migrate/migrate/v4/cmd/migrate -database "$(DB_URL)" -path db/migrations up

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

t: test
test:
	# for clear cache `-count=1`
	@GIN_MODE=test go test -short $$(go list ./... | grep -v '/mock' | grep -v '/tests/integration')

it: integration.test
integration.test:
	@GIN_MODE=test go test ./tests/integration/...

# Run both unit and integration tests
test.all:
	@GIN_MODE=test make test
	@GIN_MODE=test make integration.test

tr: test.report
test.report:
	go test -race -covermode=atomic -coverprofile=covprofile.out $$(go list ./... | grep -v '/mock')
	make tc.html

tc: test.cov
test.cov:
	go test -covermode=atomic -coverprofile=covprofile.out $$(go list ./... | grep -v '/mock')
	make test.cov.xml

c: clean
clean:
	rm -f covprofile.out covprofile.xml covprofile.html
	rm -rf tmp

f: fmt
fmt:
	go fmt ./...

g: generate
generate:
	go generate ./...

sg: sqlc-generate
sqlc-generate:
	@echo 'Generating sqlc code...'
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

b: build
build:
	go build -o apiserver ./api/cmd

# swagger
swag.init:
	swag init -g cmd/api/main.go

# auth
auth.newkey:
	openssl ecparam -name prime256v1 -genkey -noout -out ecdsa_private_key.pem
	openssl ec -in ecdsa_private_key.pem -pubout -out ecdsa_public_key.pem
