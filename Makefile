dev:
	wgo run ./cmd/main.go

start:
	go run ./cmd/main.go

infra.up:
	docker-compose up -d

setup:
	go install go.uber.org/mock/mockgen@latest
	go install github.com/axw/gocov/gocov@latest
	go install github.com/bokwoon95/wgo@latest
	make auth.newkey

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

c: clean
clean:
	rm -f covprofile.out covprofile.xml covprofile.html
	rm -rf tmp

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
	go run ./modules/auth/migrations/authMigrate.go

# swagger

swag.init:
	swag init -g cmd/main.go

# auth

auth.newkey:
	# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	# openssl rsa -in private.pem -pubout -out public.pem
	openssl ecparam -name prime256v1 -genkey -noout -out ecdsa_private_key.pem
	openssl ec -in ecdsa_private_key.pem -pubout -out ecdsa_public_key.pem
