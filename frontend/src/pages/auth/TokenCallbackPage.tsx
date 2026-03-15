/**
 * TokenCallbackPage — /auth/callback
 *
 * Handles OAuth redirects where the redirect_uri points back to auth-service itself
 * (e.g., during development or when auth-service acts as a token relay).
 *
 * Reads access_token, refresh_token, expires_in from query params,
 * stores them in localStorage, then redirects to the admin dashboard.
 */
import { useEffect } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'

export function TokenCallbackPage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()

  useEffect(() => {
    const accessToken = searchParams.get('access_token')
    const refreshToken = searchParams.get('refresh_token')

    if (accessToken) {
      localStorage.setItem('access_token', accessToken)
    }
    if (refreshToken) {
      localStorage.setItem('refresh_token', refreshToken)
    }

    navigate('/', { replace: true })
  }, [searchParams, navigate])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
      <div className="text-center">
        <div className="w-10 h-10 border-4 border-primary-500 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
        <p className="text-sm text-gray-500 dark:text-gray-400">Autenticando…</p>
      </div>
    </div>
  )
}

export default TokenCallbackPage
