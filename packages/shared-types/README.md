# Shared Types

This package contains shared type definitions and utilities that can be used across both frontend and backend applications.

## Purpose

As your monorepo grows, you may want to share type definitions between your Go backend and TypeScript frontend. This directory provides a place to:

- Define shared data structures
- Generate TypeScript types from Go structs
- Store common validation schemas
- Share constants and enums

## Usage Examples

### Option 1: Manual Type Definitions

Create parallel type definitions in both languages and keep them in sync:

```typescript
// shared-types/user.ts
export interface User {
  id: number
  email: string
  name: string
  created_at: string
  updated_at: string
}
```

### Option 2: Generate TypeScript from Go (Recommended)

Use tools like [tygo](https://github.com/gzuidhof/tygo) or [tscriptify](https://github.com/tkrajina/typescriptify-golang-structs) to automatically generate TypeScript types from your Go structs:

```bash
# Install tygo
go install github.com/gzuidhof/tygo@latest

# Create tygo.yaml config
# Generate types
tygo generate
```

Example Go struct with tygo annotations:

```go
package models

// @typescript-type
type User struct {
    ID        int       `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

This generates:

```typescript
export interface User {
  id: number
  email: string
  name: string
  created_at: string
  updated_at: string
}
```

### Option 3: JSON Schema

Use JSON Schema as a language-agnostic format:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "User",
  "type": "object",
  "properties": {
    "id": { "type": "integer" },
    "email": { "type": "string", "format": "email" },
    "name": { "type": "string" }
  },
  "required": ["id", "email", "name"]
}
```

## Current Status

This directory is currently a placeholder. As your application grows and you find yourself duplicating type definitions between frontend and backend, move them here and set up a type generation workflow.

## Recommended Tools

- [tygo](https://github.com/gzuidhof/tygo) - Generate TypeScript types from Go
- [tscriptify](https://github.com/tkrajina/typescriptify-golang-structs) - Alternative Go to TypeScript converter
- [quicktype](https://quicktype.io/) - Generate types from JSON/JSON Schema
- [zod](https://github.com/colinhacks/zod) - Runtime type validation for TypeScript
