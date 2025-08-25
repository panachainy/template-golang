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
	brew install golang-migrate

tidy:
	go mod tidy -v

## ============ Start DB ============

# make migrate.create name=<migration_name>
migrate.create:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate.up:
	go run ./cmd/migrate/main.go -action=up

# make migrate.down steps=1
migrate.down:
	go run ./cmd/migrate/main.go -action=down -steps=$(steps)

migrate.version:
	go run ./cmd/migrate/main.go -action=version

# Legacy migrate commands using migrate CLI directly
migrate.cli.up:
	migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DBNAME)?sslmode=$(DB_SSLMODE)" -path db/migrations up

migrate.cli.down:
	migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DBNAME)?sslmode=$(DB_SSLMODE)" -path db/migrations down

## ============ End DB ============

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
	go test -race -covermode=atomic -coverprofile=covprofile.out $$(go list ./... | grep -v '/mock')
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
	go tool sqlc generate

b: build
build:
	go build -o apiserver ./api/cmd

# swagger
swag.init:
	swag init -g cmd/api/main.go

# auth
auth.newkey:
	# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	# openssl rsa -in private.pem -pubout -out public.pem
	openssl ecparam -name prime256v1 -genkey -noout -out ecdsa_private_key.pem
	openssl ec -in ecdsa_private_key.pem -pubout -out ecdsa_public_key.pem
