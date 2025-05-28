dev:
	air

start:
	go run ./cmd/main.go

setup:
	go install go.uber.org/mock/mockgen@latest
	go install github.com/axw/gocov/gocov@latest

tidy:
	go mod tidy -v

t: test
test:
	go test ./...

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

clean.test:
	rm -f covprofile.out covprofile.xml covprofile.html

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
	# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	# openssl rsa -in private.pem -pubout -out public.pem
	openssl ecparam -name prime256v1 -genkey -noout -out ecdsa_private_key.pem
	openssl ec -in ecdsa_private_key.pem -pubout -out ecdsa_public_key.pem
