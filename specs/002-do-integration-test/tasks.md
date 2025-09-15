# Tasks: Integration Test for AuthHandler Interface

**Input**: Design documents from `/specs/002-do-integration-test/`
**Prerequisites**: plan.md (required)

## Execution Flow (main)
```
1. Load plan.md from feature directory
2. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: integration tests for AuthHandler endpoints
   → Core: DB migrations, test DB setup, minimal mocking
   → Integration: DB connection, logging, error handling
   → Polish: unit tests, docs
3. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
4. Number tasks sequentially (T001, T002...)
5. Generate dependency graph
6. Create parallel execution examples
7. Validate task completeness
```

## Phase 3.1: Setup
- [ ] T001 Create integration test folder and structure in `tests/integration`
- [ ] T002 Initialize Go 1.25.1 test dependencies (Testify, SQLC, Gin) in go.mod
- [ ] T003 [P] Configure linting and formatting tools (gofmt, golint)

## Phase 3.2: Tests First (TDD)
- [ ] T004 [P] Write failing integration test for `Login` endpoint in `tests/integration/auth_login_test.go`
- [ ] T005 [P] Write failing integration test for `AuthCallback` endpoint in `tests/integration/auth_callback_test.go`
- [ ] T006 [P] Write failing integration test for `Logout` endpoint in `tests/integration/auth_logout_test.go`
- [ ] T007 [P] Write failing integration test for `Example` endpoint in `tests/integration/auth_example_test.go`

## Phase 3.3: Core Implementation
- [ ] T008 Setup real PostgreSQL test database and apply migrations from `db/migrations`
- [ ] T009 Implement DB connection logic for integration tests in `tests/integration/testutils.go`
- [ ] T010 Implement minimal mocks for external dependencies (only if real not possible) in `tests/integration/mocks/`
- [ ] T011 Implement AuthHandler integration logic to use real DB in `modules/auth/handlers/authHttp.go`
- [ ] T012 Implement error handling and logging for integration tests

## Phase 3.4: Polish
- [ ] T013 [P] Add unit tests for edge cases in `tests/integration/auth_edge_test.go`
- [ ] T014 [P] Update documentation for integration test setup in `specs/002-do-integration-test/quickstart.md`
- [ ] T015 [P] Review and refactor test code for maintainability

## Dependencies
- T001-T003 before T004-T007
- T004-T007 before T008-T012
- T008 blocks T009, T011
- T010 only if needed (external dependencies)
- Implementation before polish (T013-T015)

## Parallel Example
```
# Launch T004-T007 together:
Task: "Write failing integration test for Login endpoint in tests/integration/auth_login_test.go"
Task: "Write failing integration test for AuthCallback endpoint in tests/integration/auth_callback_test.go"
Task: "Write failing integration test for Logout endpoint in tests/integration/auth_logout_test.go"
Task: "Write failing integration test for Example endpoint in tests/integration/auth_example_test.go"
```

## Validation Checklist
- [ ] All endpoints have corresponding integration tests
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task
