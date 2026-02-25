import { useState, useEffect, useCallback } from 'react'
import { Activity, RefreshCw } from 'lucide-react'
import toast from 'react-hot-toast'
import Card from '@/components/ui/Card'
import Badge from '@/components/ui/Badge'
import Button from '@/components/ui/Button'
import authService from '@/services/authService'
import type { AuthEvent, AuthEventsResponse } from '@/types/auth'

const EVENT_LABELS: Record<string, string> = {
  otp_requested: 'OTP Solicitado',
  login_success: 'Login com Sucesso',
  login_failed: 'Falha no Login',
  logout: 'Logout',
  token_refreshed: 'Token Renovado',
}

const EVENT_VARIANTS: Record<string, 'success' | 'danger' | 'info' | 'warning' | 'default'> = {
  otp_requested: 'info',
  login_success: 'success',
  login_failed: 'danger',
  logout: 'default',
  token_refreshed: 'warning',
}

const PAGE_SIZE = 50

// Keys shown in the Sistema column — hide from Details to avoid duplication
const META_HIDDEN = new Set(['source'])

const META_LABELS: Record<string, (v: unknown) => string> = {
  otp_expires_in_minutes: (v) => `OTP válido por ${v} min`,
  otp_code:               (v) => `Código: ${v}`,
  channel:                (v) => `Canal: ${v === 'email' ? 'Email' : v === 'telegram' ? 'Telegram' : v === 'whatsapp' ? 'WhatsApp' : String(v)}`,
  method:                 (v) => v === 'otp' ? 'Via OTP' : v === 'password' ? 'Via senha' : String(v),
  access_expires_in:      (v) => `Sessão: ${Math.round(Number(v) / 60)} min`,
  session_active_until:   (v) => `Até: ${formatDate(String(v))}`,
  error:                  (v) => `Erro: ${v}`,
}

function formatDate(isoDate: string): string {
  return new Date(isoDate).toLocaleString('pt-BR', {
    day: '2-digit', month: '2-digit', year: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  })
}

function formatMetaEntry(k: string, v: unknown): string {
  const formatter = META_LABELS[k]
  return formatter ? formatter(v) : `${k}: ${String(v)}`
}


export function AuditPage() {
  const [response, setResponse] = useState<AuthEventsResponse | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [filterType, setFilterType] = useState('')
  const [filterEmail, setFilterEmail] = useState('')
  const [offset, setOffset] = useState(0)

  const load = useCallback(async (off = 0) => {
    setIsLoading(true)
    try {
      const data = await authService.getEvents({
        event_type: filterType || undefined,
        email: filterEmail || undefined,
        limit: PAGE_SIZE,
        offset: off,
      })
      setResponse(data)
      setOffset(off)
    } catch {
      toast.error('Erro ao carregar eventos.')
    } finally {
      setIsLoading(false)
    }
  }, [filterType, filterEmail])

  useEffect(() => {
    load(0)
  }, [load])

  const events: AuthEvent[] = response?.events ?? []
  const total = response?.total ?? 0
  const totalPages = Math.ceil(total / PAGE_SIZE)
  const currentPage = Math.floor(offset / PAGE_SIZE) + 1

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
          <Activity className="w-6 h-6" />
          Auditoria de Acessos
        </h1>
        <Button variant="outline" size="sm" onClick={() => load(offset)} disabled={isLoading}>
          <RefreshCw className={`w-4 h-4 mr-1.5 ${isLoading ? 'animate-spin' : ''}`} />
          Atualizar
        </Button>
      </div>

      {/* Filters */}
      <Card padding="md" className="mb-4">
        <div className="flex flex-wrap gap-3 items-end">
          <div>
            <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Tipo de evento</label>
            <select
              value={filterType}
              onChange={(e) => { setFilterType(e.target.value); load(0) }}
              className="rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-1.5 text-gray-900 dark:text-gray-100"
            >
              <option value="">Todos</option>
              {Object.entries(EVENT_LABELS).map(([key, label]) => (
                <option key={key} value={key}>{label}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Email</label>
            <input
              type="text"
              placeholder="Filtrar por email..."
              value={filterEmail}
              onChange={(e) => setFilterEmail(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && load(0)}
              className="rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-1.5 text-gray-900 dark:text-gray-100 w-60"
            />
          </div>
          <Button size="sm" onClick={() => load(0)} disabled={isLoading}>
            Filtrar
          </Button>
          <span className="text-xs text-gray-400 ml-auto self-center">
            {total} evento{total !== 1 ? 's' : ''} encontrado{total !== 1 ? 's' : ''}
          </span>
        </div>
      </Card>

      {/* Table */}
      <Card padding="none">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700">
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Data / Hora</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Evento</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Sistema</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Usuário / Email</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">IP</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Detalhes</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-gray-700/50">
              {isLoading ? (
                <tr>
                  <td colSpan={6} className="px-4 py-8 text-center text-gray-400">Carregando...</td>
                </tr>
              ) : events.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-4 py-8 text-center text-gray-400">Nenhum evento encontrado.</td>
                </tr>
              ) : (
                events.map((ev) => (
                  <tr key={ev.id} className="hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                    <td className="px-4 py-3 text-gray-600 dark:text-gray-300 whitespace-nowrap font-mono text-xs">
                      {formatDate(ev.created_at)}
                    </td>
                    <td className="px-4 py-3">
                      <Badge variant={EVENT_VARIANTS[ev.event_type] ?? 'default'} size="sm">
                        {EVENT_LABELS[ev.event_type] ?? ev.event_type}
                      </Badge>
                    </td>
                    <td className="px-4 py-3 text-gray-600 dark:text-gray-300 text-xs font-medium">
                      {(ev.metadata as Record<string, unknown>)?.source as string || '—'}
                    </td>
                    <td className="px-4 py-3 text-gray-900 dark:text-gray-100">
                      <div>{ev.user_full_name || '—'}</div>
                      <div className="text-xs text-gray-400">{ev.email || ev.user_email || '—'}</div>
                    </td>
                    <td className="px-4 py-3 text-gray-600 dark:text-gray-300 font-mono text-xs">
                      {ev.ip_address || '—'}
                    </td>
                    <td className="px-4 py-3 text-gray-500 dark:text-gray-400 text-xs">
                      {ev.metadata && Object.keys(ev.metadata).some(k => !META_HIDDEN.has(k)) ? (
                        <div className="space-y-0.5">
                          {Object.entries(ev.metadata)
                            .filter(([k]) => !META_HIDDEN.has(k))
                            .map(([k, v]) => (
                            <div key={k}>{formatMetaEntry(k, v)}</div>
                          ))}
                        </div>
                      ) : '—'}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200 dark:border-gray-700">
            <span className="text-xs text-gray-500 dark:text-gray-400">
              Página {currentPage} de {totalPages}
            </span>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => load(offset - PAGE_SIZE)}
                disabled={offset === 0 || isLoading}
              >
                Anterior
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => load(offset + PAGE_SIZE)}
                disabled={offset + PAGE_SIZE >= total || isLoading}
              >
                Próxima
              </Button>
            </div>
          </div>
        )}
      </Card>
    </div>
  )
}

export default AuditPage
