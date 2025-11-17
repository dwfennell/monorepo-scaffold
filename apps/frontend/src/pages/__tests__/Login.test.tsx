import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '../../test/utils'
import userEvent from '@testing-library/user-event'
import { Login } from '../Login'
import { APIError } from '../../lib/api'

// Mock the AuthContext
vi.mock('../../contexts/AuthContext', () => ({
  useAuth: vi.fn(),
}))

// Mock react-router-dom
const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
    Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
      <a href={to}>{children}</a>
    ),
  }
})

import { useAuth } from '../../contexts/AuthContext'

describe('Login', () => {
  const mockLogin = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      loading: false,
      login: mockLogin,
      register: vi.fn(),
      logout: vi.fn(),
    })
  })

  it('renders login form', () => {
    render(<Login />)

    expect(screen.getByRole('heading', { name: /login/i })).toBeInTheDocument()
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument()
  })

  it('submits form with email and password', async () => {
    const user = userEvent.setup()
    mockLogin.mockResolvedValueOnce(undefined)

    render(<Login />)

    await user.type(screen.getByLabelText(/email/i), 'test@example.com')
    await user.type(screen.getByLabelText(/password/i), 'password123')
    await user.click(screen.getByRole('button', { name: /login/i }))

    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password123',
      })
    })

    expect(mockNavigate).toHaveBeenCalledWith('/')
  })

  it('displays error message on login failure', async () => {
    const user = userEvent.setup()
    const errorMessage = 'Invalid credentials'
    mockLogin.mockRejectedValueOnce(new APIError(401, errorMessage))

    render(<Login />)

    await user.type(screen.getByLabelText(/email/i), 'test@example.com')
    await user.type(screen.getByLabelText(/password/i), 'wrongpassword')
    await user.click(screen.getByRole('button', { name: /login/i }))

    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument()
    })
  })

  it('disables form inputs while loading', async () => {
    const user = userEvent.setup()
    mockLogin.mockImplementation(() => new Promise((resolve) => setTimeout(resolve, 100)))

    render(<Login />)

    const emailInput = screen.getByLabelText(/email/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const submitButton = screen.getByRole('button', { name: /login/i })

    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    await user.click(submitButton)

    expect(emailInput).toBeDisabled()
    expect(passwordInput).toBeDisabled()
    expect(submitButton).toBeDisabled()
    expect(screen.getByText(/logging in/i)).toBeInTheDocument()
  })

  it('shows link to register page', () => {
    render(<Login />)

    const registerLink = screen.getByRole('link', { name: /register/i })
    expect(registerLink).toBeInTheDocument()
    expect(registerLink).toHaveAttribute('href', '/register')
  })
})
