import { Navigate } from 'react-router-dom'
import { Shield, Sun, Moon } from 'lucide-react'
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
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 dark:bg-secondary-900 px-4">
      {/* Theme toggle */}
      <button
        onClick={toggleTheme}
        className="absolute top-4 right-4 p-2 rounded-lg text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
        title={theme === 'dark' ? 'Modo claro' : 'Modo escuro'}
      >
        {theme === 'dark' ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
      </button>

      {/* Logo */}
      <div className="flex items-center gap-3 mb-8">
        <div className="w-12 h-12 bg-primary-600 rounded-xl flex items-center justify-center">
          <Shield className="w-7 h-7 text-white" />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Auth Service</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">Autenticacao centralizada</p>
        </div>
      </div>

      {/* Login card */}
      <Card padding="lg" className="w-full max-w-md">
        <LoginForm />
      </Card>

      {/* Footer */}
      <p className="mt-6 text-xs text-gray-400 dark:text-gray-500">
        Acesso restrito — painel de administracao
      </p>
    </div>
  )
}

export default LoginPage
