# Copilot Instructions for Go Template Project

## Project Overview
This is a Go template project using Clean Architecture with modular design, Gin framework, SQLC ORM, Wire dependency injection, and comprehensive testing. It serves as a foundation for building scalable Go web applications.

## Code Style Guidelines
- Follow standard Go conventions (gofmt, golint)
- Use meaningful variable and function names
- Write clear, concise comments for exported functions
- Prefer composition over inheritance
- Handle errors explicitly using the Gin error handling pattern
- Use structured logging with labstack/gommon/log

## Project Structure
- `/cmd` - Main application entry point
- `/config` - Configuration management with Viper
- `/database` - Database interfaces and PostgreSQL implementation
- `/db/migrations` - Database migration scripts
- `/modules` - Feature modules organized by domain (auth, cockroach)
- `/server` - HTTP server implementation with Gin
- `/docs` - Swagger documentation
- `/mock` - Generated mocks for testing
- `/example` - Example implementations (e.g., goth auth)

## Module Architecture Pattern
Each module follows Clean Architecture with this structure:
- `entities/` - Domain entities
- `repositories/` - Data access layer interfaces and implementations
- `usecases/` - Business logic layer
- `handlers/` - HTTP handlers and request/response models
- `middlewares/` - HTTP middlewares (auth, CORS, etc.)
- `models/` - Request/response models for API
- `migrations/` - Database migration scripts
- `wire.go` - Wire dependency injection setup

## Dependencies
- **Gin** - HTTP web framework
- **SQLC** - ORM for database operations
- **Wire** - Compile-time dependency injection
- **Viper** - Configuration management
- **Swagger** - API documentation
- **Testify** - Testing framework
- **gomock** - Mock generation
- **PostgreSQL** - Primary database
- **JWT** - Authentication tokens
- **Goth** - OAuth authentication

## Testing Patterns
- Use table-driven tests with testify/assert
- Create in-memory SQLite databases for repository tests
- Implement mock interfaces using gomock
- Test names should describe the behavior being tested
- Separate test setup functions (e.g., `setupTestDB`)
- Test error cases and edge cases
- Use `TestProvide*` pattern for constructor tests
- Test database interface compliance

## Dependency Injection with Wire
- Use `Provide` functions for constructors
- Use sync.Once for singleton services when needed with Wire
- Create `ProviderSet` in each package's `wire.go`
- Use `wire.Bind` to bind interfaces to implementations
- Generate code with `//go:generate wire`
- Keep wire files separate with build tags

## Database Patterns
- Use `Database` interface for testability
- Implement repository pattern with interfaces
- Use SQLC for ORM operations
- Handle foreign key relationships properly
- Use soft deletes with SQLC's DeletedAt
- Log database operations with structured logs

## HTTP Handler Patterns
- Use Gin for routing and middleware
- Validate request bodies with go-playground/validator
- Return consistent JSON responses
- Handle errors with Gin's error handling
- Use Swagger annotations for API documentation
- Implement proper HTTP status codes

## Configuration
- Use Viper for configuration management
- Support environment variables and .env files
- Set defaults for all configuration values
- Use mapstructure tags for binding
- Load configuration once using sync.Once

## Error Handling
- Return errors explicitly from all functions
- Log errors with context using labstack/gommon/log
- Use Gin's error handling for HTTP responses
- Provide meaningful error messages to clients

## Code Generation
- Use `//go:generate` directives for code generation
- Generate everything via `make g`
- Use Wire for dependency injection code generation
- Keep generated files separate and clearly marked

## Mock Generation
- Use gomock for generating mocks
- Write generate config in file that need to generate mock e.g. `//go:generate mockgen -source=jwtUsecase.go -destination=../../../mock/mock_jwt_usecase.go -package=mock`
- Store every mock in the `mock` package
- If not have generated mock, must not create new one in test file but create a new one using `mockgen`
