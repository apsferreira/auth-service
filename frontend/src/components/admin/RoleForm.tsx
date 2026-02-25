import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import type { Role, RoleCreateRequest, RoleUpdateRequest } from '@/types/auth'

const schema = z.object({
  name: z.string().min(2, 'Minimo 2 caracteres'),
  description: z.string().optional(),
  level: z
    .number({ invalid_type_error: 'Nivel deve ser um numero' })
    .int('Apenas numeros inteiros')
    .min(1, 'Minimo 1')
    .max(10, 'Maximo 10'),
})

type FormValues = z.infer<typeof schema>

interface RoleFormProps {
  role?: Role | null
  onSubmit: (data: RoleCreateRequest | RoleUpdateRequest) => Promise<void>
  onCancel: () => void
}

export function RoleForm({ role, onSubmit, onCancel }: RoleFormProps) {
  const isEditing = !!role

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: isEditing
      ? {
          name: role.name,
          description: role.description || '',
          level: role.level,
        }
      : {
          name: '',
          description: '',
          level: 5,
        },
  })

  const handleFormSubmit = async (values: FormValues) => {
    const payload: RoleCreateRequest | RoleUpdateRequest = {
      name: values.name,
      description: values.description || undefined,
      level: values.level,
    }
    await onSubmit(payload)
  }

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Input
          {...register('name')}
          label="Nome"
          placeholder="manager, viewer, editor..."
          error={errors.name?.message}
          required
        />
        <Input
          {...register('level', { valueAsNumber: true })}
          label="Nivel (1–10)"
          type="number"
          min={1}
          max={10}
          placeholder="5"
          helperText="1 = maior privilegio, 10 = menor privilegio"
          error={errors.level?.message}
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1.5">
          Descricao
        </label>
        <textarea
          {...register('description')}
          placeholder="Descricao da role (opcional)"
          rows={2}
          className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder:text-gray-400 dark:placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors resize-none"
        />
      </div>

      <div className="flex items-center gap-3 pt-2">
        <Button type="submit" loading={isSubmitting}>
          {isEditing ? 'Salvar alteracoes' : 'Criar role'}
        </Button>
        <Button type="button" variant="ghost" onClick={onCancel}>
          Cancelar
        </Button>
      </div>
    </form>
  )
}

export default RoleForm
