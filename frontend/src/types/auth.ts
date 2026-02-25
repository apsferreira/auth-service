export interface User {
  id: string
  tenant_id: string
  email: string
  full_name: string
  avatar_url: string | null
  is_active: boolean
  role_id: string
  roles: Role[]
  created_at: string
  updated_at: string
  last_login_at: string | null
}

export interface Role {
  id: string
  name: string
  description: string
  level: number
  is_system: boolean
  tenant_id?: string
  permissions?: Permission[]
}

export interface Permission {
  id: string
  name: string
  resource: string
  action: string
  description: string
  service_id?: string
  service_slug?: string
}

export interface Service {
  id: string
  tenant_id: string
  name: string
  slug: string
  description: string
  redirect_urls: string[]
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
  expires_in: number
  roles: string[]
  permissions: Record<string, string[]>
}

export interface OTPResponse {
  message: string
  expires_in: number
}

export interface OTPRequest {
  email: string
}

export interface OTPVerifyRequest {
  email: string
  code: string
}

export interface RefreshRequest {
  refresh_token: string
}

export interface ValidateResponse {
  valid: boolean
  user_id: string
  tenant_id: string
  email: string
  roles: string[]
  permissions: Record<string, string[]>
}

export interface UserCreateRequest {
  email: string
  full_name: string
  role_id?: string
  tenant_id: string
}

export interface UserUpdateRequest {
  full_name?: string
  is_active?: boolean
  role_id?: string
  avatar_url?: string
}

export interface ServiceCreateRequest {
  name: string
  slug: string
  description?: string
  redirect_urls?: string[]
}

export interface ServiceUpdateRequest {
  name?: string
  description?: string
  redirect_urls?: string[]
  is_active?: boolean
}

export interface PermissionCreateRequest {
  name: string
  resource: string
  action: string
  description?: string
}

export interface RoleCreateRequest {
  name: string
  description?: string
  level: number
}

export interface RoleUpdateRequest {
  name?: string
  description?: string
  level?: number
}

export interface RolePermissionsRequest {
  permission_ids: string[]
}

export interface AuthEvent {
  id: string
  event_type: string
  user_id?: string
  email?: string
  ip_address?: string
  user_agent?: string
  metadata?: Record<string, unknown>
  created_at: string
  user_email?: string
  user_full_name?: string
}

export interface AuthEventsResponse {
  events: AuthEvent[]
  total: number
  limit: number
  offset: number
}
