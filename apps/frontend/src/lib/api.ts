import type { AuthResponse, LoginRequest, RegisterRequest, User } from '../types/auth'

const API_URL = import.meta.env.VITE_API_URL

if (!API_URL) {
  throw new Error(
    'VITE_API_URL is not set. Please copy .env.example to .env and configure your environment variables.'
  )
}

class APIError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message)
    this.name = 'APIError'
  }
}

async function fetchAPI<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('token')

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Unknown error' }))
    throw new APIError(response.status, error.error || 'Request failed')
  }

  return response.json()
}

export const api = {
  auth: {
    register: (data: RegisterRequest) =>
      fetchAPI<AuthResponse>('/api/v1/auth/register', {
        method: 'POST',
        body: JSON.stringify(data),
      }),

    login: (data: LoginRequest) =>
      fetchAPI<AuthResponse>('/api/v1/auth/login', {
        method: 'POST',
        body: JSON.stringify(data),
      }),

    getCurrentUser: () => fetchAPI<User>('/api/v1/me'),
  },

  health: () => fetchAPI<{ status: string }>('/health'),
}

export { APIError }
