# Testing Pipeline Architecture

This document describes the separated testing pipeline architecture implemented in this project.

## Overview

The CI/CD pipeline has been split into separate jobs for unit tests and integration tests to improve efficiency, provide faster feedback, and enable better resource utilization.

## Pipeline Structure

### 1. Unit Tests Job (`unit-tests`)
- **Purpose**: Run fast, isolated unit tests without external dependencies
- **Dependencies**: None (runs immediately on push/PR)
- **Services**: None required
- **Coverage**: Generates `covprofile-unit.out` with unit test coverage
- **Make targets**: `make test`, `make tc.unit`

**Characteristics:**
- Fast execution (typically < 30 seconds)
- No database or external services required
- Tests business logic in isolation with mocks
- Provides immediate feedback on code quality

### 2. Integration Tests Job (`integration-tests`)
- **Purpose**: Run tests that require real database and external services
- **Dependencies**: Waits for `unit-tests` job to pass
- **Services**: PostgreSQL, Redis
- **Coverage**: Generates `covprofile-integration.out` with integration test coverage
- **Make targets**: `make it`, `make tc.integration`

**Characteristics:**
- Slower execution due to service startup and database operations
- Tests real interactions between components
- Validates database schema and migrations
- Only runs if unit tests pass (fail-fast approach)

## Benefits

### 1. **Faster Feedback**
- Unit tests provide immediate feedback without waiting for service startup
- Developers know quickly if basic logic is broken

### 2. **Resource Efficiency**
- Unit tests don't consume database/Redis resources unnecessarily
- Integration tests only run when unit tests pass

### 3. **Clear Failure Isolation**
- Easy to distinguish between unit test failures (logic issues) and integration test failures (infrastructure/integration issues)
- Debugging is more focused and efficient

### 4. **Parallel Execution Potential**
- Unit tests start immediately
- Integration tests can be optimized independently
- Future: Could run multiple integration test suites in parallel

## Make Targets

### Testing Commands
```bash
# Run unit tests only
make test

# Run integration tests only  
make it

# Run both (sequential, for local development)
make test.all
```

### Coverage Commands
```bash
# Unit test coverage
make tc.unit          # Generates covprofile-unit.out and covprofile-unit.html

# Integration test coverage
make tc.integration   # Generates covprofile-integration.out and covprofile-integration.html

# Combined coverage (legacy)
make tc              # Generates covprofile.out for all tests
```

## CI/CD Configuration

The pipeline is configured in `.github/workflows/ci.yaml` with:

- **unit-tests**: Runs immediately, no external services
- **integration-tests**: Depends on unit-tests, includes PostgreSQL and Redis services
- **lint**: Runs in parallel with tests for code quality checks

Both jobs upload coverage reports to Codecov with appropriate flags (`unit` and `integration`) for separate tracking.

## Development Workflow

1. **Write unit tests first** - Fast feedback loop for TDD
2. **Run `make test`** - Verify unit tests pass locally
3. **Write integration tests** - For end-to-end validation
4. **Run `make it`** - Verify integration tests pass locally
5. **Push changes** - CI runs both test suites automatically

## Future Enhancements

- **Parallel integration suites**: Split integration tests by feature area
- **Test result caching**: Speed up repeated test runs
- **Dynamic test selection**: Run only tests affected by code changes
- **Performance testing**: Add separate job for performance regression tests