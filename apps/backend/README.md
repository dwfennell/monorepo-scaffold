# Backend

Go backend with Gin framework, PostgreSQL, and JWT authentication.

## Setup

**From project root:**
```bash
make install      # Install dependencies
make dev-db       # Start PostgreSQL
make migrate-up   # Run migrations
```

## Running

```bash
# From project root (recommended)
make dev-backend

# Or directly
go run main.go
air  # with hot reload
```

Server runs on http://localhost:8080

## Environment Variables

```bash
cp .env.example .env
```

Required variables:
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT tokens
- `FRONTEND_URL` - CORS allowed origin
- `PORT` - Server port (default: 8080)

## Database Migrations

```bash
make migrate-up                          # Run pending migrations
make migrate-down                        # Rollback last migration
make migrate-create NAME=migration_name  # Create new migration
```

## Testing

```bash
make test              # All tests (requires database)
make test-unit         # Unit tests only
make test-integration  # Integration tests only
make test-coverage     # With coverage report
```

**Note:** Integration tests use `monorepo_test` database, never `monorepo_dev`.
