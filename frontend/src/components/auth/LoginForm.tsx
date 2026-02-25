import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { User, Lock, Eye, EyeOff } from 'lucide-react'
import toast from 'react-hot-toast'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import { useAuth } from '@/contexts/AuthContext'

const schema = z.object({
  identifier: z.string().min(3, 'Minimo 3 caracteres'),
  password: z.string().min(8, 'Minimo 8 caracteres'),
})

type FormValues = z.infer<typeof schema>

export function LoginForm() {
  const { login } = useAuth()
  const [showPassword, setShowPassword] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
  })

  const onSubmit = async (data: FormValues) => {
    try {
      await login(data.identifier.trim(), data.password)
      toast.success('Login realizado com sucesso!')
    } catch (error: unknown) {
      const message =
        error instanceof Error && 'response' in error
          ? (error as { response?: { data?: { error?: string } } }).response?.data?.error
          : undefined
      toast.error(message || 'Credenciais invalidas.')
    }
  }

  return (
    <div className="w-full max-w-sm">
      <div className="text-center mb-6">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Entrar
        </h2>
        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Painel de administracao
        </p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <Input
          {...register('identifier')}
          label="Usuario ou email"
          placeholder="apsferreira"
          icon={<User className="w-5 h-5" />}
          error={errors.identifier?.message}
          autoComplete="username"
          autoFocus
          disabled={isSubmitting}
        />

        <div className="relative">
          <Input
            {...register('password')}
            label="Senha"
            type={showPassword ? 'text' : 'password'}
            placeholder="••••••••"
            icon={<Lock className="w-5 h-5" />}
            error={errors.password?.message}
            autoComplete="current-password"
            disabled={isSubmitting}
          />
          <button
            type="button"
            onClick={() => setShowPassword((v) => !v)}
            className="absolute right-3 top-9 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors"
            tabIndex={-1}
          >
            {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
          </button>
        </div>

        <Button type="submit" fullWidth loading={isSubmitting} className="mt-2">
          Entrar
        </Button>
      </form>
    </div>
  )
}

export default LoginForm
