## Why

The snake tournament application requires comprehensive testing for its core game service to ensure reliability and prevent regressions. Currently, there are no automated tests for the internal/services/game.go service, which handles critical game logic including snake movements, food generation, collision detection, and game state management. Implementing a test container approach will enable developers to write effective unit and integration tests, improving code quality and maintainability.

## What Changes

- Introduce test container setup for the game service
- Implement unit tests for core game logic functions in internal/services/game.go
- Add integration tests to verify game service behavior with mocked dependencies
- Create test helpers and utilities to support testing throughout the application
- Configure test execution in the CI/CD pipeline

## Capabilities

### New Capabilities
- **game-service-testing**: Defines the testing framework and approach for the game service
- **unit-tests**: Covers individual functions and methods in the game service
- **integration-tests**: Tests the service behavior with mocked external dependencies

### Modified Capabilities
- None

## Impact

This change will impact:
- The internal/services/game.go package
- The testing infrastructure in the project
- The development workflow for implementing new game features
- The overall test coverage of the application