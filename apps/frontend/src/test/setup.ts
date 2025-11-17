import { vi } from 'vitest'
import '@testing-library/jest-dom'

// Set environment variables before any tests run
vi.stubEnv('VITE_API_URL', 'http://localhost:8080')
