import { defineStore } from 'pinia'
import client from '@/api/client'

export const useExamStore = defineStore('exam', {
  state: () => ({
    participantId: localStorage.getItem('participantId') || null,
    sessionId: null,
    questions: [],
    answers: {} as Record<string, any>,
    timeRemaining: 0,
    isSubmitting: false,
  }),
  actions: {
    async joinExam(token: string) {
      const res = await client.post('/exam/join', { token })
      this.participantId = res.data.participant_id
      this.sessionId = res.data.session_id
      localStorage.setItem('participantId', res.data.participant_id)
      return res.data
    },
    async fetchQuestions() {
      if (!this.sessionId) return
      const res = await client.get(`/exam/${this.sessionId}/questions`)
      this.questions = res.data
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
      this.questions = []
      this.answers = {}
      localStorage.removeItem('participantId')
    }
  }
})
