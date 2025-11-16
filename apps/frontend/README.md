# Frontend

React + Vite + TypeScript frontend with authentication.

## Running

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```

## Environment Variables

Create a `.env` file:

```bash
VITE_API_URL=http://localhost:8080
```

## Project Structure

- `src/components/` - Reusable React components
- `src/contexts/` - React Context providers (auth, etc.)
- `src/lib/` - Utilities and API client
- `src/pages/` - Page components (routes)
- `src/types/` - TypeScript type definitions
