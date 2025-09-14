# Feature Spec: Generic Mockery for All Mocks

## Feature Requirements
- Refactor mock generation to use a generic approach for all mocks in the project.
- Ensure all modules (auth, cockroach, etc.) use the same mockery configuration and patterns.
- Remove duplicate or manual mock implementations.
- Support easy addition of new mocks via configuration, not manual code.

## User Stories
- As a developer, I want to generate mocks for any interface in the project using a single command/config.
- As a maintainer, I want to avoid manual updates to mock files when interfaces change.
- As a contributor, I want clear documentation and quickstart for generating and using mocks.

## Functional Requirements
- Centralize mockery configuration (e.g., via a config file or Makefile target).
- All mocks generated should be placed in the `mocks/` directory or module-specific `mock/` folders.
- Support Go interfaces in all modules.
- Generated mocks must be compatible with existing tests.

## Non-Functional Requirements
- Mock generation should be fast and reliable.
- Documentation for mock generation must be clear and up-to-date.
- No breaking changes to existing test workflows.

## Success Criteria
- All interfaces have corresponding mocks generated via the generic approach.
- No manual mock code remains in the codebase.
- Tests pass using the new mocks.
- Quickstart guide for mock generation is available.

## Acceptance Criteria
- Running the mockery command generates all required mocks.
- No errors or missing mocks after refactor.
- Documentation and quickstart updated.

## Technical Constraints
- Must use mockery (or compatible tool) for Go interface mocking.
- Must work with current project structure and CI/CD.
- No changes to business logic or API contracts.

## Dependencies
- mockery tool
- Go modules and interfaces
- Existing test framework (testify, gomock)

---
