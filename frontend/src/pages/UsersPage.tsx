import { useState, useEffect, useCallback } from 'react'
import { Plus, Pencil, Trash2, Loader2 } from 'lucide-react'
import toast from 'react-hot-toast'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Card from '@/components/ui/Card'
import UserForm from '@/components/admin/UserForm'
import type { User, UserCreateRequest, UserUpdateRequest } from '@/types/auth'
import userService from '@/services/userService'
import { formatDate } from '@/lib/utils'

export function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingUser, setEditingUser] = useState<User | null>(null)

  const loadUsers = useCallback(async () => {
    try {
      const data = await userService.list()
      setUsers(data)
    } catch {
      toast.error('Erro ao carregar usuarios.')
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    loadUsers()
  }, [loadUsers])

  const handleCreate = async (data: UserCreateRequest) => {
    try {
      await userService.create(data)
      toast.success('Usuario criado com sucesso!')
      setShowForm(false)
      loadUsers()
    } catch {
      toast.error('Erro ao criar usuario.')
    }
  }

  const handleUpdate = async (data: UserUpdateRequest) => {
    if (!editingUser) return
    try {
      await userService.update(editingUser.id, data)
      toast.success('Usuario atualizado!')
      setEditingUser(null)
      loadUsers()
    } catch {
      toast.error('Erro ao atualizar usuario.')
    }
  }

  const handleDelete = async (user: User) => {
    if (!confirm(`Deseja realmente desativar ${user.full_name || user.email}?`)) return
    try {
      await userService.delete(user.id)
      toast.success('Usuario desativado!')
      loadUsers()
    } catch {
      toast.error('Erro ao desativar usuario.')
    }
  }

  const roleBadgeVariant = (roleName: string) => {
    if (roleName === 'super_admin') return 'danger' as const
    if (roleName === 'admin') return 'warning' as const
    if (roleName === 'manager') return 'info' as const
    return 'default' as const
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Usuarios
        </h1>
        <Button
          icon={<Plus className="w-4 h-4" />}
          onClick={() => setShowForm(true)}
        >
          Novo usuario
        </Button>
      </div>

      {/* Form modal */}
      {(showForm || editingUser) && (
        <Card padding="lg" className="mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            {editingUser ? 'Editar usuario' : 'Novo usuario'}
          </h2>
          <UserForm
            user={editingUser}
            onSubmit={editingUser ? (data) => handleUpdate(data as UserUpdateRequest) : (data) => handleCreate(data as UserCreateRequest)}
            onCancel={() => {
              setShowForm(false)
              setEditingUser(null)
            }}
          />
        </Card>
      )}

      {/* Users table */}
      <Card padding="none">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700">
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Usuario
                </th>
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Role
                </th>
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Status
                </th>
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Ultimo login
                </th>
                <th className="text-right text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Acoes
                </th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr
                  key={user.id}
                  className="border-b border-gray-100 dark:border-gray-700/50 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
                >
                  <td className="px-4 py-3">
                    <div>
                      <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
                        {user.full_name || '-'}
                      </p>
                      <p className="text-xs text-gray-500 dark:text-gray-400">{user.email}</p>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex gap-1">
                      {user.roles?.map((role) => (
                        <Badge key={role.id || role.name} variant={roleBadgeVariant(role.name)} size="sm">
                          {role.name}
                        </Badge>
                      ))}
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <Badge variant={user.is_active ? 'success' : 'danger'} size="sm">
                      {user.is_active ? 'Ativo' : 'Inativo'}
                    </Badge>
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                    {user.last_login_at ? formatDate(user.last_login_at) : 'Nunca'}
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-2">
                      <button
                        onClick={() => setEditingUser(user)}
                        className="p-1.5 text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 transition-colors"
                        title="Editar"
                      >
                        <Pencil className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => handleDelete(user)}
                        className="p-1.5 text-gray-400 hover:text-danger-600 dark:hover:text-danger-400 transition-colors"
                        title="Desativar"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {users.length === 0 && (
                <tr>
                  <td colSpan={5} className="text-center py-8 text-gray-400">
                    Nenhum usuario encontrado.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  )
}

export default UsersPage
