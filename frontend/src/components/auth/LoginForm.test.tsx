import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import LoginForm from './LoginForm'

// Mock useAuth context
const mockLogin = vi.fn()
vi.mock('@/contexts/AuthContext', () => ({
  useAuth: () => ({ login: mockLogin }),
}))

// Mock react-hot-toast
vi.mock('react-hot-toast', () => ({
  default: {
    success: vi.fn(),
    error: vi.fn(),
  },
}))

import toast from 'react-hot-toast'

function renderLoginForm() {
  return render(<LoginForm />)
}

describe('LoginForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('renderização', () => {
    it('renderiza o título Entrar', () => {
      renderLoginForm()
      expect(screen.getByRole('heading', { name: /entrar/i })).toBeInTheDocument()
    })

    it('renderiza campo de usuario com label associada', () => {
      renderLoginForm()
      const input = screen.getByLabelText(/usuario ou email/i)
      expect(input).toBeInTheDocument()
    })

    it('renderiza campo de senha com label associada', () => {
      renderLoginForm()
      const input = screen.getByLabelText(/senha/i)
      expect(input).toBeInTheDocument()
    })

    it('renderiza botão de submit', () => {
      renderLoginForm()
      expect(screen.getByRole('button', { name: /entrar/i })).toBeInTheDocument()
    })

    it('campo de senha começa como tipo password', () => {
      renderLoginForm()
      const passwordInput = screen.getByLabelText(/senha/i)
      expect(passwordInput).toHaveAttribute('type', 'password')
    })
  })

  describe('visibilidade de senha', () => {
    it('alterna a visibilidade da senha ao clicar no botão', async () => {
      const user = userEvent.setup()
      renderLoginForm()

      const passwordInput = screen.getByLabelText(/senha/i)
      // Find toggle button by its position near the password field
      const toggleBtn = screen.getAllByRole('button').find(
        (btn) => btn !== screen.getByRole('button', { name: /entrar/i })
      )!

      expect(passwordInput).toHaveAttribute('type', 'password')
      await user.click(toggleBtn)
      expect(passwordInput).toHaveAttribute('type', 'text')
      await user.click(toggleBtn)
      expect(passwordInput).toHaveAttribute('type', 'password')
    })
  })

  describe('validação de campos', () => {
    it('exibe erro quando usuario tem menos de 3 caracteres', async () => {
      const user = userEvent.setup()
      renderLoginForm()

      await user.type(screen.getByLabelText(/usuario ou email/i), 'ab')
      await user.click(screen.getByRole('button', { name: /entrar/i }))

      await waitFor(() => {
        expect(screen.getByText(/minimo 3 caracteres/i)).toBeInTheDocument()
      })
    })

    it('exibe erro quando senha tem menos de 8 caracteres', async () => {
      const user = userEvent.setup()
      renderLoginForm()

      await user.type(screen.getByLabelText(/usuario ou email/i), 'usuario')
      await user.type(screen.getByLabelText(/senha/i), '123')
      await user.click(screen.getByRole('button', { name: /entrar/i }))

      await waitFor(() => {
        expect(screen.getByText(/minimo 8 caracteres/i)).toBeInTheDocument()
      })
    })

    it('campo inválido recebe aria-invalid=true', async () => {
      const user = userEvent.setup()
      renderLoginForm()

      await user.click(screen.getByRole('button', { name: /entrar/i }))

      await waitFor(() => {
        const identifierInput = screen.getByLabelText(/usuario ou email/i)
        expect(identifierInput).toHaveAttribute('aria-invalid', 'true')
      })
    })
  })

  describe('submissão', () => {
    it('chama login com identifier e password corretos', async () => {
      const user = userEvent.setup()
      mockLogin.mockResolvedValueOnce(undefined)
      renderLoginForm()

      await user.type(screen.getByLabelText(/usuario ou email/i), 'admin')
      await user.type(screen.getByLabelText(/senha/i), 'senha1234')
      await user.click(screen.getByRole('button', { name: /entrar/i }))

      await waitFor(() => {
        expect(mockLogin).toHaveBeenCalledWith('admin', 'senha1234')
      })
    })

    it('exibe toast de sucesso após login bem-sucedido', async () => {
      const user = userEvent.setup()
      mockLogin.mockResolvedValueOnce(undefined)
      renderLoginForm()

      await user.type(screen.getByLabelText(/usuario ou email/i), 'admin')
      await user.type(screen.getByLabelText(/senha/i), 'senha1234')
      await user.click(screen.getByRole('button', { name: /entrar/i }))

      await waitFor(() => {
        expect(toast.success).toHaveBeenCalledWith('Login realizado com sucesso!')
      })
    })

    it('exibe toast de erro quando login falha', async () => {
      const user = userEvent.setup()
      mockLogin.mockRejectedValueOnce(new Error('Unauthorized'))
      renderLoginForm()

      await user.type(screen.getByLabelText(/usuario ou email/i), 'admin')
      await user.type(screen.getByLabelText(/senha/i), 'senhaerrada')
      await user.click(screen.getByRole('button', { name: /entrar/i }))

      await waitFor(() => {
        expect(toast.error).toHaveBeenCalledWith('Credenciais invalidas.')
      })
    })

    it('exibe mensagem de erro da API quando disponível', async () => {
      const user = userEvent.setup()
      const apiError = { response: { data: { error: 'Conta desativada' } } }
      mockLogin.mockRejectedValueOnce(Object.assign(new Error(), apiError))
      renderLoginForm()

      await user.type(screen.getByLabelText(/usuario ou email/i), 'admin')
      await user.type(screen.getByLabelText(/senha/i), 'senha1234')
      await user.click(screen.getByRole('button', { name: /entrar/i }))

      await waitFor(() => {
        expect(toast.error).toHaveBeenCalledWith('Conta desativada')
      })
    })

    it('trim no identifier antes de submeter', async () => {
      const user = userEvent.setup()
      mockLogin.mockResolvedValueOnce(undefined)
      renderLoginForm()

      await user.type(screen.getByLabelText(/usuario ou email/i), '  admin  ')
      await user.type(screen.getByLabelText(/senha/i), 'senha1234')
      await user.click(screen.getByRole('button', { name: /entrar/i }))

      await waitFor(() => {
        expect(mockLogin).toHaveBeenCalledWith('admin', 'senha1234')
      })
    })
  })
})
