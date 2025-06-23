include .env
export $(shell sed 's/=.*//' .env)

dev:
	wgo run ./cmd/api/main.go

start:
	go run ./cmd/api/main.go

infra.up:
	docker-compose -f ./docker-compose.yml up -d

setup:
	go install go.uber.org/mock/mockgen@latest
	go install github.com/axw/gocov/gocov@latest
	go install github.com/bokwoon95/wgo@latest
	go install golang.org/x/tools/gopls@latest
	make auth.newkey
	brew install golang-migrate

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

tidy:
	go mod tidy -v

t: test
test:
	@GIN_MODE=test go test ./...

it: integration.test
integration.test:
	@GIN_MODE=test go test -tags=integration ./...

# Run both unit and integration tests
test.all:
	@GIN_MODE=test go test ./...
	@GIN_MODE=test go test -tags=integration ./...

# Run tests with verbose output
test.verbose:
	@GIN_MODE=test go test -v ./...

# Run integration tests with verbose output
integration.test.verbose:
	@GIN_MODE=test go test -tags=integration -v ./...

# Run all tests with verbose output
test.all.verbose:
	@GIN_MODE=test go test -v ./...
	@GIN_MODE=test go test -tags=integration -v ./...

# Test with race detection
test.race:
	@GIN_MODE=test go test -race ./...
	@GIN_MODE=test go test -race -tags=integration ./...

tr: test.html
test.html:
	go test -race -covermode=atomic -coverprofile=covprofile.out ./...
	make tc.html

tc: test.cov
test.cov:
	go test -race -covermode=atomic -coverprofile=covprofile.out ./...
	make test.cov.xml

tc.xml: test.cov.xml
test.cov.xml:
	gocov convert covprofile.out > covprofile.xml

tc.html: test.cov.html
test.cov.html:
	go tool cover -html=covprofile.out -o covprofile.html
	open covprofile.html

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

b: build
build:
	go build -o apiserver ./cmd


# swagger

swag.init:
	swag init -g cmd/api/main.go

# auth

auth.newkey:
	# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	# openssl rsa -in private.pem -pubout -out public.pem
	openssl ecparam -name prime256v1 -genkey -noout -out ecdsa_private_key.pem
	openssl ec -in ecdsa_private_key.pem -pubout -out ecdsa_public_key.pem
