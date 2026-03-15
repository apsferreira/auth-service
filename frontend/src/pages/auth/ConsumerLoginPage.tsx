import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Mail, ArrowLeft, Chrome, Sun, Moon, AlertTriangle } from 'lucide-react'
import toast from 'react-hot-toast'
import { useTheme } from '@/contexts/ThemeContext'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import OTPInput from '@/components/auth/OTPInput'
import { authService } from '@/services/authService'

type Step = 'email' | 'otp' | 'loading'

const PRODUCT_META: Record<string, { name: string; color: string; bg: string; border: string }> = {
  libri: {
    name: 'Libri',
    color: 'text-amber-800',
    bg: 'bg-amber-50',
    border: 'border-amber-300',
  },
  nitro: {
    name: 'Nitro',
    color: 'text-blue-900',
    bg: 'bg-blue-50',
    border: 'border-blue-300',
  },
  default: {
    name: 'Instituto Itinerante',
    color: 'text-gray-800',
    bg: 'bg-gray-50',
    border: 'border-gray-300',
  },
}

function detectProduct(redirectUri: string): string {
  if (!redirectUri) return 'default'
  if (redirectUri.includes('biblioteca') || redirectUri.includes('libri')) return 'libri'
  if (redirectUri.includes('focus-hub') || redirectUri.includes('nitro')) return 'nitro'
  return 'default'
}

function isValidRedirectUri(uri: string): boolean {
  try {
    const url = new URL(uri)
    const allowed = [
      'institutoitinerante.com.br',
      'localhost',
      '127.0.0.1',
    ]
    return allowed.some((h) => url.hostname === h || url.hostname.endsWith('.' + h))
  } catch {
    return false
  }
}

export function ConsumerLoginPage() {
  const [searchParams] = useSearchParams()
  const { theme, toggleTheme } = useTheme()

  const redirectUri = searchParams.get('redirect_uri') ?? ''
  const productKey = searchParams.get('product') ?? detectProduct(redirectUri)
  const product = PRODUCT_META[productKey] ?? PRODUCT_META.default

  const [step, setStep] = useState<Step>('email')
  const [email, setEmail] = useState('')
  const [otp, setOtp] = useState('')
  const [otpError, setOtpError] = useState('')
  const [emailError, setEmailError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [resendCooldown, setResendCooldown] = useState(0)

  const invalidUri = !redirectUri || !isValidRedirectUri(redirectUri)

  // Cooldown timer for resend
  useEffect(() => {
    if (resendCooldown <= 0) return
    const t = setTimeout(() => setResendCooldown((s) => s - 1), 1000)
    return () => clearTimeout(t)
  }, [resendCooldown])

  const handleRequestOTP = async (e?: React.FormEvent) => {
    e?.preventDefault()
    setEmailError('')

    if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setEmailError('Informe um e-mail válido.')
      return
    }

    setIsSubmitting(true)
    try {
      await authService.requestOTP(email)
      setStep('otp')
      setResendCooldown(60)
      toast.success(`Código enviado para ${email}`)
    } catch {
      setEmailError('Não foi possível enviar o código. Tente novamente.')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleVerifyOTP = async (e?: React.FormEvent) => {
    e?.preventDefault()
    setOtpError('')

    if (otp.length !== 6) {
      setOtpError('Informe o código de 6 dígitos.')
      return
    }

    setIsSubmitting(true)
    try {
      const resp = await authService.verifyOTP(email, otp)
      setStep('loading')

      const target = new URL(redirectUri)
      target.searchParams.set('access_token', resp.access_token)
      target.searchParams.set('refresh_token', resp.refresh_token)
      target.searchParams.set('expires_in', String(resp.expires_in))

      window.location.href = target.toString()
    } catch {
      setOtpError('Código inválido ou expirado. Verifique e tente novamente.')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleGoogleLogin = async () => {
    if (invalidUri) return
    setIsSubmitting(true)
    try {
      const resp = await authService.getGoogleAuthURL(redirectUri)
      window.location.href = resp.url
    } catch {
      toast.error('Google OAuth indisponível no momento.')
    } finally {
      setIsSubmitting(false)
    }
  }

  // ─── Invalid redirect URI guard ───────────────────────────────────────────
  if (invalidUri) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
        <div className="w-full max-w-md bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-8 text-center">
          <AlertTriangle className="w-12 h-12 text-amber-500 mx-auto mb-4" />
          <h1 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">
            Parâmetro inválido
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Esta página requer um parâmetro <code className="font-mono text-xs bg-gray-100 dark:bg-gray-700 px-1 py-0.5 rounded">redirect_uri</code> válido.
            Acesse o login pelo produto desejado.
          </p>
        </div>
      </div>
    )
  }

  // ─── Loading / redirecting ─────────────────────────────────────────────────
  if (step === 'loading') {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
        <div className="text-center">
          <div className="w-10 h-10 border-4 border-primary-500 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
          <p className="text-sm text-gray-500 dark:text-gray-400">Redirecionando…</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-iit-surface-subtle dark:bg-iit-surface-dark-base px-4">
      {/* Theme toggle */}
      <button
        onClick={toggleTheme}
        className="absolute top-4 right-4 p-2 rounded-lg text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 hover:bg-white dark:hover:bg-gray-800 transition-colors"
        title={theme === 'dark' ? 'Modo claro' : 'Modo escuro'}
      >
        {theme === 'dark' ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
      </button>

      {/* IIT Logo + product name */}
      <div className="mb-8 flex flex-col items-center gap-2">
        <img
          src="/logo-iit.png"
          alt="Instituto Itinerante"
          className="h-14 w-auto object-contain"
          draggable={false}
        />
        {productKey !== 'default' && (
          <span className={`text-sm font-semibold tracking-wide px-3 py-1 rounded-full ${product.bg} ${product.color} ${product.border} border`}>
            {product.name}
          </span>
        )}
      </div>

      {/* Card */}
      <div className="w-full max-w-md bg-white dark:bg-gray-800 rounded-2xl shadow-lg border border-iit-border-default dark:border-gray-700 p-8">

        {/* ── Step: email ─────────────────────────────────────────────────── */}
        {step === 'email' && (
          <form onSubmit={handleRequestOTP} noValidate>
            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-1">Entrar</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mb-6">
              Informe seu e-mail para receber o código de acesso.
            </p>

            <Input
              label="E-mail"
              type="email"
              placeholder="voce@exemplo.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              icon={<Mail className="w-5 h-5" />}
              error={emailError}
              autoComplete="email"
              autoFocus
              disabled={isSubmitting}
            />

            <Button
              type="submit"
              fullWidth
              loading={isSubmitting}
              className="mt-4"
            >
              Enviar código
            </Button>

            <div className="mt-6">
              <div className="relative flex items-center">
                <div className="flex-1 border-t border-gray-200 dark:border-gray-700" />
                <span className="mx-3 text-xs text-gray-400 dark:text-gray-500">ou</span>
                <div className="flex-1 border-t border-gray-200 dark:border-gray-700" />
              </div>

              <Button
                type="button"
                variant="secondary"
                fullWidth
                onClick={handleGoogleLogin}
                disabled={isSubmitting}
                className="mt-4 flex items-center justify-center gap-2"
              >
                <Chrome className="w-4 h-4" />
                Continuar com Google
              </Button>
            </div>
          </form>
        )}

        {/* ── Step: OTP ───────────────────────────────────────────────────── */}
        {step === 'otp' && (
          <form onSubmit={handleVerifyOTP} noValidate>
            <button
              type="button"
              onClick={() => { setStep('email'); setOtp(''); setOtpError('') }}
              className="flex items-center gap-1 text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 mb-6 transition-colors"
            >
              <ArrowLeft className="w-4 h-4" />
              Alterar e-mail
            </button>

            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-1">
              Código enviado
            </h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mb-6">
              Enviamos um código de 6 dígitos para{' '}
              <strong className="text-gray-700 dark:text-gray-300">{email}</strong>.
              Verifique sua caixa de entrada.
            </p>

            <OTPInput
              value={otp}
              onChange={setOtp}
              error={otpError}
              disabled={isSubmitting}
            />

            <Button
              type="submit"
              fullWidth
              loading={isSubmitting}
              disabled={otp.length !== 6}
              className="mt-6"
            >
              Verificar e entrar
            </Button>

            <div className="mt-4 text-center">
              {resendCooldown > 0 ? (
                <p className="text-xs text-gray-400 dark:text-gray-500">
                  Reenviar em {resendCooldown}s
                </p>
              ) : (
                <button
                  type="button"
                  onClick={handleRequestOTP}
                  disabled={isSubmitting}
                  className="text-xs text-primary-600 dark:text-primary-400 hover:underline disabled:opacity-50"
                >
                  Não recebeu? Reenviar código
                </button>
              )}
            </div>
          </form>
        )}
      </div>

      <p className="mt-6 text-xs text-iit-text-muted dark:text-gray-500 text-center max-w-xs">
        Ao continuar, você concorda com os{' '}
        <a
          href="https://institutoitinerante.com.br/termos"
          target="_blank"
          rel="noreferrer"
          className="underline hover:text-gray-700 dark:hover:text-gray-300"
        >
          Termos de Uso
        </a>{' '}
        e a{' '}
        <a
          href="https://institutoitinerante.com.br/privacidade"
          target="_blank"
          rel="noreferrer"
          className="underline hover:text-gray-700 dark:hover:text-gray-300"
        >
          Política de Privacidade
        </a>{' '}
        do Instituto Itinerante.
      </p>
    </div>
  )
}

export default ConsumerLoginPage
