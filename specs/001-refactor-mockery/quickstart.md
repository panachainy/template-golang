# Quickstart: Mockery Migration

## 1. Install Mockery
```sh
go install github.com/vektra/mockery/v2@latest
```

## 2. Add YAML Config
Create `mockery.yaml` in the project root:
```yaml
mockery:
  dir: ./modules
  output: ./mock
  package: mock
  case: underscore
  interfaces:
    - name: JWTUsecase
      path: modules/auth/usecases/jwtUsecase.go
      output: modules/auth/usecases/mock
      package: mock
    - name: CockroachUsecase
      path: modules/cockroach/usecases/cockroachUsecase.go
      output: modules/cockroach/usecases/mock
      package: mock
```

## 3. Update Makefile
Add a target to generate mocks:
```makefile
mockery:
	mockery --config mockery.yaml
```

## 4. Remove gomock
- Delete all gomock-generated mocks and mockgen commands.
- Remove gomock from go.mod and go.sum.

## 5. Regenerate Mocks
```sh
make mockery
```

## 6. Update Tests
- Ensure all tests use mockery-generated mocks.
- Run all tests:
```sh
go test ./...
```

## 7. CI/CD
- Ensure CI pipeline uses `make mockery` before running tests.

## 8. Documentation
- Update README and developer docs to reference mockery and YAML config.
