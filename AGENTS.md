# Agent Development Guidelines

## Build/Lint/Test Commands

### Build Commands
- `make build` - Builds the application binary
- `make bin/snake-tournament` - Creates a Linux binary for deployment
- `make serve` - Runs the development server locally using Go

### Lint Commands
- `make lint` - Runs golangci-lint for code quality checking

### Test Commands
- `go test ./...` - Run all tests in the project
- `go test -v ./...` - Run all tests with verbose output
- `go test ./path/to/package` - Run tests for a specific package
- `go test -v ./path/to/package` - Run tests for a specific package with verbose output

## Code Style Guidelines

### Imports
- Group imports in the following order:
  1. Standard library packages
  2. Third-party packages
  3. Internal packages (starting with project prefix like "snake-tournament/")
  
### Formatting
- Use gofmt for code formatting
- Follow Go's official naming conventions
- Use camelCase for variable and function names
- Use PascalCase for exported types and functions

### Types
- Use concrete types over interfaces when possible
- Define custom error types using `errors.New()` or `fmt.Errorf()`
- Prefer struct embedding over composition when appropriate

### Naming Conventions
- Package names: short, lowercase, no underscores or dashes
- Structs and interface names: PascalCase
- Variables and functions: camelCase
- Exported items (public): start with uppercase letter
- Unexported items (private): start with lowercase letter

### Error Handling
- Handle errors immediately after function calls
- Use meaningful error messages that include context
- Prefer creating new errors with `fmt.Errorf()` when wrapping existing errors
- Consider using custom error types for better error categorization

### Documentation
- All exported functions, methods and types must have documentation comments (docstrings)
- Use godoc comment format for generated documentation

### Configuration Management
- Use environment variables for configuration via `cleanenv`
- Configuration structs should be well-documented with `env` tags
- Provide sensible defaults in configuration struct fields

### Structure
- Package structure:
  - `cmd/` - Main applications
  - `internal/` - Private application code
  - `pkg/` - Reusable package code
  - `models/` - Data models

### File Naming
- Use lowercase with underscores for filenames (e.g., `user_service.go`)
- Test files should end with `_test.go`
- Use descriptive names for files and packages

### Logging
- Use logrus for logging (import as `logrus`)
- Log at appropriate levels: debug, info, warning, error
- Include context in logs when helpful

## Project-Specific Information

### API Endpoints
The application exposes the following snake tournament API endpoints:
- POST /api/v1/snake-tournament/start - Start a new game
- GET /api/v1/snake-tournament/start/:id - Enter a specific game
- GET /api/v1/snake-tournament/find-my - Get games for current user
- GET /api/v1/snake-tournament/find-active - Get active games for current user
- POST /api/v1/snake-tournament/create-record/:id - Submit game results
- GET /api/v1/snake-tournament/check-renew-record-allowed/:id - Check if renewing record is allowed
- POST /api/v1/snake-tournament/select-prize/:id - Select a prize for a game

### Data Models
- User DTO with Id, Username, Email, and Roles fields
- Game-related DTOs for creation, results, and prize selection
- Error DTO for standardized error responses

### Dependencies
- Uses httprouter for HTTP routing
- Uses MongoDB via go.mongodb.org/mongo-driver for database operations
- Uses cleanenv for configuration management
- Uses logrus for logging
- Uses ehttp middleware for request handling

## Additional Notes
This codebase uses Go 1.23 with MongoDB integration and HTTP routing via httprouter.
The project follows a clean architecture pattern separating business logic from infrastructure concerns.

## Database Patterns
- Uses MongoDB with the go.mongodb.org/mongo-driver for database operations
- Database access is implemented through a BaseStorage pattern
- Game-related data is stored in MongoDB using the "games" collection
- Database queries use MongoDB's BSON query syntax with complex filtering operations
- The database layer follows the repository pattern with methods for find operations
