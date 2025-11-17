# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Principles

**Configuration Management**:
- **NEVER fall back to hardcoded defaults** when configuration is missing
- **ALWAYS throw errors** for misconfigurations (missing environment variables, missing files, etc.)
- Fail fast and clearly - make it obvious when something is misconfigured
- Example: If `JWT_SECRET` is missing, throw an error rather than using a default value

**Implementation Standards**:
- **AVOID placeholder functionality** - do not leave TODOs or incomplete implementations
- If a task is too ambitious, **break it into smaller, completable parts**
- **ALWAYS inform the user** if you did not fully implement a feature or left any incomplete work
- Every feature you implement should be fully functional, not a stub or placeholder

**YAGNI (You Aren't Gonna Need It)**:
- **Only add what's needed now** - resist the temptation to build for hypothetical future requirements
- **Build extensible structures** - use patterns, interfaces, and abstractions that allow future growth
- **Don't add unnecessary code** - avoid extra features, options, or complexity that aren't currently required
- **Avoid over-engineering** - keep solutions simple and focused on the current problem
- Example: Don't add pagination to a list that currently has 5 items, but structure the API so pagination can be added later when needed

## Essential Commands

### Development
```bash
# Start frontend only (from root)
npm run dev

# Start backend (from root, separate terminal)
make dev-backend

# Start database
make dev-db

# Run all Turborepo tasks
npm run build        # Build all packages
npm run test         # Run all tests
npm run lint         # Lint all packages
npm run type-check   # TypeScript validation
```

### Database Migrations
```bash
make migrate-up                          # Run pending migrations
make migrate-down                        # Rollback last migration
make migrate-create NAME=migration_name  # Create new migration (uses current timestamp)
```

**IMPORTANT**: When creating migrations, use the current timestamp. The `migrate-create` command automatically handles this.

### Backend Testing
```bash
make test              # All tests (requires database)
make test-unit         # Unit tests only (no DB)
make test-integration  # Integration tests only
make test-coverage     # With coverage report
```

**Note**: Integration tests use `monorepo_test` database, never `monorepo_dev`.

### Frontend Testing
```bash
cd apps/frontend
npm run test              # Watch mode
npm run test:ui           # Vitest UI
npm run test:run          # Single run
npm run test:coverage     # With coverage
```

## Architecture

### Monorepo Structure
This is a Turborepo-based monorepo with workspace support. The `packages/types` package contains shared TypeScript types used by the frontend. Changes to `packages/types` require rebuilding before the frontend can pick them up.

### Backend Architecture (Go)
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with pgx driver
- **Module path**: `github.com/dwfennell/monorepo-scaffold`
- **Entry point**: `apps/backend/main.go`

**Key layers**:
- `internal/api/` - HTTP handlers and routing
  - `routes.go` - Defines all API endpoints
  - `*_handler.go` - Request handlers
  - `middleware.go` - JWT authentication middleware
- `internal/repository/` - Database access layer
- `internal/models/` - Domain models (User, etc.)
- `internal/auth/` - JWT and password utilities
- `internal/database/` - Database connection management
- `internal/testutil/` - Test helpers for database setup

**Authentication flow**:
1. Public routes: `/api/v1/auth/register`, `/api/v1/auth/login`
2. Protected routes use `AuthMiddleware()` which validates JWT tokens
3. JWT secret configured via `JWT_SECRET` env var
4. Tokens stored in localStorage on frontend, sent via `Authorization: Bearer {token}` header

**Database connection**:
- Main app uses `DATABASE_URL`
- Tests use `TEST_DATABASE_URL` from `.env.test`
- Connection pooling managed by pgx

### Frontend Architecture (React)
- **Framework**: React 19 + React Router v7
- **Build tool**: Vite
- **Styling**: Tailwind CSS v4
- **Testing**: Vitest + React Testing Library

**Key structure**:
- `src/contexts/AuthContext.tsx` - Global auth state management
- `src/components/ProtectedRoute.tsx` - Route guard for authenticated pages
- `src/lib/api.ts` - Centralized API client
- `src/pages/` - Page components (Login, Register, Home)

**API client pattern**:
- All API calls go through `api` object from `src/lib/api.ts`
- Automatically adds `Authorization` header from localStorage
- Throws `APIError` with status code for error handling
- Base URL from `VITE_API_URL` environment variable

### Shared Types Package
- Location: `packages/types/`
- Exports shared TypeScript types (User, AuthResponse, etc.)
- Frontend imports as `@workspace/types`
- Must run `npm run build` in types package after changes for frontend to see updates

## Environment Setup

**.env files required**:
- `apps/backend/.env` - Backend configuration (DATABASE_URL, JWT_SECRET, FRONTEND_URL, PORT)
- `apps/backend/.env.test` - Test database configuration (TEST_DATABASE_URL)
- `apps/frontend/.env` - Frontend configuration (VITE_API_URL)

Copy from `.env.example` files in each directory.

## Testing Strategy

**Backend**:
- Unit tests use `-short` flag to skip database-dependent tests
- Integration tests use testify suites and require test database
- Test database is separate from dev database for safety
- Use `make test-db-setup` to prepare test database

**Frontend**:
- Component tests in `__tests__/` directories next to components
- Use `src/test/utils.tsx` for test utilities
- Tests use jsdom environment via Vitest

## URLs
- Frontend dev server: http://localhost:5173
- Backend API: http://localhost:8080
- Health check: http://localhost:8080/health
