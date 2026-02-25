import api from './api'
import type {
  Service,
  ServiceCreateRequest,
  ServiceUpdateRequest,
  Permission,
  PermissionCreateRequest,
  Role,
  RoleCreateRequest,
  RoleUpdateRequest,
  RolePermissionsRequest,
} from '@/types/auth'

// Services CRUD
export const adminService = {
  // Services
  async listServices(): Promise<Service[]> {
    const { data } = await api.get('/admin/services')
    return data
  },

  async getService(id: string): Promise<Service> {
    const { data } = await api.get(`/admin/services/${id}`)
    return data
  },

  async createService(req: ServiceCreateRequest): Promise<Service> {
    const { data } = await api.post('/admin/services', req)
    return data
  },

  async updateService(id: string, req: ServiceUpdateRequest): Promise<Service> {
    const { data } = await api.put(`/admin/services/${id}`, req)
    return data
  },

  async deleteService(id: string): Promise<void> {
    await api.delete(`/admin/services/${id}`)
  },

  // Permissions
  async listAllPermissions(): Promise<Permission[]> {
    const { data } = await api.get('/admin/permissions')
    return data
  },

  async listServicePermissions(serviceId: string): Promise<Permission[]> {
    const { data } = await api.get(`/admin/services/${serviceId}/permissions`)
    return data
  },

  async createServicePermission(serviceId: string, req: PermissionCreateRequest): Promise<Permission> {
    const { data } = await api.post(`/admin/services/${serviceId}/permissions`, req)
    return data
  },

  async deletePermission(id: string): Promise<void> {
    await api.delete(`/admin/permissions/${id}`)
  },

  // Roles
  async listRoles(): Promise<Role[]> {
    const { data } = await api.get('/admin/roles')
    return data
  },

  async createRole(req: RoleCreateRequest): Promise<Role> {
    const { data } = await api.post('/admin/roles', req)
    return data
  },

  async updateRole(id: string, req: RoleUpdateRequest): Promise<Role> {
    const { data } = await api.put(`/admin/roles/${id}`, req)
    return data
  },

  async getRolePermissions(roleId: string): Promise<string[]> {
    const { data } = await api.get(`/admin/roles/${roleId}/permissions`)
    return data.permission_ids || []
  },

  async setRolePermissions(roleId: string, req: RolePermissionsRequest): Promise<void> {
    await api.put(`/admin/roles/${roleId}/permissions`, req)
  },
}

export default adminService
