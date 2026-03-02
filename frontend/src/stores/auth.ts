import { defineStore } from 'pinia'
import client from '@/api/client'

export interface User {
  id: string
  name: string
  email: string
  role: string | number
  role_name?: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    token: localStorage.getItem('hrd_token') || null,
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    isHR: (state) => state.user?.role_name === 'HR' || state.user?.role_name === 'Super Admin' || state.user?.role === 2 || state.user?.role === 1,
    isAdmin: (state) => state.user?.role_name === 'Super Admin' || state.user?.role === 1,
  },
  actions: {
    setAuth(user: User, token: string) {
      this.user = user
      this.token = token
      localStorage.setItem('hrd_token', token)
    },
    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('hrd_token')
      window.location.href = '/login'
    },
    async fetchUser() {
      if (!this.token) return
      try {
        const res = await client.get('/auth/me')
        this.user = {
          id: res.data.user_id,
          email: res.data.email,
          name: res.data.email, // fallback
          role: res.data.role,
          role_name: res.data.role,
        }
      } catch (err) {
        this.logout()
      }
    }
  }
})
