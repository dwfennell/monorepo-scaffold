import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

// Set up environment before importing API
vi.stubEnv('VITE_API_URL', 'http://localhost:8080')

import { api, APIError } from '../api'

// Mock fetch globally
const mockFetch = vi.fn()
globalThis.fetch = mockFetch

describe('API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('APIError', () => {
    it('creates error with status and message', () => {
      const error = new APIError(404, 'Not found')
      expect(error.status).toBe(404)
      expect(error.message).toBe('Not found')
      expect(error.name).toBe('APIError')
    })
  })

  describe('auth.login', () => {
    it('sends login request and returns auth response', async () => {
      const mockResponse = {
        token: 'test-token',
        user: { id: '1', email: 'test@example.com', name: 'Test', created_at: new Date() },
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await api.auth.login({ email: 'test@example.com', password: 'password' })

      expect(result).toEqual(mockResponse)
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/auth/login'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
          body: JSON.stringify({ email: 'test@example.com', password: 'password' }),
        })
      )
    })

    it('throws APIError on failed request', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({ error: 'Invalid credentials' }),
      })

      await expect(
        api.auth.login({ email: 'test@example.com', password: 'wrong' })
      ).rejects.toThrow(APIError)
    })
  })

  describe('auth.getCurrentUser', () => {
    it('includes authorization token when present', async () => {
      localStorage.setItem('token', 'test-token')

      const mockUser = { id: '1', email: 'test@example.com', name: 'Test', created_at: new Date() }
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockUser,
      })

      await api.auth.getCurrentUser()

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/me'),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer test-token',
          }),
        })
      )
    })
  })
})
