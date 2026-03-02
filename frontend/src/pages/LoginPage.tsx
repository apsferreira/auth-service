import { Navigate } from 'react-router-dom'
import { Sun, Moon } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { useTheme } from '@/contexts/ThemeContext'
import LoginForm from '@/components/auth/LoginForm'
import Card from '@/components/ui/Card'

export function LoginPage() {
  const { isAuthenticated, isLoading } = useAuth()
  const { theme, toggleTheme } = useTheme()

  if (isLoading) return null
  if (isAuthenticated) return <Navigate to="/" replace />

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

      {/* IIT Logo */}
      <div className="mb-8 flex flex-col items-center gap-3">
        <img
          src="/logo-iit.png"
          alt="Instituto Itinerante"
          className="h-16 w-auto object-contain"
          draggable={false}
        />
        <p className="text-sm text-iit-text-secondary dark:text-gray-400 font-medium tracking-wide uppercase">
          Autenticação Centralizada
        </p>
      </div>

      {/* Login card */}
      <Card padding="lg" className="w-full max-w-md shadow-lg border-iit-border-default">
        <LoginForm />
      </Card>

      {/* Footer */}
      <p className="mt-6 text-xs text-iit-text-muted dark:text-gray-500">
        Acesso restrito — painel de administração
      </p>
    </div>
  )
}

export default LoginPage
