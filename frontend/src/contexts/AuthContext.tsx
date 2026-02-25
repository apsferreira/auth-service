import React, { createContext, useContext, useEffect, useState, useCallback } from 'react'
import type { User, AuthResponse } from '@/types/auth'
import authService from '@/services/authService'

interface AuthContextType {
  user: User | null
  permissions: Record<string, string[]>
  roles: string[]
  isAuthenticated: boolean
  isLoading: boolean
  login: (identifier: string, password: string) => Promise<void>
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
  hasRole: (role: string) => boolean
  hasPermission: (permission: string) => boolean
  hasServicePermission: (serviceSlug: string, permission: string) => boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [permissions, setPermissions] = useState<Record<string, string[]>>({})
  const [roles, setRoles] = useState<string[]>([])
  const [isLoading, setIsLoading] = useState(true)

  const isAuthenticated = !!user

  const refreshUser = useCallback(async () => {
    try {
      const userData = await authService.getMe()
      setUser(userData)
      // Restore roles and permissions from JWT payload (they are not returned by /auth/me)
      const token = localStorage.getItem('access_token')
      if (token) {
        try {
          const payload = JSON.parse(atob(token.split('.')[1]))
          setRoles(payload.roles || [])
          setPermissions(payload.permissions || {})
        } catch {
          // ignore JWT decode errors
        }
      }
    } catch {
      setUser(null)
      setPermissions({})
      setRoles([])
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    }
  }, [])

  useEffect(() => {
    const token = localStorage.getItem('access_token')
    if (token) {
      refreshUser().finally(() => setIsLoading(false))
    } else {
      setIsLoading(false)
    }
  }, [refreshUser])

  const handleAuthResponse = (response: AuthResponse) => {
    localStorage.setItem('access_token', response.access_token)
    localStorage.setItem('refresh_token', response.refresh_token)
    setUser(response.user)
    setRoles(response.roles || [])
    setPermissions(response.permissions || {})
  }

  const login = async (identifier: string, password: string) => {
    const response = await authService.adminLogin(identifier, password)
    handleAuthResponse(response)
  }

  const logout = async () => {
    try {
      await authService.logout()
    } finally {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      setUser(null)
      setPermissions({})
      setRoles([])
    }
  }

  const hasRole = (role: string) => {
    if (roles.includes(role)) return true
    // Also check user.roles if they're Role objects
    return user?.roles?.some((r) => r.name === role) ?? false
  }

  const hasPermission = (permission: string) => {
    // Check across all services
    return Object.values(permissions).some((perms) => perms.includes(permission))
  }

  const hasServicePermission = (serviceSlug: string, permission: string) => {
    return permissions[serviceSlug]?.includes(permission) ?? false
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        permissions,
        roles,
        isAuthenticated,
        isLoading,
        login,
        logout,
        refreshUser,
        hasRole,
        hasPermission,
        hasServicePermission,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
