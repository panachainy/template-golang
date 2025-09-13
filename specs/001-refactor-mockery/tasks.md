# Tasks: Generic Mockery for All Mocks

Feature: Generic Mockery for All Mocks
Branch: 001-refactor-mockery
Docs: /specs/001-refactor-mockery/

---

## Setup Tasks

T001. Install mockery tool globally
- Command: `go install github.com/vektra/mockery/v2@latest`
- Dependency: Go installed

T002. Create mockery.yaml config in project root
- Path: /mockery.yaml
- Content: YAML config listing all interfaces to mock
- Dependency: T001

T003. Add mockery target to Makefile
- Path: /Makefile
- Content: `mockery: mockery --config mockery.yaml`
- Dependency: T002

T004. Remove gomock and mockgen dependencies
- Paths: /go.mod, /go.sum, all gomock-generated files
- Dependency: T003

---

## Test Tasks [P]

T005. Create contract test for mockery-generate-mocks
- Path: /specs/001-refactor-mockery/contracts/mockery-generate-mocks.md
- Description: Write failing test for mock generation contract
- Dependency: T003

T006. Integration test: Regenerate mocks and run all tests
- Command: `make mockery && go test ./...`
- Dependency: T005

---

## Core Tasks [P]

T007. For each interface in data-model, generate mock using mockery
- Paths: modules/auth/usecases/jwtUsecase.go, modules/cockroach/usecases/cockroachUsecase.go
- Output: modules/auth/usecases/mock/, modules/cockroach/usecases/mock/
- Dependency: T002

T008. Update all tests to use mockery-generated mocks
- Paths: all test files using old mocks
- Dependency: T007

---

## Integration Tasks

T009. Validate mockery integration in CI/CD pipeline
- Path: CI config (e.g., .github/workflows/)
- Ensure `make mockery` runs before tests
- Dependency: T006

---

## Polish Tasks [P]

T010. Update quickstart.md and README with mockery instructions
- Paths: /specs/001-refactor-mockery/quickstart.md, /README.md
- Dependency: T009

T011. Review and clean up unused mock files
- Paths: all modules and mocks/
- Dependency: T010

---

## Parallel Execution Guidance
- Tasks marked [P] (T005, T006, T007, T008, T010, T011) can be run in parallel if working on different files.
- Example: T007 (mock generation for each interface) can be run for multiple interfaces at once.
- Use: `copilot run T007 T008` for parallel execution.

---

## Dependency Notes
- Setup tasks (T001-T004) must be completed before tests and core tasks.
- Contract and integration tests (T005-T006) validate the mockery setup before updating codebase.
- Core tasks (T007-T008) implement the generic mockery approach.
- Integration and polish tasks (T009-T011) ensure CI/CD and documentation are updated.

---
