import { useState, useEffect, useCallback } from 'react'
import { Plus, Trash2, Loader2, ShieldCheck } from 'lucide-react'
import toast from 'react-hot-toast'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Card from '@/components/ui/Card'
import Select from '@/components/ui/Select'
import PermissionForm from '@/components/admin/PermissionForm'
import type { Service, Permission, PermissionCreateRequest } from '@/types/auth'
import adminService from '@/services/adminService'

export function PermissionsPage() {
  const [services, setServices] = useState<Service[]>([])
  const [selectedServiceId, setSelectedServiceId] = useState<string>('')
  const [permissions, setPermissions] = useState<Permission[]>([])
  const [isLoadingServices, setIsLoadingServices] = useState(true)
  const [isLoadingPermissions, setIsLoadingPermissions] = useState(false)
  const [showForm, setShowForm] = useState(false)

  const loadServices = useCallback(async () => {
    try {
      const data = await adminService.listServices()
      setServices(data)
      if (data.length > 0) {
        setSelectedServiceId(data[0].id)
      }
    } catch {
      toast.error('Erro ao carregar servicos.')
    } finally {
      setIsLoadingServices(false)
    }
  }, [])

  const loadPermissions = useCallback(async (serviceId: string) => {
    if (!serviceId) return
    setIsLoadingPermissions(true)
    try {
      const data = await adminService.listServicePermissions(serviceId)
      setPermissions(data)
    } catch {
      toast.error('Erro ao carregar permissoes.')
    } finally {
      setIsLoadingPermissions(false)
    }
  }, [])

  useEffect(() => {
    loadServices()
  }, [loadServices])

  useEffect(() => {
    if (selectedServiceId) {
      setShowForm(false)
      loadPermissions(selectedServiceId)
    }
  }, [selectedServiceId, loadPermissions])

  const handleCreate = async (data: PermissionCreateRequest) => {
    try {
      await adminService.createServicePermission(selectedServiceId, data)
      toast.success('Permissao criada com sucesso!')
      setShowForm(false)
      loadPermissions(selectedServiceId)
    } catch {
      toast.error('Erro ao criar permissao.')
    }
  }

  const handleDelete = async (permission: Permission) => {
    if (!confirm(`Deseja realmente excluir a permissao "${permission.name}"?`)) return
    try {
      await adminService.deletePermission(permission.id)
      toast.success('Permissao excluida!')
      loadPermissions(selectedServiceId)
    } catch {
      toast.error('Erro ao excluir permissao.')
    }
  }

  const serviceOptions = services.map((s) => ({ value: s.id, label: `${s.name} (${s.slug})` }))

  const actionVariant = (action: string) => {
    switch (action) {
      case 'manage': return 'danger' as const
      case 'delete': return 'warning' as const
      case 'create': return 'success' as const
      case 'update': return 'info' as const
      default: return 'default' as const
    }
  }

  if (isLoadingServices) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Permissoes</h1>
        {selectedServiceId && (
          <Button
            icon={<Plus className="w-4 h-4" />}
            onClick={() => setShowForm(true)}
            disabled={!selectedServiceId}
          >
            Nova permissao
          </Button>
        )}
      </div>

      {/* Service selector */}
      <Card padding="md" className="mb-4">
        <div className="flex items-center gap-4">
          <ShieldCheck className="w-5 h-5 text-gray-400 flex-shrink-0" />
          <div className="flex-1 max-w-sm">
            <Select
              label="Servico"
              value={selectedServiceId}
              onChange={(e) => setSelectedServiceId(e.target.value)}
              options={serviceOptions}
              placeholder="Selecione um servico..."
            />
          </div>
          {selectedServiceId && (
            <div className="pt-5">
              <Badge variant="info" size="sm">
                {permissions.length} permissao{permissions.length !== 1 ? 'es' : ''}
              </Badge>
            </div>
          )}
        </div>
      </Card>

      {/* Form inline */}
      {showForm && selectedServiceId && (
        <Card padding="lg" className="mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Nova permissao
          </h2>
          <PermissionForm
            onSubmit={handleCreate}
            onCancel={() => setShowForm(false)}
          />
        </Card>
      )}

      {/* Permissions table */}
      {selectedServiceId ? (
        <Card padding="none">
          {isLoadingPermissions ? (
            <div className="flex items-center justify-center py-16">
              <Loader2 className="w-6 h-6 animate-spin text-primary-600" />
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-200 dark:border-gray-700">
                    <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                      Nome
                    </th>
                    <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                      Recurso
                    </th>
                    <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                      Acao
                    </th>
                    <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                      Descricao
                    </th>
                    <th className="text-right text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                      Acoes
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {permissions.map((permission) => (
                    <tr
                      key={permission.id}
                      className="border-b border-gray-100 dark:border-gray-700/50 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
                    >
                      <td className="px-4 py-3">
                        <span className="text-sm font-mono font-medium text-gray-900 dark:text-gray-100">
                          {permission.name}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <span className="text-sm text-gray-600 dark:text-gray-300 font-mono">
                          {permission.resource}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <Badge variant={actionVariant(permission.action)} size="sm">
                          {permission.action}
                        </Badge>
                      </td>
                      <td className="px-4 py-3">
                        <span className="text-sm text-gray-500 dark:text-gray-400">
                          {permission.description || <span className="italic text-gray-300 dark:text-gray-600">—</span>}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-end">
                          <button
                            onClick={() => handleDelete(permission)}
                            className="p-1.5 text-gray-400 hover:text-danger-600 dark:hover:text-danger-400 transition-colors"
                            title="Excluir"
                          >
                            <Trash2 className="w-4 h-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                  {permissions.length === 0 && (
                    <tr>
                      <td colSpan={5} className="text-center py-10 text-gray-400 text-sm">
                        Nenhuma permissao cadastrada para este servico.
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </Card>
      ) : (
        <Card padding="lg">
          <p className="text-center text-gray-400 dark:text-gray-500 text-sm py-6">
            Selecione um servico para visualizar suas permissoes.
          </p>
        </Card>
      )}
    </div>
  )
}

export default PermissionsPage
