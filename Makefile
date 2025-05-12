dev:
	air

setup:
	go install go.uber.org/mock/mockgen@latest
	go install github.com/axw/gocov/gocov@latest

tidy:
	go mod tidy -v

t: test
test:
	go test ./...

tc: test.cov
test.cov:
	$(ENV_LOCAL_TEST) \
	go test -race -covermode=atomic -coverprofile=covprofile.out ./modules/...
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
	swag init -g cmd/main.go

# auth

auth.newkey:
	openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in private.pem -pubout -out public.pem
