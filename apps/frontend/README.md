# Frontend

React application with TypeScript, Tailwind CSS, and Vitest.

## Setup

```bash
# From project root (recommended)
npm install

# Or from this directory
npm install
```

## Running

```bash
# From project root
npm run dev

# Or from this directory
npm run dev
```

Server runs on http://localhost:5173

## Environment Variables

```bash
cp .env.example .env
```

Required:
- `VITE_API_URL` - Backend API URL (default: http://localhost:8080)

## Scripts

```bash
npm run dev          # Development server
npm run build        # Production build
npm run test         # Run tests
npm run test:ui      # Vitest UI
npm run lint         # ESLint
npm run type-check   # TypeScript validation
```

## Testing

Uses Vitest + React Testing Library. Tests are in `__tests__/` directories.

```bash
npm run test              # Watch mode
npm run test:ui           # Visual UI
npm run test:coverage     # With coverage
```
