import { defineStore } from 'pinia'
import client from '@/api/client'

export const useSessionStore = defineStore('session', {
  state: () => ({
    sessions: [] as any[],
    currentSession: null as any,
  }),
  actions: {
    async fetchSessions() {
      const res = await client.get('/sessions')
      this.sessions = res.data
    },
    async createSession(payload: any) {
      await client.post('/sessions', payload)
      await this.fetchSessions()
    },
    async updateSession(id: string, payload: any) {
      await client.put(`/sessions/${id}`, payload)
      await this.fetchSessions()
    },
    async deleteSession(id: string) {
      await client.delete(`/sessions/${id}`)
      await this.fetchSessions()
    },
    async fetchSessionDetails(id: string) {
      const res = await client.get(`/sessions/${id}`)
      this.currentSession = res.data
    }
  }
})

