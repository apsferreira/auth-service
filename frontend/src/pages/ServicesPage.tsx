import { useState, useEffect, useCallback } from 'react'
import { Plus, Pencil, Trash2, Loader2, ExternalLink } from 'lucide-react'
import toast from 'react-hot-toast'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import Card from '@/components/ui/Card'
import ServiceForm from '@/components/admin/ServiceForm'
import type { Service, ServiceCreateRequest, ServiceUpdateRequest } from '@/types/auth'
import adminService from '@/services/adminService'

export function ServicesPage() {
  const [services, setServices] = useState<Service[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingService, setEditingService] = useState<Service | null>(null)

  const loadServices = useCallback(async () => {
    try {
      const data = await adminService.listServices()
      setServices(data)
    } catch {
      toast.error('Erro ao carregar servicos.')
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    loadServices()
  }, [loadServices])

  const handleCreate = async (data: ServiceCreateRequest) => {
    try {
      await adminService.createService(data)
      toast.success('Servico criado com sucesso!')
      setShowForm(false)
      loadServices()
    } catch {
      toast.error('Erro ao criar servico.')
    }
  }

  const handleUpdate = async (data: ServiceUpdateRequest) => {
    if (!editingService) return
    try {
      await adminService.updateService(editingService.id, data)
      toast.success('Servico atualizado!')
      setEditingService(null)
      loadServices()
    } catch {
      toast.error('Erro ao atualizar servico.')
    }
  }

  const handleDelete = async (service: Service) => {
    if (!confirm(`Deseja realmente excluir o servico "${service.name}"?`)) return
    try {
      await adminService.deleteService(service.id)
      toast.success('Servico excluido!')
      loadServices()
    } catch {
      toast.error('Erro ao excluir servico.')
    }
  }

  const handleEdit = (service: Service) => {
    setShowForm(false)
    setEditingService(service)
  }

  const handleCancelForm = () => {
    setShowForm(false)
    setEditingService(null)
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
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Servicos</h1>
        <Button
          icon={<Plus className="w-4 h-4" />}
          onClick={() => {
            setEditingService(null)
            setShowForm(true)
          }}
        >
          Novo servico
        </Button>
      </div>

      {/* Form inline */}
      {(showForm || editingService) && (
        <Card padding="lg" className="mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            {editingService ? 'Editar servico' : 'Novo servico'}
          </h2>
          <ServiceForm
            service={editingService}
            onSubmit={editingService ? (data) => handleUpdate(data as ServiceUpdateRequest) : (data) => handleCreate(data as ServiceCreateRequest)}
            onCancel={handleCancelForm}
          />
        </Card>
      )}

      {/* Services table */}
      <Card padding="none">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700">
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Nome
                </th>
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Slug
                </th>
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Status
                </th>
                <th className="text-left text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  URLs de redirecionamento
                </th>
                <th className="text-right text-sm font-medium text-gray-500 dark:text-gray-400 px-4 py-3">
                  Acoes
                </th>
              </tr>
            </thead>
            <tbody>
              {services.map((service) => (
                <tr
                  key={service.id}
                  className="border-b border-gray-100 dark:border-gray-700/50 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
                >
                  <td className="px-4 py-3">
                    <div>
                      <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
                        {service.name}
                      </p>
                      {service.description && (
                        <p className="text-xs text-gray-400 dark:text-gray-500 truncate max-w-xs">
                          {service.description}
                        </p>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <span className="text-sm font-mono text-gray-600 dark:text-gray-300">
                      {service.slug}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <Badge variant={service.is_active ? 'success' : 'danger'} size="sm">
                      {service.is_active ? 'Ativo' : 'Inativo'}
                    </Badge>
                  </td>
                  <td className="px-4 py-3">
                    {service.redirect_urls && service.redirect_urls.length > 0 ? (
                      <div className="flex flex-col gap-1">
                        {service.redirect_urls.slice(0, 2).map((url) => (
                          <a
                            key={url}
                            href={url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center gap-1 text-xs text-primary-600 dark:text-primary-400 hover:underline truncate max-w-xs"
                          >
                            <ExternalLink className="w-3 h-3 flex-shrink-0" />
                            {url}
                          </a>
                        ))}
                        {service.redirect_urls.length > 2 && (
                          <span className="text-xs text-gray-400">
                            +{service.redirect_urls.length - 2} mais
                          </span>
                        )}
                      </div>
                    ) : (
                      <span className="text-xs text-gray-400">—</span>
                    )}
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-2">
                      <button
                        onClick={() => handleEdit(service)}
                        className="p-1.5 text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 transition-colors"
                        title="Editar"
                      >
                        <Pencil className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => handleDelete(service)}
                        className="p-1.5 text-gray-400 hover:text-danger-600 dark:hover:text-danger-400 transition-colors"
                        title="Excluir"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {services.length === 0 && (
                <tr>
                  <td colSpan={5} className="text-center py-10 text-gray-400 text-sm">
                    Nenhum servico cadastrado.
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

export default ServicesPage
