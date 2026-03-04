import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import OTPInput from './OTPInput'

function renderOTP(props: Partial<Parameters<typeof OTPInput>[0]> = {}) {
  const onChange = vi.fn()
  const result = render(
    <OTPInput value="" onChange={onChange} length={6} {...props} />
  )
  return { ...result, onChange }
}

describe('OTPInput', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('renderização', () => {
    it('renderiza 6 inputs por padrão', () => {
      renderOTP()
      const inputs = screen.getAllByRole('textbox')
      expect(inputs).toHaveLength(6)
    })

    it('renderiza quantidade customizada de inputs', () => {
      renderOTP({ length: 4 })
      const inputs = screen.getAllByRole('textbox')
      expect(inputs).toHaveLength(4)
    })

    it('cada input tem aria-label com posição', () => {
      renderOTP()
      expect(screen.getByLabelText('Dígito 1 de 6')).toBeInTheDocument()
      expect(screen.getByLabelText('Dígito 6 de 6')).toBeInTheDocument()
    })

    it('grupo tem label acessível', () => {
      renderOTP({ label: 'Código OTP' })
      expect(screen.getByRole('group', { name: 'Código OTP' })).toBeInTheDocument()
    })

    it('exibe valores passados via prop value', () => {
      renderOTP({ value: '123' })
      const inputs = screen.getAllByRole('textbox')
      expect(inputs[0]).toHaveValue('1')
      expect(inputs[1]).toHaveValue('2')
      expect(inputs[2]).toHaveValue('3')
      expect(inputs[3]).toHaveValue('')
    })
  })

  describe('interação', () => {
    it('chama onChange ao digitar um dígito', async () => {
      const user = userEvent.setup()
      const { onChange } = renderOTP()

      const firstInput = screen.getByLabelText('Dígito 1 de 6')
      await user.type(firstInput, '5')

      expect(onChange).toHaveBeenCalledWith('5')
    })

    it('ignora caracteres não-numéricos', async () => {
      const user = userEvent.setup()
      const { onChange } = renderOTP()

      const firstInput = screen.getByLabelText('Dígito 1 de 6')
      await user.type(firstInput, 'a')

      expect(onChange).not.toHaveBeenCalled()
    })

    it('inputs ficam desabilitados quando disabled=true', () => {
      renderOTP({ disabled: true })
      const inputs = screen.getAllByRole('textbox')
      inputs.forEach((input) => expect(input).toBeDisabled())
    })
  })

  describe('acessibilidade em estado de erro', () => {
    it('não exibe mensagem de erro quando error é undefined', () => {
      renderOTP()
      expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    })

    it('exibe mensagem de erro com role=alert', () => {
      renderOTP({ error: 'Código inválido ou expirado.' })
      const alert = screen.getByRole('alert')
      expect(alert).toHaveTextContent('Código inválido ou expirado.')
    })

    it('inputs recebem aria-invalid=true quando há erro', () => {
      renderOTP({ error: 'Código expirado.' })
      const inputs = screen.getAllByRole('textbox')
      inputs.forEach((input) => {
        expect(input).toHaveAttribute('aria-invalid', 'true')
      })
    })

    it('inputs recebem aria-invalid=false quando não há erro', () => {
      renderOTP()
      const inputs = screen.getAllByRole('textbox')
      inputs.forEach((input) => {
        expect(input).toHaveAttribute('aria-invalid', 'false')
      })
    })

    it('inputs apontam para o id do erro via aria-describedby', () => {
      renderOTP({ error: 'Código inválido.' })
      const inputs = screen.getAllByRole('textbox')
      inputs.forEach((input) => {
        expect(input).toHaveAttribute('aria-describedby', 'otp-error')
      })
      expect(screen.getByRole('alert')).toHaveAttribute('id', 'otp-error')
    })
  })

  describe('feedback de erro granular', () => {
    it('exibe erro de código inválido', () => {
      renderOTP({ error: 'Código inválido. Verifique e tente novamente.' })
      expect(screen.getByRole('alert')).toHaveTextContent('Código inválido. Verifique e tente novamente.')
    })

    it('exibe erro de código expirado', () => {
      renderOTP({ error: 'Código expirado. Solicite um novo código.' })
      expect(screen.getByRole('alert')).toHaveTextContent('Código expirado. Solicite um novo código.')
    })

    it('exibe erro de tentativas excedidas', () => {
      renderOTP({ error: 'Muitas tentativas. Aguarde e tente novamente.' })
      expect(screen.getByRole('alert')).toHaveTextContent('Muitas tentativas. Aguarde e tente novamente.')
    })
  })
})
