import { defineStore } from 'pinia'
import client from '@/api/client'

export const useExamStore = defineStore('exam', {
  state: () => ({
    participantId: localStorage.getItem('participantId') || null,
    sessionId: null,
    modules: [] as any[],
    currentModuleIndex: 0,
    questions: [] as any[],
    answers: {} as Record<string, any>,
    timeRemaining: 0,
    isSubmitting: false,
  }),
  actions: {
    async joinExam(payload: { token: string, name: string, email: string, age: number, position: string }) {
      const res = await client.post('/exam/join', payload)
      this.participantId = res.data.participant_id
      this.sessionId = res.data.session_id
      localStorage.setItem('participantId', res.data.participant_id)
      return res.data
    },
    async fetchModulesAndStart() {
      if (!this.sessionId) return
      try {
        const res = await client.get(`/exam/${this.sessionId}/modules`)
        this.modules = res.data || []
        this.currentModuleIndex = 0
        if (this.modules.length > 0) {
          await this.fetchQuestionsForCurrentModule()
        }
      } catch (e) {
        console.error("Failed to fetch modules:", e)
        this.modules = []
        this.questions = []
      }
    },
    async fetchQuestionsForCurrentModule() {
      if (!this.sessionId || this.modules.length === 0) return
      this.questions = [] // clear current questions for loading state
      const mId = this.modules[this.currentModuleIndex].module_id
      try {
        const res = await client.get(`/exam/${this.sessionId}/modules/${mId}/questions`)
        this.questions = res.data || []
      } catch (e) {
        console.error("Failed to fetch questions for module:", e)
        this.questions = []
      }
    },
    setAnswer(questionId: string, answer: any) {
      this.answers[questionId] = answer
    },
    async submitExam() {
      if (this.isSubmitting || !this.participantId) return
      this.isSubmitting = true
      try {
        const payload = {
          participant_id: this.participantId,
          answers: Object.entries(this.answers).map(([qId, ans]: [string, any]) => ({
            question_id: qId,
            selected_option_id: ans.selected_option_id || null,
            text_answer: ans.text_answer || null
          }))
        }
        await client.post(`/exam/${this.sessionId}/answers`, payload)
        this.clearExam()
        return true
      } catch (e) {
        console.error("Failed to submit", e)
        return false
      } finally {
        this.isSubmitting = false
      }
    },
    async reportViolation(type: string) {
      if (!this.participantId || !this.sessionId) return
      await client.post('/violations', {
        participant_id: this.participantId,
        session_id: this.sessionId,
        violation_type: type
      }).catch(e => console.error("Failed to report violation", e))
    },
    clearExam() {
      this.participantId = null
      this.sessionId = null
      this.modules = []
      this.currentModuleIndex = 0
      this.questions = []
      this.answers = {}
      localStorage.removeItem('participantId')
    }
  }
})
