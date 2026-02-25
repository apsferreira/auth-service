import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import type { PermissionCreateRequest } from '@/types/auth'

const ACTIONS = [
  { value: 'create', label: 'create — Criar' },
  { value: 'read', label: 'read — Ler' },
  { value: 'update', label: 'update — Atualizar' },
  { value: 'delete', label: 'delete — Deletar' },
  { value: 'manage', label: 'manage — Gerenciar (tudo)' },
]

const schema = z.object({
  resource: z.string().min(1, 'Recurso obrigatorio').regex(/^[a-z0-9_]+$/, 'Apenas letras minusculas, numeros e underscore'),
  action: z.enum(['create', 'read', 'update', 'delete', 'manage'], {
    errorMap: () => ({ message: 'Selecione uma acao' }),
  }),
  description: z.string().optional(),
})

type FormValues = z.infer<typeof schema>

interface PermissionFormProps {
  onSubmit: (data: PermissionCreateRequest) => Promise<void>
  onCancel: () => void
}

export function PermissionForm({ onSubmit, onCancel }: PermissionFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      resource: '',
      action: undefined,
      description: '',
    },
  })

  const resource = watch('resource')
  const action = watch('action')
  const generatedName = resource && action ? `${resource}.${action}` : ''

  const handleFormSubmit = async (values: FormValues) => {
    const payload: PermissionCreateRequest = {
      name: `${values.resource}.${values.action}`,
      resource: values.resource,
      action: values.action,
      description: values.description || undefined,
    }
    await onSubmit(payload)
  }

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      {generatedName && (
        <div className="px-3 py-2 rounded-lg bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700">
          <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5">Nome gerado</p>
          <p className="text-sm font-mono font-semibold text-primary-700 dark:text-primary-300">
            {generatedName}
          </p>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Input
          {...register('resource')}
          label="Recurso"
          placeholder="users, orders, reports..."
          helperText="Identificador do recurso protegido"
          error={errors.resource?.message}
          required
        />
        <Select
          {...register('action')}
          label="Acao"
          options={ACTIONS}
          placeholder="Selecione a acao..."
          error={errors.action?.message}
          required
        />
      </div>

      <Input
        {...register('description')}
        label="Descricao"
        placeholder="Permite criar novos usuarios no sistema"
        helperText="Explique o que esta permissao autoriza (opcional)"
        error={errors.description?.message}
      />

      <div className="flex items-center gap-3 pt-2">
        <Button type="submit" loading={isSubmitting}>
          Criar permissao
        </Button>
        <Button type="button" variant="ghost" onClick={onCancel}>
          Cancelar
        </Button>
      </div>
    </form>
  )
}

export default PermissionForm
