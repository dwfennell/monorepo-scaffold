# Monorepo Scaffold

Full-stack monorepo with React frontend, Go backend, PostgreSQL, and Turborepo.

## Quick Start

```bash
# Install dependencies
npm install

# Copy environment files
cp apps/backend/.env.example apps/backend/.env
cp apps/frontend/.env.example apps/frontend/.env

# Start frontend
npm run dev

# Start backend (in separate terminal)
make dev-backend
```

Frontend: http://localhost:5173
Backend: http://localhost:8080

## Prerequisites

- Node.js 18+
- Go 1.23+
- Docker & Docker Compose

## Essential Commands

```bash
# Development
npm run dev          # Start all apps
npm run build        # Build everything
npm run test         # Run all tests

# Database (via Makefile)
make dev-db          # Start PostgreSQL
make migrate-up      # Run migrations
make migrate-create NAME=name  # Create migration

# See all commands
make help
```

## Environment Variables

**Backend** (`apps/backend/.env`):
```env
DATABASE_URL=postgres://user:pass@localhost:5432/monorepo_dev
JWT_SECRET=your-secret-key
FRONTEND_URL=http://localhost:5173
PORT=8080
```

**Frontend** (`apps/frontend/.env`):
```env
VITE_API_URL=http://localhost:8080
```

## Structure

```
apps/
├── backend/      # Go + Gin + PostgreSQL
└── frontend/     # React + Vite + Tailwind
packages/
└── types/        # Shared TypeScript types
```
