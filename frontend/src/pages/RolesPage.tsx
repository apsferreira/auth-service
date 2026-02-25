import { useState, useEffect, useCallback } from 'react'
import { Plus, Pencil, ShieldCheck, Loader2 } from 'lucide-react'
import toast from 'react-hot-toast'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Card from '@/components/ui/Card'
import RoleForm from '@/components/admin/RoleForm'
import PermissionMatrix from '@/components/admin/PermissionMatrix'
import type { Role, Permission, RoleCreateRequest, RoleUpdateRequest } from '@/types/auth'
import adminService from '@/services/adminService'

type ActivePanel =
  | { type: 'none' }
  | { type: 'create' }
  | { type: 'edit'; role: Role }
  | { type: 'permissions'; role: Role }

export function RolesPage() {
  const [roles, setRoles] = useState<Role[]>([])
  const [allPermissions, setAllPermissions] = useState<Permission[]>([])
  const [currentPermissionIds, setCurrentPermissionIds] = useState<string[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isLoadingPermissions, setIsLoadingPermissions] = useState(false)
  const [activePanel, setActivePanel] = useState<ActivePanel>({ type: 'none' })

  const loadRoles = useCallback(async () => {
    try {
      const data = await adminService.listRoles()
      setRoles(data)
    } catch {
      toast.error('Erro ao carregar roles.')
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    loadRoles()
  }, [loadRoles])

  const openPermissionMatrix = async (role: Role) => {
    setActivePanel({ type: 'permissions', role })
    setIsLoadingPermissions(true)
    try {
      const [all, current] = await Promise.all([
        adminService.listAllPermissions(),
        adminService.getRolePermissions(role.id),
      ])
      setAllPermissions(all)
      setCurrentPermissionIds(current)
    } catch {
      toast.error('Erro ao carregar permissoes.')
      setActivePanel({ type: 'none' })
    } finally {
      setIsLoadingPermissions(false)
    }
  }

  const handleCreate = async (data: RoleCreateRequest) => {
    try {
      await adminService.createRole(data)
      toast.success('Role criada com sucesso!')
      setActivePanel({ type: 'none' })
      loadRoles()
    } catch {
      toast.error('Erro ao criar role.')
    }
  }

  const handleUpdate = async (data: RoleUpdateRequest) => {
    if (activePanel.type !== 'edit') return
    try {
      await adminService.updateRole(activePanel.role.id, data)
      toast.success('Role atualizada!')
      setActivePanel({ type: 'none' })
      loadRoles()
    } catch {
      toast.error('Erro ao atualizar role.')
    }
  }

  const handleSavePermissions = async (ids: string[]) => {
    if (activePanel.type !== 'permissions') return
    try {
      await adminService.setRolePermissions(activePanel.role.id, { permission_ids: ids })
      toast.success('Permissoes atualizadas!')
      setActivePanel({ type: 'none' })
    } catch {
      toast.error('Erro ao salvar permissoes.')
    }
  }

  const levelBadgeVariant = (level: number) => {
    if (level <= 2) return 'danger' as const
    if (level <= 4) return 'warning' as const
    if (level <= 6) return 'info' as const
    return 'default' as const
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  const panelTitle = () => {
    if (activePanel.type === 'create') return 'Nova role'
    if (activePanel.type === 'edit') return `Editar role: ${activePanel.role.name}`
    if (activePanel.type === 'permissions') return `Permissoes: ${activePanel.role.name}`
    return ''
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Roles</h1>
        <Button
          icon={<Plus className="w-4 h-4" />}
          onClick={() => setActivePanel({ type: 'create' })}
        >
          Nova role
        </Button>
      </div>

      {/* Inline panel */}
      {activePanel.type !== 'none' && (
        <Card padding="lg" className="mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            {panelTitle()}
          </h2>

          {(activePanel.type === 'create') && (
            <RoleForm
              onSubmit={(data) => handleCreate(data as RoleCreateRequest)}
              onCancel={() => setActivePanel({ type: 'none' })}
            />
          )}

          {activePanel.type === 'edit' && (
            <RoleForm
              role={activePanel.role}
              onSubmit={(data) => handleUpdate(data as RoleUpdateRequest)}
              onCancel={() => setActivePanel({ type: 'none' })}
            />
          )}

          {activePanel.type === 'permissions' && (
            isLoadingPermissions ? (
              <div className="flex items-center justify-center py-12">
                <Loader2 className="w-6 h-6 animate-spin text-primary-600" />
              </div>
            ) : (
              <PermissionMatrix
                roleId={activePanel.role.id}
                allPermissions={allPermissions}
                currentPermissionIds={currentPermissionIds}
                onSave={handleSavePermissions}
                onCancel={() => setActivePanel({ type: 'none' })}
              />
            )
          )}
        </Card>
      )}

      {/* Roles table */}
      <Card padding="none">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700">
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Nome
                </th>
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Descricao
                </th>
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Nivel
                </th>
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Tipo
                </th>
                <th className="text-right text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Acoes
                </th>
              </tr>
            </thead>
            <tbody>
              {roles.map((role) => (
                <tr
                  key={role.id}
                  className="border-b border-gray-100 dark:border-gray-700/50 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
                >
                  <td className="px-4 py-3">
                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                      {role.name}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <span className="text-sm text-gray-500 dark:text-gray-400">
                      {role.description || <span className="italic text-gray-300 dark:text-gray-600">—</span>}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <Badge variant={levelBadgeVariant(role.level)} size="sm">
                      Nivel {role.level}
                    </Badge>
                  </td>
                  <td className="px-4 py-3">
                    <Badge variant={role.is_system ? 'warning' : 'default'} size="sm">
                      {role.is_system ? 'Sistema' : 'Personalizada'}
                    </Badge>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-2">
                      <button
                        onClick={() => openPermissionMatrix(role)}
                        className="p-1.5 text-gray-400 hover:text-success-600 dark:hover:text-success-400 transition-colors"
                        title="Editar permissoes"
                      >
                        <ShieldCheck className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => setActivePanel({ type: 'edit', role })}
                        className="p-1.5 text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 transition-colors"
                        title="Editar role"
                        disabled={role.is_system}
                      >
                        <Pencil className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {roles.length === 0 && (
                <tr>
                  <td colSpan={5} className="text-center py-10 text-gray-400 text-sm">
                    Nenhuma role cadastrada.
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

export default RolesPage
