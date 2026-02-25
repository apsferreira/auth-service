import { User, Shield, Clock, Key } from 'lucide-react'
import Card from '@/components/ui/Card'
import Badge from '@/components/ui/Badge'
import { useAuth } from '@/contexts/AuthContext'
import { formatDate, formatDateTime, getInitials } from '@/lib/utils'

export function DashboardPage() {
  const { user, permissions, roles } = useAuth()

  if (!user) return null

  const roleBadgeVariant = (role: string) => {
    switch (role) {
      case 'super_admin':
        return 'danger' as const
      case 'admin':
        return 'warning' as const
      case 'manager':
        return 'info' as const
      default:
        return 'default' as const
    }
  }

  const serviceEntries = Object.entries(permissions)
  const totalPermissions = serviceEntries.reduce((sum, [, perms]) => sum + perms.length, 0)

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6">
        Dashboard
      </h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Profile card */}
        <Card padding="lg">
          <div className="flex items-start gap-4">
            <div className="w-16 h-16 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center text-xl font-bold text-primary-700 dark:text-primary-300">
              {user.full_name ? getInitials(user.full_name) : <User className="w-8 h-8" />}
            </div>
            <div className="flex-1">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                {user.full_name || 'Sem nome'}
              </h2>
              <p className="text-sm text-gray-500 dark:text-gray-400">{user.email}</p>
              <div className="flex gap-2 mt-2">
                {roles.map((role) => (
                  <Badge key={role} variant={roleBadgeVariant(role)} size="sm">
                    {role}
                  </Badge>
                ))}
                {roles.length === 0 && user.roles?.map((r) => (
                  <Badge key={r.name} variant={roleBadgeVariant(r.name)} size="sm">
                    {r.name}
                  </Badge>
                ))}
              </div>
            </div>
          </div>
        </Card>

        {/* Account info */}
        <Card padding="lg">
          <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-4 flex items-center gap-2">
            <Shield className="w-4 h-4" />
            Informacoes da conta
          </h3>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600 dark:text-gray-300">Status</span>
              <Badge variant={user.is_active ? 'success' : 'danger'} size="sm">
                {user.is_active ? 'Ativo' : 'Inativo'}
              </Badge>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600 dark:text-gray-300 flex items-center gap-1">
                <Clock className="w-3.5 h-3.5" />
                Ultimo login
              </span>
              <span className="text-sm text-gray-900 dark:text-gray-100">
                {user.last_login_at ? formatDateTime(user.last_login_at) : 'Primeiro acesso'}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600 dark:text-gray-300">Criado em</span>
              <span className="text-sm text-gray-900 dark:text-gray-100">
                {formatDate(user.created_at)}
              </span>
            </div>
          </div>
        </Card>

        {/* Permissions grouped by service */}
        <Card padding="lg" className="md:col-span-2">
          <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-4 flex items-center gap-2">
            <Key className="w-4 h-4" />
            Permissoes ({totalPermissions})
          </h3>
          {serviceEntries.length > 0 ? (
            <div className="space-y-4">
              {serviceEntries.map(([serviceSlug, perms]) => (
                <div key={serviceSlug}>
                  <p className="text-xs font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wide mb-2">
                    {serviceSlug}
                  </p>
                  <div className="flex flex-wrap gap-2">
                    {perms.map((perm) => (
                      <Badge key={`${serviceSlug}-${perm}`} variant="outline" size="sm">
                        {perm}
                      </Badge>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-gray-400">Nenhuma permissao atribuida.</p>
          )}
        </Card>
      </div>
    </div>
  )
}

export default DashboardPage
