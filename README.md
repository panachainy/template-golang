# Template Golang

## Supports

- [x] Clean Architecture
- [x] PostgreSQL
- [x] ~~air~~
  - [x] use `wgo` instead
- [x] Viper config
  - [x] Fix issue about `.env` & `.env.test`
- [x] Docker
- [x] Gorm
- [x] [Swagger-Gin](https://github.com/swaggo/gin-swagger)
- [x] Wire [need to refactor all code to wire]
- [x] Unit Test
- [x] Mock
- [x] [Validator](https://github.com/go-playground/validator)
- [x] [golang-migrate/migrate](https://github.com/golang-migrate/migrate/tree/master?tab=readme-ov-file)
- [ ] [sqlc](https://github.com/sqlc-dev/sqlc)
- [ ] Auth JWT
  - [ ] [goth](https://github.com/markbates/goth)
    - [x] [JWT goth](https://github.com/markbates/goth/issues/310)
    - [x] Permission [admin, staff, user]
      - [x] middleware with role
    - [ ] Refresh token (doing)
      - [x] add exp for JWT
      - [x] add func verify token (1. expired, 2. not exist, 3. valid)
      - [x] implement middleware jwt for call verify func
      - [ ] add refresh token func
    - [ ] Save db
- [ ] Redis
- [ ] Logger system ([zap](https://github.com/uber-go/zap))
- [ ] [Casbin](https://github.com/casbin/casbin)
- [ ] try [WTF](https://github.com/pallat/wtf)
  - [ ] Do corsHandler to support configurable.
- [ ] try [failover](https://github.com/wongnai/lmwn_gomeetup_failover)
- [x] husky
- [ ] do pkg for common tools
- [ ] try option function https://github.com/kazhuravlev/options-gen?tab=readme-ov-file
- [ ] unit & integration test
  - [ ] update makefile
  - [ ] clean makefile
  - [ ] split CI & integration test

## Design pattern

[clean arch](https://medium.com/@rayato159/how-to-implement-clean-architecture-in-golang-87e9f2c8c5e4)

## Pre-commit Hooks

This project uses [Husky](https://typicode.github.io/husky/) to run pre-commit hooks that ensure code quality before commits.

### Setup

Pre-commit hooks are automatically set up when you run:

```bash
npm install
```

### What runs on each commit

The pre-commit hook runs the following commands:

1. **Go Linter** (`make lint`):
   - `go vet ./...` - Examines Go source code and reports suspicious constructs
   - `go mod tidy` - Ensures go.mod matches the source code
   - `go fmt ./...` - Formats Go source code

2. **Go Format** (`make fmt`):
   - `go fmt ./...` - Formats Go source code according to gofmt style

3. **Go Unit Tests** (`make test`):
   - `go test -short ./...` - Runs all unit tests

### Manual execution

You can manually run the same checks that the pre-commit hook performs:

```bash
# Run all pre-commit checks
make lint
make fmt  
make test

# Or individual commands
make lint    # Run linter and tidy modules
make fmt     # Format code
make test    # Run unit tests
```

If any of these commands fail, the commit will be rejected and you'll need to fix the issues before committing.
