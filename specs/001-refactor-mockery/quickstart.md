# Quickstart: Generic Mockery for All Mocks

This guide shows how to use the new mockery-based mock generation system that replaces the previous gomock approach.

## What Changed

✅ **Before (gomock):**
- Manual `//go:generate mockgen` directives in each interface file
- Inconsistent mock locations (`mock/` vs `mocks/`)
- Required mockgen tool installation
- Complex gomock.Controller usage in tests

✅ **After (mockery):**
- Centralized configuration in `.mockery.yaml`
- Consistent mock locations in `mocks/` folders
- Simple testify mock usage in tests
- Single `make g` command generates all mocks

## Quick Start

### 1. Generate All Mocks
```bash
make g
```
This command:
- Generates SQLC code
- Generates all mocks using mockery
- Updates Swagger documentation

### 2. Run Tests
```bash
make test        # Unit tests only
make test.all    # Unit + integration tests
```

### 3. Add New Interface Mocks

To add mocks for a new interface, update `.mockery.yaml`:

```yaml
packages:
  template-golang/modules/yourmodule/usecases:
    interfaces:
      YourInterface:
```

Then run `make g` to generate the new mocks.

## Mock Configuration

The project uses a centralized `.mockery.yaml` configuration:

```yaml
all: false
dir: '{{.InterfaceDir}}/mocks'
filename: 'mock_{{.InterfaceName|snakecase}}.go'
template: testify
packages:
  template-golang/modules/auth/usecases:
    interfaces:
      JWTUsecase:
  template-golang/modules/auth/middlewares:
    interfaces:
      AuthMiddleware:
  template-golang/modules/cockroach/usecases:
    interfaces:
      CockroachUsecase:
```

## Using Mocks in Tests

### Before (gomock):
```go
ctrl := gomock.NewController(t)
defer ctrl.Finish()
mockJWT := mock.NewMockJWTUsecase(ctrl)
mockJWT.EXPECT().GenerateJWT("user-123").Return("token", nil)
```

### After (testify):
```go
mockJWT := mocks.NewMockJWTUsecase(t)
mockJWT.On("GenerateJWT", "user-123").Return("token", nil)
```

## File Structure

```
modules/
├── auth/
│   ├── middlewares/
│   │   ├── mocks/
│   │   │   └── mock_auth_middleware.go  # Generated
│   │   └── authMiddleware.go            # Interface definition
│   └── usecases/
│       ├── mocks/
│       │   └── mock_jwt_usecase.go      # Generated
│       └── jwtUsecase.go                # Interface definition
└── cockroach/
    └── usecases/
        ├── mocks/
        │   └── mock_cockroach_usecase.go # Generated
        └── cockroachUsecase.go           # Interface definition
```

## Benefits

1. **Centralized Configuration**: All mock generation controlled from `.mockery.yaml`
2. **Consistent Structure**: All mocks in `mocks/` folders
3. **Testify Integration**: Better assertion and expectation syntax
4. **Automatic Cleanup**: Testify automatically calls `AssertExpectations`
5. **Faster Generation**: Single command generates all mocks
6. **Better CI/CD**: No tool installation required (uses go run)

## Troubleshooting

### Mock not generated?
Check that your interface is listed in `.mockery.yaml` under the correct package path.

### Tests failing?
Ensure you're importing from the `mocks/` package (not `mock/`) and using testify syntax:
```go
import "template-golang/modules/auth/usecases/mocks"
mockJWT := mocks.NewMockJWTUsecase(t)
```

### Need to add new mocks?
1. Add interface to `.mockery.yaml`
2. Run `make g`
3. Import the new mock in your tests
