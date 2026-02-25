import api from './api'
import type { User, UserCreateRequest, UserUpdateRequest } from '@/types/auth'

export const userService = {
  async list(): Promise<User[]> {
    const { data } = await api.get<{ users: User[]; count: number }>('/users')
    return data.users || []
  },

  async getById(id: string): Promise<User> {
    const { data } = await api.get<User>(`/users/${id}`)
    return data
  },

  async create(payload: UserCreateRequest): Promise<User> {
    const { data } = await api.post<User>('/users', payload)
    return data
  },

  async update(id: string, payload: UserUpdateRequest): Promise<User> {
    const { data } = await api.put<User>(`/users/${id}`, payload)
    return data
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/users/${id}`)
  },
}

export default userService
