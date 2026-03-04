import { describe, it, expect, vi, beforeEach } from 'vitest'
import axios from 'axios'

vi.mock('axios')
const mockedAxios = vi.mocked(axios, true)

// Mock the api module that authService imports
vi.mock('./api', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn(),
  },
}))

import { authService } from './authService'
import api from './api'

const mockedApi = vi.mocked(api)

describe('authService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  describe('requestOTP', () => {
    it('posts to /auth/request-otp with email', async () => {
      const mockResponse = {
        message: 'OTP sent',
        expires_in: 600,
        channel: 'email',
      }
      mockedApi.post.mockResolvedValueOnce({ data: mockResponse })

      const result = await authService.requestOTP('user@example.com')

      expect(mockedApi.post).toHaveBeenCalledWith('/auth/request-otp', { email: 'user@example.com' })
      expect(result).toEqual(mockResponse)
    })

    it('propagates error when API call fails', async () => {
      mockedApi.post.mockRejectedValueOnce(new Error('Network error'))

      await expect(authService.requestOTP('user@example.com')).rejects.toThrow('Network error')
    })
  })

  describe('verifyOTP', () => {
    it('posts to /auth/verify-otp with email and code', async () => {
      const mockResponse = {
        access_token: 'access.jwt',
        refresh_token: 'refresh-token',
        expires_in: 900,
        user: { id: 'uuid', email: 'user@example.com' },
        roles: ['user'],
        permissions: {},
      }
      mockedApi.post.mockResolvedValueOnce({ data: mockResponse })

      const result = await authService.verifyOTP('user@example.com', '123456')

      expect(mockedApi.post).toHaveBeenCalledWith('/auth/verify-otp', {
        email: 'user@example.com',
        code: '123456',
      })
      expect(result.access_token).toBe('access.jwt')
    })

    it('propagates error for invalid OTP', async () => {
      mockedApi.post.mockRejectedValueOnce({ response: { status: 401, data: { error: 'invalid OTP' } } })

      await expect(authService.verifyOTP('user@example.com', '000000')).rejects.toBeDefined()
    })
  })

  describe('adminLogin', () => {
    it('posts to /auth/admin-login with identifier and password', async () => {
      const mockResponse = {
        access_token: 'admin.access.jwt',
        refresh_token: 'admin-refresh',
        expires_in: 900,
        user: { id: 'admin-uuid', email: 'admin@example.com' },
        roles: ['admin'],
        permissions: { users: ['read', 'write'] },
      }
      mockedApi.post.mockResolvedValueOnce({ data: mockResponse })

      const result = await authService.adminLogin('admin@example.com', 'password123')

      expect(mockedApi.post).toHaveBeenCalledWith('/auth/admin-login', {
        identifier: 'admin@example.com',
        password: 'password123',
      })
      expect(result.roles).toContain('admin')
    })
  })

  describe('refreshTokens', () => {
    it('reads refresh_token from localStorage and posts to /auth/refresh', async () => {
      localStorage.setItem('refresh_token', 'stored-refresh-token')
      const mockResponse = {
        access_token: 'new.access.jwt',
        refresh_token: 'new-refresh-token',
        expires_in: 900,
        user: { id: 'uuid', email: 'user@example.com' },
        roles: [],
        permissions: {},
      }
      mockedApi.post.mockResolvedValueOnce({ data: mockResponse })

      const result = await authService.refreshTokens()

      expect(mockedApi.post).toHaveBeenCalledWith('/auth/refresh', {
        refresh_token: 'stored-refresh-token',
      })
      expect(result.access_token).toBe('new.access.jwt')
    })

    it('passes null refresh_token when not in localStorage', async () => {
      mockedApi.post.mockResolvedValueOnce({ data: {} })

      await authService.refreshTokens()

      expect(mockedApi.post).toHaveBeenCalledWith('/auth/refresh', {
        refresh_token: null,
      })
    })
  })

  describe('logout', () => {
    it('posts to /auth/logout with refresh token from localStorage', async () => {
      localStorage.setItem('refresh_token', 'my-refresh-token')
      mockedApi.post.mockResolvedValueOnce({ data: {} })

      await authService.logout()

      expect(mockedApi.post).toHaveBeenCalledWith('/auth/logout', {
        refresh_token: 'my-refresh-token',
      })
    })
  })

  describe('getMe', () => {
    it('gets from /auth/me and returns user', async () => {
      const mockUser = {
        id: 'user-uuid',
        email: 'user@example.com',
        full_name: 'Test User',
        is_active: true,
        tenant_id: 'tenant-uuid',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      mockedApi.get.mockResolvedValueOnce({ data: mockUser })

      const result = await authService.getMe()

      expect(mockedApi.get).toHaveBeenCalledWith('/auth/me')
      expect(result.email).toBe('user@example.com')
    })
  })

  describe('getEvents', () => {
    it('gets from /admin/events with default params', async () => {
      const mockEvents = {
        events: [],
        total: 0,
        limit: 50,
        offset: 0,
      }
      mockedApi.get.mockResolvedValueOnce({ data: mockEvents })

      const result = await authService.getEvents()

      expect(mockedApi.get).toHaveBeenCalledWith('/admin/events', { params: undefined })
      expect(result.total).toBe(0)
    })

    it('gets events with filter params', async () => {
      const mockEvents = { events: [], total: 0, limit: 10, offset: 0 }
      mockedApi.get.mockResolvedValueOnce({ data: mockEvents })

      await authService.getEvents({ event_type: 'login_success', limit: 10, offset: 0 })

      expect(mockedApi.get).toHaveBeenCalledWith('/admin/events', {
        params: { event_type: 'login_success', limit: 10, offset: 0 },
      })
    })
  })
})
