import api from './api'
import type { AuthResponse, OTPResponse, User, AuthEventsResponse } from '@/types/auth'

export const authService = {
  async requestOTP(email: string): Promise<OTPResponse> {
    const { data } = await api.post<OTPResponse>('/auth/request-otp', { email })
    return data
  },

  async verifyOTP(email: string, code: string): Promise<AuthResponse> {
    const { data } = await api.post<AuthResponse>('/auth/verify-otp', { email, code })
    return data
  },

  async adminLogin(identifier: string, password: string): Promise<AuthResponse> {
    const { data } = await api.post<AuthResponse>('/auth/admin-login', { identifier, password })
    return data
  },

  async refreshTokens(): Promise<AuthResponse> {
    const refreshToken = localStorage.getItem('refresh_token')
    const { data } = await api.post<AuthResponse>('/auth/refresh', {
      refresh_token: refreshToken,
    })
    return data
  },

  async logout(): Promise<void> {
    const refreshToken = localStorage.getItem('refresh_token')
    await api.post('/auth/logout', { refresh_token: refreshToken })
  },

  async getMe(): Promise<User> {
    const { data } = await api.get<User>('/auth/me')
    return data
  },

  async getEvents(params?: { event_type?: string; email?: string; limit?: number; offset?: number }): Promise<AuthEventsResponse> {
    const { data } = await api.get<AuthEventsResponse>('/admin/events', { params })
    return data
  },
}

export default authService
