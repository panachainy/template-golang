# Research: Refactor to Use Mockery for Mocks

## Decision
- Use `mockery` for all mock generation in the project.
- Configure mockery using a YAML file for consistency and maintainability.
- Remove all gomock and mockgen usage.

## Rationale
- Mockery is widely adopted in the Go community and supports YAML configuration for mock generation.
- YAML config allows centralized and declarative mock management.
- Reduces manual steps and errors compared to gomock.
- Simplifies onboarding and maintenance for new contributors.

## Alternatives Considered
- **gomock**: Previously used, but requires manual mockgen commands and lacks YAML config support.
- **moq**: Simpler but less flexible and not as widely adopted for complex projects.
- **manual mocks**: Not maintainable for large codebases.

## Best Practices
- Use mockery's YAML config to specify interfaces and output paths.
- Integrate mockery generation into Makefile for reproducibility.
- Document mock generation in quickstart.md for contributors.
- Ensure all tests are updated to use mockery-generated mocks.

## Unknowns/Clarifications
- Confirm all interfaces to be mocked are public and compatible with mockery.
- Validate mockery integration with CI/CD pipeline.
- Ensure mockery-generated mocks are compatible with existing test patterns (testify, etc).

## References
- [mockery GitHub](https://github.com/vektra/mockery)
- [mockery YAML config docs](https://github.com/vektra/mockery#yaml-configuration)
