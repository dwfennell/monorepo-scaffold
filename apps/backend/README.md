# Backend

Go backend with Gin, PostgreSQL, and JWT authentication.

## Setup

```bash
# Install dependencies
go mod download

# Install golang-migrate (if not already installed)
brew install golang-migrate

# Start the database (from project root)
docker-compose up -d

# Run migrations (from project root)
make migrate-up
```

## Running

```bash
# Run with hot reload (requires Air)
air

# Or run directly
go run main.go
```

## Database Migrations

**Recommended: Use Makefile** (reads from `.env` automatically):
```bash
make migrate-up                          # Run pending migrations
make migrate-down                        # Rollback last migration
make migrate-create NAME=migration_name  # Create new migration
```

**Or use CLI directly** (requires explicit database URL):
```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq migration_name

# Run all pending migrations
migrate -path migrations -database "$DATABASE_URL" up

# Rollback the last migration
migrate -path migrations -database "$DATABASE_URL" down 1

# Check migration version
migrate -path migrations -database "$DATABASE_URL" version
```

## Environment Variables

**Development:**
```bash
cp .env.example .env
```

**Testing:**
```bash
cp .env.test.example .env.test
```

The `.env.test` file is used by all test commands and contains test database configuration.

## Testing

The project includes comprehensive tests using Go's built-in `testing` package and `testify` for assertions.

### Test Types

**Unit Tests** - Test individual functions without external dependencies:
- `internal/auth/password_test.go` - Password hashing tests
- `internal/auth/jwt_test.go` - JWT generation/validation tests

**Integration Tests** - Test with real database (requires Docker):
- `internal/repository/user_repository_test.go` - Database operations
- `internal/api/auth_handler_test.go` - HTTP endpoint tests

### Running Tests

**From project root using Makefile:**
```bash
make test                # Run all tests (requires database)
make test-unit           # Run only unit tests (no database needed)
make test-integration    # Run only integration tests
make test-coverage       # Run with coverage report
make test-coverage-html  # Generate HTML coverage report
```

**From apps/backend directory:**
```bash
go test ./...              # Run all tests
go test -v ./...           # Verbose output
go test -short ./...       # Skip integration tests
go test -cover ./...       # With coverage
go test ./internal/auth    # Test specific package
```

### Test Requirements

- **Unit tests**: No dependencies, run anytime
- **Integration tests**: Require test database
  ```bash
  docker-compose up -d        # Starts PostgreSQL with monorepo_dev + monorepo_test
  make test-db-setup          # Runs migrations on test database
  ```

### Test Database Management

```bash
make test-db-setup   # Set up test database with migrations
make test-db-reset   # Reset test database (drop all tables and recreate)
```

**Important:** Integration tests use `monorepo_test` database, never `monorepo_dev`.
This ensures your development data is never touched by tests.

### Writing Tests

Test files use the `_test.go` suffix and follow these patterns:

**Unit Test Example:**
```go
func TestHashPassword(t *testing.T) {
    hash, err := HashPassword("password123")
    assert.NoError(t, err)
    assert.NotEmpty(t, hash)
}
```

**Integration Test Example:**
```go
func TestUserRepositoryTestSuite(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests")
    }
    suite.Run(t, new(UserRepositoryTestSuite))
}
```

## Project Structure

- `internal/api/` - HTTP handlers and middleware
- `internal/auth/` - Authentication utilities (JWT, password hashing)
- `internal/database/` - Database connection and migrations
- `internal/models/` - Data models and request/response types
- `internal/repository/` - Database queries and data access
