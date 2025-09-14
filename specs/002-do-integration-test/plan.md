# Implementation Plan: [FEATURE]


**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]

# Implementation Plan: Integration Test for AuthHandler Interface


**Branch**: `002-do-integration-test` | **Date**: September 14, 2025 | **Spec**: [/specs/002-do-integration-test/spec.md]
## Execution Flow (/plan command scope)
**Input**: Feature specification from `/specs/002-do-integration-test/spec.md`
```
1. Load feature spec from Input path
## Execution Flow (/plan command scope)
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
4. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
5. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, or `GEMINI.md` for Gemini CLI).
6. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
7. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
8. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
    → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
    → Detect Project Type from context (web=frontend+backend, mobile=app+api)
    → Set Structure Decision based on project type
3. Evaluate Constitution Check section below
    → If violations exist: Document in Complexity Tracking
    → If no justification possible: ERROR "Simplify approach first"
    → Update Progress Tracking: Initial Constitution Check
4. Execute Phase 0 → research.md
    → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
5. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `.github/copilot-instructions.md` for GitHub Copilot).
6. Re-evaluate Constitution Check section
    → If new violations: Refactor design, return to Phase 1
    → Update Progress Tracking: Post-Design Constitution Check
7. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
8. STOP - Ready for /tasks command
```
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Primary requirement: End-to-end integration tests for the AuthHandler interface, verifying login, callback, logout, and example endpoint flows using a real PostgreSQL database and migrations. Only mock dependencies that cannot be tested with real resources.
[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context
**Language/Version**: Go 1.25.1
**Primary Dependencies**: Gin, SQLC, Testify, PostgreSQL, Wire, Viper
**Storage**: PostgreSQL (real DB, not in-memory)
**Testing**: Go test, Testify, integration tests in `tests/integration`
**Target Platform**: Linux/macOS server
**Project Type**: Backend web service (single project)
**Performance Goals**: [NEEDS CLARIFICATION: Integration test runtime target?]
**Constraints**: Use real DB and migrations, mock only what cannot be real (e.g., external providers)
**Scale/Scope**: [NEEDS CLARIFICATION: Expected number of test cases/users?]
## Constitution Check
**Testing**: [e.g., pytest, XCTest, cargo test or NEEDS CLARIFICATION]
**Simplicity**:
- Projects: 2 (main app, integration tests)
- Using Gin directly for HTTP, SQLC for DB
- Single data model per domain
- No unnecessary patterns; repository only where justified

## Constitution Check
**Architecture**:
- Features modularized as libraries (modules/auth, modules/cockroach)
- CLI via main.go for migrations/tests
- Library docs: [NEEDS CLARIFICATION: llms.txt planned?]

**Simplicity**:
**Testing (NON-NEGOTIABLE)**:
- RED-GREEN-Refactor cycle enforced
- Tests written before implementation
- Order: Contract→Integration→Unit
- Real PostgreSQL DB used for integration
- Integration tests for AuthHandler endpoints
- No implementation before failing test
- Using framework directly? (no wrapper classes)
- Single data model? (no DTOs unless serialization differs)
**Observability**:
- Structured logging via labstack/gommon/log
- Error context included in logs

**Architecture**:
**Versioning**:
- Versioning via go.mod and migration scripts
- BUILD increments tracked in CI
- Breaking changes require migration scripts
- Libraries listed: [name + purpose for each]
- CLI per library: [commands with --help/--version/--format]
## Project Structure


### Documentation (this feature)
```
specs/002-do-integration-test/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```
- RED-GREEN-Refactor cycle enforced? (test MUST fail first)
- Git commits show tests before implementation?
### Source Code (repository root)
```
modules/
├── auth/
│   ├── handlers/
│   ├── repositories/
│   ├── usecases/
│   └── ...
└── cockroach/
      └── ...
- Real dependencies used? (actual DBs, not mocks)
tests/
├── integration/
│   └── [integration tests for AuthHandler]
└── ...
db/
├── migrations/
│   └── [PostgreSQL migration scripts]
└── ...
```

**Structure Decision**: Single backend project with modular features and dedicated integration test folder.
- Integration tests for: new libraries, contract changes, shared schemas?
- FORBIDDEN: Implementation before test, skipping RED phase
## Phase 0: Outline & Research
**Observability**:
- Structured logging included?
- Frontend logs → backend? (unified stream)
- Error context sufficient?

**Versioning**:
- Version number assigned? (MAJOR.MINOR.BUILD)
- BUILD increments on every change?
- Breaking changes handled? (parallel tests, migration plan)

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
## Phase 1: Design & Contracts

### Source Code (repository root)
```
# Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# Option 3: Mobile + API (when "iOS/Android" detected)
api/
## Phase 2: Task Planning Approach

ios/ or android/
└── [platform-specific structure]
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```
## Phase 3+: Future Implementation
3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Complexity Tracking
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
## Progress Tracking
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

   - Each story → integration test scenario
   - Quickstart test = story validation steps
   - Preserve manual additions between markers
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- Each contract → contract test task [P]
- Each entity → model creation task [P]
- Each user story → integration test task
- Implementation tasks to make tests pass

**Ordering Strategy**:
- TDD order: Tests before implementation
- Dependency order: Models before services before UI
- Mark [P] for parallel execution (independent files)

**Estimated Output**: 25-30 numbered, ordered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)
**Phase 4**: Implementation (execute tasks.md following constitutional principles)
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [ ] Phase 0: Research complete (/plan command)
- [ ] Phase 1: Design complete (/plan command)
- [ ] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [ ] Initial Constitution Check: PASS
- [ ] Post-Design Constitution Check: PASS
- [ ] All NEEDS CLARIFICATION resolved
- [ ] Complexity deviations documented

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*