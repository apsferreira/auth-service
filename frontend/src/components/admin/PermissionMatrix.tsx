import { useState, useMemo } from 'react'
import { Loader2 } from 'lucide-react'
import Button from '@/components/ui/Button'
import Badge from '@/components/ui/Badge'
import type { Permission } from '@/types/auth'

interface PermissionMatrixProps {
  roleId: string
  allPermissions: Permission[]
  currentPermissionIds: string[]
  onSave: (ids: string[]) => Promise<void>
  onCancel: () => void
}

type PermissionsByService = {
  serviceSlug: string
  permissions: Permission[]
}

const ACTION_ORDER = ['create', 'read', 'update', 'delete', 'manage']

function groupByService(permissions: Permission[]): PermissionsByService[] {
  const map = new Map<string, Permission[]>()

  for (const perm of permissions) {
    const key = perm.service_slug || 'global'
    if (!map.has(key)) map.set(key, [])
    map.get(key)!.push(perm)
  }

  return Array.from(map.entries()).map(([serviceSlug, perms]) => ({
    serviceSlug,
    permissions: [...perms].sort((a, b) => {
      const aOrder = ACTION_ORDER.indexOf(a.action)
      const bOrder = ACTION_ORDER.indexOf(b.action)
      if (aOrder !== bOrder) return aOrder - bOrder
      return a.resource.localeCompare(b.resource)
    }),
  }))
}

export function PermissionMatrix({
  allPermissions,
  currentPermissionIds,
  onSave,
  onCancel,
}: PermissionMatrixProps) {
  const [selected, setSelected] = useState<Set<string>>(new Set(currentPermissionIds))
  const [isSaving, setIsSaving] = useState(false)

  const groups = useMemo(() => groupByService(allPermissions), [allPermissions])

  const toggle = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }

  const toggleGroup = (perms: Permission[]) => {
    const allSelected = perms.every((p) => selected.has(p.id))
    setSelected((prev) => {
      const next = new Set(prev)
      for (const p of perms) {
        if (allSelected) {
          next.delete(p.id)
        } else {
          next.add(p.id)
        }
      }
      return next
    })
  }

  const toggleAll = () => {
    if (selected.size === allPermissions.length) {
      setSelected(new Set())
    } else {
      setSelected(new Set(allPermissions.map((p) => p.id)))
    }
  }

  const handleSave = async () => {
    setIsSaving(true)
    try {
      await onSave(Array.from(selected))
    } finally {
      setIsSaving(false)
    }
  }

  const allChecked = selected.size === allPermissions.length
  const indeterminate = selected.size > 0 && !allChecked

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <label className="flex items-center gap-2 cursor-pointer select-none">
            <input
              type="checkbox"
              checked={allChecked}
              ref={(el) => {
                if (el) el.indeterminate = indeterminate
              }}
              onChange={toggleAll}
              className="w-4 h-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
            />
            <span className="text-sm font-medium text-gray-700 dark:text-gray-200">
              Selecionar todas
            </span>
          </label>
          <Badge variant="info" size="sm">
            {selected.size} / {allPermissions.length} selecionadas
          </Badge>
        </div>
      </div>

      <div className="border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
        {groups.map((group, gi) => {
          const groupAllSelected = group.permissions.every((p) => selected.has(p.id))
          const groupIndeterminate =
            group.permissions.some((p) => selected.has(p.id)) && !groupAllSelected

          return (
            <div key={group.serviceSlug}>
              {/* Group header */}
              <div className="flex items-center gap-3 px-4 py-2.5 bg-gray-50 dark:bg-gray-900/50 border-b border-gray-200 dark:border-gray-700">
                <input
                  type="checkbox"
                  checked={groupAllSelected}
                  ref={(el) => {
                    if (el) el.indeterminate = groupIndeterminate
                  }}
                  onChange={() => toggleGroup(group.permissions)}
                  className="w-4 h-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                />
                <span className="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                  {group.serviceSlug}
                </span>
                <Badge variant="default" size="sm">
                  {group.permissions.filter((p) => selected.has(p.id)).length} / {group.permissions.length}
                </Badge>
              </div>

              {/* Permissions rows */}
              {group.permissions.map((perm, pi) => {
                const isChecked = selected.has(perm.id)
                const isLast =
                  gi === groups.length - 1 && pi === group.permissions.length - 1

                return (
                  <label
                    key={perm.id}
                    className={`flex items-center gap-4 px-4 py-3 cursor-pointer transition-colors hover:bg-gray-50 dark:hover:bg-gray-800/50 ${
                      !isLast ? 'border-b border-gray-100 dark:border-gray-700/50' : ''
                    }`}
                  >
                    <input
                      type="checkbox"
                      checked={isChecked}
                      onChange={() => toggle(perm.id)}
                      className="w-4 h-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500 flex-shrink-0"
                    />
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 flex-wrap">
                        <span className="text-sm font-mono font-medium text-gray-900 dark:text-gray-100">
                          {perm.name}
                        </span>
                        <Badge
                          variant={
                            perm.action === 'manage'
                              ? 'danger'
                              : perm.action === 'delete'
                              ? 'warning'
                              : perm.action === 'create' || perm.action === 'update'
                              ? 'info'
                              : 'default'
                          }
                          size="sm"
                        >
                          {perm.action}
                        </Badge>
                      </div>
                      {perm.description && (
                        <p className="text-xs text-gray-400 dark:text-gray-500 mt-0.5 truncate">
                          {perm.description}
                        </p>
                      )}
                    </div>
                    <span className="text-xs text-gray-400 dark:text-gray-500 font-mono flex-shrink-0">
                      {perm.resource}
                    </span>
                  </label>
                )
              })}
            </div>
          )
        })}

        {allPermissions.length === 0 && (
          <div className="text-center py-10 text-gray-400 dark:text-gray-500 text-sm">
            Nenhuma permissao disponivel.
          </div>
        )}
      </div>

      <div className="flex items-center gap-3 pt-2">
        <Button onClick={handleSave} loading={isSaving}>
          {isSaving ? (
            <>
              <Loader2 className="w-4 h-4 animate-spin" />
              Salvando...
            </>
          ) : (
            'Salvar permissoes'
          )}
        </Button>
        <Button variant="ghost" onClick={onCancel} disabled={isSaving}>
          Cancelar
        </Button>
      </div>
    </div>
  )
}

export default PermissionMatrix
