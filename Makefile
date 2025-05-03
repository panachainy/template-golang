dev:
	air

setup:
	go install go.uber.org/mock/mockgen@latest

tidy:
	go mod tidy -v

t: test
test:
	go test ./...

tc: test.cov
test.cov:
	$(ENV_LOCAL_TEST) \
	go test -race -covermode=atomic -coverprofile=covprofile.out ./internal/core/... ./internal/feature/...
	make test.cov.xml

tc.xml: test.cov.xml
test.cov.xml:
	gocov convert covprofile.out > covprofile.xml

f: fmt
fmt:
	go fmt ./...

w: wire
wire:
	wire ./...

g: generate
generate:
	go generate ./...

b: build
build:
	go build -o apiserver ./cmd

migrate:
	go run ./modules/cockroach/migrations/cockroachMigrate.go

# swagger

swag.init:
	# cd ./cmd && swag init
	swag init -g cmd/main.go

# swag init -g http/api.go
# swag fmt
