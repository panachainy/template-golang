# Data Model: Mockery Migration

## Entities
- **Mock Config**
  - Fields: interface name, output path, package, case, etc.
  - Relationships: Each config entry maps to a Go interface to be mocked.

## Validation Rules
- All interfaces to be mocked must be public.
- Output paths must match project structure.
- YAML config must be valid and cover all required mocks.

## State Transitions
- Old gomock mocks removed → mockery-generated mocks added
- Makefile updated to use mockery
- Tests updated to use new mocks

---
# Example YAML Config (mockery.yaml)
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
