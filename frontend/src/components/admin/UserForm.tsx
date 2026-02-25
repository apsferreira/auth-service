import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import Select from '@/components/ui/Select'
import type { User, UserCreateRequest, UserUpdateRequest, Role } from '@/types/auth'
import adminService from '@/services/adminService'

const createSchema = z.object({
  email: z.string().email('Email invalido'),
  full_name: z.string().min(2, 'Minimo 2 caracteres'),
  role_id: z.string().uuid('Selecione uma role'),
})

const updateSchema = z.object({
  full_name: z.string().min(2, 'Minimo 2 caracteres').optional(),
  role_id: z.string().uuid('Selecione uma role').optional(),
  is_active: z.boolean().optional(),
})

interface UserFormProps {
  user: User | null
  onSubmit: (data: UserCreateRequest | UserUpdateRequest) => Promise<void>
  onCancel: () => void
}

export function UserForm({ user, onSubmit, onCancel }: UserFormProps) {
  const isEditing = !!user
  const [roleOptions, setRoleOptions] = useState<{ value: string; label: string }[]>([])
  const [loadingRoles, setLoadingRoles] = useState(true)

  useEffect(() => {
    adminService
      .listRoles()
      .then((data: Role[]) => {
        setRoleOptions(
          data.map((r) => ({
            value: r.id,
            label: r.name.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase()),
          }))
        )
      })
      .catch(() => {
        setRoleOptions([])
      })
      .finally(() => setLoadingRoles(false))
  }, [])

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm({
    resolver: zodResolver(isEditing ? updateSchema : createSchema),
    defaultValues: isEditing
      ? {
          full_name: user.full_name || '',
          role_id: user.role_id || '',
          is_active: user.is_active,
        }
      : {
          email: '',
          full_name: '',
          role_id: '',
        },
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {!isEditing && (
          <Input
            {...register('email')}
            label="Email"
            type="email"
            placeholder="usuario@email.com"
            error={errors.email?.message as string}
            required
          />
        )}
        <Input
          {...register('full_name')}
          label="Nome completo"
          placeholder="Nome Sobrenome"
          error={errors.full_name?.message as string}
          required={!isEditing}
        />
        <Select
          {...register('role_id')}
          label="Role"
          options={roleOptions}
          placeholder={loadingRoles ? 'Carregando...' : 'Selecione...'}
          error={errors.role_id?.message as string}
          required={!isEditing}
          disabled={loadingRoles}
        />
        {isEditing && (
          <div className="flex items-end">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                {...register('is_active')}
                className="w-4 h-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
              />
              <span className="text-sm text-gray-700 dark:text-gray-200">
                Usuario ativo
              </span>
            </label>
          </div>
        )}
      </div>

      <div className="flex items-center gap-3 pt-2">
        <Button type="submit" loading={isSubmitting}>
          {isEditing ? 'Salvar' : 'Criar usuario'}
        </Button>
        <Button type="button" variant="ghost" onClick={onCancel}>
          Cancelar
        </Button>
      </div>
    </form>
  )
}

export default UserForm
