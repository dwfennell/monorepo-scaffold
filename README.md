# Monorepo Scaffold

A production-ready monorepo template for building full-stack web applications with strong typing and modern tooling.

## Features

- **Backend**: Go with Gin framework
  - RESTful API architecture
  - JWT-based authentication
  - PostgreSQL database with pgx driver
  - Automatic migrations
  - Type-safe request/response handling

- **Frontend**: React + Vite + TypeScript
  - Modern React with hooks
  - Type-safe API client
  - Authentication context and protected routes
  - Clean, responsive UI components

- **Development Experience**
  - Docker Compose for local PostgreSQL
  - Hot reload for both frontend and backend
  - Strongly typed end-to-end
  - Makefile for common tasks

## Prerequisites

- [Go](https://golang.org/dl/) 1.21 or later
- [Node.js](https://nodejs.org/) 18 or later
- [Docker](https://www.docker.com/get-started) and Docker Compose
- [Air](https://github.com/cosmtrek/air) (optional, for Go hot reload): `go install github.com/cosmtrek/air@latest`

## Project Structure

```
.
├── apps/
│   ├── backend/          # Go backend
│   │   ├── internal/
│   │   │   ├── api/      # HTTP handlers and routes
│   │   │   ├── auth/     # JWT and password utilities
│   │   │   ├── database/ # Database connection and migrations
│   │   │   ├── models/   # Data models and DTOs
│   │   │   └── repository/ # Database queries
│   │   └── main.go
│   └── frontend/         # React frontend
│       ├── src/
│       │   ├── components/  # Reusable components
│       │   ├── contexts/    # React contexts (auth, etc.)
│       │   ├── lib/         # API client and utilities
│       │   ├── pages/       # Page components
│       │   └── types/       # TypeScript types
│       └── package.json
├── packages/
│   └── shared-types/     # Shared type definitions (future use)
├── docker-compose.yml
└── Makefile
```

## Quick Start

### 1. Clone and Setup

```bash
# Install dependencies
make install

# Copy environment files
cp apps/backend/.env.example apps/backend/.env
cp apps/frontend/.env.example apps/frontend/.env

# Update apps/backend/.env with your preferred settings
```

### 2. Start Development

**Option A: Using Make (recommended)**

```bash
# Start database
make dev-db

# In a second terminal, start backend
make dev-backend

# In a third terminal, start frontend
make dev-frontend
```

**Option B: Manual**

```bash
# Start database
docker-compose up -d

# Start backend (with hot reload if Air is installed)
cd apps/backend
air  # or: go run main.go

# Start frontend
cd apps/frontend
npm run dev
```

### 3. Access the Application

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Health check: http://localhost:8080/health

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }
  ```

- `POST /api/v1/auth/login` - Login
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

### Protected Routes

- `GET /api/v1/me` - Get current user (requires auth header)
  ```
  Authorization: Bearer <token>
  ```

## Development

### Backend

The Go backend uses:
- **Gin** for HTTP routing
- **pgx** for PostgreSQL database access
- **JWT** for authentication
- **bcrypt** for password hashing

To add new endpoints:
1. Define models in `internal/models/`
2. Create repository methods in `internal/repository/`
3. Add handlers in `internal/api/`
4. Register routes in `internal/api/routes.go`

### Frontend

The React frontend uses:
- **Vite** for fast development and building
- **React Router** for navigation
- **Context API** for state management

To add new pages:
1. Create component in `src/pages/`
2. Add route in `src/App.tsx`
3. Update API client in `src/lib/api.ts` if needed

### Database Migrations

Migrations run automatically on startup. To add new migrations:
1. Edit `apps/backend/internal/database/migrations.go`
2. Add your SQL statements to the `migrations` slice
3. Restart the backend

For production, consider using a migration tool like [golang-migrate](https://github.com/golang-migrate/migrate).

## Building for Production

### Backend

```bash
cd apps/backend
go build -o server main.go
./server
```

### Frontend

```bash
cd apps/frontend
npm run build
# Serve the dist/ folder with your preferred web server
```

## Environment Variables

### Backend (.env)

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT signing (change in production!)
- `FRONTEND_URL` - Frontend URL for CORS

### Frontend (.env)

- `VITE_API_URL` - Backend API URL

## Common Tasks

```bash
make help           # Show all available commands
make install        # Install all dependencies
make dev-db         # Start database only
make dev-backend    # Start backend server
make dev-frontend   # Start frontend server
make stop           # Stop all Docker services
make clean          # Clean up generated files
```

## Adding More Features

This scaffold provides a solid foundation. Consider adding:

- [ ] Email verification
- [ ] Password reset
- [ ] User profiles
- [ ] Role-based access control
- [ ] API rate limiting
- [ ] Request logging
- [ ] Error tracking (Sentry, etc.)
- [ ] Testing (Go: testify, Frontend: Vitest)
- [ ] CI/CD pipeline
- [ ] API documentation
- [ ] Shared TypeScript types generation from Go models

## License

MIT
